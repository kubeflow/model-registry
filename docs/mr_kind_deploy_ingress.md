# Model Registry - Kind Ingress Guide

## Create a Kind cluster ready for the ingress controller

1. Create a file named `kind-config.yaml` with the following content:

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 3080
    hostPort: 3080
    protocol: TCP
  - containerPort: 30443
    hostPort: 30443
    protocol: TCP
```

> 📖 **NOTE**
>
> ContainerPorts 3080 and 30443 are customisable, you can change them to any other port number, but make sure to update the port number in the following kubectl patch commands.

2. Run the following command `kind create cluster --config=kind-config.yaml`

## Install the ingress controller (nginx) on the cluster

1. Install the controller by using `kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml`


2. Patch the ports inside the controller's deployment, by running the following commands:

```shell
kubectl patch deployment -n ingress-nginx ingress-nginx-controller   --type='json'   -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/ports/0/hostPort", "value": 3080}]'

kubectl patch deployment -n ingress-nginx ingress-nginx-controller   --type='json'   -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/ports/1/hostPort", "value": 30443}]'
```

## Install model registry on the cluster

`kubectl create namespace kubeflow`

`kubectl apply -k "https://github.com/kubeflow/model-registry/manifests/kustomize/overlays/db"`

`kubectl wait --for=condition=available -n kubeflow deployment/model-registry-deployment --timeout=1m`

## Apply the ingress

1. Create a file named `mr-ingress.yaml` with the following content:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: model-registry
spec:
  rules:
  - host: "model-registry.io" # choose a name of your liking
    http:
      paths:
      - pathType: Prefix
        path: "/"
        backend:
          service:
            name: model-registry-service
            port:
              number: 8080
```

2. Run the following command `kubectl apply -f mr-ingress.yaml -n kubeflow`

3. Add the following line to the file `/etc/hosts`:

`127.0.0.1 model-registry.io`

## Test the ingress

Run `curl http://model-registry.io:3080/api/model_registry/v1alpha3/registered_models`, you should see and output similar to this:

```json
{"items":[],"nextPageToken":"","pageSize":0,"size":0}
```

## Teardown

`kind delete cluster`
