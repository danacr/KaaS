package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/cluster", kaas)

	http.HandleFunc("/favicon.ico", faviconHandler)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Printf("Starting kaas 🧀\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
