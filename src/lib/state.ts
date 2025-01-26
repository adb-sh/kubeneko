import * as k8s from "@kubernetes/client-node";
import type { KubernetesObject } from "@kubernetes/client-node";
import type { Provider } from "./provider.js";

export class State {
  private _state: Map<
    Provider,
    { namespaces: Set<string>; resources: Map<string, KubernetesObject> }
  >;
  private subscribers: Set<Function>;

  constructor() {
    this._state = new Map();
    this.subscribers = new Set();
  }

  public subscribeToNewNamespace(
    fn: (namespace: string, provider: Provider) => void
  ) {
    this.subscribers.add(fn);
  }

  private notify(namespace: string, provider: Provider) {
    this.subscribers.forEach((fn) => fn(namespace, provider));
  }

  public add(res: KubernetesObject, provider: Provider) {
    if (!this._state.has(provider)) {
      this._state.set(provider, {
        namespaces: new Set(),
        resources: new Map(),
      });
    }
    const s = this._state.get(provider);
    if (!s.namespaces.has(res.metadata.namespace)) {
      this.notify(res.metadata.namespace, provider);
    }
    s.namespaces.add(res.metadata.namespace);
    s.resources.set(
      `${res.apiVersion}/${res.kind}/${res.metadata.namespace}/${res.metadata.name}`,
      res
    );
  }

  public remove(res: KubernetesObject, provider: Provider) {
    this._state
      .get(provider)
      .resources.delete(
        `${res.apiVersion}/${res.kind}/${res.metadata.namespace}/${res.metadata.name}`
      );
  }
}
