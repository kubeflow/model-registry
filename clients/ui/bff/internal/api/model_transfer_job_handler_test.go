package api

import (
	"net/http"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestModelTransferJob", func() {
	var requestIdentity kubernetes.RequestIdentity

	BeforeEach(func() {
		requestIdentity = kubernetes.RequestIdentity{
			UserID: "user@example.com",
		}
	})

	Context("fetching model transfer jobs", func() {
		It("GET ALL returns 200", func() {
			_, rs, err := setupApiTest[ModelTransferJobListEnvelope](
				http.MethodGet,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				nil,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
		})

		It("GET returns 400 when namespace is missing", func() {
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodGet,
				"/api/v1/model_registry/model-registry/model_transfer_jobs",
				nil,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"", // empty namespace
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})

	Context("creating model transfer job", func() {
		It("POST returns 201 on success", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job-create",
					Source: models.ModelTransferJobSource{
						Type:               models.ModelTransferJobSourceTypeS3,
						Bucket:             "test-bucket",
						Key:                "models/test",
						AwsAccessKeyId:     "test-key",
						AwsSecretAccessKey: "test-secret",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Registry: "quay.io",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[ModelTransferJobEnvelope](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusCreated))
		})

		It("POST returns 400 when data is null", func() {
			payload := ModelTransferJobEnvelope{
				Data: nil,
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 404 when model registry not found in namespace", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job-404",
					Source: models.ModelTransferJobSource{
						Type:               models.ModelTransferJobSourceTypeS3,
						Bucket:             "test-bucket",
						Key:                "models/test",
						AwsAccessKeyId:     "test-key",
						AwsSecretAccessKey: "test-secret",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Registry: "quay.io",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=no-namespace",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"no-namespace",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("POST returns 400 for missing job name", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					// Name is missing
					Source: models.ModelTransferJobSource{
						Type:               models.ModelTransferJobSourceTypeS3,
						Bucket:             "test-bucket",
						Key:                "models/test",
						AwsAccessKeyId:     "test-key",
						AwsSecretAccessKey: "test-secret",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for invalid job name (too long)", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "this-job-name-is-way-too-long-and-exceeds-the-fifty-character-limit-for-dns",
					Source: models.ModelTransferJobSource{
						Type:               models.ModelTransferJobSourceTypeS3,
						Bucket:             "test-bucket",
						Key:                "models/test",
						AwsAccessKeyId:     "test-key",
						AwsSecretAccessKey: "test-secret",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for invalid job name (invalid characters)", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "INVALID_NAME!!!",
					Source: models.ModelTransferJobSource{
						Type:               models.ModelTransferJobSourceTypeS3,
						Bucket:             "test-bucket",
						Key:                "models/test",
						AwsAccessKeyId:     "test-key",
						AwsSecretAccessKey: "test-secret",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for missing source type", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						// Type is missing
						Bucket:             "test-bucket",
						Key:                "models/test",
						AwsAccessKeyId:     "test-key",
						AwsSecretAccessKey: "test-secret",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for S3 source missing bucket", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type: models.ModelTransferJobSourceTypeS3,
						// Bucket is missing
						Key:                "models/test",
						AwsAccessKeyId:     "test-key",
						AwsSecretAccessKey: "test-secret",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for S3 source missing AWS credentials", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type:   models.ModelTransferJobSourceTypeS3,
						Bucket: "test-bucket",
						Key:    "models/test",
						// AWS credentials missing
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for URI source missing URI", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type: models.ModelTransferJobSourceTypeURI,
						// URI is missing
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for missing destination credentials", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type:               models.ModelTransferJobSourceTypeS3,
						Bucket:             "test-bucket",
						Key:                "models/test",
						AwsAccessKeyId:     "test-key",
						AwsSecretAccessKey: "test-secret",
					},
					Destination: models.ModelTransferJobDestination{
						Type: models.ModelTransferJobDestinationTypeOCI,
						URI:  "quay.io/test/model:v1",
						// Username and Password missing
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for missing upload intent", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type:               models.ModelTransferJobSourceTypeS3,
						Bucket:             "test-bucket",
						Key:                "models/test",
						AwsAccessKeyId:     "test-key",
						AwsSecretAccessKey: "test-secret",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					// UploadIntent is missing
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for create_model intent missing model name", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type:               models.ModelTransferJobSourceTypeS3,
						Bucket:             "test-bucket",
						Key:                "models/test",
						AwsAccessKeyId:     "test-key",
						AwsSecretAccessKey: "test-secret",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent: models.ModelTransferJobUploadIntentCreateModel,
					// RegisteredModelName is missing
					ModelVersionName: "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for create_version intent missing model ID", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type:               models.ModelTransferJobSourceTypeS3,
						Bucket:             "test-bucket",
						Key:                "models/test",
						AwsAccessKeyId:     "test-key",
						AwsSecretAccessKey: "test-secret",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent: models.ModelTransferJobUploadIntentCreateVersion,
					// RegisteredModelId is missing
					ModelVersionName: "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for update_artifact intent missing artifact ID", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type:               models.ModelTransferJobSourceTypeS3,
						Bucket:             "test-bucket",
						Key:                "models/test",
						AwsAccessKeyId:     "test-key",
						AwsSecretAccessKey: "test-secret",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent: models.ModelTransferJobUploadIntentUpdateArtifact,
					// ModelArtifactId is missing
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 201 for URI source type", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "uri-source-job",
					Source: models.ModelTransferJobSource{
						Type: models.ModelTransferJobSourceTypeURI,
						URI:  "https://huggingface.co/test/model.safetensors",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Registry: "quay.io",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "URI Source Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[ModelTransferJobEnvelope](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusCreated))
		})

		It("POST returns 201 for create_version intent with model ID", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "version-job",
					Source: models.ModelTransferJobSource{
						Type: models.ModelTransferJobSourceTypeURI,
						URI:  "https://test.com/model.bin",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v2",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:      models.ModelTransferJobUploadIntentCreateVersion,
					RegisteredModelId: "1",
					ModelVersionName:  "v2.0.0",
				},
			}
			_, rs, err := setupApiTest[ModelTransferJobEnvelope](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusCreated))
		})

		It("POST returns 201 for update_artifact intent with artifact ID", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "artifact-job",
					Source: models.ModelTransferJobSource{
						Type: models.ModelTransferJobSourceTypeURI,
						URI:  "https://test.com/model.bin",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v3",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:    models.ModelTransferJobUploadIntentUpdateArtifact,
					ModelArtifactId: "1",
				},
			}
			_, rs, err := setupApiTest[ModelTransferJobEnvelope](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusCreated))
		})

		It("POST returns 400 for invalid source type", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type: "invalid_type",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for invalid destination type", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type: models.ModelTransferJobSourceTypeURI,
						URI:  "https://test.com/model",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     "invalid_type",
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for invalid upload intent", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type: models.ModelTransferJobSourceTypeURI,
						URI:  "https://test.com/model",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        "invalid_intent",
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for S3 source missing key", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type:   models.ModelTransferJobSourceTypeS3,
						Bucket: "test-bucket",
						// Key is missing
						AwsAccessKeyId:     "test-key",
						AwsSecretAccessKey: "test-secret",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for OCI destination missing URI", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type: models.ModelTransferJobSourceTypeURI,
						URI:  "https://test.com/model",
					},
					Destination: models.ModelTransferJobDestination{
						Type: models.ModelTransferJobDestinationTypeOCI,
						// URI is missing
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for OCI destination missing password", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type: models.ModelTransferJobSourceTypeURI,
						URI:  "https://test.com/model",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						// Password is missing
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for create_model intent missing version name", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type: models.ModelTransferJobSourceTypeURI,
						URI:  "https://test.com/model",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					// ModelVersionName is missing
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST returns 400 for create_version intent missing version name", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "test-job",
					Source: models.ModelTransferJobSource{
						Type: models.ModelTransferJobSourceTypeURI,
						URI:  "https://test.com/model",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:      models.ModelTransferJobUploadIntentCreateVersion,
					RegisteredModelId: "1",
					// ModelVersionName is missing
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("POST extracts registry from destination URI when not provided", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "registry-extract-job",
					Source: models.ModelTransferJobSource{
						Type: models.ModelTransferJobSourceTypeURI,
						URI:  "https://test.com/model",
					},
					Destination: models.ModelTransferJobDestination{
						Type: models.ModelTransferJobDestinationTypeOCI,
						URI:  "docker.io/myrepo/model:v1",
						// Registry is not provided - should be extracted from URI
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[ModelTransferJobEnvelope](
				http.MethodPost,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusCreated))
		})

	})

	Context("updating model transfer job", func() {
		It("PATCH returns 200 on success", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "new-job-name",
				},
			}
			_, rs, err := setupApiTest[ModelTransferJobEnvelope](
				http.MethodPatch,
				"/api/v1/model_registry/model-registry/model_transfer_jobs/transfer-job-001?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
		})

		It("PATCH returns 404 for non-existent job", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "new-job-name",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPatch,
				"/api/v1/model_registry/model-registry/model_transfer_jobs/does-not-exist?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("PATCH returns 400 for missing new job name", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					// Name is missing
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPatch,
				"/api/v1/model_registry/model-registry/model_transfer_jobs/transfer-job-001?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("PATCH returns 400 for invalid new job name", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "INVALID_NAME!!!",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPatch,
				"/api/v1/model_registry/model-registry/model_transfer_jobs/transfer-job-001?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("PATCH with deleteOldJob=true returns 200 on success", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "new-job-after-delete",
					Source: models.ModelTransferJobSource{
						Type: models.ModelTransferJobSourceTypeURI,
						URI:  "https://test.com/model.bin",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[ModelTransferJobEnvelope](
				http.MethodPatch,
				"/api/v1/model_registry/model-registry/model_transfer_jobs/transfer-job-001?namespace=kubeflow&deleteOldJob=true",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
		})

		It("PATCH returns 400 when namespace is missing", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "new-job",
					Source: models.ModelTransferJobSource{
						Type: models.ModelTransferJobSourceTypeURI,
						URI:  "https://test.com/model.bin",
					},
					Destination: models.ModelTransferJobDestination{
						Type:     models.ModelTransferJobDestinationTypeOCI,
						URI:      "quay.io/test/model:v1",
						Username: "user",
						Password: "pass",
					},
					UploadIntent:        models.ModelTransferJobUploadIntentCreateModel,
					RegisteredModelName: "Test Model",
					ModelVersionName:    "v1.0.0",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPatch,
				"/api/v1/model_registry/model-registry/model_transfer_jobs/transfer-job-001",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("PATCH returns 400 when data is null", func() {
			payload := ModelTransferJobEnvelope{
				Data: nil,
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPatch,
				"/api/v1/model_registry/model-registry/model_transfer_jobs/transfer-job-001?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("PATCH returns 400 when new job name equals old job name", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "transfer-job-001",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPatch,
				"/api/v1/model_registry/model-registry/model_transfer_jobs/transfer-job-001?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("PATCH returns 404 when job exists but belongs to different registry", func() {
			payload := ModelTransferJobEnvelope{
				Data: &models.ModelTransferJob{
					Name: "new-job-name",
				},
			}
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodPatch,
				"/api/v1/model_registry/other-registry/model_transfer_jobs/transfer-job-001?namespace=kubeflow",
				payload,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	Context("deleting model transfer job", func() {
		It("DELETE returns 200 on success", func() {
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodDelete,
				"/api/v1/model_registry/model-registry/model_transfer_jobs/transfer-job-001?namespace=kubeflow",
				nil,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
		})

		It("DELETE returns 404 for non-existent job", func() {
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodDelete,
				"/api/v1/model_registry/model-registry/model_transfer_jobs/does-not-exist?namespace=kubeflow",
				nil,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("DELETE returns 404 when job exists but belongs to different registry", func() {
			_, rs, err := setupApiTest[Envelope[any, any]](
				http.MethodDelete,
				"/api/v1/model_registry/other-registry/model_transfer_jobs/transfer-job-001?namespace=kubeflow",
				nil,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})

var _ = Describe("TestModelTransferJob registry filtering", func() {
	var requestIdentity kubernetes.RequestIdentity

	BeforeEach(func() {
		requestIdentity = kubernetes.RequestIdentity{
			UserID: "user@example.com",
		}
	})

	Context("GET list filtered by registry", func() {
		It("GET list for other registry returns 200 with empty items", func() {
			envelope, rs, err := setupApiTest[ModelTransferJobListEnvelope](
				http.MethodGet,
				"/api/v1/model_registry/other-registry/model_transfer_jobs?namespace=kubeflow",
				nil,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(envelope.Data).NotTo(BeNil())
			Expect(envelope.Data.Items).To(BeEmpty())
		})
	})
})
