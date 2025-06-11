package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("TestRoleBindingHandlers", func() {
	Context("Role Binding operations in standalone mode", Ordered, func() {
		var testApp App

		BeforeAll(func() {
			testApp = App{
				config:                  config.EnvConfig{DeploymentMode: config.DeploymentModeStandalone},
				kubernetesClientFactory: kubernetesMockedStaticClientFactory,
				repositories:            repositories.NewRepositories(mockMRClient),
				logger:                  logger,
			}
		})

		It("should retrieve all role bindings", func() {
			req, err := http.NewRequest(http.MethodGet, RoleBindingListPath+"?namespace=kubeflow", nil)
			Expect(err).NotTo(HaveOccurred())

			ctx := context.WithValue(req.Context(), constants.RequestIdentityKey, &kubernetes.RequestIdentity{
				UserID: "user@example.com",
			})
			ctx = context.WithValue(ctx, constants.NamespaceHeaderParameterKey, "kubeflow")
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			testApp.GetRoleBindingsHandler(rr, req, nil)

			rs := rr.Result()
			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			Expect(err).NotTo(HaveOccurred())

			var actual RoleBindingListEnvelope
			err = json.Unmarshal(body, &actual)
			Expect(err).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(actual.Data.Items).To(HaveLen(5))

			// Check the first two stub role bindings
			Expect(actual.Data.Items[0].Name).To(Equal("stub-rb-1"))
			Expect(actual.Data.Items[1].Name).To(Equal("stub-rb-2"))

			// Check model-registry permissions role binding
			Expect(actual.Data.Items[2].Name).To(Equal("model-registry-permissions"))
			Expect(actual.Data.Items[2].Labels["app.kubernetes.io/name"]).To(Equal("model-registry"))
			Expect(actual.Data.Items[2].Subjects[0].Name).To(Equal("admin-user"))

			// Check dora permissions role binding
			Expect(actual.Data.Items[3].Name).To(Equal("model-registry-dora-permissions"))
			Expect(actual.Data.Items[3].Labels["app.kubernetes.io/name"]).To(Equal("model-registry-dora"))
			Expect(actual.Data.Items[3].Subjects[0].Name).To(Equal("dora-user"))

			// Check bella permissions role binding
			Expect(actual.Data.Items[4].Name).To(Equal("model-registry-bella-permissions"))
			Expect(actual.Data.Items[4].Labels["app.kubernetes.io/name"]).To(Equal("model-registry-bella"))
			Expect(actual.Data.Items[4].Subjects[0].Name).To(Equal("bella-team"))
			Expect(actual.Data.Items[4].Subjects[0].Kind).To(Equal("Group"))
		})

		It("should create a role binding", func() {
			newRoleBinding := models.RoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-role-binding",
				},
				Subjects: []rbacv1.Subject{
					{
						Kind: "User",
						Name: "test-user@example.com",
					},
				},
				RoleRef: rbacv1.RoleRef{
					Kind:     "ClusterRole",
					Name:     "view",
					APIGroup: "rbac.authorization.k8s.io",
				},
			}

			reqBody, err := json.Marshal(RoleBindingEnvelope{Data: newRoleBinding})
			Expect(err).NotTo(HaveOccurred())

			req, err := http.NewRequest(http.MethodPost, RoleBindingListPath+"?namespace=kubeflow", bytes.NewReader(reqBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			ctx := context.WithValue(req.Context(), constants.RequestIdentityKey, &kubernetes.RequestIdentity{
				UserID: "user@example.com",
			})
			ctx = context.WithValue(ctx, constants.NamespaceHeaderParameterKey, "kubeflow")
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			testApp.CreateRoleBindingHandler(rr, req, nil)

			rs := rr.Result()
			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			Expect(err).NotTo(HaveOccurred())

			var actual RoleBindingEnvelope
			err = json.Unmarshal(body, &actual)
			Expect(err).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusCreated))
			Expect(actual.Data.Name).To(Equal("test-role-binding"))
			Expect(actual.Data.Subjects[0].Name).To(Equal("test-user@example.com"))
		})

		It("should patch a role binding", func() {
			updatedRoleBinding := models.RoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-role-binding",
				},
				Subjects: []rbacv1.Subject{
					{
						Kind: "User",
						Name: "updated-user@example.com",
					},
				},
				RoleRef: rbacv1.RoleRef{
					Kind:     "ClusterRole",
					Name:     "edit",
					APIGroup: "rbac.authorization.k8s.io",
				},
			}

			reqBody, err := json.Marshal(RoleBindingEnvelope{Data: updatedRoleBinding})
			Expect(err).NotTo(HaveOccurred())

			req, err := http.NewRequest(http.MethodPatch, RoleBindingListPath+"/test-role-binding?namespace=kubeflow", bytes.NewReader(reqBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			ctx := context.WithValue(req.Context(), constants.RequestIdentityKey, &kubernetes.RequestIdentity{
				UserID: "user@example.com",
			})
			ctx = context.WithValue(ctx, constants.NamespaceHeaderParameterKey, "kubeflow")
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			// We need to simulate the router params since we're calling the handler directly
			ps := httprouter.Params{
				{Key: RoleBindingNameParam, Value: "test-role-binding"},
			}
			testApp.PatchRoleBindingHandler(rr, req, ps)

			rs := rr.Result()
			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			Expect(err).NotTo(HaveOccurred())

			var actual RoleBindingEnvelope
			err = json.Unmarshal(body, &actual)
			Expect(err).NotTo(HaveOccurred())

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(actual.Data.Name).To(Equal("test-role-binding"))
			Expect(actual.Data.Subjects[0].Name).To(Equal("updated-user@example.com"))
			Expect(actual.Data.RoleRef.Name).To(Equal("edit"))
		})

		It("should delete a role binding", func() {
			req, err := http.NewRequest(http.MethodDelete, RoleBindingListPath+"/test-role-binding?namespace=kubeflow", nil)
			Expect(err).NotTo(HaveOccurred())

			ctx := context.WithValue(req.Context(), constants.RequestIdentityKey, &kubernetes.RequestIdentity{
				UserID: "user@example.com",
			})
			ctx = context.WithValue(ctx, constants.NamespaceHeaderParameterKey, "kubeflow")
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			// We need to simulate the router params since we're calling the handler directly
			ps := httprouter.Params{
				{Key: RoleBindingNameParam, Value: "test-role-binding"},
			}
			testApp.DeleteRoleBindingHandler(rr, req, ps)

			rs := rr.Result()
			defer rs.Body.Close()

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})
	})
})
