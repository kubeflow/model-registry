package mocks

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kubernetes ControllerRuntimeClient Test", func() {
	Context("with existing services", Ordered, func() {

		It("should retrieve the get all service successfully", func() {

			By("getting service details")
			services, err := k8sClient.GetServiceDetails(NewMockSessionContextNoParent(), "kubeflow")
			Expect(err).NotTo(HaveOccurred(), "Failed to create HTTP request")

			By("checking that all services have the modified ClusterIP and HTTPPort")
			for _, service := range services {
				Expect(service.ClusterIP).To(Equal("127.0.0.1"), "ClusterIP should be set to 127.0.0.1")
				Expect(service.HTTPPort).To(Equal(int32(8080)), "HTTPPort should be set to 8080")

			}

			By("checking that that a specific service exists")
			foundService := false
			for _, service := range services {
				if service.Name == "model-registry" {
					foundService = true
					Expect(service.DisplayName).To(Equal("Model Registry"))
					Expect(service.Description).To(Equal("Model Registry Description"))
					break
				}
			}
			Expect(foundService).To(Equal(true), "Expected to find service 'model-registry'")
		})

		It("should retrieve the service details by name", func() {

			By("getting service by name")
			service, err := k8sClient.GetServiceDetailsByName(NewMockSessionContextNoParent(), "dora-namespace", "model-registry-dora")
			Expect(err).NotTo(HaveOccurred(), "Failed to create k8s request")

			By("checking that service details are correct")
			Expect(service.Name).To(Equal("model-registry-dora"))
			Expect(service.Description).To(Equal("Model Registry Dora description"))
			Expect(service.DisplayName).To(Equal("Model Registry Dora"))
		})

		It("should retrieve the services names", func() {

			By("getting service by name")
			services, err := k8sClient.GetServiceNames(NewMockSessionContextNoParent(), "kubeflow")
			Expect(err).NotTo(HaveOccurred(), "Failed to create HTTP request")

			By("checking that service details are correct")
			Expect(services).To(ConsistOf("model-registry", "model-registry-one"))

		})
	})

})

var _ = Describe("KubernetesNativeClient SAR Test", func() {
	Context("Subject Access Review", func() {

		It("should allow allowed user to access services", func() {
			By("performing SAR for Kubeflow User ID")
			allowed, err := k8sClient.PerformSARonGetListServicesByNamespace(KubeflowUserIDHeaderValue, []string{}, "kubeflow")
			Expect(err).NotTo(HaveOccurred(), "Failed to perform SAR for Kubeflow User ID\"")
			Expect(allowed).To(BeTrue(), "Expected Kubeflow User ID to have access")
		})

		It("check dora access to namespaces", func() {
			By("performing SAR for dora user")
			allowed, err := k8sClient.PerformSARonGetListServicesByNamespace(DoraNonAdminUser, []string{}, "kubeflow")
			Expect(err).NotTo(HaveOccurred(), "Failed to perform SAR for unauthorized-dora@example.com")
			Expect(allowed).To(BeFalse(), "Expected doraNonAdmin@example.com to be denied access")

			allowed, err = k8sClient.PerformSARonGetListServicesByNamespace(DoraNonAdminUser, []string{}, "dora-namespace")
			Expect(err).NotTo(HaveOccurred(), "Failed to perform SAR for unauthorized-dora@example.com")
			Expect(allowed).To(BeTrue(), "Expected doraNonAdmin@example.com ID to have access")
		})

		It("should deny access for another user", func() {
			By("performing SAR for another user")
			allowed, err := k8sClient.PerformSARonGetListServicesByNamespace("unauthorized-dora@example.com", []string{}, "kubeflow")
			Expect(err).NotTo(HaveOccurred(), "Failed to perform SAR for unauthorized-dora@example.com")
			Expect(allowed).To(BeFalse(), "Expected unauthorized-dora@example.com to be denied access")
		})

	})
})

var _ = Describe("KubernetesClient PerformSARonSpecificService Group Tests", func() {
	Context("checking access using group memberships", func() {
		const (
			namespace    = "dora-namespace"
			serviceName  = "model-registry-dora"
			existingUser = "bentoOnlyGroupAccess@example.com"
		)

		It("should deny access for a group that does not exist", func() {
			groups := []string{"non-existent-group"}

			allowed, err := k8sClient.PerformSARonSpecificService(existingUser, groups, namespace, serviceName)
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeFalse(), "Access should be denied for a non-existent group")
		})

		It("should allow service access for the DoraServiceGroup", func() {
			groups := []string{DoraServiceGroup}

			allowed, err := k8sClient.PerformSARonSpecificService(existingUser, groups, namespace, serviceName)
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeTrue(), "Access should be allowed for the DoraServiceGroup group")
		})

		It("should allow access when one group exists and the other does not", func() {
			groups := []string{DoraServiceGroup, "non-existent-group"}

			allowed, err := k8sClient.PerformSARonSpecificService(existingUser, groups, namespace, serviceName)
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeTrue(), "Access should be allowed if any group in the list has access")
		})

		It("should allow access only when I've service access and namespace access", func() {
			groups := []string{DoraServiceGroup}

			allowed, err := k8sClient.PerformSARonSpecificService(existingUser, groups, namespace, serviceName)
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeTrue(), "Access should be allowed for the DoraServiceGroup group")

			allowed, err = k8sClient.PerformSARonGetListServicesByNamespace(existingUser, groups, namespace)
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeFalse(), "Access should not be allowed for only DoraServiceGroup group")

			allGroups := []string{DoraServiceGroup, DoraNamespaceGroup}

			allowed, err = k8sClient.PerformSARonGetListServicesByNamespace(existingUser, allGroups, namespace)
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeTrue(), "Access should be allowed for both groups")

			allowed, err = k8sClient.PerformSARonSpecificService(existingUser, allGroups, namespace, serviceName)
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(BeTrue(), "Access should be allowed for for both groups")

		})
	})
})

var _ = Describe("KubernetesClient isClusterAdmin Test", func() {
	Context("checking cluster admin status", func() {
		It("should confirm that user@example.com(KubeflowUserIDHeaderValue) is a cluster-admin", func() {
			isAdmin, err := k8sClient.IsClusterAdmin(KubeflowUserIDHeaderValue)
			Expect(err).NotTo(HaveOccurred())
			Expect(isAdmin).To(BeTrue())
		})
		It("should confirm that doraNonAdmin@example.com(DoraNonAdminUser) is a not cluster-admin", func() {
			isAdmin, err := k8sClient.IsClusterAdmin(DoraNonAdminUser)
			Expect(err).NotTo(HaveOccurred())
			Expect(isAdmin).To(BeFalse())
		})
		It("should confirm that non existent user is not a cluster-admin", func() {
			isAdmin, err := k8sClient.IsClusterAdmin("bella@non-existent.com")
			Expect(err).NotTo(HaveOccurred())
			Expect(isAdmin).To(BeFalse())
		})
	})
})
