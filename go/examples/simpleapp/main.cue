package simpleapp

import "github.com/adb-sh/kubeneko/core"


config: core.#Config & {
  components: {
    app1: #genSimpleApp & {
      in: {
        name: "hello-world"
        image: "nginx"
        port: 80
      }
    },
    // app2: #genSimpleApp & {
    //   in: {
    //     name: "whoami"
    //     image: "traefik/whoami"
    //     port: app1.out.resources.pod.spec.containers[0].ports[0].containerPort
    //   }
    // }
  }
}
