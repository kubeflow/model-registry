# Kind Cluster Teardown

Stop all running model-registry dev environment processes. Run each step sequentially.

## Arguments

Options (space-separated after the command). If none specified, stop everything.

- **keep-cluster** — Stop BFF, Frontend, and Tilt but keep the Kind cluster and Colima running

## Stop Frontend (port 9000)
```
lsof -ti:9000 | xargs kill -9 2>/dev/null || true
```

## Stop BFF (port 4000)
```
lsof -ti:4000 | xargs kill -9 2>/dev/null || true
```

## Stop Tilt (port 10350)
```
cd devenv && make tilt-down
```

## Delete Kind cluster (skip if `keep-cluster`)
```
kind delete cluster --name model-registry
```

## Stop Colima (skip if `keep-cluster`)
```
colima stop
```

## Examples

- `/kind-cluster-teardown` — Stop everything and delete the cluster
- `/kind-cluster-teardown keep-cluster` — Stop processes but keep cluster for quick restart
