# Go Kubernetes Pod Watcher Application

A simple Go application that watches Pod events in a Kubernetes cluster and logs Added, Modified, and Deleted events.

## Overview

This application demonstrates:

- Using the Kubernetes `client-go` library
- Watching resources (Pods) across all namespaces witbh label filtering.
- Handling watch events (Added, Modified, Deleted)
- Automatic reconnection on failures
- Both in-cluster and out-of-cluster configurations

## Prerequisites

- Go 1.23 or later
- kubectl access to a Kubernetes cluster
- Docker (for building container image)

## Running Locally

### 1. Install Dependencies

```bash
cd applications/pod-monitor
go mod download
```

### 2. Run Against Your Cluster

```bash
# Ensure kubectl is configured
kubectl cluster-info

# Run the watcher
go run main.go
```

You'll see output like:

```sh
2026/02/01 12:25:55 Starting pod watcher...
2026/02/01 12:25:55 [CREATED] database/mysql-0 | Phase=Running
2026/02/01 12:25:55 [CREATED] database/mysql-1 | Phase=Running
2026/02/01 12:25:55 [CREATED] web-server/web-server-ddbb87f44-7mkww | Phase=Running
2026/02/01 12:25:55 [CREATED] web-server/web-server-ddbb87f44-fklml | Phase=Running
2026/02/01 12:25:55 [CREATED] web-server/web-server-ddbb87f44-rg4kx | Phase=Running
```

## Building

### Build Binary

```bash
go build -o pod-watcher main.go
./pod-watcher
```

### Build Docker Image

```bash
docker build -t ng-voice/pod-watcher:v1.0 .
```

## Deploying to Kubernetes

### 1. Create Service Account and RBAC

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: pod-watcher
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: pod-watcher
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: pod-watcher
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: pod-watcher
subjects:
- kind: ServiceAccount
  name: pod-watcher
  namespace: default
EOF
```

### 2. Deploy the Controller

```bash
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pod-watcher
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: pod-watcher
  template:
    metadata:
      labels:
        app: pod-watcher
    spec:
      serviceAccountName: pod-watcher
      containers:
      - name: controller
        image: farhaanshamsee/pod-watcher:v1.0
        imagePullPolicy: IfNotPresent
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "100m"
EOF
```

### 3. View Logs

```bash
kubectl logs -f deployment/pod-watcher
```

## Code Structure

```go
main()
├── getConfig()           // Load kubeconfig (in-cluster or local)
├── watchPods()           // Watch pod events
│   └── handleEvent()     // Process each event
└── Reconnection loop     // Automatic reconnection
```

### Key Functions

**getConfig():**
- Tries in-cluster config first (when running in Kubernetes)
- Falls back to local kubeconfig (~/.kube/config)

**watchPods():**
- Creates a watcher for all pods across all namespaces filtered by label `project=ng-voice`
- Listens for Added, Modified, Deleted events
- Returns on error (handled by reconnection loop)

## How It Works

1. **Configuration**: Loads Kubernetes client configuration
2. **Watch Setup**: Creates a watch on pods resource
3. **Event Loop**: Processes events as they arrive
4. **Reconnection**: Automatically reconnects if watch fails
