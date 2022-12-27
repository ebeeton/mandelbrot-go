package main

import (
	"log"
	"net/http"

	"github.com/ebeeton/mandelbrot-go/mandelbrot"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := mandelbrot.GetQueryParams(r)
		img := mandelbrot.PlotImage(p)
		mandelbrot.WriteImage(w, img)
	})
	if e := http.ListenAndServe(":3000", nil); e != nil {
		log.Fatal(e)
	}
}
