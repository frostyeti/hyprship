package core

import (
	"fmt"
	"strings"
)

// BuildListQuery takes a base query like "SELECT * FROM users" and appends
// WHERE, ORDER BY, LIMIT, OFFSET clauses based on ListOptions.
// fieldMap allows translating API field names to database column names safely.
func BuildListQuery(dialect string, baseQuery string, opts ListOptions, fieldMap map[string]string) (string, []any) {
	var sb strings.Builder
	sb.WriteString(baseQuery)

	var args []any
	var whereClauses []string

	for _, filter := range opts.Filter {
		colName, ok := fieldMap[filter.Field]
		if !ok {
			continue // skip unknown fields for safety
		}

		op := "="
		switch strings.ToLower(filter.Operator) {
		case "eq":
			op = "="
		case "neq":
			op = "!="
		case "gt":
			op = ">"
		case "lt":
			op = "<"
		case "gte":
			op = ">="
		case "lte":
			op = "<="
		case "like":
			op = "LIKE"
		}

		whereClauses = append(whereClauses, fmt.Sprintf("%s %s ?", colName, op))
		args = append(args, filter.Value)
	}

	if len(whereClauses) > 0 {
		if strings.Contains(strings.ToUpper(baseQuery), " WHERE ") {
			sb.WriteString(" AND ")
		} else {
			sb.WriteString(" WHERE ")
		}
		sb.WriteString(strings.Join(whereClauses, " AND "))
	}

	hasOrder := false
	if len(opts.Sort) > 0 {
		sb.WriteString(" ORDER BY ")
		var sortClauses []string
		for _, s := range opts.Sort {
			colName, ok := fieldMap[s.Field]
			if !ok {
				continue
			}
			dir := "ASC"
			if s.Desc {
				dir = "DESC"
			}
			sortClauses = append(sortClauses, fmt.Sprintf("%s %s", colName, dir))
		}
		if len(sortClauses) > 0 {
			sb.WriteString(strings.Join(sortClauses, ", "))
			hasOrder = true
		} else {
			// default to id if mapped
			if idCol, ok := fieldMap["id"]; ok {
				sb.WriteString(fmt.Sprintf("%s ASC", idCol))
				hasOrder = true
			}
		}
	} else {
		if idCol, ok := fieldMap["id"]; ok {
			sb.WriteString(fmt.Sprintf(" ORDER BY %s ASC", idCol))
			hasOrder = true
		}
	}

	if opts.Limit > 0 {
		if dialect == "mssql" || dialect == "sqlserver" {
			if !hasOrder && !strings.Contains(strings.ToUpper(sb.String()), " ORDER BY ") {
				// Fallback sort for mssql offset support
				sb.WriteString(" ORDER BY (SELECT NULL)")
			}
			sb.WriteString(" OFFSET ? ROWS FETCH NEXT ? ROWS ONLY")
			args = append(args, opts.Offset, opts.Limit)
		} else {
			sb.WriteString(" LIMIT ? OFFSET ?")
			args = append(args, opts.Limit, opts.Offset)
		}
	}

	return sb.String(), args
}
