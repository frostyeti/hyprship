package core

// ListOptions provides common filtering, pagination, and sorting for repositories.
type ListOptions struct {
	Page     int
	PageSize int
	Sort     string
	SortDesc bool
	Filter   string // Basic text filter
}

// ListResult wraps a slice of items with the total count.
type ListResult[T any] struct {
	Items []T
	Total int64
}

// Exporter is an interface for stores that support bulk export.
type Exporter[T any] interface {
	Export(opts ListOptions) ([]T, error)
}

// Importer is an interface for stores that support bulk import.
type Importer[T any] interface {
	Import(items []T) error
}
