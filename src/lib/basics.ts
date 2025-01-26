import * as k8s from "@kubernetes/client-node";
import { effect, createSignal } from "./reactive.js";
import type { KubernetesObject } from "@kubernetes/client-node";
import { z, ZodSchema } from "zod";
import merge from "deepmerge";
import dotenv from "dotenv";
import type { Provider } from "./provider.js";
import { State } from "./state.js";

dotenv.config();
const stateId = process.env.KUBENEKO_STATE_ID;
const group = "kubeneko.adb.sh";

const state = new State();

state.subscribeToNewNamespace((namespace, provider) => {
  provider.watch.watch(
    `/api/v1/namespaces/${namespace}/events`,
    {},
    (type, event) => {
      console.log(`Type: ${type}`);
      console.log(`Event: ${JSON.stringify(event, null, 2)}`);
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
      _template.metadata?.namespace,
      _template.kind,
      _template.metadata?.name
    );
    template.value = _template;
  });
  return template as ReturnType<typeof createSignal<T>>;
};

const applyTemplate = async (
  template: KubernetesObject,
  { provider }: { provider: Provider }
) => {
  const current = await provider.objectApi.read(template).catch(() => null);

  if (current) {
    if (
      current.metadata.labels?.["app.kubernetes.io/managed-by"] !== "kubeneko"
    ) {
      throw new Error(
        "Resource already exists and is not managed by kubeneko!"
      );
    }
    if (current.metadata.annotations?.[`${group}/state-id`] !== stateId) {
      throw new Error("Resource is managed by another kubeneko state!");
    }
  }

  const active =
    template.metadata?.annotations?.[`${group}/active`] !== "false";

  if (active && current) {
    try {
      if (current) {
        const res = await provider.objectApi.replace(template);
        state.add(res, provider);
        return res;
      }
    } catch (e) {
      console.error(e);
    }
  } else if (active) {
    try {
      const res = await provider.objectApi.create(template);
      state.add(res, provider);
      return res;
    } catch (e) {
      console.error(e);
    }
  } else if (current) {
    try {
      const res = await provider.objectApi.delete(template);
      state.remove(res, provider);
      return null;
    } catch (e) {
      console.error(e);
    }
  }
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
        resource.value = await applyTemplate(newTemp, { provider });
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
        resource.value = await applyTemplate(newTemp, { provider });
      } else {
        resource.value = await applyTemplate(newTemp, { provider });
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
