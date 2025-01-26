import * as k8s from "@kubernetes/client-node";
import { effect, createSignal } from "./reactive.js";
import type { KubernetesObject } from "@kubernetes/client-node";
import { z, ZodSchema } from "zod";
import merge from "deepmerge";
import dotenv from "dotenv";
import type { Provider } from "./provider.js";
import { State } from "./state.js";
import { applyTemplate } from "./applyTemplate.js";

dotenv.config();
const stateId = process.env.KUBENEKO_STATE_ID;
const group = "kubeneko.adb.sh";

const state = new State();

state.subscribeToNewNamespace((namespace, provider) => {
  provider.watch.watch(
    `/api/v1/namespaces/${namespace}/events`,
    {},
    async (type, event) => {
      // console.log(`Type: ${type}`);
      // console.log(`Event: ${JSON.stringify(event, null, 2)}`);
      const res = event.involvedObject;
      if (!res) return;
      if (type === "MODIFIED") {
        const stateEntry = state.state
          .get(provider)
          .resources.get(
            `${res.apiVersion}/${res.kind}/${res.namespace}/${res.name}`
          );
        if (!stateEntry) return;
        console.log(
          "remote resource changed:",
          stateEntry.resource.value.apiVersion,
          stateEntry.resource.value.kind,
          stateEntry.resource.value.metadata?.namespace,
          stateEntry.resource.value.metadata?.name
        );
        const resFromRemote = await provider.objectApi.read(
          stateEntry.resource.value
        );

        const isEqual =
          JSON.stringify((resFromRemote as any).spec) ===
          JSON.stringify((stateEntry.resource.value as any).spec);
        if (isEqual) {
          stateEntry.resource.value = resFromRemote;
        } else {
          console.log("yikes!");
          stateEntry.resource.value = await applyTemplate(
            stateEntry.template.value,
            { provider, group, stateId }
          );
        }
      }
    },
    (err) => {
      console.error("Error watching events:", err);
    }
  );
});

export const createTemplate = <T extends KubernetesObject | KubernetesObject>(
  getManifest: () => T
) => {
  const template = createSignal<T | null>(null);
  effect(() => {
    const _template = merge(getManifest(), {
      metadata: {
        labels: {
          "app.kubernetes.io/managed-by": "kubeneko",
        },
        annotations: {
          [`${group}/state-id`]: stateId,
        },
      },
    }) as T;
    console.log(
      `template updated:`,
      _template.apiVersion,
      _template.kind,
      _template.metadata?.namespace,
      _template.metadata?.name
    );
    template.value = _template;
  });
  return template as ReturnType<typeof createSignal<T>>;
};

export const createResource = <T extends KubernetesObject | KubernetesObject>(
  getTemplate: () => ReturnType<typeof createTemplate<T>>,
  { provider }: { provider: Provider }
) => {
  return new Promise<ReturnType<typeof createSignal<T>>>((resolve) => {
    const template = getTemplate();
    const resource = createSignal(null);

    effect(async () => {
      const oldTemp = template.oldValue;
      const newTemp = template.value;

      if (!oldTemp) {
        resource.value = await applyTemplate(newTemp, {
          provider,
          group,
          stateId,
        });
        state.add({ template, resource }, provider);
        resolve(resource);
        return;
      }

      if (
        oldTemp.apiVersion !== newTemp.apiVersion ||
        oldTemp.kind !== newTemp.kind
      ) {
        throw Error("cannot transform api version!");
      }
      if (
        oldTemp.metadata.name !== newTemp.metadata.name ||
        oldTemp.metadata.namespace !== newTemp.metadata.namespace
      ) {
        const res = await provider.objectApi.delete(oldTemp);
        state.remove(res, provider);
        resource.value = await applyTemplate(newTemp, {
          provider,
          group,
          stateId,
        });
        state.add({ template, resource }, provider);
      } else {
        resource.value = await applyTemplate(newTemp, {
          provider,
          group,
          stateId,
        });
        state.add({ template, resource }, provider);
      }
    });
  });
};

// export const createComponent = async <T>(
//   name: string,
//   schema: ZodSchema<T>,
//   getTemplates: (
//     data: ReturnType<typeof createSignal<T>>
//   ) => Array<ReturnType<typeof createTemplate>>,
//   { provider }: { provider: Provider }
// ) => {
//   const group = "kubeneko.adb.sh";
//   const crd = await createResource(
//     () => {
//       return createTemplate(() => ({
//         apiVersion: "apiextensions.k8s.io/v1",
//         kind: "CustomResourceDefinition",
//         metadata: {
//           name: `${name}.${group}`,
//         },
//         spec: {
//           group: group,
//           scope: "Namespaced",
//           names: {
//             plural: `${name}s`,
//             singular: name,
//             kind: `${name.charAt(0).toUpperCase()}${name.slice(1)}`,
//           },
//           version: [
//             {
//               name: "v1",
//               served: true,
//               storage: true,
//               schema: {
//                 openAPIV3Schema: schema,
//               },
//             },
//           ],
//         },
//       }));
//     },
//     { provider }
//   );

//   // return resource;
// };
