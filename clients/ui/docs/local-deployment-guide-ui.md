[Model registry server set up]: ../../bff/docs/dev-guide.md

# Deploying the Model Registry UI in a local cluster

For this guide, we will be using kind for locally deploying our cluster. See
the [Model registry server set up] guide for prerequisites on setting up kind
and deploying the model registry server.

## Setup

### 1. Create a kind cluster

Create a local cluster for running the MR UI using the following command:

```shell
kind create cluster
```

### 2. Create kubeflow namespace

Create a namespace for model registry to run in, by default this is kubeflow, run:

```shell
kubectl create namespace kubeflow
```

### 3. Build a standalone image for the UI

Right now, the default image is targeted for the KF Central Dashboard. To build a standalone image for the UI, run:

```shell
make docker-build-standalone
make docker-push-standalone
```

**Note: You will need to set up `IMG_UI_STANDALONE` in your .env.local file to push the image to your own registry.**

### 4. Deploy Model Registry UI to cluster

You can now deploy the UI and BFF to your newly created cluster using the kustomize configs in the root manifest directory:

First you need to set up your new image

```shell
cd manifests/kustomize/options/ui/base
kustomize edit set image model-registry-ui=${IMG_UI_STANDALONE}
```

Now you can set the namespace to kubeflow and apply the manifests:

```shell
cd manifests/kustomize/options/ui/overlays/standalone
kustomize edit set namespace kubeflow
kubectl apply -k .
```

After a few seconds you should see 1 pod running:

```shell
kubectl get pods -n kubeflow
```

```shell
NAME                                  READY   STATUS    RESTARTS   AGE
model-registry-ui-58755c4754-zdrnr    1/1     Running   0          11s
```

### 5. Access the Model Registry UI running in the cluster

Now that the pods are up and running you can access the UI.

First you will need to port-forward the UI service by running the following in it's own terminal:

```shell
kubectl port-forward service/model-registry-ui-service 8080:8080 -n kubeflow
```

You can then access the UI running in your cluster locally at http://localhost:8080/

You can now make API requests to the BFF endpoints like:

```shell
curl http://localhost:8080/api/v1/model-registry
```

```json
{
    "model_registry": null
}
```

## Troubleshooting

### Running on macOS

When running locally on macOS you may find the pods fail to deploy, with one or more stuck in the `pending` state. This is usually due to insufficient memory allocated to your docker / podman virtual machine. You can verify this by running:

```shell
kubectl describe pods -n kubeflow
```

If you're experiencing this issue you'll see an output containing something similar to the following:

```shell
Events:
  Type     Reason            Age   From               Message
  ----     ------            ----  ----               -------
  Warning  FailedScheduling  29s   default-scheduler  0/1 nodes are available: 1 Insufficient memory. preemption: 0/1 nodes are available: 1 No preemption victims found for incoming pod.
```

To fix this, you'll need to increase the amount of memory available to the VM. This can be done through either the Podman Desktop or Docker Desktop GUI. 6-8GB of memory is generally a sufficient amount to use.

# Running with Kubeflow and Istio

Alternatively, if you'd like to run the UI and BFF pods with an Istio configuration for the KF Central Dashboard, you can apply the manifests by running:

```shell
kubectl apply -k manifests/kustomize/options/ui/overlays/istio -n kubeflow
```
