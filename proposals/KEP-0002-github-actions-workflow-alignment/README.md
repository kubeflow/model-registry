# KEP-0002: GitHub Actions Workflow Alignment for Container Image Signature and Attestation

## Summary

Standardize GitHub Actions workflows to use `docker/build-push-action` for multi-arch builds with metadata, `anchore/sbom-action` for SPDX SBOM generation, and cosign attestation for signed SBOMs.

## Motivation

Current workflows use inconsistent approaches for building container images, generating SBOMs, and signing artifacts. Some images use `docker/build-push-action` while others use custom scripts. SBOM generation varies in format and signing status. This creates:

- Inconsistent metadata across images
- Unsigned or incorrectly formatted SBOMs that are not cosign-compatible
- Manual tagging steps that could leverage automation

### Goals

- Use `docker/build-push-action` for all multi-arch container builds in CI/CD (GitHub Action workflows)
- Generate consistent image metadata via `docker/metadata-action`
- Produce SPDX SBOMs using `anchore/sbom-action`
- Sign images, and attest SBOMs using `cosign` (not attach)
- Sign images by digest reference, not tag (as best practice)

### Non-Goals

- Changing existing image registries or repositories
- Modifying local development workflows
- Altering image versioning schemes

## Proposal

Update all GitHub Actions workflows that build and push container images to:

1. Use `docker/build-push-action` for multi-arch builds (linux/arm64, linux/amd64)
2. Generate metadata with `docker/metadata-action` for tags and labels
3. Sign images using cosign with digest-based references (e.g., `image@sha256:...`)
4. Generate SPDX SBOMs with `anchore/sbom-action`
5. Attest (not attach) SBOMs using cosign to ensure signature with `cosign attest`

### Risks and Mitigations

**Risk**: Workflow changes may temporarily break image builds.

**Mitigation**: Test the workflow changes as possible in the PR created for this KEP before merging. Workflows already exist in PR #1790](https://github.com/kubeflow/model-registry/pull/1790#pullrequestreview-3374876556) and [PR #1588](https://github.com/kubeflow/model-registry/pull/1588) as partial implementations. Merge the PR for this KEP right after a release, so to mitigate potential unforeseen issues while integrating this strategy for all workflows (i.e.: not jeopardize the release process).

## Design Details

Workflows shall follow this pattern:

1. Set up Docker Buildx and QEMU for multi-arch
2. Generate metadata for tags using `docker/metadata-action`
3. Build and push multi-arch images using `docker/build-push-action`
4. Sign image by digest using `cosign sign image@sha256:...`
5. Generate SBOM using `anchore/sbom-action`
6. Attest SBOM using `cosign attest`

This ensures:
- Rich metadata embedded in images
- Multi-arch support (amd64, arm64, etc.)
- Signed images with digest references
- Signed SBOMs in cosign-compatible SPDX format

### Execution Details

Align all workflows to the alignment pattern:

- [.github/workflows/build-and-push-image.yml](.github/workflows/build-and-push-image.yml) (model-registry server)
  - Replace custom script with `docker/build-push-action`
  - Add `docker/metadata-action` for rich metadata
  - Add multi-arch support (already had QEMU/Buildx)
  - Add image signing by digest
  - Change from `cosign attach` to `cosign attest`

- [.github/workflows/build-and-push-async-upload.yml](.github/workflows/build-and-push-async-upload.yml) (async-upload)
  - Add QEMU for multi-arch support
  - Add `docker/metadata-action`
  - Add `platforms: linux/arm64,linux/amd64`
  - Replace built-in SBOM with `anchore/sbom-action`
  - Add `cosign` signing and attestation

- [.github/workflows/build-and-push-csi-image.yml](.github/workflows/build-and-push-csi-image.yml) (storage-initializer)
  - Replace Makefile-based build with `docker/build-push-action`
  - Add QEMU and Buildx for multi-arch support
  - Add `platforms: linux/arm64,linux/amd64`
  - Add `docker/metadata-action`
  - Add image signing by digest
  - Change from `cosign attach` to `cosign attest`

- [.github/workflows/build-and-push-ui-images.yml](.github/workflows/build-and-push-ui-images.yml) (UI)
  - Add QEMU for multi-arch support
  - Add `platforms: linux/arm64,linux/amd64`
  - Replace built-in SBOM with `anchore/sbom-action`
  - Added `cosign` signing and attestation

- [.github/workflows/build-and-push-ui-images-standalone.yml](.github/workflows/build-and-push-ui-images-standalone.yml) (UI standalone)
  - Add QEMU for multi-arch support
  - Add `platforms: linux/arm64,linux/amd64`
  - Replace built-in SBOM with `anchore/sbom-action`
  - Add `cosign` signing and attestation

So that we can ensure for all workflows:
- Multi-arch builds: `linux/arm64`, `linux/amd64` with QEMU setup (extending it in the future as needed)
- `docker/metadata-action`: Rich metadata generation
- `docker/build-push-action`: Standardized build action
- Image signing: Using `cosign sign` with _digest_ references
- SBOM generation: Using `anchore/sbom-action` for SPDX format
- SBOM attestation: Using `cosign attest` (instead of `cosign attach`)
- Digest-based references: All signing/attestation uses `@sha256:...`
- Permissions:
   - Add `id-token: write` as needed by Cosig
   - Add `actions: read`, `contents: write` as needed by `anchore/sbom-action`

What is not chaning is the structure of container image Tags, i.e. we will keep the principle that:

```yaml
tags: |
  type=raw,value=${{ env.VERSION }}                                 # e.g., v0.3.x or main-a1b2c3d
  type=raw,value=latest,enable=${{ env.BUILD_CONTEXT == 'main' }}   # latest=main branch
  type=raw,value=main,enable=${{ env.BUILD_CONTEXT == 'main' }}     # explicit main Tag
```

### Test Plan

We already had previous work done in this areas:

- [x] Existing workflows in [PR #1790](https://github.com/kubeflow/model-registry/pull/1790#pullrequestreview-3374876556) demonstrate multi-arch builds
- [x] Existing workflows in [PR #1588](https://github.com/kubeflow/model-registry/pull/1588) demonstrate SBOM generation
- [ ] Implement this solution right after a release, so to mitigate potential unforeseen problems while integrating all workflows with this strategy

#### Unit Tests

N/A - workflow changes do not affect application code.

#### E2E tests

Workflow execution in GitHub Actions validates functionality. Manual verification:
- `cosign verify` to check image signatures
- `cosign verify-attestation` to check SBOM attestations

### Graduation Criteria

N/A - this is an infrastructure alignment, not a feature.

## Implementation History

- 2025-12-18: KEP creation
- Previous work:
  - see [this comment](https://github.com/kubeflow/model-registry/pull/1790#pullrequestreview-3374876556) in introduction of multi-arch support
  - see [this PR](https://github.com/kubeflow/model-registry/pull/1588) for introduction of SBOM

## Drawbacks

Requires updating multiple workflow definitions. Increases workflow complexity slightly compared to simple docker build and push actions, but benefit with providing signature and SBOM (signed SBOM).

## Alternatives

1. Continue with current inconsistent approaches
   - Not viable: we are required by OpenSSF best practices to sign container images and provide SBOM in a standard way.

2. Use `cosign attach` instead of `cosign attest` for SBOMs
   - Not preferred: while we adopted this approach in [PR #1588](https://github.com/kubeflow/model-registry/pull/1588), the Attachment does not sign the SBOM, only associates it with the image, while Attesting does.

3. Use docker buildx directly in scripts instead of GitHub Actions
   - Not preferred: more complex automation; we can keep local development focused on DevEX while using standard workflows for the project
