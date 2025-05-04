package simpleapp

import "github.com/adb-sh/kubeneko/core"

let app1 = #genSimpleApp & {
  name: "app1"
  in: {
    name: "hello-world"
    image: "nginx"
    port: 80
  }
}

let app2 = #genSimpleApp & {
  name: "app2"
  in: {
    name: "whoami"
    image: "traefik/whoami"
    port: 80
  }
}


config: core.#Config & {
  components: [
    app1, app2
  ]
}
