package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/karthiksgd/go-in-cluster/k8sclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var ClientSet *kubernetes.Clientset

func init() {
	// Establish K8s Connectivity
	ClientSet = k8sclient.ConnectToK8s()

}

type Pods struct {
	PodName string `json:"pod_name"`
}

type Namespace struct {
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
}

func GetPods(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")

	// List all Pods
	pods, err := ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal("Unable to list all Pods", err)
	}

	var podList []Pods

	if len(pods.Items) > 0 {
		// List All Pods
		for _, pod := range pods.Items {

			podList = append(podList, Pods{pod.GetName()})
		}
	}
	json.NewEncoder(w).Encode(podList)
}

func GetNamespaces(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")

	// List all Namespaces
	ns, err := ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal("Unable to list all Namespaces", err)
	}

	var nsList []Namespace

	if len(ns.Items) > 0 {
		fmt.Println("No of Namespaces in the Cluster : ", len(ns.Items))
		for _, ns := range ns.Items {
			nsList = append(nsList, Namespace{ns.GetName(), ns.Kind})
		}
	}
	json.NewEncoder(w).Encode(nsList)
}
