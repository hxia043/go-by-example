package main

import (
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/core"
)

func main() {
	pod := &core.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind: "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{"name": "foo"},
		},
	}

	fmt.Println(*pod)

	obj := runtime.Object(pod)
	pod2, ok := obj.(*core.Pod)
	if !ok {
		panic("unexpected runtime object")
	}

	if !reflect.DeepEqual(pod, pod2) {
		panic("unexpected")
	}

	fmt.Println(*pod2)
}
