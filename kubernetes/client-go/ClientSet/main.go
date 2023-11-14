package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	podClient := clientset.CoreV1().Pods("kube-system")

	list, err := podClient.List(context.TODO(), metav1.ListOptions{Limit: 500})
	if err != nil {
		panic(err)
	}

	for _, d := range list.Items {
		fmt.Printf("NAMESPACE: %v \t NAME: %v \t STATUS:%+v\n", d.Namespace, d.Name, d.Status.Phase)
	}
}
