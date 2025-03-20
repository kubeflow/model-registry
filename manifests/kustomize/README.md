# Install Kubeflow Model Registry

This folder contains [Kubeflow Model Registry](https://www.kubeflow.org/docs/components/model-registry/installation/) Kustomize manifests

## Installation

To install Kubeflow Model Registry, follow [Kubeflow Model Registry deployment documentation](https://www.kubeflow.org/docs/components/model-registry/installation/).

> **⚠️ Warning:** To install Model Registry in a way that's compatible with the Model Registry UI, please follow the instructions in the [Installing Model Registry on a Kubeflow Profile](#installing-model-registry-on-a-kubeflow-profile) section below.

The following instructions will summarize how to deploy Model Registry as separate component in the context of a default Kubeflow >=1.9 installation.
Ensure you are running these commands from the directory containing this README.md file (e.g.: you could check with `pwd`).

```bash
kubectl apply -k overlays/db
```

As the default Kubeflow installation provides an Istio mesh, apply the necessary manifests:

```bash
kubectl apply -k options/istio
```

Check everything is up and running:

```bash
kubectl wait --for=condition=available -n kubeflow deployment/model-registry-deployment --timeout=2m
kubectl logs -n kubeflow deployment/model-registry-deployment
```

Optionally, you can also port-forward the REST API container port of Model Registry to interact with it from your terminal:

```bash
kubectl port-forward svc/model-registry-service -n kubeflow 8081:8080
```

And then, from another terminal:

```bash
curl -sX 'GET' \
  'http://localhost:8081/api/model_registry/v1alpha3/registered_models?pageSize=100&orderBy=ID&sortOrder=DESC' \
  -H 'accept: application/json' | jq
```

### UI Installation

To install the Model Registry UI as a Kubeflow component, you need first to deploy the Model Registry UI:

```bash
kubectl apply -k options/ui/overlays/istio -n kubeflow
```

And then to make it accessible through Kubeflow Central Dashboard, you need to edit the `centraldashboard-config` ConfigMap to add the Model Registry UI link to the Central Dashboard by running the following command:

```bash
kubectl get configmap centraldashboard-config -n kubeflow -o json | jq '.data.links |= (fromjson | .menuLinks += [{"icon": "assignment", "link": "/model-registry/", "text": "Model Registry", "type": "item"}] | tojson)' | kubectl apply -f - -n kubeflow
```

Alternatively, you can edit the ConfigMap manually by running:

```bash
kubectl edit configmap -n kubeflow centraldashboard-config
```

```yaml
apiVersion: v1
data:
  links: |-
    {
        "menuLinks": [
            {
                "icon": "assignment",
                "link": "/model-registry/",
                "text": "Model Registry",
                "type": "item"
            },
            ...
```

#### Installing Model Registry on a Kubeflow Profile (namespace)

Kubeflow Central Dashboard uses [Profiles](https://www.kubeflow.org/docs/components/central-dash/profiles/) to handle user namespaces and permissions. By default, the manifests deploy the Model Registry instance in the `kubeflow` namespace, to install a compatible version of Model Registry for Kubeflow, you should deploy a separate instance of Model Registry in the profile's namespace. For that just run:

```sh
PROFILE_NAME=<your-profile>
for DIR in options/istio overlays/db ; do (cd $DIR; kustomize edit set namespace $PROFILE_NAME; kubectl apply -k .); done

## Usage

For a basic usage of the Kubeflow Model Registry, follow the [Kubeflow Model Registry getting started documentation](https://www.kubeflow.org/docs/components/model-registry/getting-started/)

## Uninstall

To uninstall the Kubeflow Model Registry run:

```bash
# Delete istio options
kubectl delete -k options/istio

# Delete model registry db and deployment
kubectl delete -k overlays/db
```
