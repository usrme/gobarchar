package gobarchar

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultTmpl))
}

var tpl *template.Template

var defaultTmpl string = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>{{.Title}}</title>
</head>
<style>
	* {
		box-sizing: border-box;
	}
	body {
		font-family: sans-serif;
		line-height: 1.33;
		margin: 0 auto;
		max-width: 650px;
		padding: 1rem;
	}
	pre {
		overflow: auto;
		user-select: all;
	}
	a {
		word-break: break-all;
	}
</style>
<body>
	<h1>GoBarChar</h1>
	<p><strong>The graphing solution that might not suit you 📊</strong></p>
	<hr />
	<p>What is this? This is a small <a href="https://github.com/usrme/gobarchar">project</a> to generate ASCII bar charts using just query parameters.</p>
	<br/>
	<pre>{{.Chart}}</pre>
	<br/>
	<hr>
	<p>Link used to generate the current chart: <a href="{{.ChartUrl}}">{{.ChartUrl}}</a></p>
	<hr>
</body>
<footer>
	<small>Hosted on <a href="https://fly.io">Fly</a>.</small>
</footer>
</html>`

type pageData struct {
	Title    string
	Chart    string
	ChartUrl string
}

type entry struct {
	Label string
	Value int
}

// Create custom type that implements 'sort.Interface' for easier sorting
type chartData []entry

func (a chartData) Len() int           { return len(a) }
func (a chartData) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a chartData) Less(i, j int) bool { return a[i].Value < a[j].Value }

func PresentBarChart(w http.ResponseWriter, r *http.Request) {
	// Check if there are any query parameters; if not, add random key-value pairs
	if len(r.URL.Query()) == 0 {
		encodeRandomQuery(r)
	}

	chart := createBarChart(r)

	// Skip all templating if user is requesting through 'curl' or 'wget'
	agent := r.UserAgent()
	if strings.HasPrefix(agent, "curl") || strings.HasPrefix(agent, "Wget") {
		w.Write([]byte(chart))
		return
	}

	pageData := pageData{
		Title:    "GoBarChar",
		Chart:    chart,
		ChartUrl: r.URL.String(),
	}

	err := tpl.Execute(w, pageData)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "something went wrong...", http.StatusInternalServerError)
	}
}

func encodeRandomQuery(r *http.Request) {
	params := r.URL.Query()

	for i := 0; i < 6; i++ {
		var month string
		for {
			month = randomMonth()
			if !params.Has(month) {
				break
			}
		}
		value := rand.Intn(101)
		params.Add(month, strconv.Itoa(value))
	}
	r.URL.RawQuery = params.Encode()
}

func randomMonth() string {
	months := []string{
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	}
	return months[rand.Intn(len(months))]
}

func createBarChart(r *http.Request) string {
	// Use custom type to make sorting easier
	entries := make([]entry, 0)
	maxValue, total := 0, 0
	// Use 'RawQuery' instead of something like 'Query()' or other automatic
	// methods to parse query parameters as those get turned into maps that
	// then lose their order of appearance in the original query.
	//
	// The order make intuitive sense when presenting, otherwise every request
	// with the same data shows the results with a slightly different order if
	// the 'sort' query parameter isn't given.
	raw := r.URL.RawQuery
	orderedParams := strings.Split(raw, "&")

	sortIndex := -1
	for i, pair := range orderedParams {
		kv := strings.Split(pair, "=")
		// TODO: Add logic for 'layout' query parameter to support vertical layout too
		if kv[0] == "sort" {
			// Store where the 'sort' query parameter was first found
			sortIndex = i
			// Don't parse 'sort' query parameter as that is not part of the data
			continue
		}
		count, err := strconv.Atoi(kv[1])
		if err == nil {
			entries = append(entries, entry{Label: kv[0], Value: count})
			if count > maxValue {
				maxValue = count
			}
			total += count
		}
	}

	// Only remove the 'sort' query parameter if it was found.
	//
	// This is at O(n) complexity, but I'm more interested in keeping the order.
	if sortIndex != -1 {
		orderedParams = append(orderedParams[:sortIndex], orderedParams[sortIndex+1:]...)
	}

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = ""
	}

	if sortOrder == "asc" {
		sort.Sort(chartData(entries))
	} else if sortOrder == "desc" {
		sort.Sort(sort.Reverse(chartData(entries)))
	}

	avg := math.Round(float64(total) / float64(len(entries)))

	entries = append(entries, entry{Label: "Avg.", Value: int(avg)})
	orderedParams = append(orderedParams, fmt.Sprintf("%s=%d", "Avg.", int(avg)))

	entries = append(entries, entry{Label: "Total", Value: total})
	orderedParams = append(orderedParams, fmt.Sprintf("%s=%d", "Total", total))

	increment := float64(maxValue) / 25.0

	// Find the longest label to determine padding later on
	longestLabelLength := 0
	for _, entry := range entries {
		if len(entry.Label) > longestLabelLength {
			longestLabelLength = len(entry.Label)
		}
	}

	var chartContent strings.Builder
	maximumBarChunk := 0
	for i := range orderedParams {
		// Skip parsing the total for now to not interfere with calculating
		// the maximum number of bar chunks for all of the labels
		if entries[i].Label == "Total" {
			continue
		}
		barChunks := int(float64(entries[i].Value) * 8 / increment)
		remainder := barChunks % 8
		barChunks /= 8

		if barChunks > maximumBarChunk {
			maximumBarChunk = barChunks
		}

		bar := calculateBars(barChunks, remainder)
		chartContent.WriteString(
			fmt.Sprintf(
				"%s %4d %s\n",
				padRight(entries[i].Label, longestLabelLength), entries[i].Value, bar,
			),
		)
	}
	bar := calculateBars(maximumBarChunk, 0)
	chartContent.WriteString(fmt.Sprintf("%s %4d %s\n", padRight("Total", longestLabelLength), total, bar))
	return chartContent.String()
}

func calculateBars(count, remainder int) string {
	// First draw the full width chunks
	bar := strings.Repeat("█", count)

	// Then add the fractional part
	if remainder > 0 {
		bar += string(rune('█' + (8 - remainder)))
	}

	// If the bar is empty (i.e. a value of 0 was given), add a left one-eighth block
	if bar == "" {
		bar = "▏"
	}

	return bar
}

func padRight(str string, length int) string {
	if len(str) >= length {
		return str
	}
	return str + strings.Repeat(" ", length-len(str))
}
