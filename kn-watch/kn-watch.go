package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/configmap/informer"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	NS_FILE = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	CM_NAME = "clusters-config"
)

var (
	namespace string
)

func init() {
	nsBytes, err := ioutil.ReadFile(NS_FILE)
	if err != nil {
		panic(fmt.Sprintf("Unable to read namespace file at %s", NS_FILE))
	}

	namespace = string(nsBytes)
}

func main() {
	var (
		watcher         configmap.Watcher
		currentEndpoint string
	)

	watcherStopChan := make(chan struct{})

	// Get things set up for watching - we need a valid k8s client
	clientCfg, err := rest.InClusterConfig()
	if err != nil {
		panic("Unable to get our client configuration")
	}

	clientset, err := kubernetes.NewForConfig(clientCfg)
	if err != nil {
		panic("Unable to create our clientset")
	}

	// Create our watcher
	req, _ := labels.NewRequirement("watcherManaged", selection.Equals, []string{"yes"})
	watcher = informer.NewInformedWatcher(clientset, namespace, *req)

	// Specify our callback for the configmap with the name stored in CM_NAME
	watcher.Watch(CM_NAME, func(updated *corev1.ConfigMap) {
		if endpointKey, ok := updated.Data["current.target"]; ok {
			if endpoint, ok := updated.Data[endpointKey]; ok {
				currentEndpoint = endpoint
				fmt.Println("Endpoint updated")
			}
		}
	})

	// Start watching
	watcher.Start(watcherStopChan)
	fmt.Println("Watcher started...")

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		body := []byte(fmt.Sprintf(`{"current_endpoint": "%s"}`, currentEndpoint))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})

	fmt.Printf("Listening on port 8080\n")
	http.ListenAndServe(":8080", nil)
}
