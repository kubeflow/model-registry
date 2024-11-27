package api

import (
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("TestGetRegisteredModelHandler", func() {
	Context("testing registered models by id", Ordered, func() {

		It("should retrieve a registered model", func() {
			By("fetching all model registries")
			data := mocks.GetRegisteredModelMocks()[0]
			expected := RegisteredModelEnvelope{Data: &data}
			actual, rs, err := setupApiTest[RegisteredModelEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/registered_models/1", nil, k8sClient, mocks.KubeflowUserIDHeaderValue)
			Expect(err).NotTo(HaveOccurred())
			By("should match the expected model registry")
			//TODO assert the full structure, I couldn't get unmarshalling to work for the full customProperties values
			// this issue is in the test only
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data.Name).To(Equal(expected.Data.Name))
		})

		It("should retrieve all registered models", func() {
			By("fetching all registered models")
			data := mocks.GetRegisteredModelListMock()
			expected := RegisteredModelListEnvelope{Data: &data}
			actual, rs, err := setupApiTest[RegisteredModelListEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/registered_models", nil, k8sClient, mocks.KubeflowUserIDHeaderValue)
			Expect(err).NotTo(HaveOccurred())
			By("should match the expected model registry")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data.Size).To(Equal(expected.Data.Size))
			Expect(actual.Data.PageSize).To(Equal(expected.Data.PageSize))
			Expect(actual.Data.NextPageToken).To(Equal(expected.Data.NextPageToken))
			Expect(len(actual.Data.Items)).To(Equal(len(expected.Data.Items)))
		})

		It("creating registered models", func() {
			By("post to registered models")
			data := mocks.GetRegisteredModelMocks()[0]
			expected := RegisteredModelEnvelope{Data: &data}
			body := RegisteredModelEnvelope{Data: openapi.NewRegisteredModel("Model One")}
			actual, rs, err := setupApiTest[RegisteredModelEnvelope](http.MethodPost, "/api/v1/model_registry/model-registry/registered_models", body, k8sClient, mocks.KubeflowUserIDHeaderValue)
			Expect(err).NotTo(HaveOccurred())

			By("should do a successful post")
			Expect(rs.StatusCode).To(Equal(http.StatusCreated))
			Expect(actual.Data.Name).To(Equal(expected.Data.Name))
			Expect(rs.Header.Get("location")).To(Equal("/api/v1/model_registry/model-registry/registered_models/1"))
		})

		It("updating registered models", func() {
			By("path to registered models")
			data := mocks.GetRegisteredModelMocks()[0]
			expected := RegisteredModelEnvelope{Data: &data}
			reqData := openapi.RegisteredModelUpdate{
				Description: openapi.PtrString("This is a new description"),
			}
			body := RegisteredModelUpdateEnvelope{Data: &reqData}
			actual, rs, err := setupApiTest[RegisteredModelEnvelope](http.MethodPatch, "/api/v1/model_registry/model-registry/registered_models/1", body, k8sClient, mocks.KubeflowUserIDHeaderValue)
			Expect(err).NotTo(HaveOccurred())

			By("should do a successful patch")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data.Description).To(Equal(expected.Data.Description))
		})

		It("get all model versions for registered model", func() {
			By("get to registered models versions")
			data := mocks.GetModelVersionListMock()
			expected := ModelVersionListEnvelope{Data: &data}

			actual, rs, err := setupApiTest[ModelVersionListEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/registered_models/1/versions", nil, k8sClient, mocks.KubeflowUserIDHeaderValue)
			Expect(err).NotTo(HaveOccurred())

			By("should get all items")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data.Size).To(Equal(expected.Data.Size))
			Expect(actual.Data.PageSize).To(Equal(expected.Data.PageSize))
			Expect(actual.Data.NextPageToken).To(Equal(expected.Data.NextPageToken))
			Expect(len(actual.Data.Items)).To(Equal(len(expected.Data.Items)))
		})

		It("create model version for registered model", func() {
			By("doing a post to registered model versions")
			data := mocks.GetModelVersionMocks()[0]
			expected := ModelVersionEnvelope{Data: &data}

			body := ModelVersionEnvelope{Data: openapi.NewModelVersion("Version Fifty", "")}
			actual, rs, err := setupApiTest[ModelVersionEnvelope](http.MethodPost, "/api/v1/model_registry/model-registry/registered_models/1/versions", body, k8sClient, mocks.KubeflowUserIDHeaderValue)
			Expect(err).NotTo(HaveOccurred())

			By("should successfully create it")
			Expect(rs.StatusCode).To(Equal(http.StatusCreated))
			Expect(actual.Data.Name).To(Equal(expected.Data.Name))
			Expect(rs.Header.Get("Location")).To(Equal("/api/v1/model_registry/model-registry/model_versions/1"))

		})

		It("should return 403 when not using the correct KubeflowUserIDHeaderValue", func() {
			By("making a request with an incorrect username")
			wrongUserIDHeader := "bella@dora.com" // Incorrect username header value

			// Test: GET /registered_models/1
			_, rs, err := setupApiTest[RegisteredModelEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/registered_models/1", nil, k8sClient, wrongUserIDHeader)
			Expect(err).NotTo(HaveOccurred())
			By("should return a 403 Forbidden response for GET registered model by ID")
			Expect(rs.StatusCode).To(Equal(http.StatusForbidden))

			// Test: GET /registered_models
			_, rs, err = setupApiTest[RegisteredModelListEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/registered_models", nil, k8sClient, wrongUserIDHeader)
			Expect(err).NotTo(HaveOccurred())
			By("should return a 403 Forbidden response for GET all registered models")
			Expect(rs.StatusCode).To(Equal(http.StatusForbidden))

			// Test: POST /registered_models
			body := RegisteredModelEnvelope{Data: openapi.NewRegisteredModel("Model One")}
			_, rs, err = setupApiTest[RegisteredModelEnvelope](http.MethodPost, "/api/v1/model_registry/model-registry/registered_models", body, k8sClient, wrongUserIDHeader)
			Expect(err).NotTo(HaveOccurred())
			By("should return a 403 Forbidden response for POST create registered model")
			Expect(rs.StatusCode).To(Equal(http.StatusForbidden))

			// Test: PATCH /registered_models/1
			reqData := openapi.RegisteredModelUpdate{
				Description: openapi.PtrString("This is a new description"),
			}
			body2 := RegisteredModelUpdateEnvelope{Data: &reqData}
			_, rs, err = setupApiTest[RegisteredModelEnvelope](http.MethodPatch, "/api/v1/model_registry/model-registry/registered_models/1", body2, k8sClient, wrongUserIDHeader)
			Expect(err).NotTo(HaveOccurred())
			By("should return a 403 Forbidden response for PATCH update registered model")
			Expect(rs.StatusCode).To(Equal(http.StatusForbidden))

			// Test: GET /registered_models/1/versions
			_, rs, err = setupApiTest[ModelVersionListEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/registered_models/1/versions", nil, k8sClient, wrongUserIDHeader)
			Expect(err).NotTo(HaveOccurred())
			By("should return a 403 Forbidden response for GET model versions of registered model")
			Expect(rs.StatusCode).To(Equal(http.StatusForbidden))

			// Test: POST /registered_models/1/versions
			body3 := ModelVersionEnvelope{Data: openapi.NewModelVersion("Version Fifty", "")}
			_, rs, err = setupApiTest[ModelVersionEnvelope](http.MethodPost, "/api/v1/model_registry/model-registry/registered_models/1/versions", body3, k8sClient, wrongUserIDHeader)
			Expect(err).NotTo(HaveOccurred())
			By("should return a 403 Forbidden response for POST create model version for registered model")
			Expect(rs.StatusCode).To(Equal(http.StatusForbidden))
		})
	})
})
