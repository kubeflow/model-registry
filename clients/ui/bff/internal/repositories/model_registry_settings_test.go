package repositories

import (
	"context"

	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestModelRegistrySettingsRepository", func() {
	Context("fetching groups from the stub", Ordered, func() {

		It("should return the mocked group list", func() {
			By("initializing the repository and client")
			repo := NewModelRegistrySettingsRepository()
			k8sClient, err := kubernetesMockedStaticClientFactory.GetClient(mocks.NewMockSessionContextNoParent())
			Expect(err).NotTo(HaveOccurred())

			By("fetching groups")
			groups, err := repo.GetGroups(context.Background(), k8sClient)
			Expect(err).NotTo(HaveOccurred())

			By("verifying the returned group models")
			Expect(groups).To(HaveLen(2))

			// Verify first group
			firstGroup := groups[0]
			Expect(firstGroup.Metadata.Name).To(Equal("dora-group-mock"))
			Expect(firstGroup.Users).To(ConsistOf("dora-user@example.com", "dora-admin@example.com"))
			Expect(*firstGroup.APIVersion).To(Equal("user.openshift.io/v1"))
			Expect(*firstGroup.Kind).To(Equal("Group"))
			Expect(firstGroup).To(BeAssignableToTypeOf(models.Group{}))

			// Verify second group
			secondGroup := groups[1]
			Expect(secondGroup.Metadata.Name).To(Equal("bella-group-mock"))
			Expect(secondGroup.Users).To(ConsistOf("bella-user@example.com", "bella-maintainer@example.com"))
			Expect(*secondGroup.APIVersion).To(Equal("user.openshift.io/v1"))
			Expect(*secondGroup.Kind).To(Equal("Group"))
		})
	})

	Context("creating model registry with default database", Ordered, func() {
		It("should create ModelRegistry CR without Secret", func() {
			By("initializing the repository and client")
			repo := NewModelRegistrySettingsRepository()
			k8sClient, err := kubernetesMockedStaticClientFactory.GetClient(mocks.NewMockSessionContextNoParent())
			Expect(err).NotTo(HaveOccurred())

			By("preparing payload for default Postgres database")
			generateDeployment := true
			payload := models.ModelRegistrySettingsPayload{
				ModelRegistry: models.ModelRegistryKind{
					Metadata: models.Metadata{
						Name:      "test-registry-default",
						Namespace: "default",
					},
					Spec: models.ModelRegistrySpec{
						GRPC: models.EmptyObject{},
						REST: models.EmptyObject{},
						Istio: models.IstioConfig{
							Gateway: models.GatewayConfig{
								GRPC: models.GRPCConfig{TLS: models.EmptyObject{}},
								REST: models.RESTConfig{TLS: models.EmptyObject{}},
							},
						},
						Postgres: &models.PostgresConfig{
							Database:           "model_registry",
							GenerateDeployment: &generateDeployment,
						},
					},
				},
				DatabasePassword: nil, // No password for default database
			}

			By("creating the model registry")
			result, err := repo.CreateModelRegistry(context.Background(), k8sClient, "default", payload)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())

			By("verifying the created ModelRegistry")
			Expect(result.Metadata.Name).To(Equal("test-registry-default"))
			Expect(result.Spec.Postgres).NotTo(BeNil())
			Expect(result.Spec.Postgres.GenerateDeployment).NotTo(BeNil())
			Expect(*result.Spec.Postgres.GenerateDeployment).To(BeTrue())
			Expect(result.Spec.Postgres.PasswordSecret).To(BeNil(), "Default database should not have password secret")
		})
	})

	Context("creating model registry with external MySQL", Ordered, func() {
		It("should create ModelRegistry CR and Secret", func() {
			By("initializing the repository and client")
			repo := NewModelRegistrySettingsRepository()
			k8sClient, err := kubernetesMockedStaticClientFactory.GetClient(mocks.NewMockSessionContextNoParent())
			Expect(err).NotTo(HaveOccurred())

			By("preparing payload for external MySQL database")
			port := 3306
			password := "my-secret-password"
			payload := models.ModelRegistrySettingsPayload{
				ModelRegistry: models.ModelRegistryKind{
					Metadata: models.Metadata{
						Name:      "test-registry-mysql",
						Namespace: "default",
					},
					Spec: models.ModelRegistrySpec{
						GRPC: models.EmptyObject{},
						REST: models.EmptyObject{},
						Istio: models.IstioConfig{
							Gateway: models.GatewayConfig{
								GRPC: models.GRPCConfig{TLS: models.EmptyObject{}},
								REST: models.RESTConfig{TLS: models.EmptyObject{}},
							},
						},
						MySQL: &models.MySQLConfig{
							Host:     "mysql.example.com",
							Database: "model_registry",
							Username: "admin",
							Port:     &port,
						},
					},
				},
				DatabasePassword: &password,
			}

			By("creating the model registry")
			result, err := repo.CreateModelRegistry(context.Background(), k8sClient, "default", payload)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())

			By("verifying the created ModelRegistry")
			Expect(result.Metadata.Name).To(Equal("test-registry-mysql"))
			Expect(result.Spec.MySQL).NotTo(BeNil())
			Expect(result.Spec.MySQL.Host).To(Equal("mysql.example.com"))
			Expect(result.Spec.MySQL.Username).To(Equal("admin"))

			By("verifying the password secret reference")
			Expect(result.Spec.MySQL.PasswordSecret).NotTo(BeNil())
			Expect(result.Spec.MySQL.PasswordSecret.Name).To(Equal("test-registry-mysql-database-password"))
			Expect(result.Spec.MySQL.PasswordSecret.Key).To(Equal("database-password"))
		})
	})

	Context("creating model registry with external PostgreSQL", Ordered, func() {
		It("should create ModelRegistry CR and Secret", func() {
			By("initializing the repository and client")
			repo := NewModelRegistrySettingsRepository()
			k8sClient, err := kubernetesMockedStaticClientFactory.GetClient(mocks.NewMockSessionContextNoParent())
			Expect(err).NotTo(HaveOccurred())

			By("preparing payload for external PostgreSQL database")
			port := 5432
			password := "postgres-secret-password"
			payload := models.ModelRegistrySettingsPayload{
				ModelRegistry: models.ModelRegistryKind{
					Metadata: models.Metadata{
						Name:      "test-registry-postgres",
						Namespace: "default",
					},
					Spec: models.ModelRegistrySpec{
						GRPC: models.EmptyObject{},
						REST: models.EmptyObject{},
						Istio: models.IstioConfig{
							Gateway: models.GatewayConfig{
								GRPC: models.GRPCConfig{TLS: models.EmptyObject{}},
								REST: models.RESTConfig{TLS: models.EmptyObject{}},
							},
						},
						Postgres: &models.PostgresConfig{
							Host:     "postgres.example.com",
							Database: "model_registry",
							Username: "postgres",
							Port:     &port,
							SSLMode:  "require",
						},
					},
				},
				DatabasePassword: &password,
			}

			By("creating the model registry")
			result, err := repo.CreateModelRegistry(context.Background(), k8sClient, "default", payload)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())

			By("verifying the created ModelRegistry")
			Expect(result.Metadata.Name).To(Equal("test-registry-postgres"))
			Expect(result.Spec.Postgres).NotTo(BeNil())
			Expect(result.Spec.Postgres.Host).To(Equal("postgres.example.com"))
			Expect(result.Spec.Postgres.Username).To(Equal("postgres"))
			Expect(result.Spec.Postgres.SSLMode).To(Equal("require"))

			By("verifying the password secret reference")
			Expect(result.Spec.Postgres.PasswordSecret).NotTo(BeNil())
			Expect(result.Spec.Postgres.PasswordSecret.Name).To(Equal("test-registry-postgres-database-password"))
			Expect(result.Spec.Postgres.PasswordSecret.Key).To(Equal("database-password"))
		})
	})

	Context("error handling", Ordered, func() {
		It("should return error when registry name is empty", func() {
			By("initializing the repository and client")
			repo := NewModelRegistrySettingsRepository()
			k8sClient, err := kubernetesMockedStaticClientFactory.GetClient(mocks.NewMockSessionContextNoParent())
			Expect(err).NotTo(HaveOccurred())

			By("preparing payload with empty name")
			payload := models.ModelRegistrySettingsPayload{
				ModelRegistry: models.ModelRegistryKind{
					Metadata: models.Metadata{
						Name:      "", // Empty name
						Namespace: "default",
					},
					Spec: models.ModelRegistrySpec{},
				},
			}

			By("attempting to create the model registry")
			result, err := repo.CreateModelRegistry(context.Background(), k8sClient, "default", payload)

			By("verifying error is returned")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("model registry name is required"))
			Expect(result).To(BeNil())
		})
	})
})
