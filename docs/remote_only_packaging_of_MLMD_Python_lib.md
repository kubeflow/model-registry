# Remote-only packaging of MLMD Python lib

## Context and Problem statement

Google’s ML Metadata (MLMD) is a project composed of a C++ server, and a Python client library.
The server exposes a gRPC interface, and is only distributed for x86-64 architectures.
It is embedded in the client's wheel binary, providing an additional convenience [method for running the server locally (in memory)](https://www.tensorflow.org/tfx/guide/mlmd#metadata_storage_backends_and_store_connection_configuration), 
whilst also making it [architecture specific](https://pypi.org/project/ml-metadata/1.14.0/#files).

The [Model Registry project](https://docs.google.com/document/d/1G-pjdGaS2kLELsB5kYk_D4AmH-fTfnCnJOhJ8xENjx0/edit?usp=sharing) (MR) is built on top of MLMD.
The Go implementation interfaces with the MLMD server via gRPC, typically available as a Docker container.
The [MR Python client](https://github.com/opendatahub-io/model-registry/tree/main/clients/python#readme) wraps the MLMD client.

As the MLMD client is architecture specific, so is the MR Python client, which **severely limits the targets it can run on**, as it only supports x86-64.
This **poses many challenges to contributors** using other CPU architectures, specially ARM, as that's become more prevalent in recent years.

Since the Model Registry python client does _not_ leverage the embedded MLMD server's in-memory connection, we present options for a “soft-fork” and re-distribution of the MLMD client as a pure Python library, making it platform/architecture agnostic.

The resulting artifact, a MLMD python wheel supporting only remote gRPC connections at the same time being a “pure library” hence not requiring specific CPU architectures, OSes, or CPython versions, could be as well published on PyPI under the ODH org, and become the dependency of [MR python client](/clients/python/README.md).


## Goals

* consider required changes to “soft-fork” MLMD wheel to support only remote gRPC connections
* repackaging as a pure library


## Non-Goals

* “hard”-fork of MLMD
* maintaining original MLMD tests (those bound to embedded server)


## Proposed solution

Refer to the conclusions section for the motivations behind selecting:
1. soft-fork upstream repo, modify pip+bazel build, to produce the distributable Python client ("Alternative B", below)
2. create a `ml-metadata-remote` or similarly named package on PyPI based on the distributable wheel and files from the step1 above ("Packaging Option1", below)

For documentation purposes, the exploration of the different alternatives is reported below.


## Alternative A: repackage the resulting wheel

This solution explores the steps required to simply repackage the existing and distributed (by Google) MLMD python wheel, removing all embedded binaries, and repackaging the wheel as a pure library.

This has been experimented with success on this repository: ([link](https://github.com/tarilabs/ml-metadata-remote))

The steps required are recorded here: ([link](https://github.com/tarilabs/ml-metadata-remote/commits/v1.14.0))

and mainly consists of:

1. Download one platform-specific MLMD v1.14.0 wheel and extract its content ([reference](https://github.com/tarilabs/ml-metadata-remote/commit/39dd0c7dcd063e0440a6354017445dada8423f0c#diff-b335630551682c19a781afebcf4d07bf978fb1f8ac04c6bf87428ed5106870f5))
2. Remove embedded code, apply required code changes ([reference](https://github.com/tarilabs/ml-metadata-remote/commit/bcb1f0ffd37600e056342aff39e154bb35422668#diff-f363c85a1cf3536a48a7b721b02a6999b80a08b9c305d185327e87e2769b6f21))
3. Recompute dist-info checksums before repackaging ([reference](https://github.com/tarilabs/ml-metadata-remote/commit/fda125fb742ab8ecf4a7153705717d8b50f59326#diff-53bdc596caf062825dbb42b65e5b2305db70d2e533c03bc677b13cc8c7cfd236))
4. Repackage the directories as a new pure library wheel ([reference](https://github.com/tarilabs/ml-metadata-remote/commit/5d199f808eea0cb7ba78a0702be8de3306477df8))

The resulting artifact has been [tested](https://github.com/tarilabs/ml-metadata-remote#readme:~:text=Testing%20with%20launching%20a%20local%20server) locally with a gRPC connection to MLMD server made available via Docker. The resulting artifact is directly available for local download: ([link](https://github.com/tarilabs/ml-metadata-remote/releases/tag/1.14.0))


## Alternative B: build by soft-fork upstream repo

This solution explores how to use the upstream MLMD repo by Google and by making necessary code changes, so to directly produce with the pip+bazel build the wheel as a pure library.

This has been experimented with success on this fork: ([link](https://github.com/tarilabs/ml-metadata/commits/remote-r1.14.0))

The steps required mainly consists of:

1. Make changes to the bazel BUILD file ([reference](https://github.com/tarilabs/ml-metadata/commit/079aeb3a9da69eb960e428a7866e279d0bfb533b#diff-c8858dec4f58c1d8a280af8c117ff8480f7ed4ae863b96e1ba20b52f83222aab))
2. Make changes to sh build script ([reference](https://github.com/tarilabs/ml-metadata/commit/079aeb3a9da69eb960e428a7866e279d0bfb533b#diff-125a2f247ce39f711e1c8a77f430bd5b1b865cd10b5c5fef0d9140d276c617f2))
3. Make changes to setup.py build file ([reference](https://github.com/tarilabs/ml-metadata/commit/079aeb3a9da69eb960e428a7866e279d0bfb533b#diff-60f61ab7a8d1910d86d9fda2261620314edcae5894d5aaa236b821c7256badd7))
4. Apply required code changes analogously to “Alternative A” (see other changes in [this commit](https://github.com/tarilabs/ml-metadata/commit/079aeb3a9da69eb960e428a7866e279d0bfb533b))

The resulting artifact has been [tested](https://github.com/tarilabs/ml-metadata/commit/794ec39d97e3ac70db2ca18fcf5807c44f339f0b) locally with a gRPC connection to MLMD server made available via Docker, similar to instructions provided in “Alternative A”.


## Packaging

In this section we consider packaging and delivery options for the resulting artifact from the alternative selected above.


### Packaging Option1: separate repo on ODH

This delivery option considers having a separate repo on ODH, called “ml-metadata-remote” (or the likes). Repeat the exercise from the alternative selected above on this repo. Then deliver this as a package on PyPI.

Pros:

* Well isolated dependency
    * also, if one day upstream Google MLMD resolves to be platform/arch agnostic, is just a matter of changing again the consumed dependency from MR python client
* Google code (copyright header) is isolated from Model Registry code
* The resulting artifact could also be re-used by other communities/users

Cons:

* Additional artifact to publish on PyPI


### Packaging Option2: mix resulting artifact inside Model Registry repo

This delivery option considers placing the resulting artifact by executing the exercise from the alternative selected above and placing it directly inside the Model Registry repo, with the python client source [location](https://github.com/opendatahub-io/model-registry/tree/main/clients/python). (for analogy, this is similar to “shading”/”uberjar” in Java world for those familiar with the concept)

Pros:

* Only one artifact to publish on PyPI

Cons:

* Google code (copyright header) is mixed with Model Registry code
    * at this stage is not clear if any implications with uplifting the MR project in KF community
* The resulting artifact cannot be re-used by other communities/users
* If one day upstream Google MLMD resolves to be platform/arch agnostic, changing back the MR python client to use the original ml-metadata could require extra work and effort


## Conclusion

Based on analysis of the alternatives provided, we decided to further pursue:
- the repackaging by **Alternative B** because makes it actually easier to demonstrate the steps and modifications required using as baseline the upstream repo.
- the distribution by **Packaging Option1** because it will make it easier to "revert" to the upstream `ml-metadata` dependency if upstream will publish for all architectures, OSes, etc. and as the pros outweight considered cons.

MR python client [tests](https://github.com/opendatahub-io/model-registry/blob/259b39320953bf05942dcec1fb5ec74f7eb5d4a7/clients/python/tests/conftest.py#L19) should be rewritten using Testcontainers, and not leveraging the embedded server (already done with [this PR](https://github.com/opendatahub-io/model-registry/pull/225)).

The group concur this is a sensible approach ([recorded here](https://redhat-internal.slack.com/archives/C05LGBNUK9C/p1700763823505259?thread_ts=1700427888.670999&cid=C05LGBNUK9C)).

This change would also better align the test approach used for the MR python client, by aligning with the same strategy of the MR core Go layer test framework, which already makes use of Testcontainers for Go ([reference](https://github.com/opendatahub-io/model-registry/blob/259b39320953bf05942dcec1fb5ec74f7eb5d4a7/internal/testutils/test_container_utils.go#L59)).

This would allow to update the constraint on the version for the `attrs` dependency as part of this activity.
