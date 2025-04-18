package api

import (
	"io"
	"net/http"
	httptest "net/http/httptest"

	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Static File serving Test", func() {
	var (
		server *httptest.Server
		client *http.Client
	)

	Context("serving static files at /", Ordered, func() {

		BeforeAll(func() {
			envConfig := config.EnvConfig{
				StaticAssetsDir: resolveStaticAssetsDirOnTests(),
			}
			app := &App{
				kubernetesClientFactory: kubernetesMockedStaticClientFactory,
				repositories:            repositories.NewRepositories(mockMRClient),
				logger:                  logger,
				config:                  envConfig,
			}

			server = httptest.NewServer(app.Routes())
			client = server.Client()
		})

		AfterAll(func() {
			server.Close()
		})

		It("should serve index.html from the root path", func() {
			resp, err := client.Get(server.URL + "/")
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(ContainSubstring("BFF Stub Page"))
		})

		It("should serve subfolders from the root path", func() {
			resp, err := client.Get(server.URL + "/sub/test.html")
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(ContainSubstring("BFF Stub Subfolder Page"))
		})

		It("should return index.html for a non-existent static file", func() {
			resp, err := client.Get(server.URL + "/non-existent.html")
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			//BFF Stub page is the context of index.html
			Expect(string(body)).To(ContainSubstring("BFF Stub Page"))
		})

	})
})
