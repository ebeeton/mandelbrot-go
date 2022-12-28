// Package mandelbrot provides functionality to plot images and write them to
// HTTP responses.
package mandelbrot

import (
	"net/url"
	"strconv"
)

type params struct {
	width,
	height,
	maxIterations int
}

func newParams(width, height, iterations int) *params {
	p := params{width: width, height: height, maxIterations: iterations}
	return &p
}

// GetQueryParams gets plotting parameters from a map of query parameters.
func GetQueryParams(q url.Values) *params {
	const (
		defaultWidth      int = 1024
		defaultHeight     int = 768
		defaultIterations int = 512
	)

	width, err := strconv.Atoi(q.Get("width"))
	if err != nil {
		width = defaultWidth
	}
	height, err := strconv.Atoi(q.Get("height"))
	if err != nil {
		height = defaultHeight
	}
	iterations, err := strconv.Atoi(q.Get("iterations"))
	if err != nil {
		iterations = defaultIterations
	}

	return newParams(width, height, iterations)
}
