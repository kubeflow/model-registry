package api

import (
	"context"
	"encoding/json"
	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("TestNamespacesHandler", func() {
	Context("when running in dev mode", Ordered, func() {
		var testApp App

		BeforeAll(func() {
			By("setting up the test app in dev mode")
			testApp = App{
				config:           config.EnvConfig{DevMode: true},
				kubernetesClient: k8sClient,
				repositories:     repositories.NewRepositories(mockMRClient),
				logger:           logger,
			}
		})

		It("should return only dora-namespace for doraNonAdmin@example.com", func() {
			By("creating the HTTP request with the kubeflow-userid header")
			req, err := http.NewRequest(http.MethodGet, NamespaceListPath, nil)
			ctx := context.WithValue(req.Context(), constants.KubeflowUserIdKey, mocks.DoraNonAdminUser)
			req = req.WithContext(ctx)
			Expect(err).NotTo(HaveOccurred())
			rr := httptest.NewRecorder()

			By("calling the GetNamespacesHandler")
			testApp.GetNamespacesHandler(rr, req, nil)
			rs := rr.Result()
			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			Expect(err).NotTo(HaveOccurred())

			By("unmarshalling the response")
			var actual NamespacesEnvelope
			err = json.Unmarshal(body, &actual)
			Expect(err).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusOK))

			By("validating the response contains only dora-namespace")
			expected := []models.NamespaceModel{{Name: "dora-namespace"}}
			Expect(actual.Data).To(ConsistOf(expected))
		})

		It("should return all namespaces for user@example.com", func() {
			By("creating the HTTP request with the kubeflow-userid header")
			req, err := http.NewRequest(http.MethodGet, NamespaceListPath, nil)
			ctx := context.WithValue(req.Context(), constants.KubeflowUserIdKey, mocks.KubeflowUserIDHeaderValue)
			req = req.WithContext(ctx)
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("kubeflow-userid", "user@example.com")
			rr := httptest.NewRecorder()

			By("calling the GetNamespacesHandler")
			testApp.GetNamespacesHandler(rr, req, nil)
			rs := rr.Result()
			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			Expect(err).NotTo(HaveOccurred())

			By("unmarshalling the response")
			var actual NamespacesEnvelope
			err = json.Unmarshal(body, &actual)
			Expect(err).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusOK))

			By("validating the response contains all namespaces")
			expected := []models.NamespaceModel{
				{Name: "kubeflow"},
				{Name: "dora-namespace"},
			}
			Expect(actual.Data).To(ContainElements(expected))
		})

		It("should return no namespaces for non-existent user", func() {
			By("creating the HTTP request with a non-existent kubeflow-userid")
			req, err := http.NewRequest(http.MethodGet, NamespaceListPath, nil)
			ctx := context.WithValue(req.Context(), constants.KubeflowUserIdKey, "nonexistent@example.com")
			req = req.WithContext(ctx)
			Expect(err).NotTo(HaveOccurred())
			rr := httptest.NewRecorder()

			By("calling the GetNamespacesHandler")
			testApp.GetNamespacesHandler(rr, req, nil)
			rs := rr.Result()
			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			Expect(err).NotTo(HaveOccurred())

			By("unmarshalling the response")
			var actual NamespacesEnvelope
			err = json.Unmarshal(body, &actual)
			Expect(err).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusOK))

			By("validating the response contains no namespaces")
			Expect(actual.Data).To(BeEmpty())
		})
	})

})
