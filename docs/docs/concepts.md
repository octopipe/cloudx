---
title: Concepts
hide_next: true
---

# Concepts

Below are some conecepts that are specific for Cloudx

**SharedInfra**: The sharedinfra represents the infraestructure diagram with the steps and orchestration by dependencies, this is a CRD.

**ConnectionInterface**: The connection interface represents the outputs of a specific plugin executed from sharedinfra, this is a CRD.

**ProviderConfig**: The providerconfig is a set of configurations to create the cloudx connection with the cloud environment, this is a CRD.

**Execution**: A execution is a manifest generated from sharedinfra, its represents the execution of sharedinfra by cloudx engine.

**Webhook**: The webhook is a configration set used by the cloudx engine to call the target defined in the manifests, this is a CRD.

