package k8sclient

import (
	"flag"
)

// Create Job
func CreateJob() {

	jobNamespace := flag.String("namespace", "k8s-client-go", "Namespace of the Job")
	jobName := flag.String("jobname", "client-go-in-cluster", "Name of the Job")
	containerImage := flag.String("image", "ubuntu:latest", "Container Image")
	entryCmd := flag.String("command", "ls", "Command should run inside the container")

	flag.Parse()

	jobSpecDefined := JobsK8s{ClientSet, jobNamespace, jobName, containerImage, entryCmd}
	jobSpecDefined.createK8sJob()

}

// Edit Job

// Delete Job
