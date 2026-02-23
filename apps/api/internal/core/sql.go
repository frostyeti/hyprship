package core

import (
	"strconv"
	"strings"
)

// Rebind changes '?' placeholders in a query to driver-specific placeholders.
// e.g. for postgres: '$1', '$2', ...
// e.g. for sqlserver: '@p1', '@p2', ...
func Rebind(driverName string, query string) string {
	switch driverName {
	case "postgres", "pg":
		return rebindDollar(query)
	case "mssql", "sqlserver":
		return rebindAtP(query)
	default: // sqlite, mysql use ?
		return query
	}
}

func rebindDollar(query string) string {
	var sb strings.Builder
	idx := 1
	for _, char := range query {
		if char == '?' {
			sb.WriteRune('$')
			sb.WriteString(strconv.Itoa(idx))
			idx++
		} else {
			sb.WriteRune(char)
		}
	}
	return sb.String()
}

func rebindAtP(query string) string {
	var sb strings.Builder
	idx := 1
	for _, char := range query {
		if char == '?' {
			sb.WriteString("@p")
			sb.WriteString(strconv.Itoa(idx))
			idx++
		} else {
			sb.WriteRune(char)
		}
	}
	return sb.String()
}
