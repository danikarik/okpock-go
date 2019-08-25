package api

import (
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
func NewProject(title, name, desc string, passType PassType) *Project {
	return &Project{
		ID:               uuid.NewV4().String(),
		Title:            title,
		OrganizationName: name,
		Description:      desc,
		PassType:         passType,
	}
}

// Project holds project structure related fields.
type Project struct {
	ID string `json:"id" db:"id"`

	Title            string   `json:"title" db:"title"`
	OrganizationName string   `json:"organizationName" db:"organization_name"`
	Description      string   `json:"description" db:"description"`
	PassType         PassType `json:"passType" db:"pass_type"`

	BackgroundImage string `json:"backgroundImage" db:"background_image"`
	FooterImage     string `json:"footerImage" db:"footer_image"`
	IconImage       string `json:"iconImage" db:"icon_image"`
	StripImage      string `json:"stripImage" db:"strip_image"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// IsValid checks whether input is valid or not.
func (p *Project) IsValid() error {
	if p.Title == "" {
		return errors.New("title is empty")
	}
	if p.OrganizationName == "" {
		return errors.New("organization name is empty")
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
