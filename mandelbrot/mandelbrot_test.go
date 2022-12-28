package mandelbrot

import (
	"fmt"
	"math"
	"testing"
)

func TestLinearScale(t *testing.T) {
	// "Close enough" from https://stackoverflow.com/a/47969546
	const delta = 1e-10
	var tests = []struct {
		val, minScaleFrom, maxScaleFrom, minScaleTo, maxScaleTo, want float64
	}{
		{5, 0, 10, 0, 100, 50},
		{50, 0, 100, 0, 1, 0.5},
		{75, 0, 100, 0, 1, 0.75},
		{0, -2, 0.47, 0, 3072, 2487.449392712551},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%f,%f,%f,%f,%f", tt.val, tt.minScaleFrom, tt.maxScaleFrom, tt.minScaleTo, tt.maxScaleTo)
		t.Run(testname, func(t *testing.T) {
			result := linearScale(tt.val, tt.minScaleFrom, tt.maxScaleFrom, tt.minScaleTo, tt.maxScaleTo)
			if math.Abs(tt.want-result) > delta {
				t.Errorf("Got %f, want %f", result, tt.want)
			}
		})
	}
}
