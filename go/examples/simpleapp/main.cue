package simpleapp

import "github.com/adb-sh/kubeneko/core"

// let app1 = #genSimpleApp & {
//   name: "app1"
//   in: {
//     name: "hello-world"
//     image: "nginx"
//     port: 80
//   }
// }

// let app2 = #genSimpleApp & {
//   name: "app2"
//   in: {
//     name: "whoami"
//     image: "traefik/whoami"
//     port: app1.out.resources[0].spec.containers[0].ports[0].containerPort
//   }
// }


config: core.#Config & {
  components: {
    app1: #genSimpleApp & {
      in: {
        name: "hello-world"
        image: "nginx"
        port: 80
      }
    },
    app2: #genSimpleApp & {
      in: {
        name: "whoami"
        image: "traefik/whoami"
        port: app1.out.resources[0].spec.containers[0].ports[0].containerPort
      }
    }
  }
}
