package gobarchar

import (
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

var htmlFirstHalf string = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="icon" href="data:,">
	<title>GoBarChar</title>
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
	h1 > a, h1 > a:visited {
		color: black;
	}
</style>
</head>
<body>
	<h1><a href="/">GoBarChar</a></h1>
	<p><strong>The charting solution that might not suit you 📊</strong></p>
	<hr />
	<p>What is this? This is a small <a href="https://github.com/usrme/gobarchar">project</a> to generate ASCII bar charts using just query parameters.</p>
	<br/>
`

var htmlSecondHalf string = `<hr>
<footer>
	<small>Hosted on <a href="https://fly.io">Fly</a>.</small>
</footer>
</body>
</html>`

type entry struct {
	Label string
	Value float64
}

// Create custom type that implements 'sort.Interface' for easier sorting
type chartData []entry

func (a chartData) Len() int           { return len(a) }
func (a chartData) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a chartData) Less(i, j int) bool { return a[i].Value < a[j].Value }

func PresentBarChart(examples string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		chartUrl := r.URL.String()
		html := fmt.Sprintf(
			"%s<pre>%s</pre><br/><hr><p>Link used to generate the current chart: <a href='%s'>%s</a></p>%s%s",
			htmlFirstHalf, chart, chartUrl, chartUrl, examples, htmlSecondHalf,
		)
		w.Write([]byte(html))
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
	maxValue, total := 0.0, 0.0
	// Use 'RawQuery' instead of something like 'Query()' or other automatic
	// methods to parse query parameters as those get turned into maps that
	// then lose their order of appearance in the original query.
	//
	// The order make intuitive sense when presenting, otherwise every request
	// with the same data shows the results with a slightly different order if
	// the 'sort' query parameter isn't given.
	raw := r.URL.RawQuery
	orderedParams := strings.Split(raw, "&")
	addSpaces := slices.Contains(orderedParams, "spaces=yes")
	for _, pair := range orderedParams {
		kv := strings.Split(pair, "=")
		key := kv[0]
		// Don't parse certain query parameters as they are not part of the data
		if key == "sort" || key == "spaces" {
			continue
		}
		count, err := strconv.ParseFloat(kv[1], 64)
		if err == nil {
			if addSpaces {
				key = strings.Replace(key, "%20", " ", -1)
			}

			entries = append(entries, entry{Label: key, Value: float64(count)})
			if count > maxValue {
				maxValue = count
			}
			total += count
		}
	}

	// Only remove the 'sort' query parameter if it was found.
	sortOrder := r.URL.Query().Get("sort")

	if sortOrder == "asc" {
		sort.Sort(chartData(entries))
	} else if sortOrder == "desc" {
		sort.Sort(sort.Reverse(chartData(entries)))
	}

	avg := total / float64(len(entries))

	entries = append(entries, entry{Label: "Avg.", Value: avg})
	orderedParams = append(orderedParams, fmt.Sprintf("%s=%.2f", "Avg.", avg))

	entries = append(entries, entry{Label: "Total", Value: total})
	orderedParams = append(orderedParams, fmt.Sprintf("%s=%.2f", "Total", total))

	increment := maxValue / 25.0

	// Find the longest label to determine padding later on
	longestLabelLength := 0
	for _, entry := range entries {
		if len(entry.Label) > longestLabelLength {
			longestLabelLength = len(entry.Label)
		}
	}

	longestValueLength := 0
	for _, entry := range entries {
		valueStr := formatValue(entry.Value)
		if len(valueStr) > longestValueLength {
			longestValueLength = len(valueStr)
		}
	}

	var chartContent strings.Builder
	maximumBarChunk := 0
	for i := range entries {
		// Skip parsing the total for now to not interfere with calculating
		// the maximum number of bar chunks for all of the labels
		if entries[i].Label == "Total" {
			continue
		}
		barChunks := int(entries[i].Value * 8 / increment)
		remainder := barChunks % 8
		barChunks /= 8

		if barChunks > maximumBarChunk {
			maximumBarChunk = barChunks
		}

		bar := calculateBars(barChunks, remainder)
		valueStr := formatValue(entries[i].Value)
		chartContent.WriteString(
			fmt.Sprintf(
				"%s %s %s\n",
				padRight(entries[i].Label, longestLabelLength), padLeft(valueStr, longestValueLength), bar,
			),
		)
	}
	bar := calculateBars(maximumBarChunk, 0)
	totalStr := formatValue(total)
	chartContent.WriteString(fmt.Sprintf("%s %s %s\n", padRight("Total", longestLabelLength), padLeft(totalStr, longestValueLength), bar))
	return chartContent.String()
}

func formatValue(value float64) string {
	if value == float64(int(value)) {
		return fmt.Sprintf("%4d", int(value))
	}
	return fmt.Sprintf("%6.2f", value)
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

func padLeft(str string, length int) string {
	if len(str) >= length {
		return str
	}
	return strings.Repeat(" ", length-len(str)) + str
}
