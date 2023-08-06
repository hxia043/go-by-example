package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
	// KnownType external
	coreGV := schema.GroupVersion{Group: "", Version: "v1"}
	extensionsGV := schema.GroupVersion{Group: "extensions", Version: "v1beta1"}

	// KnownType internal
	coreInternalGV := schema.GroupVersion{Group: "", Version: runtime.APIVersionInternal}

	// UnversionedType
	Unversioned := schema.GroupVersion{Group: "", Version: "v1"}

	schema := runtime.NewScheme()
	schema.AddKnownTypes(coreGV, &corev1.Pod{})
	schema.AddKnownTypes(extensionsGV, &appsv1.DaemonSet{})
	schema.AddKnownTypes(coreInternalGV, &corev1.Pod{})
	schema.AddUnversionedTypes(Unversioned, &metav1.Status{})

	fmt.Println(*schema)

	fmt.Println(schema.KnownTypes(coreGV))
}
