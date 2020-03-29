package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/google/uuid"
)

// Cluster struct
type Cluster struct {
	Version string
	PubKey  string
	ID      string
	Minutes string
	Region  string
	Cfg     string
}

type SupportedVersions struct {
	Versions []string
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
		cluster := Cluster{}

		err := json.NewDecoder(r.Body).Decode(&cluster)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

		}
		supported, err := checkversions()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if stringInSlice(cluster.Version, supported.Versions) {
			id, err := uuid.NewUUID()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			cluster.ID = id.String()
			cluster.Minutes = "30"
			cluster.Region = "nyc3"
			cluster.Cfg = "https://storage.googleapis.com/" + cluster.ID + "/cluster-config.gpg"

			err = createcluster(cluster)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			js, err := json.Marshal(cluster.Cfg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(js)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else {
			fmt.Fprintf(w, "Unsupported version")
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
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

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}
