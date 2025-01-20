"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
const basics_ts_1 = require("./lib/basics.ts");
const reative_ts_1 = require("./lib/reative.ts");
const k8s = __importStar(require("@kubernetes/client-node"));
// Load the kubeconfig file
const kc = new k8s.KubeConfig();
kc.loadFromDefault();
// Get the current cluster
const cluster = kc.getCurrentCluster();
if (!cluster) {
    throw new Error("No cluster found in kubeconfig");
}
// Get the current user
const user = kc.getCurrentUser();
if (!user) {
    throw new Error("No user found in kubeconfig");
}
// Decode the CA, client cert, and key
// const ca = cluster.caData ? atob(cluster.caData) : undefined;
// const cert = user.certData ? atob(user.certData) : undefined;
// const key = user.keyData ? atob(user.keyData) : undefined;
// console.log("ca", ca);
// console.log("cert", cert);
// console.log("key", key);
// const agent = Deno.createHttpClient({
//   caCerts: ca ? [ca] : undefined,
//   key,
//   cert,
// });
// await kc.applyToFetchOptions({
//   ca: ca ? [ca] : undefined,
//   key,
//   cert,
// });
// const res = await fetch("https://10.10.1.0:6443/api/v1/namespaces", {client: agent}).then(body => body.json());
// console.log(res);
const provider = kc.makeApiClient(k8s.KubernetesObjectApi);
const image = (0, reative_ts_1.createSignal)("nginx");
const namespace = (0, basics_ts_1.createResource)(() => {
    const name = (0, reative_ts_1.createSignal)("kloud");
    return (0, basics_ts_1.createTemplate)(() => ({
        apiVersion: "v1",
        kind: "Namespace",
        metadata: {
            name: name.value,
            namespace: name.value,
        },
    }));
}, { provider });
const deployment = (0, basics_ts_1.createResource)(() => {
    const name = (0, reative_ts_1.createSignal)("test");
    return (0, basics_ts_1.createTemplate)(() => {
        var _a;
        return ({
            apiVersion: "v1",
            kind: "Pod",
            metadata: {
                name: name.value,
                namespace: (_a = namespace.value.metadata) === null || _a === void 0 ? void 0 : _a.name,
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
        });
    });
}, { provider });
// const deployment = createResource(
//   () => {
//     const name = createSignal("test");
//     return createTemplate(
//       () =>
//         ({
//           apiVersion: "apps/v1",
//           kind: "Deployment",
//           metadata: {
//             name: name.value,
//             namespace: namespace.value.metadata?.name,
//           },
//           spec: {
//             template: {
//               spec: {
//                 containers: [
//                   {
//                     image: image.value,
//                     name: "app",
//                     ports: [
//                       {
//                         containerPort: 80,
//                         name: "http",
//                         protocol: "TCP",
//                       },
//                     ],
//                   },
//                 ],
//               },
//             },
//           },
//         } satisfies K8sResource)
//     );
//   },
//   { provider }
// );
