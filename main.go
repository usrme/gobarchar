package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var PageTitle string = "GoBarChar"

type PageData struct {
	Title    string
	Chart    template.HTML
	ChartUrl string
}

type Entry struct {
	Label string
	Value int
}

type ByValue []Entry

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByValue) Less(i, j int) bool { return a[i].Value < a[j].Value }

func generateBarChartContent(r *http.Request) string {
	// Use custom type to make sorting easier
	data := make([]Entry, 0)
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
			data = append(data, Entry{Label: kv[0], Value: count})
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
		sort.Sort(ByValue(data))
	} else if sortOrder == "desc" {
		sort.Sort(sort.Reverse(ByValue(data)))
	}

	avg := math.Round(float64(total) / float64(len(data)))

	data = append(data, Entry{Label: "Avg.", Value: int(avg)})
	orderedParams = append(orderedParams, fmt.Sprintf("%s=%d", "Avg.", int(avg)))

	data = append(data, Entry{Label: "Total", Value: total})
	orderedParams = append(orderedParams, fmt.Sprintf("%s=%d", "Total", total))

	increment := float64(maxValue) / 25.0

	// Find the longest label to determine padding later on
	longestLabelLength := 0
	for _, entry := range data {
		if len(entry.Label) > longestLabelLength {
			longestLabelLength = len(entry.Label)
		}
	}

	var chartContent strings.Builder
	maximumBarChunk := 0
	for i := range orderedParams {
		// Skip parsing the total for now to not interfere with calculating
		// the maximum number of bar chunks for all of the labels
		if data[i].Label == "Total" {
			continue
		}
		barChunks := int(float64(data[i].Value) * 8 / increment)
		remainder := barChunks % 8
		barChunks /= 8

		if barChunks > maximumBarChunk {
			maximumBarChunk = barChunks
		}

		bar := calculateBars(barChunks, remainder)
		chartContent.WriteString(fmt.Sprintf("%s %4d %s\n", padRight(data[i].Label, longestLabelLength), data[i].Value, bar))
	}
	bar := calculateBars(maximumBarChunk, 0)
	chartContent.WriteString(fmt.Sprintf("%s %4d %s\n", padRight("Total", longestLabelLength), total, bar))
	return chartContent.String()
}

func generateBarChart(w http.ResponseWriter, r *http.Request) {
	// Check if there are any query parameters; if not, add random key-value pairs
	if len(r.URL.Query()) == 0 {
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

	chartContent := generateBarChartContent(r)

	agent := r.UserAgent()
	if strings.HasPrefix(agent, "curl") || strings.HasPrefix(agent, "Wget") {
		w.Write([]byte(chartContent))
		return
	}

	pageData := PageData{
		Title:    PageTitle,
		Chart:    template.HTML(chartContent),
		ChartUrl: r.URL.String(),
	}

	tmpl, err := template.New("layout").Parse(layout)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("error parsing template:", err)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("error executing template:", err)
		return
	}
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

func timer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Now().Sub(startTime)
		log.Println("completed in:", duration)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on:", port)
	http.Handle("/", timer(http.HandlerFunc(generateBarChart)))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var layout string = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>{{.Title}}</title>
</head>
<style>
	pre {
		user-select: all;
	}
</style>
<body>
	<pre>{{.Chart}}</pre>
	<hr>
	<p>What is this? This is a small <a href="https://github.com/usrme/gobarchar">project</a> to generate ASCII bar charts using just query parameters.</p>
	<p>Link used to generate the current chart: <a href="{{.ChartUrl}}">{{.ChartUrl}}</p>
</body>
</html>`
