package api

import (
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

var _ = Describe("TestGroupsHandler", func() {
	Context("when fetching group list in standalone mode", Ordered, func() {
		var testApp App

		BeforeAll(func() {
			testApp = App{
				config:                  config.EnvConfig{DeploymentMode: config.DeploymentModeStandalone},
				kubernetesClientFactory: kubernetesMockedStaticClientFactory,
				repositories:            repositories.NewRepositories(mockMRClient),
				logger:                  logger,
			}
		})

		It("should return the group names for user@example.com", func() {
			req, err := http.NewRequest(http.MethodGet, GroupsPath+"?namespace=kubeflow", nil)
			Expect(err).NotTo(HaveOccurred())

			ctx := context.WithValue(req.Context(), constants.RequestIdentityKey, &kubernetes.RequestIdentity{
				UserID: "user@example.com",
			})
			ctx = context.WithValue(ctx, constants.NamespaceHeaderParameterKey, "kubeflow")
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			testApp.GetGroupsHandler(rr, req, nil)

			rs := rr.Result()
			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			Expect(err).NotTo(HaveOccurred())

			var actual GroupsEnvelope
			err = json.Unmarshal(body, &actual)
			Expect(err).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(actual.Data).To(HaveLen(2))

			// Verify first group
			firstGroup := actual.Data[0]
			Expect(firstGroup.Metadata.Name).To(Equal("dora-group-mock"))
			Expect(firstGroup.Users).To(ConsistOf("dora-user@example.com", "dora-admin@example.com"))

			// Verify second group
			secondGroup := actual.Data[1]
			Expect(secondGroup.Metadata.Name).To(Equal("bella-group-mock"))
			Expect(secondGroup.Users).To(ConsistOf("bella-user@example.com", "bella-maintainer@example.com"))
		})
	})
})
