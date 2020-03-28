package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v30/github"
)

// User struct
type User struct {
	Version string
	PubKey  string
}
type SupportedVersions struct {
	Versions []string
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func kaas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		client := github.NewClient(nil)

		taglist, _, err := client.Repositories.ListTags(context.Background(), "poseidon", "typhoon", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var supported SupportedVersions
		for _, tag := range taglist {
			supported.Versions = append(supported.Versions, tag.GetName())
		}
		err = json.NewEncoder(w).Encode(supported.Versions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	case "POST":
		user := User{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

		}

		log.Println("result")

		http.StatusText(200)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	var err error
	http.HandleFunc("/get", kaas)
	http.HandleFunc("/favicon.ico", faviconHandler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	fmt.Printf("Starting kaas 🧀\n")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
