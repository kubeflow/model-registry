package api

import (
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("TestGetModelVersionHandler", func() {
	Context("testing Model Version Handler", Ordered, func() {

		It("should retrieve a model version", func() {
			By("fetching a model version")
			data := mocks.GetModelVersionMocks()[0]
			expected := ModelVersionEnvelope{Data: &data}
			actual, rs, err := setupApiTest[ModelVersionEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/model_versions/1", nil, k8sClient)
			Expect(err).NotTo(HaveOccurred())
			By("should match the expected model version")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data.Name).To(Equal(expected.Data.Name))
		})

		It("should create a model version", func() {
			By("creating a model version")
			data := mocks.GetModelVersionMocks()[0]
			expected := ModelVersionEnvelope{Data: &data}
			body := ModelVersionEnvelope{Data: openapi.NewModelVersion("Model One", "1")}
			actual, rs, err := setupApiTest[ModelVersionEnvelope](http.MethodPost, "/api/v1/model_registry/model-registry/model_versions", body, k8sClient)
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected model version created")
			Expect(rs.StatusCode).To(Equal(http.StatusCreated))
			Expect(actual.Data.Name).To(Equal(expected.Data.Name))
			Expect(rs.Header.Get("Location")).To(Equal("/api/v1/model_registry/model-registry/model_versions/1"))
		})

		It("should updated a model version", func() {
			By("updating a model version")
			data := mocks.GetModelVersionMocks()[0]
			expected := ModelVersionEnvelope{Data: &data}

			reqData := openapi.ModelVersionUpdate{
				Description: openapi.PtrString("New description"),
			}
			body := ModelVersionUpdateEnvelope{Data: &reqData}

			actual, rs, err := setupApiTest[ModelVersionEnvelope](http.MethodPatch, "/api/v1/model_registry/model-registry/model_versions/1", body, k8sClient)
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected model version updated")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data.Name).To(Equal(expected.Data.Name))
		})

		It("get all model artifacts by a model version", func() {
			By("getting a model artifacts by model version")
			data := mocks.GetModelArtifactListMock()
			expected := ModelArtifactListEnvelope{Data: &data}
			actual, rs, err := setupApiTest[ModelArtifactListEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/model_versions/1/artifacts", nil, k8sClient)
			Expect(err).NotTo(HaveOccurred())

			By("should get all expected model version artifacts")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data.Size).To(Equal(expected.Data.Size))
			Expect(actual.Data.PageSize).To(Equal(expected.Data.PageSize))
			Expect(actual.Data.NextPageToken).To(Equal(expected.Data.NextPageToken))
			Expect(len(actual.Data.Items)).To(Equal(len(expected.Data.Items)))
		})

		It("create Model Artifact By Model Version", func() {
			By("creating a model version")
			data := mocks.GetModelArtifactMocks()[0]
			expected := ModelArtifactEnvelope{Data: &data}

			artifact := openapi.ModelArtifact{
				Name:         openapi.PtrString("Artifact One"),
				ArtifactType: "ARTIFACT_TYPE_ONE",
			}
			body := ModelArtifactEnvelope{Data: &artifact}
			actual, rs, err := setupApiTest[ModelArtifactEnvelope](http.MethodPost, "/api/v1/model_registry/model-registry/model_versions/1/artifacts", body, k8sClient)
			Expect(err).NotTo(HaveOccurred())

			By("should get all expected model artifacts")
			Expect(rs.StatusCode).To(Equal(http.StatusCreated))
			Expect(actual.Data.GetArtifactType()).To(Equal(expected.Data.GetArtifactType()))
			Expect(rs.Header.Get("Location")).To(Equal("/api/v1/model_registry/model-registry/model_artifacts/1"))

		})
	})
})
