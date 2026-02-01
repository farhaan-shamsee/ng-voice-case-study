# Kubernetes Manifests

This directory contains raw Kubernetes manifests for deploying the NG-Voice application.
These manifests can be applied directly using `kubectl` without Helm.
This also contains Kustomization files for easier management of resources.

## Prerequisites

- kubectl installed and configured
- Access to a Kubernetes cluster.

## Deployment

### 1. Namespace Creation

```bash
kubectl apply -k namespaces/
```

### 2. Database Deployment

```bash
kubectl apply -k database/
```

### 3. Web Server Deployment

```bash
kubectl apply -k web-server/
```