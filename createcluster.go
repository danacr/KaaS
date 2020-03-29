package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func createcluster(cluster Cluster) error {
	fmt.Println(cluster.ID, cluster.Minutes, cluster.Region, cluster.Version, cluster.PubKey)
	var kubeconfig *string
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	kubeconfig = flag.String("kubeconfig", filepath.Join(path, "kubeconfig.yaml"), "(optional) absolute path to the kubeconfig file")

	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	jobsClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: cluster.ID,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "terraform",
							Image: "ubuntu",
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
				},
			},
		},
	}
	// Create Job
	fmt.Println("Creating job...")
	result, err := jobsClient.Create(job)
	if err != nil {
		return err
	}
	fmt.Printf("Created job %q.\n", result.GetObjectMeta().GetName())
	return nil
}
