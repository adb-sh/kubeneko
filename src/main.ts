import { createResource, createTemplate, K8sResource } from "./lib/basics.js";
import { createSignal } from "./lib/reative.js";
import * as k8s from "@kubernetes/client-node";
import { effect } from "./lib/reative.js";

// Load the kubeconfig file
const kc = new k8s.KubeConfig();
kc.loadFromDefault();
const provider = kc.makeApiClient(k8s.KubernetesObjectApi);

const image = createSignal("nginx");
const active = createSignal(true);

const namespace = await createResource(
  () => {
    const name = createSignal("kubeneko");

    return createTemplate(() => ({
      apiVersion: "v1",
      kind: "Namespace",
      metadata: {
        name: name.value,
        namespace: name.value,
      },
    }));
  },
  { provider }
);

const pod = await createResource(
  () => {
    const name = createSignal("test");

    // let i = 0;
    // setInterval(() => {
    //   name.value = `test${++i}`;
    // }, 15000);

    return createTemplate(
      () =>
        ({
          apiVersion: "v1",
          kind: "Pod",
          metadata: {
            name: name.value,
            namespace: namespace.value.metadata?.name,
            annotations: {
              "kubeneko/active": active.value.toString(),
            },
          },
          spec: {
            containers: [
              {
                name: "my-container",
                image: image.value,
                ports: [
                  {
                    containerPort: 80,
                  },
                ],
              },
            ],
          },
        } satisfies K8sResource)
    );
  },
  { provider }
);

effect(() => {
  console.log(pod.value);
});

setInterval(() => {
  active.value = !active.value;
}, 15000);
