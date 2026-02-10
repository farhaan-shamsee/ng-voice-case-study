# NG-Voice Helm Chart

A Helm chart for deploying the NG-Voice DevOps case study application with MySQL database cluster and Nginx web server.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- Storage class available for persistent volumes (default: `standard`)
- Nodes labeled with `db-node=true` for MySQL pods

## Components

This Helm chart deploys the following components:

### 1. Database (MySQL Cluster via Bitnami Chart)
- MySQL cluster with replication architecture
  - 1 Primary node (read-write)
  - 2 Secondary/Replica nodes (read-only)

### 2. Web Server (Nginx)
- Nginx deployment with 3 replicas
- Init container to generate hostname-based content
- ConfigMap for Nginx configuration
- NodePort service for external access
- Readiness probes

### 3. Networking
- Network policies to isolate database access
- Multus CNI secondary network support

## Installation

### 1. Add Bitnami Helm Repository

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

### 2. Update Dependencies

Before installing, download the MySQL chart dependency:

```bash
cd helm-charts/ng-voice
helm dependency update
```

This will download the Bitnami MySQL chart to the `charts/` directory.

### 3. Label Nodes for MySQL

MySQL pods require nodes labeled with `db-node=true`:

```bash
kubectl label nodes <node-name> db-node=true
```

### 4. Install the Chart

```bash
# Install with default values
helm install ng-voice . -n default --create-namespace

# Or with custom values
helm install ng-voice . -f custom-values.yaml -n default --create-namespace
```

### 5. Verify Installation

```bash
# Check all pods
kubectl get pods -A

# Check MySQL cluster status
kubectl get pods -n database -l app.kubernetes.io/name=mysql

# Check web server
kubectl get pods -n web-server
```

## Uninstall

```bash
helm uninstall ng-voice -n default
```

Note: PVCs are not automatically deleted. To remove them:

```bash
kubectl delete pvc -n database -l app.kubernetes.io/name=mysql
kubectl delete pvc -n database mysql-backup
```
