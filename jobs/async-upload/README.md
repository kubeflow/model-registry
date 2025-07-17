# Kubeflow Model Registry – Async Upload Job

> Lightweight, non‑root Python Job containerised for **transferring model artefacts** between storage back‑ends (S3, OCI, PVC, …) and registering them in **Kubeflow Model Registry**. Born out of the discussion in [kubeflow/model-registry #1108](https://github.com/kubeflow/model-registry/issues/1108).

---

## Quick start

```bash
# 1 – Build & tag ( reproducible build args are optional )
docker build \
  --build-arg VCS_REF=$(git rev-parse --short HEAD) \
  --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  -t quay.io/<org>/async-upload-job:$(git rev-parse --short HEAD) .

# 2 – Push to your registry
docker push quay.io/<org>/model-registry-job-async-upload:<tag>

# TODO: Run locally...
```

## Job Configuration

The async job is able to take a number of configuration parameters and environment variables which can be consumed to perform the job of synchronizing a model from a source to a destination.

When using environment variables to configure the job, you will need to provide them in the Kubernetes Job manifest in the `spec.template.spec.containers[*].env` section. From the job's perspective, these will become standard env vars that it can read. See the [samples directory](./samples/) for typical usage.

When using a parameter-based approach, the configuration variables will need to be passed in as `args` to the `command`.

When providing parameters in a mixed-fashion to the job, the job will prioritize certain sources of those parameters over others. The order of priority is below:

1. Command-line arguments (`args: []` in the manifest)
2. Environment variables (`env: { ... }` in the manifest)
3. Credentials files (read from the parameter \*\_CREDENTIALS_PATH env/arg)
4. Default values

Below is a table of the configuration that can be passed into the job

See asterisks below table for details

| Environment Variable                         | Arg                                 | Default Value     | Required | Description                                                                             |
| -------------------------------------------- | ----------------------------------- | ----------------- | -------- | --------------------------------------------------------------------------------------- |
| MODEL_SYNC_SOURCE_TYPE                       | --source-type                       | s3                | ✅       |                                                                                         |
| MODEL_SYNC_SOURCE_S3_CREDENTIALS_PATH        | --source-s3-credentials-path        |                   |          |                                                                                         |
| MODEL_SYNC_SOURCE_OCI_CREDENTIALS_PATH       | --source-oci-credentials-path       |                   |          |                                                                                         |
| MODEL_SYNC_SOURCE_AWS_BUCKET                 | --source-aws-bucket                 |                   | ✅\*     | When --source-type is "s3"                                                              |
| MODEL_SYNC_SOURCE_AWS_KEY                    | --source-aws-key                    |                   | ✅\*     | "                                                                                       |
| MODEL_SYNC_SOURCE_AWS_REGION                 | --source-aws-region                 |                   |          | "                                                                                       |
| MODEL_SYNC_SOURCE_AWS_ACCESS_KEY_ID          | --source-aws-access-key-id          |                   | ✅\*     | "                                                                                       |
| MODEL_SYNC_SOURCE_AWS_SECRET_ACCESS_KEY      | --source-aws-secret-access-key      |                   | ✅\*     | "                                                                                       |
| MODEL_SYNC_SOURCE_AWS_ENDPOINT               | --source-aws-endpoint               |                   |          | "                                                                                       |
| MODEL_SYNC_SOURCE_OCI_URI                    | --source-oci-uri                    |                   | ✅\+     | When --source-type is "oci". The tag to use when pulling the image                      |
| MODEL_SYNC_SOURCE_OCI_REGISTRY               | --source-oci-registry               |                   | ✅\+     | When --source-type is "oci". Indicates which registry the creds belong to               |
| MODEL_SYNC_SOURCE_OCI_USERNAME               | --source-oci-username               |                   | ✅\+     | "                                                                                       |
| MODEL_SYNC_SOURCE_OCI_PASSWORD               | --source-oci-password               |                   | ✅\+     | "                                                                                       |
| MODEL_SYNC_DESTINATION_TYPE                  | --destination-type                  | oci               | ✅       |                                                                                         |
| MODEL_SYNC_DESTINATION_S3_CREDENTIALS_PATH   | --destination-s3-credentials-path   |                   |          |                                                                                         |
| MODEL_SYNC_DESTINATION_OCI_CREDENTIALS_PATH  | --destination-oci-credentials-path  |                   |          |                                                                                         |
| MODEL_SYNC_DESTINATION_AWS_BUCKET            | --destination-aws-bucket            |                   | ✅\*     | When --destination-type is "s3"                                                         |
| MODEL_SYNC_DESTINATION_AWS_KEY               | --destination-aws-key               |                   | ✅\*     | "                                                                                       |
| MODEL_SYNC_DESTINATION_AWS_REGION            | --destination-aws-region            |                   |          | "                                                                                       |
| MODEL_SYNC_DESTINATION_AWS_ACCESS_KEY_ID     | --destination-aws-access-key-id     |                   | ✅\*     | "                                                                                       |
| MODEL_SYNC_DESTINATION_AWS_SECRET_ACCESS_KEY | --destination-aws-secret-access-key |                   | ✅\*     | "                                                                                       |
| MODEL_SYNC_DESTINATION_AWS_ENDPOINT          | --destination-aws-endpoint          |                   |          | "                                                                                       |
| MODEL_SYNC_DESTINATION_OCI_URI               | --destination-oci-uri               |                   | ✅\+     | When --destination-type is "oci". The tag to use when pushing the image                 |
| MODEL_SYNC_DESTINATION_OCI_REGISTRY          | --destination-oci-registry          |                   | ✅\+     | When --destination-type is "oci". Indicates which registry the creds belong to          |
| MODEL_SYNC_DESTINATION_OCI_USERNAME          | --destination-oci-username          |                   | ✅\+     | "                                                                                       |
| MODEL_SYNC_DESTINATION_OCI_PASSWORD          | --destination-oci-password          |                   | ✅\+     | "                                                                                       |
| MODEL_SYNC_DESTINATION_OCI_BASE_IMAGE        | --destination-oci-base-image        | busybox:latest    |          | When --destination-type is "oci". The image to use when pushing to an OCI registry      |
| MODEL_SYNC_DESTINATION_OCI_ENABLE_TLS_VERIFY | --destination-oci-enable-tls-verify | true              |          | When --destination-type is "oci". Specifies whether to use TLS when pushing to registry |
| MODEL_SYNC_MODEL_ID                          | --model-id                          |                   | ✅       | The `RegisteredModel.id`                                                                |
| MODEL_SYNC_VERSION_ID                        | --model-version-id                  |                   | ✅       | The `ModelVersion.id`                                                                   |
| MODEL_SYNC_ARTIFACT_ID                       | --model-artifact-id                 |                   | ✅       | The `ModelArtifact.id`                                                                  |
| MODEL_SYNC_STORAGE_PATH                      | --storage-path                      | `/tmp/model-sync` | ✅       | Temporary storage, must be large enough to hold the entire model                        |
| MODEL_SYNC_REGISTRY_SERVER_ADDRESS           | --registry-server-address           |                   | ✅       | Server address for the model-registry client to connect to                              |
| MODEL_SYNC_REGISTRY_PORT                     | --registry-port                     | 443               |          |                                                                                         |
| MODEL_SYNC_REGISTRY_IS_SECURE                | --registry-is-secure                | True              |          |                                                                                         |
| MODEL_SYNC_REGISTRY_AUTHOR                   | --registry-author                   |                   |          |                                                                                         |
| MODEL_SYNC_REGISTRY_USER_TOKEN               | --registry-user-token               |                   |          |                                                                                         |
| MODEL_SYNC_REGISTRY_USER_TOKEN_ENVVAR        | --registry-user-token-envvar        |                   |          |                                                                                         |
| MODEL_SYNC_REGISTRY_CUSTOM_CA                | --registry-custom-ca                |                   |          |                                                                                         |
| MODEL_SYNC_REGISTRY_CUSTOM_CA_ENVVAR         | --registry-custom-ca-envvar         |                   |          |                                                                                         |
| MODEL_SYNC_REGISTRY_LOG_LEVEL                | --registry-log-level                |                   |          |                                                                                         |

✅\*: Must be present in some form when the source/destination is `s3`. This might be from the parameter in the table, or from the credentials file(s) that was specified/provided.

✅\+: Must be present in some from when the source/destination is `oci`. This might be from the parameter in the table, or from the credentials file(s) that was specified/provided.

## References

- Issue thread : [https://github.com/kubeflow/model-registry/issues/1108](https://github.com/kubeflow/model-registry/issues/1108)
- OCI Image Spec : [https://github.com/opencontainers/image-spec](https://github.com/opencontainers/image-spec)
- Kubernetes Pod Security : [https://kubernetes.io/docs/concepts/security/pod-security-standards/](https://kubernetes.io/docs/concepts/security/pod-security-standards/)
