package main

import (
	"context"
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
		return err
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
							Image: "danacr/stk:latest",
							Env: []apiv1.EnvVar{
								apiv1.EnvVar{
									Name:  "TF_VAR_cluster_id",
									Value: cluster.ID},
								apiv1.EnvVar{
									Name:  "TF_VAR_cluster_region",
									Value: cluster.Region},
								apiv1.EnvVar{
									Name:  "HOW_LONG",
									Value: cluster.Minutes},
								apiv1.EnvVar{
									Name:  "CLUSTER_VERSION",
									Value: cluster.Version},
								apiv1.EnvVar{
									Name:  "pubkey",
									Value: cluster.PubKey},
							},
							VolumeMounts: []apiv1.VolumeMount{
								apiv1.VolumeMount{
									Name:      "do-token",
									MountPath: "/home/terraform/config/do_token",
									SubPath:   "do_token",
								},
								apiv1.VolumeMount{
									Name:      "service-account",
									MountPath: "/home/terraform/config/service-account-key.json",
									SubPath:   "service-account-key.json",
								},
							},
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
					Volumes: []apiv1.Volume{
						apiv1.Volume{
							Name: "do-token",
							VolumeSource: apiv1.VolumeSource{
								Secret: &apiv1.SecretVolumeSource{
									SecretName: "do-token",
									Items: []apiv1.KeyToPath{
										apiv1.KeyToPath{
											Key:  "do_token",
											Path: "do_token",
										},
									},
								},
							},
						},
						apiv1.Volume{
							Name: "service-account",
							VolumeSource: apiv1.VolumeSource{
								Secret: &apiv1.SecretVolumeSource{
									SecretName: "service-account",
									Items: []apiv1.KeyToPath{
										apiv1.KeyToPath{
											Key:  "service-account-key.json",
											Path: "service-account-key.json",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	// Create Job
	fmt.Println("Creating job...")
	result, err := jobsClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created job %q.\n", result.GetObjectMeta().GetName())
	return nil
}
