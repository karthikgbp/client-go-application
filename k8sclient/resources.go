package k8sclient

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	// Establish K8s Connectivity
	ClientSet = ConnectToK8s()

}

type Object struct {
	ObjType string
}

func ListPods() {

	// List all Pods
	pods, err := ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal("Unable to list all Pods", err)
	}

	if len(pods.Items) > 0 {
		fmt.Println("No of Pods in the Cluster: ", len(pods.Items))

		// List All Pods
		for i, pod := range pods.Items {
			fmt.Println("Pod ", i, " - ", pod.GetName(), " ;  Pod Lable : ", pod.Labels)
		}
	}
}

func ListNamespaces() {

	// List All Namespaces
	ns, err := ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal("Unable to list all Namespaces", err)
	}

	if len(ns.Items) > 0 {
		fmt.Println("No of Namespaces in the Cluster : ", len(ns.Items))
		for i, ns := range ns.Items {
			fmt.Println("Namespace :", i, " - ", ns.GetName())
			fmt.Println("Label :", i, " - ", ns.Labels["name"])
			// fmt.Println("Object Meta ", i, " - ", ns.ObjectMeta)
		}
	}
}

func GetNamespaces() []string {

	var nsList []string

	ns, err := ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal("Unable to list all Namespaces", err)
	}

	for _, ns := range ns.Items {
		nsList = append(nsList, ns.GetName())
	}
	return nsList
}

func ListUserDefNSLabels(uD string) map[string]string {

	var nsLabel = make(map[string]string)

	// ns, err := ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	ns, err := ClientSet.CoreV1().Namespaces().Get(context.TODO(), uD, metav1.GetOptions{})

	if err != nil {
		log.Println("Unable to list all Namespaces", err)
	}

	// fmt.Println(ns.GetLabels())

	for k, val := range ns.GetLabels() {

		if k == "kubernetes.io/metadata.name" {
			continue
		}
		nsLabel[k] = val

	}

	return nsLabel
}

func DeleteNamespace(ns string) bool {

	err := ClientSet.CoreV1().Namespaces().Delete(context.TODO(), ns, metav1.DeleteOptions{})

	if err != nil {
		log.Fatal(err)
	}

	return true
}
