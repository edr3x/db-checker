# DB checker

App to check if your db connection is successful or not

## Requirements

### environment variables 

- `DB_URL` 
- `REDIS_URL` 

## Run via Docker

```sh
docker run --name db-checker -e DB_URL=postgresql://username:password@host:5432/database -e REDIS_URL=redis://0.0.0.0:6379 -d edr3x/db-check
```

## Run via k8s

### configmap.yaml
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-env-config
data:
  DB_URL: postgresql://user:pass@host:5432/app
  REDIS_URL: redis://host:6379
---
```

### deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: db-checker-deployment
  labels:
    app: db-checker
spec:
  replicas: 3
  selector:
    matchLabels:
      app: db-checker 
  template:
    metadata:
      labels:
        app: db-checker 
    spec:
      containers:
      - name: db-checker
        image: edr3x/db-check:0.2
        ports:
        - containerPort: 8080
        envFrom:
          - configMapRef:
              name: app-env-config
---
```

### service.yaml
```yaml
apiVersion: v1
kind: Service
metadata:
  name: db-checker-service
spec:
  selector:
    app: db-checker
  ports:
  - protocol: TCP
    port: 8080
    targetPort: 8080
---
```

Run all of them with `kubectl` and do the following

```sh
kubectl port-forward service/db-checker-service 8080:8080
```
