[frontend requirements]: ./frontend/docs/dev-setup.md#requirements
[BFF requirements]: ./bff/README.md#pre-requisites
[frontend dev setup]: ./frontend/docs/dev-setup.md#development
[BFF dev setup]: ./bff/README.md#development

# Model Registry UI

## Overview

The Model Registry UI is a standalone web app for Kubeflow Model Registry. In this repository, you will find the frontend and backend for the Model Registry UI.

## Prerequisites

* [Frontend requirements]
* [BFF requirements]

## Set Up

### Development

To run the a mocked dev environment you can either:

* Use the makefile command to install dependencies `make dev-install-dependencies` and then start the dev environment `make dev-start`.

* Or follow the [frontend dev setup] and [BFF dev setup].

### Docker deployment

To build the Model Registry UI container, run the following command:

```shell
make docker-compose
```

### Kubernetes Deployment

For a in-depth guide on how to deploy the Model Registry UI, please refer to the [local kubernetes deployment](./bff/docs/dev-guide.md) documentation.

To quickly enable the Model Registry UI in your Kind cluster, you can use the following command:

```shell
make kind-deployment
```

## OpenAPI Specification

You can find the OpenAPI specification for the Model Registry UI in the [openapi](./api/openapi) directory.
A live version of the OpenAPI specification can be found [here](https://editor.swagger.io/?url=https://raw.githubusercontent.com/kubeflow/model-registry/main/clients/ui/api/openapi/mod-arch.yaml).
