---
title: Getting Started
hide_next: true
---

# Getting Started

## Requirements

- You have a cluster that uses Kubernetes v1.24 or newer.
- You have installed the kubectl CLI.
- Your Kubernetes cluster must have access to the internet, because Kubernetes needs to be able to fetch images

## Install

To install the core components of cloudx running the command:
```
kubectl create ns cloudx-system
kubectl apply -n cloudx-system -f https://raw.githubusercontent.com/octopipe/cloudx/main/install/install-crds.yaml
kubectl apply -n cloudx-system -f https://raw.githubusercontent.com/octopipe/cloudx/main/install/install-base.yaml
```

## Install UI
Intall the UI running the command:
```
kubectl apply -n cloudx-system -f https://raw.githubusercontent.com/octopipe/cloudx/main/install/install-ui.yaml
```

To test the ui installation you can handle the application running the command:

```
kubectl port-forward service/cloudx-ui 3000:80 -n cloudx-system
```
