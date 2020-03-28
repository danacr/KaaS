package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tidwall/gjson"
)

// User struct
type User struct {
	Version string
	PubKey  string
}
type SupportedVersions struct {
	Array []string
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func kaas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		resp, err := http.Get("https://api.github.com/repos/poseidon/typhoon/tags")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		value := gjson.Get(bodyString, "#.name")
		versions := value.Array()
		js, err := json.Marshal(versions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(js)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case "POST":
		user := User{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
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
