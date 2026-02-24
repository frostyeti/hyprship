package core

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
)

var ErrUnsupportedFormat = errors.New("unsupported import/export format")

// Collection is an interface that allows batch operations across different types
type Collection interface {
	Name() string
	Headers() []string
	Export(ctx context.Context, opts ListOptions) (any, error)
	ExportRows(ctx context.Context, opts ListOptions) ([][]string, error)
	Import(ctx context.Context, data any) ([]ImportError, error)
	ImportRows(ctx context.Context, rows [][]string) ([]ImportError, error)
}

// CollectionDef defines how a specific collection is exported/imported.
type CollectionDef[T any] struct {
	NameStr     string // e.g. "users", "roles"
	HeaderNames []string
	ToRow       func(T) []string
	FromRow     func([]string) (T, error)
	GetItems    func(ctx context.Context, opts ListOptions) ([]T, error)
	SaveItems   func(ctx context.Context, items []T) ([]ImportError, error)
}

func (c CollectionDef[T]) Name() string { return c.NameStr }

func (c CollectionDef[T]) Headers() []string { return c.HeaderNames }

func (c CollectionDef[T]) Export(ctx context.Context, opts ListOptions) (any, error) {
	return c.GetItems(ctx, opts)
}

func (c CollectionDef[T]) ExportRows(ctx context.Context, opts ListOptions) ([][]string, error) {
	items, err := c.GetItems(ctx, opts)
	if err != nil {
		return nil, err
	}
	rows := make([][]string, len(items))
	for i, item := range items {
		rows[i] = c.ToRow(item)
	}
	return rows, nil
}

func (c CollectionDef[T]) Import(ctx context.Context, data any) ([]ImportError, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var items []T
	if err := json.Unmarshal(b, &items); err != nil {
		return nil, err
	}
	return c.SaveItems(ctx, items)
}

func (c CollectionDef[T]) ImportRows(ctx context.Context, rows [][]string) ([]ImportError, error) {
	var items []T
	var errs []ImportError

	for i, row := range rows {
		item, err := c.FromRow(row)
		if err != nil {
			errs = append(errs, ImportError{RowNumber: i + 2, Error: err}) // +2 for header offset
			continue
		}
		items = append(items, item)
	}

	if len(items) > 0 {
		saveErrs, err := c.SaveItems(ctx, items)
		if err != nil {
			return errs, err
		}
		errs = append(errs, saveErrs...)
	}

	return errs, nil
}

// BatchExport handles exporting multiple collections to a single output format.
func BatchExport(ctx context.Context, w io.Writer, format ExportFormat, opts ListOptions, collections ...Collection) error {
	switch format {
	case FormatJSON:
		return exportJSON(ctx, w, opts, collections...)
	case FormatYAML:
		return exportYAML(ctx, w, opts, collections...)
	case FormatCSV:
		return exportCSV(ctx, w, opts, collections...)
	case FormatExcel:
		return exportExcel(ctx, w, opts, collections...)
	default:
		return ErrUnsupportedFormat
	}
}

func exportJSON(ctx context.Context, w io.Writer, opts ListOptions, collections ...Collection) error {
	result := make(map[string]any)
	for _, c := range collections {
		items, err := c.Export(ctx, opts)
		if err != nil {
			return err
		}
		result[c.Name()] = items
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func exportYAML(ctx context.Context, w io.Writer, opts ListOptions, collections ...Collection) error {
	result := make(map[string]any)
	for _, c := range collections {
		items, err := c.Export(ctx, opts)
		if err != nil {
			return err
		}
		result[c.Name()] = items
	}

	return yaml.NewEncoder(w).Encode(result)
}

func exportCSV(ctx context.Context, w io.Writer, opts ListOptions, collections ...Collection) error {
	if len(collections) == 1 {
		return writeCsvCollection(ctx, w, opts, collections[0])
	}

	zw := zip.NewWriter(w)
	for _, c := range collections {
		f, err := zw.Create(c.Name() + ".csv")
		if err != nil {
			return err
		}
		if err := writeCsvCollection(ctx, f, opts, c); err != nil {
			return err
		}
	}
	return zw.Close()
}

func writeCsvCollection(ctx context.Context, w io.Writer, opts ListOptions, collection Collection) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(collection.Headers()); err != nil {
		return err
	}

	rows, err := collection.ExportRows(ctx, opts)
	if err != nil {
		return err
	}

	for _, row := range rows {
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func exportExcel(ctx context.Context, w io.Writer, opts ListOptions, collections ...Collection) error {
	f := excelize.NewFile()
	first := true
	for _, c := range collections {
		name := c.Name()
		if first {
			f.SetSheetName("Sheet1", name)
			first = false
		} else {
			f.NewSheet(name)
		}

		headers := c.Headers()
		headerRow := make([]interface{}, len(headers))
		for i, h := range headers {
			headerRow[i] = h
		}
		f.SetSheetRow(name, "A1", &headerRow)

		rows, err := c.ExportRows(ctx, opts)
		if err != nil {
			return err
		}

		for i, strRow := range rows {
			row := make([]interface{}, len(strRow))
			for j, val := range strRow {
				row[j] = val
			}
			f.SetSheetRow(name, fmt.Sprintf("A%d", i+2), &row)
		}
	}
	_, err := f.WriteTo(w)
	return err
}

// BatchImport handles importing multiple collections from a single input format.
func BatchImport(ctx context.Context, r io.Reader, format ExportFormat, collections ...Collection) ([]ImportError, error) {
	switch format {
	case FormatJSON:
		return importJSON(ctx, r, collections...)
	case FormatYAML:
		return importYAML(ctx, r, collections...)
	case FormatCSV:
		return importCSV(ctx, r, collections...)
	case FormatExcel:
		return importExcel(ctx, r, collections...)
	default:
		return nil, ErrUnsupportedFormat
	}
}

func importJSON(ctx context.Context, r io.Reader, collections ...Collection) ([]ImportError, error) {
	var payload map[string]any
	if err := json.NewDecoder(r).Decode(&payload); err != nil {
		return nil, err
	}

	var allErrs []ImportError
	for _, c := range collections {
		if data, ok := payload[c.Name()]; ok {
			errs, err := c.Import(ctx, data)
			if err != nil {
				return allErrs, err
			}
			allErrs = append(allErrs, errs...)
		}
	}
	return allErrs, nil
}

func importYAML(ctx context.Context, r io.Reader, collections ...Collection) ([]ImportError, error) {
	var payload map[string]any
	if err := yaml.NewDecoder(r).Decode(&payload); err != nil {
		return nil, err
	}

	var allErrs []ImportError
	for _, c := range collections {
		if data, ok := payload[c.Name()]; ok {
			errs, err := c.Import(ctx, data)
			if err != nil {
				return allErrs, err
			}
			allErrs = append(allErrs, errs...)
		}
	}
	return allErrs, nil
}

func importCSV(ctx context.Context, r io.Reader, collections ...Collection) ([]ImportError, error) {
	if len(collections) == 0 {
		return nil, errors.New("no collections provided for import")
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err == nil {
		var allErrs []ImportError
		for _, c := range collections {
			for _, file := range zipReader.File {
				if file.Name == c.Name()+".csv" {
					f, err := file.Open()
					if err != nil {
						return allErrs, err
					}
					cr := csv.NewReader(f)
					records, err := cr.ReadAll()
					f.Close()
					if err != nil {
						return allErrs, err
					}
					if len(records) > 1 {
						errs, err := c.ImportRows(ctx, records[1:])
						if err != nil {
							return allErrs, err
						}
						allErrs = append(allErrs, errs...)
					}
				}
			}
		}
		return allErrs, nil
	}

	c := collections[0]
	cr := csv.NewReader(bytes.NewReader(data))
	records, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) > 1 {
		return c.ImportRows(ctx, records[1:])
	}
	return nil, nil
}

func importExcel(ctx context.Context, r io.Reader, collections ...Collection) ([]ImportError, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var allErrs []ImportError
	for _, c := range collections {
		rows, err := f.GetRows(c.Name())
		if err != nil {
			if len(collections) == 1 {
				rows, err = f.GetRows(f.GetSheetName(0))
				if err != nil {
					continue
				}
			} else {
				continue
			}
		}

		if len(rows) > 1 {
			errs, err := c.ImportRows(ctx, rows[1:])
			if err != nil {
				return allErrs, err
			}
			allErrs = append(allErrs, errs...)
		}
	}
	return allErrs, nil
}
