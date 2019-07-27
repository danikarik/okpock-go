package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// M is an alias for map.
type M map[string]interface{}

// IsValid checks whether input is valid or not.
func (m *M) IsValid() error {
	return nil
}

// String returns string representation of struct.
func (m *M) String() string {
	data, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(data)
}

// Validator is common interface used for input validation.
type Validator interface {
	String() string
	IsValid() error
}

// SyntaxError holds error and source data.
type SyntaxError struct {
	*json.SyntaxError
	input []byte
}

func (e SyntaxError) Error() string {
	if e.input == nil || len(e.input) == 0 {
		return "empty json body"
	}
	return fmt.Sprintf("syntax error near: `%s`", string(e.input[e.Offset-1:]))
}

func readJSON(r *http.Request, v Validator) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.Unmarshal(data, &v)
	if err != nil {
		if se, ok := err.(*json.SyntaxError); ok {
			return SyntaxError{se, data}
		}
		return err
	}

	err = v.IsValid()
	if err != nil {
		return err
	}

	return err
}

func sendJSON(w http.ResponseWriter, code int, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("json: encoding json response: %v", err)
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("json: writing json response: %v", err)
	}

	return nil
}
