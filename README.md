Kubeneko
===

kubeneko is a POC reactive Kubernetes framework using, built around CUE and Go. It aims to provide a more declarative, type-safe way to define and manage Kubernetes resources via schemas and signals.

## What is this

This is an early-stage proof of concept.

The go-cue branch combines Go and CUE to validate and generate Kubernetes manifests based on schema constraints, making manifest definitions safer and less error-prone.

The idea: define resource shape and constraints with CUE, then leverage code generation / validation in Go for rendering or verifying Kubernetes YAML / manifest files.

Using CUE for Kubernetes manifests is a known pattern: the CUE ecosystem offers Kubernetes schema support for types and CRDs, enabling validation and structured config generation.

## Repository Structure
```
/go   → go playground (controllers / helpers / generators) 
/js   → legacy js playground
```
*Note: Branch is currently go-cue.*

## Why CUE + Kubernetes

CUE is a data-constraint language that works well with structured configuration instead of free-form YAML.

With this setup you can define strong typing / validation for Kubernetes manifests. This reduces runtime errors compared to plain YAML.

Good foundation for building a reactive framework around Kubernetes, where "signals" or config changes generate or reconcile resources automatically.


## When / Where to Use

This is a good fit if you:

Want schema-driven manifest definitions rather than manually writing raw YAML.

Prefer type-safety and validation upfront (compile-time / config-time) rather than discovering errors at runtime.

Build larger Kubernetes setups where configuration drift or misconfigurations are risky.

Want to experiment with a GitOps / declarative + programmatic approach using Go + CUE.

## Status & What’s Missing

This is not production-ready — it’s a POC / early prototype.

No published releases or package distributions yet.

Limited documentation and usage examples so far.
