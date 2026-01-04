# Siteforge

Siteforge is a small Go-based control plane that automates tenant isolation and application runtime deployment on Kubernetes.

This project was built as a hands-on learning exercise to understand how platform teams manage multi-tenant runtimes, enforce isolation, and reconcile desired state in production systems.

---

## Motivation

Modern multi-tenant platforms run thousands of customer applications on shared infrastructure. To do this safely and reliably, platform teams need to:

- isolate tenants
- manage application lifecycles
- enforce resource limits
- gate traffic based on readiness
- recover cleanly from failures
- reconcile desired state continuously

Siteforge explores these ideas in a minimal, local environment using Go and Kubernetes.

---

## What This Project Does

- Provisions tenants as Kubernetes namespaces
- Deploys application runtimes (WordPress) per tenant
- Manages runtimes using Deployments and Services
- Applies CPU and memory requests/limits for isolation
- Uses readiness probes to control traffic safely
- Reconciles updates instead of relying on one-time creation
- Recovers from deleted resources and cluster restarts

---

## High-Level Architecture

```
CLI (Go)
  ↓
Control Plane Logic
  ↓
Kubernetes API
  ↓
Namespace (Tenant)
   ├── Deployment (Runtime)
   │     └── Pod (Container)
   └── Service (Networking)
```

- **Tenant** → Kubernetes Namespace  
- **Runtime** → Deployment + Service  
- **Control Plane** → Go program using `client-go`

---

## What I Learned

- How control planes interact with the Kubernetes API
- The roles of Pods, Deployments, and Services in production
- Why readiness probes matter and how they can fail
- The difference between idempotent creation and reconciliation
- How Kubernetes recovers from failure by design
- How to structure Go code for long-running infrastructure services

---

## Running Locally

### Requirements
- Go 1.21+
- Docker
- kind
- kubectl

### Steps

```bash
kind create cluster --name siteforge
go run ./cmd/siteforge --tenant blog-aku
kubectl port-forward svc/wordpress -n blog-aku 8080:80
```

Then open: http://localhost:8080

---

## Future Improvements

- Persistent storage and database integration
- Runtime status reporting from the control plane
- Structured logging and basic observability
- CI/CD-driven runtime updates
