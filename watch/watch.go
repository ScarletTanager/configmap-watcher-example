package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	NS_FILE          = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	CM_NAME          = "clusters-config"
	DEFAULT_ENDPOINT = "https://wazanga.partytime"
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
		currentEndpoint string
		mutex           sync.Mutex
	)

	// Let's make sure we don't forget to set a default
	currentEndpoint = DEFAULT_ENDPOINT

	// Get things set up for watching - we need a valid k8s client
	clientCfg, err := rest.InClusterConfig()
	if err != nil {
		panic("Unable to get our client configuration")
	}

	clientset, err := kubernetes.NewForConfig(clientCfg)
	if err != nil {
		panic("Unable to create our clientset")
	}

	watcher, err := clientset.CoreV1().ConfigMaps(namespace).Watch(context.TODO(),
		metav1.SingleObject(metav1.ObjectMeta{Name: CM_NAME, Namespace: namespace}))
	if err != nil {
		panic("Unable to create watcher")
	}

	go updateCurrentEndpoint(watcher.ResultChan(), &currentEndpoint, mutex)

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		body := []byte(fmt.Sprintf(`{"current_endpoint": "%s"}`, currentEndpoint))
		mutex.Unlock()
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})

	fmt.Printf("Listening on port 8080\n")
	http.ListenAndServe(":8080", nil)
}

func updateCurrentEndpoint(eventChannel <-chan watch.Event, endpoint *string, mutex sync.Mutex) {
	for event := range eventChannel {
		switch event.Type {
		case watch.Added:
			fallthrough
		case watch.Modified:
			mutex.Lock()
			// Update our endpoint
			if updatedMap, ok := event.Object.(*corev1.ConfigMap); ok {
				if endpointKey, ok := updatedMap.Data["current.target"]; ok {
					if targetEndpoint, ok := updatedMap.Data[endpointKey]; ok {
						*endpoint = targetEndpoint
					}
				}
			}
			mutex.Unlock()
		case watch.Deleted:
			mutex.Lock()
			// Fall back to the default value
			*endpoint = DEFAULT_ENDPOINT
			mutex.Unlock()
		default:
			// Do nothing
		}
	}
}
