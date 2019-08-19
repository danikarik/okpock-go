package api

import (
	"encoding/json"
	"errors"
	"time"
)

// NewOrganization returns a new instance of organization.
func NewOrganization(title, desc string, data map[string]interface{}) (*Organization, error) {
	return &Organization{
		Title:       title,
		Description: desc,
		MetaData:    data,
	}, nil
}

// Organization holds company information.
type Organization struct {
	ID     int64 `json:"id" db:"id"`
	UserID int64 `json:"-" db:"user_id"`

	Title       string  `json:"title" db:"title"`
	Description string  `json:"description" db:"description"`
	MetaData    JSONMap `json:"metaData" db:"raw_metadata"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// IsValid checks whether input is valid or not.
func (o *Organization) IsValid() error {
	if o.Title == "" {
		return errors.New("title is empty")
	}
	if o.Description == "" {
		return errors.New("description is empty")
	}
	return nil
}

// String returns string representation of struct.
func (o *Organization) String() string {
	data, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(data)
}
