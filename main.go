package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/karthiksgd/go-in-cluster/k8sclient"
	"github.com/karthiksgd/go-in-cluster/rest"
	kubernetes "k8s.io/client-go/kubernetes"
)

var ClientSet *kubernetes.Clientset

func init() {
	// Establish K8s Connectivity
	ClientSet = k8sclient.ConnectToK8s()

}

func main() {

	fmt.Println("Loading from Main ....")

	// Create a Job
	k8sclient.CreateJob()

	// Create Deployments
	k8sclient.CreateDeployment()

	// List Resources

	// Pods
	k8sclient.ListPods()

	// Namespaces
	k8sclient.ListNamespaces()

	//Serve Http
	router := mux.NewRouter()
	router.HandleFunc("/api/pods", rest.GetPods).Methods("Get")
	router.HandleFunc("/api/namespaces", rest.GetNamespaces).Methods("Get")
	http.ListenAndServe(":5050", router)
}
