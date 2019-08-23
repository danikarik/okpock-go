package api

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
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

// NewProject returns a new instance of project.
func NewProject(orgID, desc string, passType PassType) *Project {
	return &Project{
		ID:             uuid.NewV4().String(),
		OrganizationID: orgID,
		Description:    desc,
		PassType:       passType,
	}
}

// Project holds project structure related fields.
type Project struct {
	ID             string `json:"id" db:"id"`
	OrganizationID string `json:"-" db:"organization_id"`

	Description string   `json:"description" db:"description"`
	PassType    PassType `json:"passType" db:"pass_type"`

	BackgroundImage string `json:"backgroundImage" db:"background_image"`
	FooterImage     string `json:"footerImage" db:"footer_image"`
	IconImage       string `json:"iconImage" db:"icon_image"`
	StripImage      string `json:"stripImage" db:"strip_image"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// IsValid checks whether input is valid or not.
func (p *Project) IsValid() error {
	if p.OrganizationID == "" {
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

// ProjectStore implements method for project logic.
type ProjectStore interface {
	// IsProjectExists ...
	// TODO: description
	IsProjectExists(ctx context.Context, orgID, desc string, passType PassType) (bool, error)

	// SaveNewProject ...
	// TODO: description
	SaveNewProject(ctx context.Context, proj *Project) error

	// LoadProject ...
	// TODO: description
	LoadProject(ctx context.Context, id string) (*Project, error)

	// LoadProjects ...
	// TODO: description
	LoadProjects(ctx context.Context, userID string) ([]*Project, error)

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
