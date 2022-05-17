package k8sclient

import "flag"

func CreateDeployment() {

	deploymentNS := flag.String("namespaceDeployment", "k8s-client-go", "Namespace of the Deployment")
	label := flag.String("labelNs", "enabled", "Label for Namespace")
	deploymentName := flag.String("deploymentName", "client-go-cluster-deployment", "Name of Deployment")
	replicas := flag.Int("replicas", 1, "No of Replicas")
	appName := flag.String("appName", "app-in-go-deployment", "App Name")
	podName := flag.String("podName", "app-in-go-deployment", "Container Image")
	image := flag.String("imageDeployment", "alpine:3.15.4", "Container Image")
	entryCmd := flag.String("commandDeployment", "sleep 1800", "Command should run inside the container")
	serviceName := flag.String("serviceName", "serv-deployement", "Service as Load balancer")

	flag.Parse()

	deploymentDefined := DeploymentK8s{ClientSet, deploymentNS, label, deploymentName, replicas, appName, podName, image, entryCmd, serviceName}

	deploymentDefined.createK8sDeployment()
}
