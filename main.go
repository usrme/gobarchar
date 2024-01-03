package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	Label string
	Value int
}

type ByValue []Entry

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByValue) Less(i, j int) bool { return a[i].Value < a[j].Value }

func generateBarChart(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintf(w, "%s %4d %s\n", padRight(data[i].Label, longestLabelLength), data[i].Value, bar)
	}
	bar := calculateBars(maximumBarChunk, 0)
	fmt.Fprintf(w, "%s %4d %s\n", padRight("Total", longestLabelLength), total, bar)
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
