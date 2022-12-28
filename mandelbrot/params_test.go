package mandelbrot

import (
	"fmt"
	"testing"
)

func TestGetQueryParams(t *testing.T) {

	var tests = []struct {
		query map[string][]string
		want  params
	}{
		{
			map[string][]string{"width": {"100"}, "height": {"75"}, "iterations": {"500"}},
			params{width: 100, height: 75, maxIterations: 500},
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.query)
		t.Run(testname, func(t *testing.T) {
			result := GetQueryParams(tt.query)
			if *result != tt.want {
				t.Errorf("Got %v, want %v", *result, tt.want)
			}
		})
	}
}
