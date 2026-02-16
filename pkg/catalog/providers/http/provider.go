// Package http provides a base HTTP provider for catalog data.
// It handles HTTP fetching, polling, authentication, and rate limiting,
// while delegating entity-specific conversion to user-provided functions.
package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kubeflow/model-registry/pkg/catalog"
)

// Config configures an HTTP provider.
type Config[E any, A any] struct {
	// BaseURLKey is the property key in source.Properties for the base URL.
	// Defaults to "url" if empty.
	BaseURLKey string

	// DefaultBaseURL is used if no URL is provided in source.Properties.
	DefaultBaseURL string

	// SyncIntervalKey is the property key for sync interval.
	// Defaults to "syncInterval" if empty.
	SyncIntervalKey string

	// DefaultSyncInterval is the default polling interval.
	// Defaults to 24 hours if zero.
	DefaultSyncInterval time.Duration

	// HTTPClient is the HTTP client to use.
	// Defaults to a client with 30 second timeout if nil.
	HTTPClient *http.Client

	// FetchRecords fetches records from the HTTP API.
	// This is the main function that catalog-specific providers must implement.
	FetchRecords func(ctx context.Context, client *http.Client, baseURL string, source *catalog.Source) ([]catalog.Record[E, A], error)

	// GetAuthHeader returns the authorization header name and value.
	// If nil or returns empty strings, no auth header is added.
	GetAuthHeader func(source *catalog.Source) (name, value string)

	// ValidateCredentials validates the API credentials before fetching.
	// If nil, no validation is performed.
	ValidateCredentials func(ctx context.Context, client *http.Client, baseURL string, source *catalog.Source) error

	// Logger for logging messages (optional).
	Logger Logger

	// UserAgent is the User-Agent header value.
	// Defaults to "model-registry-catalog" if empty.
	UserAgent string
}

// Logger is an interface for logging.
type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Warningf(format string, args ...any)
}

type noopLogger struct{}

func (noopLogger) Infof(format string, args ...any)    {}
func (noopLogger) Errorf(format string, args ...any)   {}
func (noopLogger) Warningf(format string, args ...any) {}

// Provider is an HTTP-based data provider with periodic polling.
type Provider[E any, A any] struct {
	config       Config[E, A]
	client       *http.Client
	baseURL      string
	syncInterval time.Duration
	source       *catalog.Source
	filter       *catalog.ItemFilter
	logger       Logger
}

// NewProvider creates a new HTTP provider with the given configuration.
func NewProvider[E any, A any](config Config[E, A], source *catalog.Source) (*Provider[E, A], error) {
	// Parse base URL
	baseURLKey := config.BaseURLKey
	if baseURLKey == "" {
		baseURLKey = "url"
	}

	baseURL := config.DefaultBaseURL
	if url, ok := source.Properties[baseURLKey].(string); ok && url != "" {
		baseURL = url
	}
	if baseURL == "" {
		return nil, fmt.Errorf("missing base URL (property %s or default)", baseURLKey)
	}

	// Parse sync interval
	syncIntervalKey := config.SyncIntervalKey
	if syncIntervalKey == "" {
		syncIntervalKey = "syncInterval"
	}

	syncInterval := config.DefaultSyncInterval
	if syncInterval == 0 {
		syncInterval = 24 * time.Hour
	}
	if intervalStr, ok := source.Properties[syncIntervalKey].(string); ok && intervalStr != "" {
		if parsed, err := time.ParseDuration(intervalStr); err == nil {
			syncInterval = parsed
		}
	}

	// HTTP client
	client := config.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	// Build filter from source configuration
	filter, err := catalog.NewItemFilterFromSource(source, nil, nil)
	if err != nil {
		return nil, err
	}

	logger := config.Logger
	if logger == nil {
		logger = noopLogger{}
	}

	return &Provider[E, A]{
		config:       config,
		client:       client,
		baseURL:      baseURL,
		syncInterval: syncInterval,
		source:       source,
		filter:       filter,
		logger:       logger,
	}, nil
}

// Records starts fetching data and returns a channel of records.
// The channel is closed when the context is canceled.
// The provider polls periodically based on syncInterval.
func (p *Provider[E, A]) Records(ctx context.Context) (<-chan catalog.Record[E, A], error) {
	// Validate credentials if configured
	if p.config.ValidateCredentials != nil {
		if err := p.config.ValidateCredentials(ctx, p.client, p.baseURL, p.source); err != nil {
			return nil, fmt.Errorf("credential validation failed: %w", err)
		}
	}

	// Fetch initial data to catch errors early
	records, err := p.fetch(ctx)
	if err != nil {
		return nil, err
	}

	ch := make(chan catalog.Record[E, A])
	go func() {
		defer close(ch)

		// Send initial records
		p.emit(ctx, records, ch)

		// Set up periodic polling
		ticker := time.NewTicker(p.syncInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				p.logger.Infof("Periodic sync: fetching records for source %s", p.source.ID)
				records, err := p.fetch(ctx)
				if err != nil {
					p.logger.Errorf("Failed to fetch records: %v", err)
					continue
				}
				p.emit(ctx, records, ch)
			}
		}
	}()

	return ch, nil
}

func (p *Provider[E, A]) fetch(ctx context.Context) ([]catalog.Record[E, A], error) {
	if p.config.FetchRecords == nil {
		return nil, fmt.Errorf("FetchRecords function not configured")
	}
	return p.config.FetchRecords(ctx, p.client, p.baseURL, p.source)
}

func (p *Provider[E, A]) emit(ctx context.Context, records []catalog.Record[E, A], out chan<- catalog.Record[E, A]) {
	done := ctx.Done()
	for _, record := range records {
		select {
		case out <- record:
		case <-done:
			return
		}
	}

	// Send an empty record to indicate batch completion
	var zero catalog.Record[E, A]
	select {
	case out <- zero:
	case <-done:
	}
}

// NewProviderFunc creates a ProviderFunc that can be registered with a ProviderRegistry.
func NewProviderFunc[E any, A any](config Config[E, A]) catalog.ProviderFunc[E, A] {
	return func(ctx context.Context, source *catalog.Source, reldir string) (<-chan catalog.Record[E, A], error) {
		provider, err := NewProvider(config, source)
		if err != nil {
			return nil, err
		}
		return provider.Records(ctx)
	}
}

// Request represents an HTTP request configuration.
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    io.Reader
}

// DoRequest performs an HTTP request with standard error handling.
func DoRequest[T any](ctx context.Context, client *http.Client, req Request) (*T, error) {
	method := req.Method
	if method == "" {
		method = "GET"
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, req.URL, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// DoRequestRaw performs an HTTP request and returns the raw response body.
func DoRequestRaw(ctx context.Context, client *http.Client, req Request) ([]byte, error) {
	method := req.Method
	if method == "" {
		method = "GET"
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, req.URL, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}

// PaginatedFetcher handles paginated API responses.
type PaginatedFetcher[T any] struct {
	Client    *http.Client
	Headers   map[string]string
	BuildURL  func(cursor string) string
	ParseNext func(response *http.Response, items []T) (cursor string, hasMore bool)
	MaxItems  int
}

// FetchAll fetches all items from a paginated API.
func (f *PaginatedFetcher[T]) FetchAll(ctx context.Context) ([]T, error) {
	var allItems []T
	cursor := ""

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if f.MaxItems > 0 && len(allItems) >= f.MaxItems {
			break
		}

		url := f.BuildURL(cursor)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		for key, value := range f.Headers {
			req.Header.Set(key, value)
		}

		resp, err := f.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
		}

		var items []T
		if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		nextCursor, hasMore := f.ParseNext(resp, items)
		_ = resp.Body.Close()

		allItems = append(allItems, items...)

		if !hasMore || nextCursor == "" {
			break
		}
		cursor = nextCursor
	}

	return allItems, nil
}

// RateLimiter provides simple rate limiting for API calls.
type RateLimiter struct {
	interval time.Duration
	lastCall time.Time
}

// NewRateLimiter creates a rate limiter with the given minimum interval between calls.
func NewRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{interval: interval}
}

// Wait waits until the next API call is allowed.
func (r *RateLimiter) Wait(ctx context.Context) error {
	now := time.Now()
	elapsed := now.Sub(r.lastCall)
	if elapsed < r.interval {
		select {
		case <-time.After(r.interval - elapsed):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	r.lastCall = time.Now()
	return nil
}
