import type { Provider } from "./provider.js";
import type { KubernetesObject } from "@kubernetes/client-node";

export const applyTemplate = async (
  template: KubernetesObject,
  {
    provider,
    group,
    stateId,
  }: { provider: Provider; group: string; stateId: string }
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
        return await provider.objectApi.replace(template);
      }
    } catch (e) {
      console.error(e);
    }
  } else if (active) {
    try {
      return await provider.objectApi.create(template);
    } catch (e) {
      console.error(e);
    }
  } else if (current) {
    try {
      await provider.objectApi.delete(template);
      return null;
    } catch (e) {
      console.error(e);
    }
  }
};
