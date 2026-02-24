package core

import (
	"context"
	"io"
)

// ListOptions provides common filtering, pagination, and sorting for repositories.
type ListOptions struct {
	Limit  int            `form:"limit"`
	Offset int            `form:"offset"`
	Sort   []SortOption   `form:"sort"`
	Filter []FilterOption `form:"filter"`
	Expand []string       `form:"expand"`
}

type SortOption struct {
	Field string
	Desc  bool
}

type FilterOption struct {
	Field    string
	Operator string // e.g., "eq", "neq", "like", "gt", "lt"
	Value    any
}

// ListResult wraps a slice of items with the total count.
type ListResult[T any] struct {
	Items      []T     `json:"items"`
	TotalCount int     `json:"totalCount"`
	NextPage   *string `json:"nextPage,omitempty"`
}

type ExportFormat string

const (
	FormatJSON  ExportFormat = "json"
	FormatCSV   ExportFormat = "csv"
	FormatExcel ExportFormat = "excel"
	FormatYAML  ExportFormat = "yaml"
)

type ImportError struct {
	RowNumber int
	Error     error
}

// Exporter is an interface for stores that support bulk export.
type Exporter[T any] interface {
	Export(ctx context.Context, w io.Writer, format ExportFormat, opts ListOptions) error
}

// Importer is an interface for stores that support bulk import.
type Importer[T any] interface {
	Import(ctx context.Context, r io.Reader, format ExportFormat) ([]ImportError, error)
}
