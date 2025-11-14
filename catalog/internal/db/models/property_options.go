package models

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

type PropertyOptionType int

const (
	ContextPropertyOptionType PropertyOptionType = iota
	ArtifactPropertyOptionType
)

type PropertyOption struct {
	TypeID           int32    `json:"type_id"`
	Name             string   `json:"name"`
	IsCustomProperty bool     `json:"is_custom_property"`
	StringValue      []string `json:"string_value"`
	ArrayValue       []string `json:"array_value"`
	MinDoubleValue   *float64 `json:"min_double_value"`
	MaxDoubleValue   *float64 `json:"max_double_value"`
	MinIntValue      *int64   `json:"min_int_value"`
	MaxIntValue      *int64   `json:"max_int_value"`
}

const (
	StringValueField = "string_value"
	DoubleValueField = "double_value"
	IntValueField    = "int_value"
	ArrayValueField  = "array_value"
)

// ValueField returns string_value, array_value, double_value or int_value
// depending on which fields are non-nil.
func (po *PropertyOption) ValueField() string {
	switch {
	case po.MinDoubleValue != nil || po.MaxDoubleValue != nil:
		return DoubleValueField
	case po.MinIntValue != nil || po.MaxIntValue != nil:
		return IntValueField
	case len(po.ArrayValue) > 0:
		return ArrayValueField
	}

	return StringValueField
}

// FullName returns the complete name of the property to pass to filterQuery. Prefix is optional.
func (po *PropertyOption) FullName(prefix string) string {
	parts := make([]string, 0, 3)

	if prefix != "" {
		parts = append(parts, prefix)
	}

	parts = append(parts, po.Name)

	if po.IsCustomProperty {
		parts = append(parts, po.ValueField())
	}

	return strings.Join(parts, ".")
}

type PropertyOptionsRepository interface {
	// Refresh rebuilds the materialized view.
	Refresh(t PropertyOptionType) error
	// List returns all the options for a type. If typeID is 0, all types are returned.
	List(t PropertyOptionType, typeID int32) ([]PropertyOption, error)
}

// PropertyOptionsRefresher refreshes the materialized views after a short
// delay to prevent unnecessary duplicate refreshes.
type PropertyOptionsRefresher struct {
	ticker *time.Ticker
	delay  time.Duration
	mu     sync.Mutex
}

func NewPropertyOptionsRefresher(ctx context.Context, repo PropertyOptionsRepository, delay time.Duration) *PropertyOptionsRefresher {
	ticker := time.NewTicker(time.Hour)
	ticker.Stop()

	r := &PropertyOptionsRefresher{
		ticker: ticker,
		delay:  delay,
	}
	go r.bg(ctx, repo)
	return r
}

func (r *PropertyOptionsRefresher) Trigger() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ticker.Reset(r.delay)
}

func (r *PropertyOptionsRefresher) bg(ctx context.Context, repo PropertyOptionsRepository) {
	done := ctx.Done()
	for {
		select {
		case <-done:
			return
		case <-r.ticker.C:
			// Fallthrough
		}

		r.mu.Lock()
		r.ticker.Stop()
		r.mu.Unlock()

		err := repo.Refresh(ContextPropertyOptionType)
		if err != nil {
			glog.Warningf("Failed to refresh context property options: %v", err)
		}

		err = repo.Refresh(ArtifactPropertyOptionType)
		if err != nil {
			glog.Warningf("Failed to refresh artifact property options: %v", err)
		}
	}
}
