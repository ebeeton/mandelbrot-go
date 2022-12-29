package main

import (
	"log"
	"net/http"

	"github.com/ebeeton/mandelbrot-go/mandelbrot"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := mandelbrot.GetParams(r.URL.Query())
		img := mandelbrot.PlotImage(p)
		if err := mandelbrot.WriteImage(w, img); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
