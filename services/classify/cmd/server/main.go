// Classify service stub — production traffic handled by gateway defer endpoint.
package main

import "net/http"

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })
	http.ListenAndServe(":8082", nil)
}
