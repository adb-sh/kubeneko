import { createResource, createTemplate } from "./lib/basics.js";
import { createSignal, effect } from "./lib/reactive.js";
import * as k8s from "@kubernetes/client-node";
import { z } from "zod";
import type { V1Pod } from "@kubernetes/client-node";
import { Provider } from "./lib/provider.js";

// Load the kubeconfig file
const kc = new k8s.KubeConfig();
kc.loadFromDefault();
const provider = new Provider(kc);

const image = createSignal("nginx");
const active = createSignal(true);
const name = createSignal("test");
const namespaceName = createSignal("kubeneko");

const service = await createResource(
  () => {
    return createTemplate<k8s.V1Service>(() => ({
      apiVersion: "v1",
      kind: "Service",
      metadata: {
        name: `${name.value}-svc`,
        namespace: namespaceName.value,
      },
      spec: {
        type: "LoadBalancer",
        selector: {
          name: name.value,
        },
        ports: [
          {
            protocol: "TCP",
            port: 80,
            targetPort: 80,
          },
        ],
      },
    }));
  },
  { provider }
);

const configMap = await createResource(
  () => {
    return createTemplate<k8s.V1ConfigMap>(() => {
      const ip = service.value.status.loadBalancer.ingress[0].ip;
      return {
        apiVersion: "v1",
        kind: "ConfigMap",
        metadata: {
          name: `${name.value}-cm`,
          namespace: namespaceName.value,
        },
        data: {
          "index.html": `my ip is "${ip}"`,
        },
      };
    });
  },
  { provider }
);

const pod = await createResource(
  () => {
    const name = createSignal("test");

    return createTemplate<k8s.V1Deployment>(() => ({
      apiVersion: "apps/v1",
      kind: "Deployment",
      metadata: {
        name: name.value,
        namespace: namespaceName.value,
        annotations: {
          "kubeneko.adb.sh/active": active.value.toString(),
        },
        labels: {
          name: name.value,
        },
      },
      spec: {
        replicas: 1,
        selector: {
          matchLabels: {
            name: name.value,
          },
        },
        template: {
          metadata: {
            labels: {
              name: name.value,
            },
          },
          spec: {
            containers: [
              {
                name: "nginx",
                image: image.value,
                ports: [
                  {
                    containerPort: 80,
                  },
                ],
                volumeMounts: [
                  {
                    name: "html",
                    mountPath: "/usr/share/nginx/html/ip/",
                    readOnly: true,
                  },
                ],
              },
            ],
            volumes: [
              {
                name: "html",
                configMap: {
                  name: configMap.value.metadata.name,
                  items: [
                    {
                      key: "index.html",
                      path: "index.html",
                    },
                  ],
                },
              },
            ],
          },
        },
      },
    }));
  },
  { provider }
);

// const ip = createSignal(null);
// effect(() => {
//   const _ip = service.value.status.loadBalancer.ingress[0].ip;
//   console.log("LB IP:", _ip);
//   ip.value = _ip;
// });

// effect(async () => {
//   const res = await fetch(`http://${ip.value}/ip/`);
//   console.log(await res.text());
// });

// setInterval(() => {
//   active.value = !active.value;
// }, 15000);

// const app = createComponent(
//   "app",
//   z.object({
//     name: z.string(),
//     namespace: z.string(),
//     image: z.string(),
//     port: z.number(),
//   }),
//   (data) => {
//     const ns = createTemplate(() => ({
//       apiVersion: "v1",
//       kind: "Namespace",
//       metadata: {
//         name: data.value.namespace,
//         namespace: data.value.namespace,
//       },
//     }));
//     const pod = createTemplate(
//       () =>
//         ({
//           apiVersion: "v1",
//           kind: "Pod",
//           metadata: {
//             name: data.value.name,
//             namespace: ns.value.metadata.name,
//           },
//           spec: {
//             containers: [
//               {
//                 name: "my-container",
//                 image: data.value.image,
//                 ports: [
//                   {
//                     containerPort: data.value.port,
//                   },
//                 ],
//               },
//             ],
//           },
//         } satisfies K8sResource)
//     );
//     return [ns, pod];
//   },
//   { provider }
// );
