package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"knative.dev/pkg/configmap"
)

const (
	CM_NAME = "clusters-config"
)

var (
	namespace string
)

func main() {
	cfg, err := configmap.Load(CM_NAME)
	if err != nil {
		panic("Unable to load config map!")
	} else {
		http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
			body, err := json.Marshal(cfg)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(body)
			}
		})
	}

	fmt.Printf("Listening on port 8080\n")
	http.ListenAndServe(":8080", nil)
}
