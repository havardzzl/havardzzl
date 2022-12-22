package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
)

func K8sGo() {
	x := resource.MustParse("1Gi")
	fmt.Println(x.MilliValue())
	fmt.Println(x.String())

	x2 := resource.MustParse("1Mi")
	fmt.Println(x2.MilliValue())
	fmt.Println(x2.String())
}

func mainK8() {
	K8sGo()
}
