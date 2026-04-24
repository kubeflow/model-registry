package api

import (
	"net/http"

	"github.com/kubeflow/hub/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/hub/ui/bff/internal/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestMcpServersCatalogHandler", func() {
	Context("testing MCP Servers Catalog Handler", Ordered, func() {

		It("should retrieve all MCP servers", func() {
			By("fetching all MCP servers")
			requestIdentity := kubernetes.RequestIdentity{
				UserID: "user@example.com",
			}

			actual, rs, err := setupApiTest[McpServerListEnvelope](http.MethodGet, "/api/v1/mcp_catalog/mcp_servers?namespace=kubeflow", nil, kubernetesMockedStaticClientFactory, requestIdentity, "kubeflow")
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected MCP server list")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data.Size).To(Equal(int32(10)))
			Expect(actual.Data.PageSize).To(Equal(int32(10)))
			Expect(actual.Data.NextPageToken).NotTo(BeEmpty())
			Expect(len(actual.Data.Items)).To(Equal(10))
		})

		It("should retrieve MCP server filter options", func() {
			By("fetching MCP server filter options")
			data := mocks.GetMcpFilterOptionsListMock()
			requestIdentity := kubernetes.RequestIdentity{
				UserID: "user@example.com",
			}

			expected := McpServerFilterOptionsListEnvelope{Data: &data}
			actual, rs, err := setupApiTest[McpServerFilterOptionsListEnvelope](http.MethodGet, "/api/v1/mcp_catalog/mcp_servers_filter_options?namespace=kubeflow", nil, kubernetesMockedStaticClientFactory, requestIdentity, "kubeflow")
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected filter options")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data).NotTo(BeNil())
			Expect(actual.Data).To(Equal(expected.Data))
		})

		It("should retrieve a single MCP server by id", func() {
			By("fetching MCP server by server_id")
			data := mocks.GetMcpServerMocks()[0]
			requestIdentity := kubernetes.RequestIdentity{
				UserID: "user@example.com",
			}

			actual, rs, err := setupApiTest[McpServerEnvelope](http.MethodGet, "/api/v1/mcp_catalog/mcp_servers/1?namespace=kubeflow", nil, kubernetesMockedStaticClientFactory, requestIdentity, "kubeflow")
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected MCP server")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data).NotTo(BeNil())
			Expect(actual.Data.Name).To(Equal(data.Name))
			Expect(actual.Data.ID).To(Equal(data.ID))
		})

		It("should retrieve MCP server tools", func() {
			By("fetching MCP server tools")
			data := mocks.GetMcpToolListMock()
			requestIdentity := kubernetes.RequestIdentity{
				UserID: "user@example.com",
			}

			expected := McpServerToolsListEnvelope{Data: &data}
			actual, rs, err := setupApiTest[McpServerToolsListEnvelope](http.MethodGet, "/api/v1/mcp_catalog/mcp_servers/1/tools?namespace=kubeflow", nil, kubernetesMockedStaticClientFactory, requestIdentity, "kubeflow")
			Expect(err).NotTo(HaveOccurred())

			By("should match the expected tool list")
			Expect(rs.StatusCode).To(Equal(http.StatusOK))
			Expect(actual.Data).NotTo(BeNil())
			Expect(actual.Data.Size).To(Equal(expected.Data.Size))
			Expect(len(actual.Data.Items)).To(Equal(len(expected.Data.Items)))
		})
	})
})
