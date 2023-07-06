---
title: Getting Started
hide_next: true
---

## Requirements

## Install
```
kubectl create ns cloudx-system
kubectl apply -n cloudx-system -f https://raw.githubusercontent.com/octopipe/cloudx/main/install/install-crds.yaml
kubectl apply -n cloudx-system -f https://raw.githubusercontent.com/octopipe/cloudx/main/install/install-base.yaml
```

## Install UI
```
kubectl apply -n cloudx-system -f https://raw.githubusercontent.com/octopipe/cloudx/main/install/install-ui.yaml
```