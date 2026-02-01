# NG-Voice Helm Chart

A Helm chart for deploying the NG-Voice case study application with MySQL database and Nginx web server.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- Storage class available for persistent volumes (default: `local-path`)

## Components

This Helm chart deploys the following components:

### 1. Database (MySQL)
- MySQL 8.0 StatefulSet with 2 replicas
- Headless service for stable network identities
- Persistent storage for data
- Automated backups (optional)
  - Internal backup to PVC
  - S3 backup (configurable)
- Network policy for access control

### 2. Web Server (Nginx)
- Nginx deployment with 3 replicas
- Init container to generate hostname-based content
- ConfigMap for Nginx configuration
- NodePort service for external access
- Readiness probes

### 3. Networking
- Network policies to isolate database access
- Multus CNI secondary network support
