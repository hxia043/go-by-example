package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/apps"
)

func main() {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1beta1.SchemeGroupVersion, &appsv1beta1.Deployment{})
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &appsv1.Deployment{})
	scheme.AddKnownTypes(apps.SchemeGroupVersion, &apps.Deployment{})
	metav1.AddToGroupVersion(scheme, appsv1beta1.SchemeGroupVersion)
	metav1.AddToGroupVersion(scheme, appsv1.SchemeGroupVersion)
	//metav1.AddToGroupVersion(scheme, apps.SchemeGroupVersion)

	v1beta1Deployment := &appsv1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1beta1", // Corrected APIVersion here
		},
	}

	// v1beta1 -> __internal
	objInternal, err := scheme.ConvertToVersion(v1beta1Deployment, appsv1beta1.SchemeGroupVersion)
	if err != nil {
		panic(err)
	}

	fmt.Println("GVK: ", objInternal.GetObjectKind().GroupVersionKind().String())

	// __internal -> v1
	objV1, err := scheme.ConvertToVersion(objInternal, apps.SchemeGroupVersion)
	if err != nil {
		panic(err)
	}

	v1Deployment, ok := objV1.(*appsv1.Deployment)
	if !ok {
		panic("wrong Deployment type")
	}

	fmt.Println("GVK: ", v1Deployment.GetObjectKind().GroupVersionKind().String())

}
