import * as k8s from "@kubernetes/client-node";

export class Provider {
  private _objectApi: k8s.KubernetesObjectApi;
  private _watch: k8s.Watch;

  constructor(kc: k8s.KubeConfig) {
    this._objectApi = kc.makeApiClient(k8s.KubernetesObjectApi);
    this._watch = new k8s.Watch(kc);
  }

  get objectApi() {
    return this._objectApi;
  }
  get watch() {
    return this._watch;
  }
}
