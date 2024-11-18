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

## Environment Variables

The following environment variables are used to configure the deployment and development environment for the Model Registry UI. These variables should be defined in a `.env.local` file in the `clients/ui` directory of the project. **This values will affect the build and push commands**.

### `CONTAINER_TOOL`

* **Description**: Specifies the container tool to be used for building and running containers.
* **Default Value**: `docker`
* **Possible Values**: `docker`, `podman`, etc.
* **Example**: `CONTAINER_TOOL=docker`

### `IMG_BFF`

* **Description**: Specifies the image name and tag for the Backend For Frontend (BFF) service.
* **Default Value**: `model-registry-bff:latest`
* **Example**: `IMG_BFF=model-registry-bff:latest`

### `IMG_FRONTEND`

* **Description**: Specifies the image name and tag for the frontend service.
* **Default Value**: `model-registry-frontend:latest`
* **Example**: `IMG_FRONTEND=model-registry-frontend:latest`

### Example `.env.local` File

Here is an example of what your `.env.local` file might look like:

```shell
CONTAINER_TOOL=docker
IMG_BFF=model-registry-bff:latest
IMG_FRONTEND=model-registry-frontend:latest
```

## Build and Push Commands

The following Makefile targets are used to build and push the Docker images for the Backend For Frontend (BFF) and frontend services. These targets utilize the environment variables defined in the `.env.local` file.

### Build Commands

* **`build-bff`**: Builds the Docker image for the BFF service.
  * Command: `make build-bff`
  * This command uses the `CONTAINER_TOOL` and `IMG_BFF` environment variables to build the image.

* **`build-frontend`**: Builds the Docker image for the frontend service.
  * Command: `make build-frontend`
  * This command uses the `CONTAINER_TOOL` and `IMG_FRONTEND` environment variables to build the image.

* **`build`**: Builds the Docker images for both the BFF and frontend services.
  * Command: `make build`
  * This command runs both `build-bff` and `build-frontend` targets.

### Push Commands

* **`push-bff`**: Pushes the Docker image for the BFF service to the container registry.
  * Command: `make push-bff`
  * This command uses the `CONTAINER_TOOL` and `IMG_BFF` environment variables to push the image.

* **`push-frontend`**: Pushes the Docker image for the frontend service to the container registry.
  * Command: `make push-frontend`
  * This command uses the `CONTAINER_TOOL` and `IMG_FRONTEND` environment variables to push the image.

* **`push`**: Pushes the Docker images for both the BFF and frontend services to the container registry.
  * Command: `make push`
  * This command runs both `push-bff` and `push-frontend` targets.
