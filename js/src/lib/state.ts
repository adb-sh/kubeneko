import * as k8s from "@kubernetes/client-node";
import type { KubernetesObject } from "@kubernetes/client-node";
import type { Provider } from "./provider.js";
import type { Signal } from "./reactive.js";

type StateEntry = {
  template: Signal<KubernetesObject>;
  resource: Signal<KubernetesObject>;
};

export class State {
  private _state: Map<
    Provider,
    { namespaces: Set<string>; resources: Map<string, StateEntry> }
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

  public add(entry: StateEntry, provider: Provider) {
    if (!this._state.has(provider)) {
      this._state.set(provider, {
        namespaces: new Set(),
        resources: new Map(),
      });
    }
    const res = entry.template.value;
    const s = this._state.get(provider);
    if (!s.namespaces.has(res.metadata.namespace)) {
      this.notify(res.metadata.namespace, provider);
    }
    s.namespaces.add(res.metadata.namespace);
    s.resources.set(
      `${res.apiVersion}/${res.kind}/${res.metadata.namespace}/${res.metadata.name}`,
      entry
    );
  }

  public remove(res: KubernetesObject, provider: Provider) {
    this._state
      .get(provider)
      .resources.delete(
        `${res.apiVersion}/${res.kind}/${res.metadata.namespace}/${res.metadata.name}`
      );
  }

  get state() {
    return this._state;
  }
}
