package api

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONMap is an alias for raw json.
type JSONMap map[string]interface{}

// Value is a value that drivers must be able to handle.
func (j JSONMap) Value() (driver.Value, error) {
	data, err := json.Marshal(j)
	if err != nil {
		return driver.Value(""), err
	}
	return driver.Value(string(data)), nil
}

// Scan value from database.
func (j JSONMap) Scan(src interface{}) error {
	var source []byte
	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	case sql.NullString:
		source = []byte("")
	default:
		return errors.New("invalid data type for JSONMap")
	}

	if len(source) == 0 {
		source = []byte("{}")
	}
	return json.Unmarshal(source, &j)
}
