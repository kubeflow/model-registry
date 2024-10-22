package mocks

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kubernetes Client Test", func() {
	Context("with existing services", Ordered, func() {

		It("should retrieve the get all service successfully", func() {

			By("getting service details")
			services, err := k8sClient.GetServiceDetails()
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
			service, err := k8sClient.GetServiceDetailsByName("model-registry-dora")
			Expect(err).NotTo(HaveOccurred(), "Failed to create k8s request")

			By("checking that service details are correct")
			Expect(service.Name).To(Equal("model-registry-dora"))
			Expect(service.Description).To(Equal("Model Registry Dora description"))
			Expect(service.DisplayName).To(Equal("Model Registry Dora"))
		})

		It("should retrieve the services names", func() {

			By("getting service by name")
			services, err := k8sClient.GetServiceNames()
			Expect(err).NotTo(HaveOccurred(), "Failed to create HTTP request")

			By("checking that service details are correct")
			Expect(services[0]).To(Equal("model-registry"))
			Expect(services[1]).To(Equal("model-registry-bella"))
			Expect(services[2]).To(Equal("model-registry-dora"))
		})
	})

})
