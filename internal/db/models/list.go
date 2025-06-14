package models

type ListWrapper[T any] struct {
	Items         []T
	NextPageToken string
	PageSize      int32
	Size          int32
}
