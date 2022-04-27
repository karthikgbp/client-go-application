package k8sclient

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	kubernetes "k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

var ClientSet *kubernetes.Clientset

func init() {
	// Establish K8s Connectivity
	ClientSet = ConnectToK8s()

}

// Connect to K8s
func ConnectToK8s() *kubernetes.Clientset {

	home, exist := os.LookupEnv("HOME")
	if !exist {
		home = "/root"
	}

	fmt.Println("Home directory is ", home)

	configPath := filepath.Join(home, ".kube", "config")

	fmt.Println("File Path : ", configPath)

	// Create K8s Config
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Fatalln("Failed to Create k8s config", err)
	}

	// Create K8s Client Set
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Failes to create Client Set", err)
	}

	return clientSet

}
