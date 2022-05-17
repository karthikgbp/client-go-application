package k8sclient

import (
	"context"
	"log"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type JobsK8s struct {
	ClientSet    *kubernetes.Clientset
	Namespace    *string
	JobName      *string
	Image        *string
	EntryCommand *string
}

type DeploymentK8s struct {
	ClientSet      *kubernetes.Clientset
	Namespace      *string
	WatchEnabled   *string
	DeploymentName *string
	Replicas       *int
	AppName        *string
	PodName        *string
	Image          *string
	EntryCommand   *string
	ServiceName    *string
}

// Create Job in K8s

func (launch *JobsK8s) createK8sJob() {

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

	result, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatal("Failed to create K8s job.", err)
	}

	log.Println("Create K8s job successfully : ", result.GetObjectMeta().GetName())

}

// Create Deployment

func (launch *DeploymentK8s) createK8sDeployment() {

	if _, err := ClientSet.CoreV1().Namespaces().Get(context.TODO(), *launch.Namespace, metav1.GetOptions{}); errors.IsNotFound(err) {

		nsSpec := &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: *launch.Namespace,
				Labels: map[string]string{
					"label": *launch.WatchEnabled,
				},
			},
		}
		_, err := ClientSet.CoreV1().Namespaces().Create(context.TODO(), nsSpec, metav1.CreateOptions{})
		if err != nil {
			log.Println("Error while creating Namespace", err.Error())
		}

	}

	// Set Deployment Spec

	deployment := ClientSet.AppsV1().Deployments(*launch.Namespace)

	deploymentSpec := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *launch.DeploymentName,
			Namespace: *launch.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(int32(*launch.Replicas)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"component": *launch.AppName,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"component": *launch.AppName,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    *launch.PodName,
							Image:   *launch.Image,
							Command: strings.Split(*launch.EntryCommand, " "),
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: 8081,
								},
							},
						},
					},
					// ServiceAccountName: *launch.ServiceName,
					RestartPolicy: v1.RestartPolicyAlways,
				},
			},
		},
	}

	// Create Deployment
	result, err := deployment.Create(context.TODO(), deploymentSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatal("Error occurred while creating Deployment", err)
	}

	log.Println("Successfully Created Deployment :", result.GetObjectMeta().GetName())

	// Set Service Spec
	service := ClientSet.CoreV1().Services(*launch.Namespace)

	serviceSpec := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *launch.ServiceName,
			Namespace: *launch.Namespace,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"component": *launch.AppName,
			},
			Type: "ClusterIP",
			Ports: []v1.ServicePort{
				{
					Port: 8081,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 8081,
					},
				},
			},
		},
	}

	// Attaching Service
	result2, err := service.Create(context.TODO(), serviceSpec, metav1.CreateOptions{})

	if err != nil {
		log.Fatal("Error occurred while attaching Service", err)
	}

	log.Println("Successfully Attached Service :", result2.GetObjectMeta().GetName())

}

func int32Ptr(i int32) *int32 {
	log.Println("Replica Set Value : ", &i)
	return &i
}

func createk8sInformer() {

	// factory := informers.NewSharedInformerFactory(ClientSet, 1000000000)
	factory := informers.NewFilteredSharedInformerFactory(ClientSet, 1000000000, "client-namespace", func(lo *metav1.ListOptions) {
		// lo.LabelSelector = "component=web"
		lo.LabelSelector = "name=enabled"

	})

	nsinformer := factory.Core().V1().Namespaces()

	Informer := nsinformer.Informer()

	ch := make(chan struct{})

	defer close(ch)

	// Kubernetes provides a utility to handle API crashes
	defer runtime.HandleCrash()

	// Start Informer
	go factory.Start(ch)

	Informer.AddEventHandler(cache.ResourceEventHandlerFuncs{

		AddFunc: func(obj interface{}) {

			addObj := obj.(metav1.Object).GetName()
			log.Println("Object is : ", addObj)
		},

		// UpdateFunc: func(interface{}, interface{}) { panic("") },

		// DeleteFunc: func(interface{}) { panic("") },
	})

	// podInformer.Run(ch)
	<-ch
}
