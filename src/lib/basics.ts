import * as k8s from "@kubernetes/client-node";
import { effect, createSignal } from "./reative.js";
import type { KubernetesObject } from "@kubernetes/client-node";

export type K8sResource =
  | k8s.V1Pod
  | k8s.V1Service
  | k8s.V1Deployment
  | k8s.V1ReplicaSet
  | k8s.V1StatefulSet
  | k8s.V1DaemonSet
  | k8s.V1Job
  | k8s.V1CronJob
  | k8s.V1ConfigMap
  | k8s.V1Secret
  | k8s.V1PersistentVolume
  | k8s.V1PersistentVolumeClaim
  | k8s.V1Namespace
  | k8s.V1Ingress
  | k8s.V1NetworkPolicy
  // | k8s.V1Event
  | k8s.V1LimitRange
  | k8s.V1ResourceQuota
  | k8s.V1PodTemplate;

export const createTemplate = (getManifest: () => KubernetesObject) => {
  const template = createSignal<K8sResource | null>(null);
  effect(() => {
    const _template = getManifest();
    console.log(
      `template updated:`,
      _template.apiVersion,
      _template.kind,
      _template.metadata?.namespace,
      _template.metadata?.name
    );
    template.value = _template;
  });
  return template as ReturnType<typeof createSignal<K8sResource>>;
};

const applyTemplate = async (
  template: K8sResource,
  { provider }: { provider: k8s.KubernetesObjectApi }
) => {
  const active =
    template.metadata?.annotations?.["kubeneko/active"] !== "false";
  if (active) {
    try {
      await provider.read(template);
      return await provider.patch(template);
    } catch (e) {
      return await provider.create(template);
    }
  } else {
    try {
      await provider.read(template);
      return await provider.delete(template);
    } catch (e) {
      console.error(e);
    }
  }
};

export const createResource = async <T>(
  getTemplate: () => ReturnType<typeof createTemplate>,
  { provider }: { provider: k8s.KubernetesObjectApi }
) => {
  const template = getTemplate();
  const resource = createSignal(
    (await applyTemplate(template.value, { provider })) as any
  );

  effect(async () => {
    const oldTemp = template.oldValue;
    const newTemp = template.value;

    if (!oldTemp) {
      resource.value = await applyTemplate(newTemp, { provider });
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
      await provider.delete(oldTemp);
      resource.value = await applyTemplate(newTemp, { provider });
    } else {
      resource.value = await applyTemplate(newTemp, { provider });
    }
  });

  return resource;
};
