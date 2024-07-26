package gobarchar

import (
	"net/http"
	"strings"
	"testing"
)

func TestGenerateBarChart(t *testing.T) {
	testCases := []struct {
		name        string
		queryParams string
		expected    string
	}{
		{
			name:        "No 'sort' parameter",
			queryParams: "A=10&B=20&C=15",
			expected: strings.TrimSpace(`
A       10 ████████████▌
B       20 █████████████████████████
C       15 ██████████████████▊
Avg.    15 ██████████████████▊
Total   45 █████████████████████████
			`),
		},
		{
			name:        "Sort by ascending",
			queryParams: "A=10&B=20&C=15&sort=asc",
			expected: strings.TrimSpace(`
A       10 ████████████▌
C       15 ██████████████████▊
B       20 █████████████████████████
Avg.    15 ██████████████████▊
Total   45 █████████████████████████
			`),
		},
		{
			name:        "Sort by descending",
			queryParams: "A=10&B=20&C=15&sort=desc",
			expected: strings.TrimSpace(`
B       20 █████████████████████████
C       15 ██████████████████▊
A       10 ████████████▌
Avg.    15 ██████████████████▊
Total   45 █████████████████████████
			`),
		},
		{
			name:        "Invalid value for 'sort' parameter",
			queryParams: "A=10&B=20&C=15&sort=invalid",
			expected: strings.TrimSpace(`
A       10 ████████████▌
B       20 █████████████████████████
C       15 ██████████████████▊
Avg.    15 ██████████████████▊
Total   45 █████████████████████████
			`),
		},
		{
			name:        "Several input values without 'sort' parameter",
			queryParams: "2012=8&2013=6&2014=8&2015=14&2016=8&2017=6&2018=0&2019=24&2020=17&2021=21&2022=17&2023=13",
			expected: strings.TrimSpace(`
2012       8 ████████▎
2013       6 ██████▎
2014       8 ████████▎
2015      14 ██████████████▌
2016       8 ████████▎
2017       6 ██████▎
2018       0 ▏
2019      24 █████████████████████████
2020      17 █████████████████▋
2021      21 █████████████████████▉
2022      17 █████████████████▋
2023      13 █████████████▌
Avg.   11.83 ████████████▎
Total    142 █████████████████████████
			`),
		},
		{
			name:        "Keep HTML entity '%20' when no 'spaces' parameter",
			queryParams: "Year%202024=10&Year%202023=8",
			expected: strings.TrimSpace(`
Year%202024   10 █████████████████████████
Year%202023    8 ████████████████████
Avg.           9 ██████████████████████▌
Total         18 █████████████████████████
			`),
		},
		{
			name:        "Replace HTML entity '%20' with spaces",
			queryParams: "Year%202024=10&Year%202023=8&spaces=yes",
			expected: strings.TrimSpace(`
Year 2024   10 █████████████████████████
Year 2023    8 ████████████████████
Avg.         9 ██████████████████████▌
Total       18 █████████████████████████
			`),
		},
		{
			name:        "Invalid value for 'spaces' parameter",
			queryParams: "Year%202024=10&Year%202023=8&spaces=invalid",
			expected: strings.TrimSpace(`
Year%202024   10 █████████████████████████
Year%202023    8 ████████████████████
Avg.           9 ██████████████████████▌
Total         18 █████████████████████████
			`),
		},
		{
			name:        "Add 'title' parameter with literal spaces",
			queryParams: "A=10&B=20&C=15&title=A descriptive title",
			expected: strings.TrimSpace(`
A descriptive title

A       10 ████████████▌
B       20 █████████████████████████
C       15 ██████████████████▊
Avg.    15 ██████████████████▊
Total   45 █████████████████████████
			`),
		},
		{
			name:        "Add 'title' parameter with HTML entity spaces",
			queryParams: "A=10&B=20&C=15&title=A%20descriptive%20title",
			expected: strings.TrimSpace(`
A descriptive title

A       10 ████████████▌
B       20 █████████████████████████
C       15 ██████████████████▊
Avg.    15 ██████████████████▊
Total   45 █████████████████████████
			`),
		},
		{
			name:        "Add 'title' parameter as first parameter",
			queryParams: "title=A descriptive title&A=10&B=20&C=15",
			expected: strings.TrimSpace(`
A descriptive title

A       10 ████████████▌
B       20 █████████████████████████
C       15 ██████████████████▊
Avg.    15 ██████████████████▊
Total   45 █████████████████████████
			`),
		},
		{
			name:        "Add 'title' parameter with other parameters",
			queryParams: "A=10&B=20&C=15&sort=desc&title=A descriptive title",
			expected: strings.TrimSpace(`
A descriptive title

B       20 █████████████████████████
C       15 ██████████████████▊
A       10 ████████████▌
Avg.    15 ██████████████████▊
Total   45 █████████████████████████
			`),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r, err := http.NewRequest("GET", "/?"+testCase.queryParams, nil)
			if err != nil {
				t.Fatal(err)
			}
			actual := strings.TrimSpace(createBarChart(r))
			if actual != testCase.expected {
				t.Errorf("Unexpected output:\nExpected:\n%s\nGot:\n%s", testCase.expected, actual)
			}
		})
	}
}
