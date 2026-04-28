package filter

import (
	platformfilter "github.com/kubeflow/hub/internal/platform/db/filter"
)

// Value type constants
const (
	StringValueType = platformfilter.StringValueType
	DoubleValueType = platformfilter.DoubleValueType
	IntValueType    = platformfilter.IntValueType
	BoolValueType   = platformfilter.BoolValueType
	ArrayValueType  = platformfilter.ArrayValueType
)

// AST types
type WhereClause = platformfilter.WhereClause
type Expression = platformfilter.Expression
type OrExpression = platformfilter.OrExpression
type AndExpression = platformfilter.AndExpression
type Term = platformfilter.Term
type Comparison = platformfilter.Comparison
type PropertyRef = platformfilter.PropertyRef
type Value = platformfilter.Value
type ValueList = platformfilter.ValueList
type SingleValue = platformfilter.SingleValue

// Public API types
type FilterExpression = platformfilter.FilterExpression
type PropertyReference = platformfilter.PropertyReference

// Parse parses a filter query string and returns the root expression
var Parse = platformfilter.Parse
