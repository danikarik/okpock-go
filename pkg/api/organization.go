package api

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

// NewOrganization returns a new instance of organization.
func NewOrganization(userID, title, desc string, data map[string]interface{}) *Organization {
	return &Organization{
		ID:          uuid.NewV4().String(),
		UserID:      userID,
		Title:       title,
		Description: desc,
		MetaData:    data,
	}
}

// Organization holds company information.
type Organization struct {
	ID     string `json:"id" db:"id"`
	UserID string `json:"-" db:"user_id"`

	Title       string  `json:"title" db:"title"`
	Description string  `json:"description" db:"description"`
	MetaData    JSONMap `json:"metaData" db:"raw_metadata"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// IsValid checks whether input is valid or not.
func (o *Organization) IsValid() error {
	if o.UserID == "" {
		return errors.New("user id is empty")
	}
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
	IsOrganizationExists(ctx context.Context, userID, title string) (bool, error)

	// SaveNewOrganization ...
	// TODO: description
	SaveNewOrganization(ctx context.Context, org *Organization) error

	// LoadOrganization ...
	// TODO: description
	LoadOrganization(ctx context.Context, id string) (*Organization, error)

	// LoadOrganizations ...
	// TODO: description
	LoadOrganizations(ctx context.Context, userID string) ([]*Organization, error)

	// UpdateOrganizationDescription ...
	// TODO: description
	UpdateOrganizationDescription(ctx context.Context, desc string, org *Organization) error

	// UpdateOrganizationMetaData ...
	// TODO: description
	UpdateOrganizationMetaData(ctx context.Context, data map[string]interface{}, org *Organization) error
}
