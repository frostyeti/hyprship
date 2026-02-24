package core

import (
	"archive/zip"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
)

var ErrUnsupportedFormat = errors.New("unsupported import/export format")

// CollectionDef defines how a specific collection is exported/imported.
type CollectionDef[T any] struct {
	Name      string // e.g. "users", "roles"
	Headers   []string
	ToRow     func(T) []string
	FromRow   func([]string) (T, error)
	GetItems  func(ctx context.Context, opts ListOptions) ([]T, error)
	SaveItems func(ctx context.Context, items []T) ([]ImportError, error)
}

// BatchExport handles exporting multiple collections to a single output format.
func BatchExport(ctx context.Context, w io.Writer, format ExportFormat, opts ListOptions, collections ...any) error {
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

func exportJSON(ctx context.Context, w io.Writer, opts ListOptions, collections ...any) error {
	result := make(map[string]any)
	for _, c := range collections {
		v := reflect.ValueOf(c)
		name := v.FieldByName("Name").String()
		getItems := v.FieldByName("GetItems")
		res := getItems.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(opts)})
		if !res[1].IsNil() {
			return res[1].Interface().(error)
		}
		result[name] = res[0].Interface()
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func exportYAML(ctx context.Context, w io.Writer, opts ListOptions, collections ...any) error {
	result := make(map[string]any)
	for _, c := range collections {
		v := reflect.ValueOf(c)
		name := v.FieldByName("Name").String()
		getItems := v.FieldByName("GetItems")
		res := getItems.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(opts)})
		if !res[1].IsNil() {
			return res[1].Interface().(error)
		}
		result[name] = res[0].Interface()
	}

	return yaml.NewEncoder(w).Encode(result)
}

func exportCSV(ctx context.Context, w io.Writer, opts ListOptions, collections ...any) error {
	if len(collections) == 1 {
		// Single file
		return writeCsvCollection(ctx, w, opts, collections[0])
	}

	// Multiple collections require a zip file containing multiple CSVs
	zw := zip.NewWriter(w)
	for _, c := range collections {
		v := reflect.ValueOf(c)
		name := v.FieldByName("Name").String()
		f, err := zw.Create(name + ".csv")
		if err != nil {
			return err
		}
		if err := writeCsvCollection(ctx, f, opts, c); err != nil {
			return err
		}
	}
	return zw.Close()
}

func writeCsvCollection(ctx context.Context, w io.Writer, opts ListOptions, collection any) error {
	cw := csv.NewWriter(w)
	v := reflect.ValueOf(collection)
	headers := v.FieldByName("Headers").Interface().([]string)
	if err := cw.Write(headers); err != nil {
		return err
	}

	getItems := v.FieldByName("GetItems")
	res := getItems.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(opts)})
	if !res[1].IsNil() {
		return res[1].Interface().(error)
	}

	items := res[0]
	toRow := v.FieldByName("ToRow")
	for i := 0; i < items.Len(); i++ {
		rowRes := toRow.Call([]reflect.Value{items.Index(i)})
		row := rowRes[0].Interface().([]string)
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func exportExcel(ctx context.Context, w io.Writer, opts ListOptions, collections ...any) error {
	f := excelize.NewFile()
	first := true
	for _, c := range collections {
		v := reflect.ValueOf(c)
		name := v.FieldByName("Name").String()
		if first {
			f.SetSheetName("Sheet1", name)
			first = false
		} else {
			f.NewSheet(name)
		}

		headers := v.FieldByName("Headers").Interface().([]string)
		headerRow := make([]interface{}, len(headers))
		for i, h := range headers {
			headerRow[i] = h
		}
		f.SetSheetRow(name, "A1", &headerRow)

		getItems := v.FieldByName("GetItems")
		res := getItems.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(opts)})
		if !res[1].IsNil() {
			return res[1].Interface().(error)
		}

		items := res[0]
		toRow := v.FieldByName("ToRow")
		for i := 0; i < items.Len(); i++ {
			rowRes := toRow.Call([]reflect.Value{items.Index(i)})
			strRow := rowRes[0].Interface().([]string)
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

// Implement Import similarly for batch imports...
