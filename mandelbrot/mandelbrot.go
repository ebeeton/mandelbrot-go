// Package mandelbrot provides functionality to plot images and write them to
// HTTP responses.
package mandelbrot

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"net/http"
	"strconv"
	"sync"
)

// PlotImage plots a Mandelbrot image given a set of parameters and returns
// an image.RGBA.
func PlotImage(p *params) *image.RGBA {
	const (
		min float64 = -2.0
		max float64 = 2.0
	)

	img := image.NewRGBA(image.Rect(0, 0, p.width, p.height))

	aspectRatio := float64(p.height) / float64(p.width)
	minI := aspectRatio * min
	maxI := aspectRatio * max

	var wg sync.WaitGroup

	for y := 0; y < p.height; y++ {
		i := linearScale(float64(y), 0, float64(p.height), minI, maxI)
		// Plot each scan line concurrently.
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
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
		}(y)
	}

	wg.Wait()
	return img
}

// WriteImage writes an image.RGBA to an http.ResponseWriter.
func WriteImage(w http.ResponseWriter, img *image.RGBA) error {
	buf := new(bytes.Buffer)

	if err := png.Encode(buf, img); err != nil {
		return err
	}

	w.Header().Set("Content-type", "image/png")
	w.Header().Set("Content-length", strconv.Itoa(buf.Len()))
	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
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
