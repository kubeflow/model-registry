package api

import (
	"context"
	"encoding/json"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	KubeflowUserIDHeaderValue = "user@example.com"
	DoraNonAdminUser          = "doraNonAdmin@example.com"
)

var _ = Describe("TestUserHandler", func() {
	Context("fetching user details", Ordered, func() {
		var testApp App

		BeforeAll(func() {
			By("creating the test app")
			testApp = App{
				kubernetesClient: k8sClient,
				repositories:     repositories.NewRepositories(mockMRClient),
				logger:           logger,
			}
		})

		It("should show that KubeflowUserIDHeaderValue (user@example.com) is a cluster-admin", func() {
			By("creating the http request")
			req, err := http.NewRequest(http.MethodGet, UserPath, nil)
			ctx := context.WithValue(req.Context(), KubeflowUserIdKey, mocks.KubeflowUserIDHeaderValue)
			req = req.WithContext(ctx)
			Expect(err).NotTo(HaveOccurred())

			By("creating the http test infrastructure")
			rr := httptest.NewRecorder()

			By("invoking the UserHandler")
			testApp.UserHandler(rr, req, nil)
			rs := rr.Result()
			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			Expect(err).NotTo(HaveOccurred())

			By("unmarshalling the user response")
			var actual UserEnvelope
			err = json.Unmarshal(body, &actual)
			Expect(err).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusOK))

			By("checking that the user is cluster-admin")
			Expect(actual.Data.UserID).To(Equal(KubeflowUserIDHeaderValue))
			Expect(actual.Data.ClusterAdmin).To(BeTrue(), "Expected this user to be cluster-admin")
		})

		It("should show that DoraNonAdminUser (doraNonAdmin@example.com) is not a cluster-admin", func() {
			By("creating the http request")
			req, err := http.NewRequest(http.MethodGet, UserPath, nil)
			ctx := context.WithValue(req.Context(), KubeflowUserIdKey, DoraNonAdminUser)
			req = req.WithContext(ctx)
			Expect(err).NotTo(HaveOccurred())

			By("creating the http test infrastructure")
			rr := httptest.NewRecorder()

			By("invoking the UserHandler")
			testApp.UserHandler(rr, req, nil)
			rs := rr.Result()
			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			Expect(err).NotTo(HaveOccurred())

			By("unmarshalling the user response")
			var actual UserEnvelope
			err = json.Unmarshal(body, &actual)
			Expect(err).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusOK))

			By("checking that the user is not cluster-admin")
			Expect(actual.Data.UserID).To(Equal(DoraNonAdminUser))
			Expect(actual.Data.ClusterAdmin).To(BeFalse(), "Expected this user to not be cluster-admin")
		})

		It("should show that a random non-existent user is not a cluster-admin", func() {
			randomUser := "bellaUser@example.com"

			By("creating the http request")
			req, err := http.NewRequest(http.MethodGet, UserPath, nil)
			ctx := context.WithValue(req.Context(), KubeflowUserIdKey, randomUser)
			req = req.WithContext(ctx)
			Expect(err).NotTo(HaveOccurred())

			By("creating the http test infrastructure")
			rr := httptest.NewRecorder()

			By("invoking the UserHandler")
			testApp.UserHandler(rr, req, nil)
			rs := rr.Result()
			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			Expect(err).NotTo(HaveOccurred())

			By("unmarshalling the user response")
			var actual UserEnvelope
			err = json.Unmarshal(body, &actual)
			Expect(err).NotTo(HaveOccurred())
			Expect(rr.Code).To(Equal(http.StatusOK))

			By("checking that the user is not cluster-admin")
			Expect(actual.Data.UserID).To(Equal(randomUser))
			Expect(actual.Data.ClusterAdmin).To(BeFalse(), "Expected this user to not be cluster-admin")
		})
	})

})
