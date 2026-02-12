package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CheckNamespaceRegistryAccessHandler", func() {
	Context("when identity and body are valid", Ordered, func() {
		var testApp App

		BeforeAll(func() {
			testApp = App{
				config:                  config.EnvConfig{DevMode: true},
				kubernetesClientFactory: kubernetesMockedStaticClientFactory,
				repositories:            repositories.NewRepositories(mockMRClient, mockModelCatalogClient),
				logger:                  logger,
			}
		})

		It("should return hasAccess true when namespace default SA has access to the registry", func() {
			body := CheckNamespaceRegistryAccessRequestEnvelope{
				Data: CheckNamespaceRegistryAccessRequest{
					Namespace:         "dora-namespace",
					RegistryName:      "model-registry-dora",
					RegistryNamespace: "dora-namespace",
				},
			}
			bodyBytes, err := json.Marshal(body)
			Expect(err).NotTo(HaveOccurred())

			req, err := http.NewRequest(http.MethodPost, CheckNamespaceRegistryAccessPath, bytes.NewReader(bodyBytes))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			reqIdentity := &kubernetes.RequestIdentity{UserID: KubeflowUserIDHeaderValue}
			ctx := context.WithValue(req.Context(), constants.RequestIdentityKey, reqIdentity)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			testApp.CheckNamespaceRegistryAccessHandler(rr, req, nil)

			rs := rr.Result()
			defer rs.Body.Close()
			Expect(rs.StatusCode).To(Equal(http.StatusOK))

			respBody, err := io.ReadAll(rs.Body)
			Expect(err).NotTo(HaveOccurred())
			var envelope CheckNamespaceRegistryAccessEnvelope
			Expect(json.Unmarshal(respBody, &envelope)).To(Succeed())
			Expect(envelope.Data.HasAccess).To(BeTrue())
		})

		It("should return hasAccess false when namespace default SA has no access to the registry", func() {
			body := CheckNamespaceRegistryAccessRequestEnvelope{
				Data: CheckNamespaceRegistryAccessRequest{
					Namespace:         "bella-namespace",
					RegistryName:      "model-registry-dora",
					RegistryNamespace: "dora-namespace",
				},
			}
			bodyBytes, err := json.Marshal(body)
			Expect(err).NotTo(HaveOccurred())

			req, err := http.NewRequest(http.MethodPost, CheckNamespaceRegistryAccessPath, bytes.NewReader(bodyBytes))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			reqIdentity := &kubernetes.RequestIdentity{UserID: KubeflowUserIDHeaderValue}
			ctx := context.WithValue(req.Context(), constants.RequestIdentityKey, reqIdentity)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			testApp.CheckNamespaceRegistryAccessHandler(rr, req, nil)

			rs := rr.Result()
			defer rs.Body.Close()
			Expect(rs.StatusCode).To(Equal(http.StatusOK))

			respBody, err := io.ReadAll(rs.Body)
			Expect(err).NotTo(HaveOccurred())
			var envelope CheckNamespaceRegistryAccessEnvelope
			Expect(json.Unmarshal(respBody, &envelope)).To(Succeed())
			Expect(envelope.Data.HasAccess).To(BeFalse())
		})
	})

	Context("when request is invalid", Ordered, func() {
		var testApp App

		BeforeAll(func() {
			testApp = App{
				config:                  config.EnvConfig{DevMode: true},
				kubernetesClientFactory: kubernetesMockedStaticClientFactory,
				repositories:            repositories.NewRepositories(mockMRClient, mockModelCatalogClient),
				logger:                  logger,
			}
		})

		It("should return 400 when identity is missing", func() {
			body := CheckNamespaceRegistryAccessRequestEnvelope{
				Data: CheckNamespaceRegistryAccessRequest{
					Namespace:         "dora-namespace",
					RegistryName:      "model-registry-dora",
					RegistryNamespace: "dora-namespace",
				},
			}
			bodyBytes, _ := json.Marshal(body)
			req, _ := http.NewRequest(http.MethodPost, CheckNamespaceRegistryAccessPath, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			// no RequestIdentityKey in context

			rr := httptest.NewRecorder()
			testApp.CheckNamespaceRegistryAccessHandler(rr, req, nil)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 400 when required body fields are missing", func() {
			body := CheckNamespaceRegistryAccessRequestEnvelope{
				Data: CheckNamespaceRegistryAccessRequest{
					Namespace: "dora-namespace",
				},
			}
			bodyBytes, _ := json.Marshal(body)
			req, _ := http.NewRequest(http.MethodPost, CheckNamespaceRegistryAccessPath, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			reqIdentity := &kubernetes.RequestIdentity{UserID: KubeflowUserIDHeaderValue}
			ctx := context.WithValue(req.Context(), constants.RequestIdentityKey, reqIdentity)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			testApp.CheckNamespaceRegistryAccessHandler(rr, req, nil)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})
})
