package core

import (
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

#KubeRef: {
  apiVersion: string
  kind:       string
  metadata: {
    metav1.#ObjectMeta
    name: string
  }
  // Allow any other fields
  ...
}

#Component: {
  in: {
    name: string
    _
    ...
  }
  out: #Config
}

#Config: {
  components?: [string]: #Component
  resources?: [string]: #KubeRef
  ...
}
