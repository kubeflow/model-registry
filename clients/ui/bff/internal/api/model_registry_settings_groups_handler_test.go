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
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestGroupsHandler", func() {
	Context("when fetching group list in standalone mode", Ordered, func() {
		var testApp App

		BeforeAll(func() {
			testApp = App{
				config:                  config.EnvConfig{StandaloneMode: true},
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

			expected := []models.GroupModel{
				{Name: "dora-group-mock"},
				{Name: "bella-group-mock"},
			}

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(actual.Data).To(ConsistOf(expected))
		})
	})
})
