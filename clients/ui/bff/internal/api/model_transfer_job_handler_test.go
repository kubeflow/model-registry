package api

import (
	"net/http"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Model transfer job handler", func() {
	requestIdentity := kubernetes.RequestIdentity{
		UserID: "user@example.com",
	}

	Describe("GET list", func() {
		It("should list model transfer jobs for namespace", func() {
			actual, rs, err := setupApiTest[ModelTransferJobListEnvelope](
				http.MethodGet,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=bella-namespace",
				nil,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"bella-namespace",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data).NotTo(BeNil())
			Expect(actual.Data.Items).NotTo(BeEmpty())

			// Envtest seeds transfer-job-001 in bella-namespace with label job-id=001
			var found bool
			for _, item := range actual.Data.Items {
				if item.Id == "001" && item.Name == "transfer-job-001" {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue(), "expected at least one job with id 001 and name transfer-job-001")
		})

		It("should return list for kubeflow namespace", func() {
			actual, rs, err := setupApiTest[ModelTransferJobListEnvelope](
				http.MethodGet,
				"/api/v1/model_registry/model-registry/model_transfer_jobs?namespace=kubeflow",
				nil,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"kubeflow",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data).NotTo(BeNil())
			Expect(actual.Data.Items).NotTo(BeEmpty())
		})
	})

	Describe("DELETE", func() {
		It("should delete model transfer job by job name and return 200", func() {
			// Path param is job name (K8s Job resource name from list response). BFF deletes by name.
			// BFF returns 200 with JSON body so the frontend can parse the response (204 No Content breaks response.json()).
			actual, rs, err := setupApiTest[ModelTransferJobOperationStatusEnvelope](
				http.MethodDelete,
				"/api/v1/model_registry/model-registry/model_transfer_jobs/transfer-job-001?namespace=bella-namespace",
				nil,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"bella-namespace",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data.Status).To(Equal("deleted"))
		})

		It("should return 404 when deleting non-existent job", func() {
			_, rs, err := setupApiTest[ModelTransferJobListEnvelope](
				http.MethodDelete,
				"/api/v1/model_registry/model-registry/model_transfer_jobs/nonexistent-job?namespace=bella-namespace",
				nil,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"bella-namespace",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should return 400 when namespace is missing", func() {
			// Omit namespace query param so AttachNamespace middleware returns 400
			_, rs, err := setupApiTest[ModelTransferJobListEnvelope](
				http.MethodDelete,
				"/api/v1/model_registry/model-registry/model_transfer_jobs/transfer-job-001",
				nil,
				kubernetesMockedStaticClientFactory,
				requestIdentity,
				"",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(rs.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})
})
