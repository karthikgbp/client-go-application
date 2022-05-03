package k8sclient

import (
	"context"
	"log"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
)

type LaunchK8s struct {
	ClientSet    *kubernetes.Clientset
	Namespace    *string
	JobName      *string
	Image        *string
	EntryCommand *string
}

// Create Job in K8s

func (launch *LaunchK8s) createK8sJob() {

	if _, err := ClientSet.CoreV1().Namespaces().Get(context.TODO(), *launch.Namespace, metav1.GetOptions{}); errors.IsNotFound(err) {

		nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: *launch.Namespace}}
		_, err := ClientSet.CoreV1().Namespaces().Create(context.TODO(), nsSpec, metav1.CreateOptions{})
		if err != nil {
			log.Println("Error while creating Namespace", err.Error())
		}
	}

	jobs := launch.ClientSet.BatchV1().Jobs(*launch.Namespace)
	var backOffLimit int32 = 0

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *launch.JobName,
			Namespace: *launch.Namespace,
		},

		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    *launch.JobName,
							Image:   *launch.Image,
							Command: strings.Split(*launch.EntryCommand, " "),
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	_, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatal("Failed to create K8s job.", err)
	}

	log.Println("Create K8s job successfully")

}
