package main

import (
	"flag"
	"fmt"

	"net/http"

)

// This assumes that there is a worker running
// To test it:
// curl -X POST -H "Content-Type: application/json" -d '{"img_url":"http://localhost:8081/img","engine":0}' http://localhost:8081/ocr

func init() {

}

func main() {

	var http_port int
	flagFunc := func() {
		flag.IntVar(
			&http_port,
			"http_port",
			8888,
			"The http port to listen on, eg, 8081",
		)

	}

	// any requests to root, just redirect to main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		text := `<h1>OpenOCR is running!<h1> Need <a href="http://www.openocr.net">docs</a>?`
		fmt.Fprintf(w, text)
	})

//	http.Handle("/ocr", func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintf(w, "OCR")
//	})

	listenAddr := fmt.Sprintf(":%d", http_port)
	http.ListenAndServe(listenAddr , nil)

}
