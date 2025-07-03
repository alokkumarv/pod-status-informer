# 🚦 Pod Health Dashboard Controller

A lightweight Kubernetes controller built in Go that watches Pods in the cluster and exposes a live REST API to display their health status.

This project uses Kubernetes **client-go informers** to efficiently cache and monitor real-time pod updates, making it scalable and production-ready.

---

## 🎯 Features

- ✅ Watches all Pods in the cluster using informers
- ✅ Supports real-time caching of pod health info
- ✅ Exposes a REST API:
  - `GET /pods?labelSelector=app=nginx` — List matching pods with status
  - `GET /summary?labelSelector=app=nginx` — Count of pods by status
- ✅ Query-time filtering using label selectors

---

## 🚀 Demo




Sample Output:

`/pods`
```json
[
  {
    "name": "nginx-deploy-68f4f7b4f5-6rszn",
    "namespace": "default",
    "status": "Running"
  },
  {
    "name": "nginx-deploy-68f4f7b4f5-vmx78",
    "namespace": "default",
    "status": "Pending"
  }
]


{
  "Running": 2,
  "Pending": 1,
  "Failed": 0
}

