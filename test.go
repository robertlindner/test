package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var kubeconfig *string
var clientSet *kubernetes.Clientset

func main() {
	var config *rest.Config
	var err error

	var kubeconfig *string
	home := homedir.HomeDir()
	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	// Creates the clientset
	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	repNo := int32(2)
	configMapName := "config1"
	secretName := "secret1"

	//This containerList is used for both Statefulset and Deployment
	containerList := []apiv1.Container{
		{
			Name:            "postgres",
			Image:           "postgres:13",
			ImagePullPolicy: "IfNotPresent",
			Ports: []apiv1.ContainerPort{{
				ContainerPort: int32(5432),
				Name:          "postgredb",
			}},
			EnvFrom: []apiv1.EnvFromSource{
				{
					ConfigMapRef: &apiv1.ConfigMapEnvSource{
						LocalObjectReference: apiv1.LocalObjectReference{
							Name: configMapName,
						},
					},

					SecretRef: &apiv1.SecretEnvSource{
						LocalObjectReference: apiv1.LocalObjectReference{
							Name: secretName,
						},
					},
				},
			},
		},
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "d1",
		},

		Spec: appsv1.DeploymentSpec{
			Replicas: &repNo,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "postgres",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "postgres",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: containerList,
				},
			},
		},
	}

	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "s1",
		},

		Spec: appsv1.StatefulSetSpec{
			Replicas: &repNo,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "postgres",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "postgres",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: containerList,
				},
			},
		},
	}

	result, err1 := clientSet.AppsV1().StatefulSets(apiv1.NamespaceDefault).Create(context.TODO(), statefulset, metav1.CreateOptions{})
	if err1 != nil {
		fmt.Printf(" Error below exist!")
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Stateful ...... %q created!\n", result.GetObjectMeta().GetName())
	}

	result2, err := clientSet.AppsV1().Deployments(apiv1.NamespaceDefault).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("Error below exist!")
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Deployment...... %q created!\n", result2.GetObjectMeta().GetName())
	}
}
