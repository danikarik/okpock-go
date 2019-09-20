package api

import (
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

// ImageSize is an alias for image size.
type ImageSize string

const (
	// ImageSize1x original image size.
	ImageSize1x = ImageSize("1x")
	// ImageSize2x retina image size.
	ImageSize2x = ImageSize("2x")
	// ImageSize3x super retina image size.
	ImageSize3x = ImageSize("3x")
)

// NewProject returns a new instance of project.
func NewProject(title, name, desc string, passType PassType) *Project {
	return &Project{
		Title:            title,
		OrganizationName: name,
		Description:      desc,
		PassType:         passType,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

// Project holds project structure related fields.
type Project struct {
	ID int64 `json:"id" db:"id"`

	Title            string   `json:"title" db:"title"`
	OrganizationName string   `json:"organizationName" db:"organization_name"`
	Description      string   `json:"description" db:"description"`
	PassType         PassType `json:"passType" db:"pass_type"`

	BackgroundImage   string `json:"backgroundImage" db:"background_image"`
	BackgroundImage2x string `json:"backgroundImage2x" db:"background_image_2x"`
	BackgroundImage3x string `json:"backgroundImage3x" db:"background_image_3x"`

	FooterImage   string `json:"footerImage" db:"footer_image"`
	FooterImage2x string `json:"footerImage2x" db:"footer_image_2x"`
	FooterImage3x string `json:"footerImage3x" db:"footer_image_3x"`

	IconImage   string `json:"iconImage" db:"icon_image"`
	IconImage2x string `json:"iconImage2x" db:"icon_image_2x"`
	IconImage3x string `json:"iconImage3x" db:"icon_image_3x"`

	LogoImage   string `json:"logoImage" db:"logo_image"`
	LogoImage2x string `json:"logoImage2x" db:"logo_image_2x"`
	LogoImage3x string `json:"logoImage3x" db:"logo_image_3x"`

	StripImage   string `json:"stripImage" db:"strip_image"`
	StripImage2x string `json:"stripImage2x" db:"strip_image_2x"`
	StripImage3x string `json:"stripImage3x" db:"strip_image_3x"`

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
