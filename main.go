package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(http.ResponseWriter, *http.Request) {
		log.Print("Hello World")
	})

	http.Handle("/goodbye", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Oooopsy", http.StatusBadRequest)
			return
		}

		log.Printf("Data %s", data)
		fmt.Fprintf(w, "Bye %s\n", string(data))

	}))

	http.ListenAndServe(":9090", nil)
}
