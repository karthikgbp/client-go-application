package k8sclient

import (
	"flag"
	"fmt"
)

func CreateJob() {

	jobName := flag.String("jobname", "test-job", "Name of the Job")
	containerImage := flag.String("image", "ubuntu:latest", "Container Image")
	entryCmd := flag.String("command", "ls", "Command should run inside the container")

	flag.Parse()

	fmt.Printf("Args are %s %s %s \n", *jobName, *containerImage, *entryCmd)

	jobSpecDefined := LaunchK8s{ClientSet, jobName, containerImage, entryCmd}
	jobSpecDefined.createK8sJob()

}
