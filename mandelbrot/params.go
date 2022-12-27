package mandelbrot

import (
	"net/http"
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

func GetQueryParams(r *http.Request) *params {
	const (
		defaultWidth      int = 1024
		defaultHeight     int = 768
		defaultIterations int = 512
	)

	width, err := strconv.Atoi(r.URL.Query().Get("width"))
	if err != nil {
		width = defaultWidth
	}
	height, err := strconv.Atoi(r.URL.Query().Get("height"))
	if err != nil {
		height = defaultHeight
	}
	iterations, err := strconv.Atoi(r.URL.Query().Get("iterations"))
	if err != nil {
		iterations = defaultIterations
	}

	return newParams(width, height, iterations)
}
