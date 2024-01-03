package main

import (
	"net/http"
	"net/http/httptest"
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
2012     8 ████████▎
2013     6 ██████▎
2014     8 ████████▎
2015    14 ██████████████▌
2016     8 ████████▎
2017     6 ██████▎
2018     0 ▏
2019    24 █████████████████████████
2020    17 █████████████████▋
2021    21 █████████████████████▉
2022    17 █████████████████▋
2023    13 █████████████▌
Avg.    12 ████████████▌
Total  142 █████████████████████████
			`),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/?"+testCase.queryParams, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(generateBarChart)

			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Expected status code 200, got %d", rr.Code)
			}

			actual := strings.TrimSpace(rr.Body.String())
			if actual != testCase.expected {
				t.Errorf("Unexpected output:\nExpected:\n%s\nGot:\n%s", testCase.expected, actual)
			}
		})
	}
}
