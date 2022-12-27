package mandelbrot

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
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

func PlotImage(p *params) *image.RGBA {
	const (
		min float64 = -2.0
		max float64 = 2.0
	)

	img := image.NewRGBA(image.Rect(0, 0, p.width, p.height))

	aspectRatio := float64(p.height) / float64(p.width)
	minI := aspectRatio * min
	maxI := aspectRatio * max

	for y := 0; y < p.height; y++ {
		i := linearScale(float64(y), 0, float64(p.height), minI, maxI)
		for x := 0; x < p.width; x++ {
			r := linearScale(float64(x), 0, float64(p.width), min, max)

			var gray uint8
			if isInSet, iter := isInMandelbrotSet(complex(r, i), p.maxIterations); isInSet {
				// Leave points in the set black.
				gray = 0
			} else {
				gray = uint8(float64(iter) / float64(p.maxIterations) * math.MaxUint8)
			}

			img.Set(x, y, color.NRGBA{
				R: gray,
				G: gray,
				B: gray,
				A: uint8(color.Opaque.A),
			})
		}
	}

	return img
}

func WriteImage(w http.ResponseWriter, img *image.RGBA) {
	buf := new(bytes.Buffer)

	if err := png.Encode(buf, img); err != nil {
		log.Printf("Failed to encode image. %s", err)
		return
	}

	w.Header().Set("Content-type", "image/png")
	w.Header().Set("Content-length", strconv.Itoa(buf.Len()))
	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Printf("Failed to write image. %s", err)
	}
}

func isInMandelbrotSet(c complex128, maxIterations int) (bool, int) {
	const bailout float64 = 2
	z := c
	for i := 0; i < maxIterations; i++ {
		if math.Abs(real(z)) > bailout || math.Abs(imag(z)) > bailout {
			return false, i
		}
		z = z*z + c
	}

	return true, maxIterations
}

func linearScale(val, minScaleFrom, maxScaleFrom, minScaleTo, maxScaleTo float64) float64 {
	return (val-minScaleFrom)/(maxScaleFrom-minScaleFrom)*(maxScaleTo-minScaleTo) + minScaleTo
}
