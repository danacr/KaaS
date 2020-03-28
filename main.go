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
		versions, err := checkversions()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(versions.Versions)
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
		supported, err := checkversions()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if stringInSlice(user.Version, supported.Versions) {
			http.StatusText(200)
			fmt.Println("found")

		} else {
			fmt.Println("not found")
			http.StatusText(400)
		}

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

func checkversions() (SupportedVersions, error) {
	client := github.NewClient(nil)

	taglist, _, err := client.Repositories.ListTags(context.Background(), "poseidon", "typhoon", nil)
	if err != nil {
		return SupportedVersions{}, err
	}
	var supported SupportedVersions
	for _, tag := range taglist {
		supported.Versions = append(supported.Versions, tag.GetName())
	}
	if err != nil {
		return SupportedVersions{}, err
	}
	return supported, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
