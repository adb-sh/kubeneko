package simpleapp

import (
  "k8s.io/api/core/v1"
  netv1 "k8s.io/api/networking/v1"
  "github.com/adb-sh/kubeneko/core"
)

#SimpleAppInput: {
  name: string
  image: string
  port: int
}

#genSimpleApp: S=core.#Component & {
  in: #SimpleAppInput

  _pod: v1.#Pod & {
    kind: "Pod"
    apiVersion: "v1"
    metadata: {
      name: in.name
      labels: {
        app: in.name
      }
    }
    spec: {
      containers: [{
        name:  "my-container"
        image: in.image
        ports: [{
          containerPort: in.port
        }]
      }]
    }
  }

  _service: v1.#Service & {
    kind: "Service"
    apiVersion: "v1"
    metadata: {
      name: in.name
    }
    spec: {
      selector: {
        app: _pod.metadata.labels.app
      }
      ports: [{
        port:       in.port
        targetPort: _pod.spec.containers[0].ports[0].containerPort
      }]
    }
  }

  _ingress: netv1.#Ingress & {
    kind: "Ingress"
    apiVersion: "networking.k8s.io/v1"
    metadata: {
      name: in.name
    }
    spec: {
      rules: [{
        host: "example.com"
        http: {
          paths: [{
            path: "/"
            pathType: "Prefix"
            backend: {
              service: {
                name: in.name
                port: {
                  number: in.port
                }
              }
            }
          }]
        }
      }]
    }
  }
}


#PostgresDBInput: {
  name: string
  image: string
  port: int
}

#genPostgresDB: core.#Component & {
  in: #PostgresDBInput

  out: {
    resources: {
      pod: v1.#Pod & {
        kind: "Pod"
        apiVersion: "v1"
        metadata: {
          name: in.name
          labels: {
            app: in.name
          }
        }
        spec: {
          containers: [{
            name:  "my-container"
            image: in.image
            ports: [{
              containerPort: in.port
            }]
          }]
        }
      },
      service: v1.#Service & {
        kind: "Service"
        apiVersion: "v1"
        metadata: {
          name: in.name
        }
        spec: {
          selector: {
            app: in.name
          }
          ports: [{
            port:       in.port
            targetPort: in.port
          }]
        }
      },
    },
    db: {
      host: resources.service.metadata.name
      port: resources.service.spec.ports[0].port
    }
  }
}
