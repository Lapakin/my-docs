package rest

import (
	f "github.com/lapotkin/file-storage/internal/app/filter"
)

type Query struct {
	Param        string
	ValidateFunc func(string) error
	Required     bool
}

type Queries []*Query

func CreateFiltersFromQuery(urlParams map[string][]string, queries Queries) (f.Filters, error) {
	filters := make(f.Filters)
	for _, q := range queries {
		values, exists := urlParams[q.Param]
		if !exists || len(values) == 0 || values[0] == "" {
			if q.Required {
				return nil, ErrEmptyValue
			}
			continue
		}
		v := values[0]
		if q.ValidateFunc != nil {
			if err := q.ValidateFunc(v); err != nil {
				return nil, ErrParseQuery
			}
		}
		filters[q.Param] = v
	}
	return filters, nil
}
