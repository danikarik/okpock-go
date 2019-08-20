package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

// PassType refers to `Style Keys` type.
type PassType string

const (
	// BoardingPass holds information specific to a boarding pass.
	BoardingPass = PassType("boardingPass")
	// Coupon holds information specific to a coupon.
	Coupon = PassType("coupon")
	// EventTicket holds information specific to an event ticket.
	EventTicket = PassType("eventTicket")
	// Generic holds information specific to a generic pass.
	Generic = PassType("generic")
	// StoreCard holds information specific to a store card.
	StoreCard = PassType("storeCard")
)

// ProjectField is an alias for project field.
type ProjectField int

const (
	_ ProjectField = iota
	// BackgroundImage refers to `project.BackgroundImage`.
	BackgroundImage
	// FooterImage refers to `project.FooterImage`.
	FooterImage
	// IconImage refers to `project.IconImage`.
	IconImage
	// StripImage refers to `project.StripImage`.
	StripImage
)

// NewProject returns a new instance of project.
func NewProject(desc string, orgID int64, passType PassType) (*Project, error) {
	return &Project{
		Description:    desc,
		OrganizationID: orgID,
		PassType:       passType,
	}, nil
}

// Project holds project structure related fields.
type Project struct {
	ID             int64 `json:"id" db:"id"`
	OrganizationID int64 `json:"-" db:"organization_id"`

	Description string   `json:"description" db:"description"`
	PassType    PassType `json:"passType" db:"pass_type"`

	BackgroundImage sql.NullString `json:"backgroundImage" db:"background_image"`
	FooterImage     sql.NullString `json:"footerImage" db:"footer_image"`
	IconImage       sql.NullString `json:"iconImage" db:"icon_image"`
	StripImage      sql.NullString `json:"stripImage" db:"strip_image"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// IsValid checks whether input is valid or not.
func (p *Project) IsValid() error {
	if p.OrganizationID == 0 {
		return errors.New("organization id is empty")
	}
	if p.Description == "" {
		return errors.New("description is empty")
	}
	if p.PassType == "" {
		return errors.New("pass type is empty")
	}
	return nil
}

// String returns string representation of struct.
func (p *Project) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(data)
}

// SetField sets string value of project field.
func (p *Project) SetField(field ProjectField, value string) {
	switch field {
	case BackgroundImage:
		p.BackgroundImage.Valid = true
		p.BackgroundImage.String = value
		break
	case FooterImage:
		p.FooterImage.Valid = true
		p.FooterImage.String = value
		break
	case IconImage:
		p.IconImage.Valid = true
		p.IconImage.String = value
		break
	case StripImage:
		p.StripImage.Valid = true
		p.StripImage.String = value
		break
	}
}

// GetField gets string value of project field.
func (p *Project) GetField(field ProjectField) string {
	switch field {
	case BackgroundImage:
		if p.BackgroundImage.Valid {
			return p.BackgroundImage.String
		}
	case FooterImage:
		if p.FooterImage.Valid {
			return p.FooterImage.String
		}
	case IconImage:
		if p.IconImage.Valid {
			return p.IconImage.String
		}
	case StripImage:
		if p.StripImage.Valid {
			return p.StripImage.String
		}
	}
	return ""
}

// ProjectStore implements method for project logic.
type ProjectStore interface {
	// IsProjectExists ...
	// TODO: description
	IsProjectExists(ctx context.Context, desc string, orgID int64, passType PassType) (bool, error)

	// SaveNewProject ...
	// TODO: description
	SaveNewProject(ctx context.Context, proj *Project) error

	// LoadProject ...
	// TODO: description
	LoadProject(ctx context.Context, id int64) (*Project, error)

	// LoadProjects ...
	// TODO: description
	LoadProjects(ctx context.Context, userID int64) ([]*Project, error)

	// UpdateProjectDescription ...
	// TODO: description
	UpdateProjectDescription(ctx context.Context, desc string, proj *Project) error

	// SetBackgroundImage ...
	// TODO: description
	SetBackgroundImage(ctx context.Context, key string, proj *Project) error

	// SetFooterImage ...
	// TODO: description
	SetFooterImage(ctx context.Context, key string, proj *Project) error

	// SetIconImage ...
	// TODO: description
	SetIconImage(ctx context.Context, key string, proj *Project) error

	// SetStripImage ...
	// TODO: description
	SetStripImage(ctx context.Context, key string, proj *Project) error
}
