package main

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

func main() {
	http.HandleFunc("/", mandelbrotHandler)
	e := http.ListenAndServe(":3000", nil)
	if e != nil {
		log.Fatal(e)
	}
}

func mandelbrotHandler(w http.ResponseWriter, r *http.Request) {
	p := getQueryParams(r)
	img := plotImage(p)
	writeImage(w, img)
}

type params struct {
	width         int
	height        int
	maxIterations int
}

func newParams(width int, height int, iterations int) *params {
	p := params{width: width, height: height, maxIterations: iterations}
	return &p
}

func getQueryParams(r *http.Request) *params {
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

func plotImage(p *params) *image.RGBA {
	const (
		min float64 = -2.0
		max float64 = 2.0
	)

	img := image.NewRGBA(image.Rect(0, 0, p.width, p.height))

	var aspectRatio = float64(p.height) / float64(p.width)
	var minI = aspectRatio * min
	var maxI = aspectRatio * max

	for y := 0; y < p.height; y++ {
		i := linearScale(float64(y), 0, float64(p.height), minI, maxI)
		for x := 0; x < p.width; x++ {
			r := linearScale(float64(x), 0, float64(p.width), min, max)
			isInSet, iter := isInMandelbrotSet(complex(r, i), p.maxIterations)
			var gray uint8
			if isInSet {
				// Leave points in the set black.
				gray = 0
			} else {
				gray = uint8(float64(iter) / float64(p.maxIterations) * 255)
			}

			img.Set(x, y, color.NRGBA{
				R: gray,
				G: gray,
				B: gray,
				A: 255,
			})
		}
	}

	return img
}

func writeImage(w http.ResponseWriter, img *image.RGBA) {
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

func linearScale(val float64, minScaleFrom float64, maxScaleFrom float64, minScaleTo float64, maxScaleTo float64) float64 {
	return (val-minScaleFrom)/(maxScaleFrom-minScaleFrom)*(maxScaleTo-minScaleTo) + minScaleTo
}
