package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

// User struct
type User struct {
	Version string
	PubKey  string
}

// Cluster struct
type Cluster struct {
	Creation time.Time
	Ready    bool
	Version  string
}

// ClusterID struct
type ClusterID struct {
	ID string
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
			id, err := createcluster(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, id)
		} else {
			fmt.Fprintf(w, "Unsupported version")
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	var err error
	if err = serviceAccount(); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/get", kaas)

	http.HandleFunc("/cluster", clusterHandler)

	http.HandleFunc("/favicon.ico", faviconHandler)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Printf("Starting kaas 🧀\n")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func createcluster(user User) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	time := time.Now()
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "k8stfw")
	if err != nil {
		return "", err
	}
	collection := client.Collection("Clusters")
	clusterid := id.String()
	document := collection.Doc(clusterid)
	wr, err := document.Create(ctx, Cluster{
		Creation: time,
		Ready:    false,
		Version:  user.Version,
	})
	if err != nil {
		return "", err
	}
	err = terraformcluster(clusterid, user)
	if err != nil {
		return "", err
	}
	fmt.Println(wr)
	return clusterid, nil
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func clusterHandler(w http.ResponseWriter, r *http.Request) {
	clusterid := ClusterID{}

	err := json.NewDecoder(r.Body).Decode(&clusterid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, "k8stfw")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	collection := client.Collection("Clusters")
	document := collection.Doc(clusterid.ID)
	data, err := document.Get(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	datastring, err := json.Marshal(data.Data())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(datastring))

}

func terraformcluster(id string, user User) error {
	f, err := os.Create("static/" + id)
	if err != nil {
		return err
	}
	_, err = f.WriteString("Hello World")
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	err = cfgencrypt(user.PubKey, "static/"+id)
	if err != nil {
		return err
	}
	return nil
}
