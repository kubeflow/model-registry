# Set up local development with kubeflow

**Note: this should only be needed in edge cases in which we need to test a local integration with the kubeflow dashboard.**

## Prerequisites

- [Kubeflow repo](https://github.com/kubeflow/kubeflow/tree/master/components/centraldashboard#development)
- [Model Registry repo](../README.md)

## Setup

### Kubeflow repo (under centraldashboard)

1. Change the [webpack config](https://github.com/kubeflow/kubeflow/blob/master/components/centraldashboard/webpack.config.js#L186) proxies to allow Model Registry:

```javascript
        proxy: {
            '/api/v1': {
                target: 'http://localhost:4000',
            },
            '/api': {
                target: 'http://localhost:8082',
            },
            '/jupyter': {
                target: 'http://localhost:8085',
                pathRewrite: {'^/jupyter': ''},
            },
            // NOTE: this makes `/notebook` requests fail with a 504 error
            '/notebook': {
                target: 'http://localhost:8086',
                pathRewrite: {
                    '^/notebook/(.*?)/(.*?)/(.*)':
                        '/$1/services/$2/proxy/notebook/$1/$2/$3',
                },
            },
            '/pipeline': {
                target: 'http://localhost:8087',
                pathRewrite: {'^/pipeline': ''},
            },
            '/model-registry': {
                target: 'http://localhost:9000',
                changeOrigin: true,
                ws: true,
                secure: false,
            },
        },
```

2. Run the centraldashboard:

```shell
npm run dev
```

### Model Registry repo

1. Just run the repo in kubeflow dev mode

```shell
make dev-start-kubeflow
```

### Access the cluster

You need to have a kubeflow cluster up and running, to get the Model Registry working you'll need to port-forward these two services:

```shell
kubectl port-forward service/model-registry-service 8085:8080 -n <targeted-namespace-of-the-mr-service>
```

```shell
kubectl port-forward svc/profiles-kfam 8081:8081 -n kubeflow
```
