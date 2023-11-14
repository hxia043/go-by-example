package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		panic(err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	unstructObj, err := dynamicClient.Resource(gvr).Namespace("kube-system").List(context.Background(), metav1.ListOptions{Limit: 500})
	if err != nil {
		panic(err)
	}

	podList := &corev1.PodList{}
	if err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), podList); err != nil {
		panic(err)
	}

	for _, d := range podList.Items {
		fmt.Printf("NAMESPACE: %v \t NAME: %v \t STATUS:%+v\n", d.Namespace, d.Name, d.Status.Phase)
	}
}
