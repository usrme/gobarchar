package gobarchar

import (
	"fmt"
	"math"
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
	<p><strong>The graphing solution that might not suit you 📊</strong></p>
	<hr />
	<p>What is this? This is a small <a href="https://github.com/usrme/gobarchar">project</a> to generate ASCII bar charts using just query parameters.</p>
	<br/>
`

var exampleQueries = map[string]string{
	"List of presidents of the United States by age at start of presidency": "George%20Washington=57&John%20Adams=61&Thomas%20Jefferson=57&James%20Madison=57&James%20Monroe=58&John%20Quincy%20Adams=57&Andrew%20Jackson=61&Martin%20Van%20Buren=54&William%20Henry%20Harrison=68&John%20Tyler=51&James%20K.%20Polk=49&Zachary%20Taylor=64&Millard%20Fillmore=50&Franklin%20Pierce=48&James%20Buchanan=65&Abraham%20Lincoln=52&Andrew%20Johnson=56&Ulysses%20S.%20Grant=46&Rutherford%20B.%20Hayes=54&James%20A.%20Garfield=49&Chester%20A.%20Arthur=51&Grover%20Cleveland=55&Benjamin%20Harrison=55&William%20McKinley=54&Theodore%20Roosevelt=42&William%20Howard%20Taft=51&Woodrow%20Wilson=56&Warren%20G.%20Harding=55&Calvin%20Coolidge=51&Herbert%20Hoover=54&Franklin%20D.%20Roosevelt=51&Harry%20S.%20Truman=60&Dwight%20D.%20Eisenhower=62&John%20F.%20Kennedy=43&Lyndon%20B.%20Johnson=55&Richard%20Nixon=56&Gerald%20Ford=61&Jimmy%20Carter=52&Ronald%20Reagan=69&George%20H.%20W.%20Bush=64&Bill%20Clinton=46&George%20W.%20Bush=54&Barack%20Obama=47&Donald%20Trump=70&Joe%20Biden=78&spaces=yes&sort=asc",
}

var htmlSecondHalf string = `<hr>
<footer>
	<small>Hosted on <a href="https://fly.io">Fly</a>.</small>
</footer>
</body>
</html>`

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

	chartUrl := r.URL.String()
	html := fmt.Sprintf(
		"%s<pre>%s</pre><br/><hr><p>Link used to generate the current chart: <a href='%s'>%s</a></p>%s%s",
		htmlFirstHalf, chart, chartUrl, chartUrl, createListItems("/", exampleQueries), htmlSecondHalf,
	)
	w.Write([]byte(html))
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
	addSpaces := slices.Contains(orderedParams, "spaces=yes")
	for _, pair := range orderedParams {
		kv := strings.Split(pair, "=")
		key := kv[0]
		// Don't parse certain query parameters as they are not part of the data
		if key == "sort" || key == "spaces" {
			continue
		}
		count, err := strconv.Atoi(kv[1])
		if err == nil {
			if addSpaces {
				key = strings.Replace(key, "%20", " ", -1)
			}

			entries = append(entries, entry{Label: key, Value: count})
			if count > maxValue {
				maxValue = count
			}
			total += count
		}
	}

	// Only remove the 'sort' query parameter if it was found.
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder != "" {
		i := slices.Index(orderedParams, fmt.Sprintf("sort=%s", sortOrder))
		orderedParams = slices.Delete(orderedParams, i, i+1)
	}

	if sortOrder == "asc" {
		sort.Sort(chartData(entries))
	} else if sortOrder == "desc" {
		sort.Sort(sort.Reverse(chartData(entries)))
	}

	// Only remove the 'spaces' query parameter if it was found.
	spaces := r.URL.Query().Get("spaces")
	if spaces != "" {
		i := slices.Index(orderedParams, fmt.Sprintf("spaces=%s", spaces))
		orderedParams = slices.Delete(orderedParams, i, i+1)
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

func createListItems(url string, elements map[string]string) string {
	start := `<hr/>
	<p>More example queries:</p>
	<ul>
`
	end := `
	</ul>
`
	var listItems strings.Builder
	for k, v := range elements {
		listItems.WriteString(
			fmt.Sprintf("<li><a href='%s?%s'>%s</a></li>", url, v, k),
		)
	}
	return fmt.Sprintf("%s%s%s", start, listItems.String(), end)
}
