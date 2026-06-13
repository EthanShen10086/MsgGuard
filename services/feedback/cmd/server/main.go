// Feedback service stub — production traffic handled by gateway.
package main

import "net/http"

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })
	http.ListenAndServe(":8084", nil)
}
