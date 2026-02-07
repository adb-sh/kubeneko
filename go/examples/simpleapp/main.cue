package simpleapp

import (
  "github.com/adb-sh/kubeneko/core"
  "k8s.io/api/core/v1"
)

app1: #genSimpleApp & {
  name: "app1"
  in: {
    name: "hello-world"
    image: "nginx"
    port: 80
  }
}

app2: #genSimpleApp & {
  name: "app2"
  in: {
    name: "whoami"
    image: "traefik/whoami"
    port: 80
  }
}

myService: v1.#Service & {
  kind: "Service"
  apiVersion: "v1"
  metadata: {
    name: "lol"
  }
  spec: {
    selector: {
      app: "lol"
    }
    ports: [{
      port:       80
      targetPort: 80
    }]
  }
}
myService2: v1.#Service & {
  kind: "Service"
  apiVersion: "v1"
  metadata: {
    name: "lol2"
    labels: {
      other: myService.metadata.name
    }
  }
  spec: {
    selector: {
      app: "lol"
    }
    ports: [{
      port:       80
      targetPort: myService.spec.ports[0].targetPort
    }]
  }
}


config: core.#Config & {
  components: [
    app1, app2
  ]
  resources: [
    myService, myService2
  ]
}
