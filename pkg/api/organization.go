package api

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

// NewOrganization returns a new instance of organization.
func NewOrganization(userID int64, title, desc string, data map[string]interface{}) (*Organization, error) {
	return &Organization{
		UserID:      userID,
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

// OrganizationStore implements method for organization logic.
type OrganizationStore interface {
	// IsOrganizationExists ...
	// TODO: description
	IsOrganizationExists(ctx context.Context, title string, userID int64) (bool, error)

	// SaveNewOrganization ...
	// TODO: description
	SaveNewOrganization(ctx context.Context, org *Organization) error

	// LoadOrganization ...
	// TODO: description
	LoadOrganization(ctx context.Context, id int64) (*Organization, error)

	// LoadOrganizations ...
	// TODO: description
	LoadOrganizations(ctx context.Context, userID int64) ([]*Organization, error)

	// UpdateOrganizationDescription ...
	// TODO: description
	UpdateOrganizationDescription(ctx context.Context, desc string, org *Organization) error

	// UpdateOrganizationMetaData ...
	// TODO: description
	UpdateOrganizationMetaData(ctx context.Context, data map[string]interface{}, org *Organization) error
}
