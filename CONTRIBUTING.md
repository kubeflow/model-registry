# Kubeflow Contributor Guide

Welcome to the Kubeflow Model Registry project! We'd love to accept your patches and 
contributions to this project. Please read the 
[contributor's guide in our docs](https://www.kubeflow.org/docs/about/contributing/).

The contributor's guide:

* explains how to use the Developer Certificate of Origin (DCO) with git commit system
  * like explained in [this guide](https://wiki.linuxfoundation.org/dco), or [this example](https://github.com/kubeflow/community/tree/master/dco-signoff-hook#signing-off-commits)
* references the [Code of Conduct](https://www.kubeflow.org/docs/about/contributing/#follow-the-code-of-conduct)
* helps you get started with your first contribution to Kubeflow
  * like how to identify [good-first-issue](https://github.com/kubeflow/model-registry/labels/good%20first%20issue) tickets
* and [describes](https://www.kubeflow.org/docs/about/contributing/#owners-files-and-pr-workflow) the pull request and review workflow in detail, including the
  OWNERS files and automated workflow tool

The [Kubeflow Community calendar](https://www.kubeflow.org/docs/about/community/#kubeflow-community-calendar) provides an overview of the community calls and specifically the bi-weekly Kubeflow Model Registry meeting.

<hr/>

The remainder of this document focuses on technical aspects while contributing to the Model Registry project specifically.

The [Model Registry documentation](https://github.com/kubeflow/model-registry?tab=readme-ov-file#documentation-links) provides details about: what is a Model Registry, how to Install, logical model Concepts, how the MR python client works, Tutorials, FAQs and most importantly the [Technical References](https://github.com/kubeflow/model-registry?tab=readme-ov-file#documentation-links:~:text=FAQs-,Development,-introduction%20to%20local)!

# Contributing to Model Registry using Apple-silicon/ARM-based computers

Some limitations apply when developing on this project, specifically using Apple-silicon and Mac OSX.
The content from this guide might also be applicable in part for general ARM-based developers/users, beyond Mac OSX.

## Consider installing coreutils

While we strived to keep the repository scripts cross-compatible with GNU and BSD utilities, you might want to consider installing the `coreutils` from Homebrew, [following these instructions](https://formulae.brew.sh/formula/coreutils).
You will need to ensure the gnubin directory is added to your shell path. 

## Makefile

The make command shipped with Mac OSX (at the time of writing) is a bit old:

```
% make --version
GNU Make 3.81
This program built for i386-apple-darwin11.3.0
```

 and known to cause problems when using this project's Makefile:

```
% make build
openapi-generator-cli validate -i api/openapi/model-registry.yaml
make: openapi-generator-cli: No such file or directory
make: *** [openapi/validate] Error 1
```

i.e. failing to locate the `bin/` executables managed in the Makefile.

The solution is to use and updated version of `make`.

You can install it with Homebrew:

```
% brew install make
...
==> Pouring make--4.4.1.arm64_ventura.bottle.tar.gz
==> Caveats
GNU "make" has been installed as "gmake".
If you need to use it as "make", you can add a "gnubin" directory
to your PATH from your bashrc like:

    PATH="/opt/homebrew/opt/make/libexec/gnubin:$PATH"
...
```

and now you can substitute `gmake` every time the make command is mentioned in guides (or perform the path management per the caveat).

## Local kubernetes deployment of Model Registry

To test the Model Registry locally without mocking the k8s calls, the Model Registry backend can be deployed using Kind.

### Prerequisites

The following tools need to be installed in your local environment:

- [Podman](https://podman.io/) (Docker should also work)
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

Choose the networking setup that fits your needs, either port-forwarding or Ingress.

### Port-forwarding guide

Create a Kind cluster with the following command:

```sh
kind create cluster
```

and then follow the steps from the [Installation guide](https://www.kubeflow.org/docs/components/model-registry/installation/#standalone-installation) on the Kubeflow website, to set up the port-forwarding and deploy the Model Registry on the cluster.

### Ingress guide

Follow the [Ingress guide](docs/mr_kind_deploy_ingress.md) to set up the Ingress controller and deploy the Model Registry on the cluster.

## Docker engine

Several options of docker engines are available for Mac.
Having Docker installed is also helpful for Testcontainers.

### Podman

[Podman](https://podman.io/) also supports Rosetta (Apple specific) emulation which is handy since the Google MLMD project dependency is x86 specific.

We recommend setting up the Podman machine with root privileges,
and setting the environment variable

```sh
export TESTCONTAINERS_RYUK_PRIVILEGED=true
```

when running TestContainer-based Model Registry Python tests (for more information, see [here](https://pypi.org/project/testcontainers/#:~:text=TESTCONTAINERS_RYUK_PRIVILEGED)).

#### Fedora

Fedora requires further setup to make testcontainers work with Podman, the following steps are required:

- start the podman socket service

```sh
systemctl --user start podman.socket
```

- set the environment variable

```sh
export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/podman/podman.sock
```

If you need more information, please refer to the [Testcontainers using Podman](https://golang.testcontainers.org/system_requirements/using_podman/).

### Colima

Colima offers Rosetta (Apple specific) emulation which is handy since the Google MLMD project dependency is x86 specific.
You can install Colima (and Docker) with Homebrew.

You can create a Colima "docker context" focusing on x86 emulation with:

```
colima start --vz-rosetta --vm-type vz --arch x86_64 --cpu 4 --memory 8
```

To use with *Testcontainers for Go* you can use these commands:

```
export DOCKER_HOST="unix://${HOME}/.colima/default/docker.sock"
export TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE="/var/run/docker.sock"
```

as instructed in [this guide](https://golang.testcontainers.org/system_requirements/using_colima/#:~:text=Set%20the%20DOCKER_HOST%20environment).

This colima setups allows to:
- launch Integration tests in Go (used in Core go layer) with Testcontainers for Go
- launch DevContainer to be able to install MLMD python wheel dependency (which is x86 specific)

## DevContainer

Using a [DevContainer](https://containers.dev) is helpful to develop with the Model Registry Python client, since it needs to wrap MLMD python dependency (which is [x86 specific](https://pypi.org/project/ml-metadata/#files)).

This allows for instance with [VSCode DevContainer extension](https://code.visualstudio.com/docs/devcontainers/containers) to re-open VSCode window "inside" an x86 emulated docker container.
The experience is very similar to when on GitHub.com you press dot `.` and you get a VSCode "inside something", except it is local to your machine.
It's not super fast because x86 is emulated via Rosetta, but works "good enough" to complete most tasks without requiring remote connection to a real x86 server.

To use DevContainer as a framework directly, a command-line tool and an SDK is available as well on the upstream project: https://containers.dev.

Don't forget you will need a Docker context for x86 emulation, for instance with colima (see previous step) this can be achieved with:

```
colima start --vz-rosetta --vm-type vz --arch x86_64 --cpu 4 --memory 8
```

Define this `.devcontainer/devcontainer.json` file :
```jsonc
{
	"name": "Python 3",
	// "image": "mcr.microsoft.com/devcontainers/python:1-3.10-bullseye"
	"build": {
		"dockerfile": "Dockerfile"
	},
	"runArgs": [
		"--network=host"
	],
	"customizations": {
		"vscode": {
			"extensions": [
				// does not work well in DevContainer: "robocorp.robotframework-lsp",
				"d-biehl.robotcode"
			]
		}
	},
//...
}
```

The `network=host` allow from _inside_ the devcontainer to reach any "service" exposed on your computer (host).
This is helpful if other containers are started on your computer (eg: a PostgreSQL or DB in another container).

The `customizations.vscode.extensions` pre-loads additional extensions needed in VSCode to be executing from *inside* the DevContainer.

Define this `.devcontainer/Dockerfile` file:

```docker
FROM mcr.microsoft.com/devcontainers/python:1-3.10-bullseye

# Here I use the USER from the FROM image
ARG USERNAME=vscode
ARG GROUPNAME=vscode

# Here I use the UID/GID from _my_ computer, as I'm _not_ using Docker Desktop
ARG USER_UID=501
ARG USER_GID=20

RUN groupmod --gid $USER_GID -o $GROUPNAME \
    && usermod --uid $USER_UID --gid $USER_GID $USERNAME \
    && chown -R $USER_UID:$USER_GID /home/$USERNAME

# General setup which is "cached" for convenience
RUN pip install -U pip setuptools
RUN pip install -U poetry
RUN pip install -U "ml-metadata==1.14.0"
RUN pip install -U robotframework
RUN pip install -U robotframework-requests
RUN pip install -U PyYAML
```

The group/user is needed as on Mac anything _but_ Docker Desktop will need to set correct FS permissions. (more details here: `https://github.com/devcontainers/spec/issues/325`).

The RUN pip install "caches" the local installation (inside the DevContainer) of some wheels which are almost always needed, for convenience.

Please notice the line `RUN pip install -U "ml-metadata==1.14.0"`: that pip install would otherwise fail on Apple-silicon/ARM-based machines.

E.g. when issued on bare-metal-Apple-silicon fails with:
```
% pip install -U "ml-metadata==1.14.0"
ERROR: Could not find a version that satisfies the requirement ml-metadata==1.14.0 (from versions: 0.12.0.dev0, 0.13.0.dev0, 0.13.1.dev0)
ERROR: No matching distribution found for ml-metadata==1.14.0
```

As all the wheels are x86 specific: https://pypi.org/project/ml-metadata/1.14.0/#files

So it's not even possible to receive code assists. However, after clicking to re-open the project inside an (emulated) DevContainer:

![](docs/Screenshot%202023-11-29%20at%2014.08.12.png)

Then with the given setup MLMD is already installed inside the DevContainer:

![](docs/Screenshot%202023-11-29%20at%2014.10.14.png)

At this point Poetry is already installed as well and can be used to build and run test of the Model Registry Python client.

<!-- to be continued with explanation of this "hack": https://github.com/tarilabs/ml-metadata-remote#readme -->

# FAQ

## Error `docker.errors.NotFound: 404 Client Error for http+docker://localhost/v1.41/containers ...`

This happens on Mac OSX when running Testcontainers-based Python tests, and [Ryuk](https://github.com/testcontainers/moby-ryuk) does not have the correct access priviledges.

```
-------------------------------------------- Captured stderr setup ---------------------------------------------
Pulling image testcontainers/ryuk:0.7.0
Container started: 4af385c2d670
---------------------------------------------- Captured log setup ----------------------------------------------
WARNING  root:config.py:59 DOCKER_AUTH_CONFIG is experimental, see testcontainers/testcontainers-python#566
INFO     testcontainers.core.container:container.py:88 Pulling image testcontainers/ryuk:0.7.0
INFO     testcontainers.core.container:container.py:101 Container started: 4af385c2d670
=========================================== short test summary info ============================================
ERROR tests/store/test_wrapper.py::test_get_undefined_artifact_type_id - docker.errors.NotFound: 404 Client Error for http+docker://localhost/v1.41/containers/4af385c2d67067875f6f0...
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! stopping after 1 failures !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
=============================================== 1 error in 0.50s ===============================================
```

Solution:
- Use `export TESTCONTAINERS_RYUK_PRIVILEGED=true` (we recommend this option)
- Use `export TESTCONTAINERS_RYUK_DISABLED=true` to disable Ryuk entirely; the tests are configured to close all container resources, but we do not recommend this option

For more information:
- section [Docker engine](https://github.com/kubeflow/model-registry/blob/main/CONTRIBUTING.md#docker-engine) in this document
- https://pypi.org/project/testcontainers/#:~:text=the%20database%20version.-,Configuration,-Env%20Variable
