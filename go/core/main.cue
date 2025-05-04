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
  name: string
  in: _
  out: #Config
}

// #Config: {
//   components?: [...#Component]
//   resources?: [...#KubeRef]
// }

#Config: {
  components?: [Name=string]: #Component & {
    name: string | Name
  }
  resources?: [...#KubeRef]
}
