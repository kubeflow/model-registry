package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetOrderBy_NilReturnsDefault verifies that a nil OrderBy returns the default.
func TestGetOrderBy_NilReturnsDefault(t *testing.T) {
	p := &Pagination{}
	assert.Equal(t, DefaultOrderBy, p.GetOrderBy(), "nil OrderBy should return default")
}

// TestGetOrderBy_EmptyStringReturnsDefault verifies that a pointer to an empty
// string still returns the default OrderBy. This is the root cause of
// RHOAIENG-57798: the API layer sets OrderBy to a pointer to "" which bypasses
// the nil check and results in no ORDER BY clause, breaking cursor-based pagination.
func TestGetOrderBy_EmptyStringReturnsDefault(t *testing.T) {
	empty := ""
	p := &Pagination{OrderBy: &empty}
	assert.Equal(t, DefaultOrderBy, p.GetOrderBy(),
		"empty string OrderBy should return default 'id', not empty string")
}

// TestGetOrderBy_ExplicitValuePreserved verifies that an explicit value is returned as-is.
func TestGetOrderBy_ExplicitValuePreserved(t *testing.T) {
	val := "CREATE_TIME"
	p := &Pagination{OrderBy: &val}
	assert.Equal(t, "CREATE_TIME", p.GetOrderBy())
}

// TestGetSortOrder_NilReturnsDefault verifies that a nil SortOrder returns the default.
func TestGetSortOrder_NilReturnsDefault(t *testing.T) {
	p := &Pagination{}
	assert.Equal(t, DefaultSortOrder, p.GetSortOrder(), "nil SortOrder should return default")
}

// TestGetSortOrder_EmptyStringReturnsDefault verifies that a pointer to an empty
// string still returns the default SortOrder. Same root cause as RHOAIENG-57798.
func TestGetSortOrder_EmptyStringReturnsDefault(t *testing.T) {
	empty := ""
	p := &Pagination{SortOrder: &empty}
	assert.Equal(t, DefaultSortOrder, p.GetSortOrder(),
		"empty string SortOrder should return default 'ASC', not empty string")
}

// TestGetSortOrder_ExplicitValuePreserved verifies that an explicit value is returned as-is.
func TestGetSortOrder_ExplicitValuePreserved(t *testing.T) {
	val := "DESC"
	p := &Pagination{SortOrder: &val}
	assert.Equal(t, "DESC", p.GetSortOrder())
}
