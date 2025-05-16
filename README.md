# kube-chache-api

A Go backend that provides a user-scoped, cached, paginated, and searchable API for Kubernetes pods. Built with [Gin](https://github.com/gin-gonic/gin) and [client-go](https://github.com/kubernetes/client-go), this project demonstrates efficient backend patterns for Kubernetes dashboards and developer tools.

---

## Features

- **User-scoped caching**: Each user (by `X-User-ID` header) gets isolated cached data.
- **Kubernetes integration**: Fetches real pod data from any K8s cluster (kind, minikube, cloud, etc).
- **TTL-based cache**: Reduces API load and latency, with automatic expiry (default: 60s).
- **Pagination**: Request only the pods you need per page.
- **Search**: Fast substring search on pod name, namespace, or status.
- **Stateless, easy to run**: No database required, just Go and access to a kubeconfig.

---

## Architecture

- **Gin** HTTP server
- **client-go** for Kubernetes API
- **In-memory cache** per user
- **Endpoints**:
  - `GET /pods` — paginated pod list
  - `GET /search` — search pods by query

---

## Setup & Running

### Prerequisites
- Go 1.18+
- Docker Desktop (for kind)
- [kind](https://kind.sigs.k8s.io/) (Kubernetes IN Docker)
- kubectl

### 1. Start a Local Kubernetes Cluster
```sh
kind create cluster --name demo
kubectl create deployment nginx --image=nginx
kubectl get pods
```
Wait until the nginx pod is `Running`.

### 2. Run the Backend
Clone this repo and run:
```sh
go run main.go cache.go handlers.go
```

The backend will use your current kubeconfig (e.g., kind, minikube, or any cluster).

---

## Usage

### Authentication
- Every request must include an `X-User-ID` header (simulate different users by changing its value).

### List Pods (Paginated)
```sh
curl -H "X-User-ID: alice" "http://localhost:8080/pods?page=1&limit=5"
```

### Search Pods
```sh
curl -H "X-User-ID: alice" "http://localhost:8080/search?q=nginx"
```

### Simulate Another User
```sh
curl -H "X-User-ID: bob" "http://localhost:8080/pods"
```

### Example Response
```json
{
  "limit": 5,
  "page": 1,
  "pods": [
    {"name": "nginx-123", "namespace": "default", "status": "Running"},
    ...
  ],
  "total": 10
}
```

---

## How Caching Works
- On first request per user, pods are fetched from Kubernetes and cached for 60 seconds.
- Subsequent requests within 60s are served from cache (fast).
- After 60s, the next request fetches fresh data and updates the cache.
- Each user’s cache is isolated (by `X-User-ID`).

---

## How Pagination & Search Work
- **Pagination**: `/pods?page=2&limit=5` returns the 2nd page of 5 pods.
- **Search**: `/search?q=nginx` returns pods whose name, namespace, or status contains `nginx` (case-insensitive).

---

## Troubleshooting
- **Pod stuck in Pending?**
  - Check Docker Desktop is running and has enough resources.
  - Run `kubectl describe pod <pod-name>` for details.
- **No pods returned?**
  - Make sure you have running pods in your cluster (`kubectl get pods -A`).
- **Kubernetes API errors?**
  - Ensure your kubeconfig is valid and points to the right cluster.
- **Cache not updating?**
  - Wait 60 seconds or restart the backend to force a refresh.

---

## Extending This Project
- Add support for other Kubernetes resources (Deployments, Services, etc).
- Add manual cache invalidation endpoints.
- Add metrics (cache hits/misses, request latency)
- Add authentication/authorization.
- Build a frontend to consume the API.

---

## License
MIT
