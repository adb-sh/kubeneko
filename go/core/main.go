package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubeRef struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   metav1.ObjectMeta `json:"metadata"`
	// Extra fields allowed — use map or `runtime.RawExtension` for unstructured data
	Spec map[string]interface{} `json:"spec"`
}

type Component struct {
	Name string      `json:"name"`
	In   interface{} `json:"in"`  // flexible generic input
	Out  Config      `json:"out"` // replace with your actual Config type
}

type Config struct {
	Components []Component `json:"components"`
	Resources  []KubeRef   `json:"resources"`
}
