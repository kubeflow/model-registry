package integrations

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type HTTPClientInterface interface {
	GetModelRegistryID() (modelRegistryService string)
	GET(url string) ([]byte, error)
	POST(url string, body io.Reader) ([]byte, error)
	PATCH(url string, body io.Reader) ([]byte, error)
}

type HTTPClient struct {
	client          *http.Client
	baseURL         string
	ModelRegistryID string
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type HTTPError struct {
	StatusCode int `json:"-"`
	ErrorResponse
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s - %s", e.StatusCode, e.Code, e.Message)
}

func NewHTTPClient(modelRegistryID string, baseURL string) (HTTPClientInterface, error) {

	return &HTTPClient{
		client: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}},
		baseURL:         baseURL,
		ModelRegistryID: modelRegistryID,
	}, nil
}

func (c *HTTPClient) GetModelRegistryID() string {
	return c.ModelRegistryID
}

func (c *HTTPClient) GET(url string) ([]byte, error) {
	fullURL := c.baseURL + url
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return body, nil
}

func (c *HTTPClient) POST(url string, body io.Reader) ([]byte, error) {
	fullURL := c.baseURL + url
	fmt.Println(fullURL)
	req, err := http.NewRequest("POST", fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if response.StatusCode != http.StatusCreated {
		var errorResponse ErrorResponse
		if err := json.Unmarshal(responseBody, &errorResponse); err != nil {
			return nil, fmt.Errorf("error unmarshalling error response: %w", err)
		}
		httpError := &HTTPError{
			StatusCode:    response.StatusCode,
			ErrorResponse: errorResponse,
		}
		//Sometimes the code comes empty from model registry API
		//also not all error codes are correctly implemented
		//see https://github.com/kubeflow/model-registry/issues/95
		if httpError.ErrorResponse.Code == "" {
			httpError.ErrorResponse.Code = strconv.Itoa(response.StatusCode)
		}
		return nil, httpError
	}

	return responseBody, nil
}

func (c *HTTPClient) PATCH(url string, body io.Reader) ([]byte, error) {
	fullURL := c.baseURL + url
	req, err := http.NewRequest(http.MethodPatch, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.Unmarshal(responseBody, &errorResponse); err != nil {
			return nil, fmt.Errorf("error unmarshalling error response: %w", err)
		}
		httpError := &HTTPError{
			StatusCode:    response.StatusCode,
			ErrorResponse: errorResponse,
		}
		//Sometimes the code comes empty from model registry API
		//also not all error codes are correctly implemented
		//see https://github.com/kubeflow/model-registry/issues/95
		if httpError.ErrorResponse.Code == "" {
			httpError.ErrorResponse.Code = strconv.Itoa(response.StatusCode)
		}
		return nil, httpError
	}
	return responseBody, nil
}
