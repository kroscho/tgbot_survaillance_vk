// Package stommer converts golang struct to map and allows to get fields and values
package stommer

import (
	"database/sql/driver"
	"sort"

	"github.com/elgris/stom"
	"github.com/pkg/errors"
)

type Stommer struct {
	Columns []string
	Values  []interface{}
}

func New(o interface{}, omittedFields ...string) (*Stommer, error) {
	fields, err := stom.MustNewStom(o).ToMap(o)
	if err != nil {
		return nil, err
	}

	var (
		columns     []string
		fieldValues = map[string]interface{}{}
	)
	for name, value := range fields {
		if isOmitted(name, omittedFields...) {
			continue
		}

		tmpVal := value
		if valuer, ok := value.(driver.Valuer); ok {
			tmpVal, err = valuer.Value()
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}

		fieldValues[name] = tmpVal
		columns = append(columns, name)
	}

	sort.Strings(columns)
	var values []interface{}
	for _, columnName := range columns {
		values = append(values, fieldValues[columnName])
	}

	return &Stommer{
		Columns: columns,
		Values:  values,
	}, nil
}

func isOmitted(key string, omitted ...string) bool {
	for _, omitted := range omitted {
		if key == omitted {
			return true
		}
	}

	return false
}
