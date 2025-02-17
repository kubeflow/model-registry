package helper

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"io"
	"log/slog"
	"net/http"
	"slices"
)

func GetContextLoggerFromReq(r *http.Request) *slog.Logger {
	return GetContextLogger(r.Context())
}

func GetContextLogger(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(constants.TraceLoggerKey).(*slog.Logger)

	if !ok {
		logger = slog.New(slog.Default().Handler())
		logger.Warn("Unable to get context logger for tracing, falling back to default")
	}

	return logger
}

var sensitiveHeaders = []string{
	"Authorization",
	"Cookie",
	"Set-Cookie",
	"Proxy-Authorization",
}

func isSensitiveHeader(h string) bool {
	return slices.Contains(sensitiveHeaders, http.CanonicalHeaderKey(h))
}

type HeaderLogValuer struct {
	Header http.Header
}

func (h HeaderLogValuer) LogValue() slog.Value {
	var values []slog.Attr

	for k, v := range h.Header {
		if len(v) == 0 {
			values = append(values, slog.String(k, ""))
			continue
		}

		if isSensitiveHeader(k) {
			values = append(values, slog.String(k, "[REDACTED]"))
			continue
		}

		values = append(values, slog.String(k, v[0]))
	}

	return slog.GroupValue(values...)
}

func CloneBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("no body provided")
	}
	buf, _ := io.ReadAll(r.Body)
	readerCopy := io.NopCloser(bytes.NewBuffer(buf))
	readerOriginal := io.NopCloser(bytes.NewBuffer(buf))
	r.Body = readerOriginal

	defer readerCopy.Close()
	cloneBody, err := io.ReadAll(readerCopy)

	return cloneBody, err
}

type RequestLogValuer struct {
	Request *http.Request
}

func (r RequestLogValuer) LogValue() slog.Value {
	body := ""

	if r.Request.Body != nil {
		cloneBody, err := CloneBody(r.Request)
		if err != nil {
			body = fmt.Sprintf("error: %v", err)
		} else {
			body = string(cloneBody)
		}
	}

	return slog.GroupValue(
		slog.String("method", r.Request.Method),
		slog.String("url", r.Request.URL.String()),
		slog.String("body", body),
		slog.Any("headers", HeaderLogValuer{Header: r.Request.Header}))
}

type ResponseLogValuer struct {
	Response *http.Response
	Body     []byte
}

func (r ResponseLogValuer) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("status_code", r.Response.StatusCode),
		slog.String("status", r.Response.Status),
		slog.Any("body", r.Body),
		slog.Any("headers", HeaderLogValuer{Header: r.Response.Header}))
}
