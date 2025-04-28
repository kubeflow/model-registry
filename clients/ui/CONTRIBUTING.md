[frontend requirements]: ./frontend/docs/dev-setup.md#requirements
[BFF requirements]: ./bff/README.md#pre-requisites
[frontend dev setup]: ./frontend/docs/dev-setup.md#development
[BFF dev setup]: ./bff/README.md#development
[issue]: https://github.com/kubeflow/model-registry/issues/new/choose
[contributing guidelines]: https://github.com/kubeflow/model-registry/blob/main/CONTRIBUTING.md
# Contributing

Individual bug fixes are welcome. Please open an [issue] to track the fix you are planning to implement. If you are unsure how best to solve it, start by opening the issue and note your desire to contribute.
We have [contributing guidelines] available for you to follow.

## Requirements

To review the requirements, please refer to:

* [Frontend requirements]
* [BFF requirements]

## Set Up

### Development

To run the mocked development environment you can either:

* Use the makefile command to install dependencies `make dev-install-dependencies`, and then start the dev environment with `make dev-start`.

* Or follow the steps in the [frontend dev setup] and [BFF dev setup] guides.

### Kubernetes Deployment

For an in-depth guide on how to deploy the Model Registry UI, please refer to the [local kubernetes deployment](./docs/local-deployment-guide.md) documentation.

To quickly enable the Model Registry UI in your Kind cluster, you can use the following command:

```shell
make kind-deployment
```

## Debugging and Testing

See [frontend testing guidelines](frontend/docs/testing.md) for testing the frontend.