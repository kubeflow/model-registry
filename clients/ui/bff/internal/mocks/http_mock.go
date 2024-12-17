package mocks

import (
	"github.com/stretchr/testify/mock"
	"io"
)

type MockHTTPClient struct {
	mock.Mock
}

func (c *MockHTTPClient) GetModelRegistryID() string {
	return "model-registry"
}

func (m *MockHTTPClient) GET(url string) ([]byte, error) {
	args := m.Called(url)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockHTTPClient) POST(url string, body io.Reader) ([]byte, error) {
	args := m.Called(url, body)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockHTTPClient) PATCH(url string, body io.Reader) ([]byte, error) {
	args := m.Called(url, body)
	return args.Get(0).([]byte), args.Error(1)
}
