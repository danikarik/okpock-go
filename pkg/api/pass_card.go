package api

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/danikarik/okpock/pkg/secure"
	uuid "github.com/satori/go.uuid"
)

const (
	// PKDataDetectorTypePhoneNumber refers to phone number.
	PKDataDetectorTypePhoneNumber string = "PKDataDetectorTypePhoneNumber"
	// PKDataDetectorTypeLink refers to link.
	PKDataDetectorTypeLink string = "PKDataDetectorTypeLink"
	// PKDataDetectorTypeAddress refers to address.
	PKDataDetectorTypeAddress string = "PKDataDetectorTypeAddress"
	// PKDataDetectorTypeCalendarEvent refers to calendar event.
	PKDataDetectorTypeCalendarEvent string = "PKDataDetectorTypeCalendarEvent"
)

const (
	// PKTextAlignmentLeft formats by left side.
	PKTextAlignmentLeft string = "PKTextAlignmentLeft"
	// PKTextAlignmentCenter formats into center.
	PKTextAlignmentCenter string = "PKTextAlignmentCenter"
	// PKTextAlignmentRight formats by right side.
	PKTextAlignmentRight string = "PKTextAlignmentRight"
	// PKTextAlignmentNatural formats by default.
	PKTextAlignmentNatural string = "PKTextAlignmentNatural"
)

const (
	// PKDateStyleNone formats into "no style".
	PKDateStyleNone string = "PKDateStyleNone"
	// PKDateStyleShort formats into "11/23/37" or "3:30 PM".
	PKDateStyleShort string = "PKDateStyleShort"
	// PKDateStyleMedium formats into "Nov 23, 1937" or "3:30:32 PM".
	PKDateStyleMedium string = "PKDateStyleMedium"
	// PKDateStyleLong formats into "November 23, 1937" or "3:30:32 PM PST".
	PKDateStyleLong string = "PKDateStyleLong"
	// PKDateStyleFull formats into "Tuesday, April 12, 1952 AD" or "3:30:42 PM Pacific Standard Time".
	PKDateStyleFull string = "PKDateStyleFull"
)

const (
	// PKNumberStyleDecimal formats into "1,234.568".
	PKNumberStyleDecimal string = "PKNumberStyleDecimal"
	// PKNumberStylePercent formats into "12%".
	PKNumberStylePercent string = "PKNumberStylePercent"
	// PKNumberStyleScientific formats into "1.2345678E3".
	PKNumberStyleScientific string = "PKNumberStyleScientific"
	// PKNumberStyleSpellOut formats into "one hundred twenty-three".
	PKNumberStyleSpellOut string = "PKNumberStyleSpellOut"
)

const (
	// PKTransitTypeAir refers to air.
	PKTransitTypeAir string = "PKTransitTypeAir"
	// PKTransitTypeBoat refers to boat.
	PKTransitTypeBoat string = "PKTransitTypeBoat"
	// PKTransitTypeBus refers to bus.
	PKTransitTypeBus string = "PKTransitTypeBus"
	// PKTransitTypeGeneric refers to generic.
	PKTransitTypeGeneric string = "PKTransitTypeGeneric"
	// PKTransitTypeTrain refers to train.
	PKTransitTypeTrain string = "PKTransitTypeTrain"
)

const (
	// PKBarcodeFormatQR refers to QR code.
	PKBarcodeFormatQR string = "PKBarcodeFormatQR"
	// PKBarcodeFormatPDF417 refers to PDF 417.
	PKBarcodeFormatPDF417 string = "PKBarcodeFormatPDF417"
	// PKBarcodeFormatAztec refers to Aztec.
	PKBarcodeFormatAztec string = "PKBarcodeFormatAztec"
	// PKBarcodeFormatCode128 refers to 128.
	PKBarcodeFormatCode128 string = "PKBarcodeFormatCode128"
)

// FieldType is an alias for `Field` type.
type FieldType string

const (
	// AuxiliaryFieldsType refers to `AuxiliaryFields` fields.
	AuxiliaryFieldsType = FieldType("auxiliaryFields")
	// BackFieldsType refers to `BackFields` fields.
	BackFieldsType = FieldType("backFields")
	// HeaderFieldsType refers to `HeaderFields` fields.
	HeaderFieldsType = FieldType("headerFields")
	// PrimaryFieldsType refers to `PrimaryFields` fields.
	PrimaryFieldsType = FieldType("primaryFields")
	// SecondaryFieldsType refers to `SecondaryFields` fields.
	SecondaryFieldsType = FieldType("primaryFields")
)

const w3cDate string = time.RFC3339

// Beacon refers to `Beacon Dictionary Keys`.
type Beacon struct {
	Major         uint16 `json:"major,omitempty"`
	Minor         uint16 `json:"minor,omitempty"`
	ProximityUUID string `json:"proximityUUID"`
	RelevantText  string `json:"relevantText,omitempty"`
}

// Location refers to `Location Dictionary Keys`.
type Location struct {
	Altitude     float64 `json:"altitude,omitempty"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	RelevantText string  `json:"relevantText,omitempty"`
}

// Barcode refers to `Barcode Dictionary Keys`.
type Barcode struct {
	AltText         string `json:"altText,omitempty"`
	Format          string `json:"format"`
	Message         string `json:"message"`
	MessageEncoding string `json:"messageEncoding"`
}

// NFC refers to `NFC Dictionary Keys`.
type NFC struct {
	Message             string `json:"message"`
	EncryptionPublicKey string `json:"encryptionPublicKey,omitempty"`
}

// Field refers to `Field Dictionary Keys`.
type Field struct {
	// Standard Field Dictionary Keys
	AttributedValue   string      `json:"attributedValue,omitempty"`
	ChangeMessage     string      `json:"changeMessage,omitempty"`
	DataDetectorTypes []string    `json:"dataDetectorTypes,omitempty"`
	Key               string      `json:"key"`
	Label             string      `json:"label,omitempty"`
	TextAlignment     string      `json:"textAlignment,omitempty"`
	Value             interface{} `json:"value"`

	// Date Style Keys
	DateStyle       string `json:"dateStyle,omitempty"`
	IgnoresTimeZone bool   `json:"ignoresTimeZone,omitempty"`
	IsRelative      bool   `json:"isRelative,omitempty"`
	TimeStyle       string `json:"timeStyle,omitempty"`

	// Number Style Keys
	CurrencyCode string `json:"currencyCode,omitempty"`
	NumberStyle  string `json:"numberStyle,omitempty"`
}

// PassStructure refers to `Pass Structure Dictionary Keys`.
type PassStructure struct {
	AuxiliaryFields []*Field `json:"auxiliaryFields,omitempty"`
	BackFields      []*Field `json:"backFields,omitempty"`
	HeaderFields    []*Field `json:"headerFields,omitempty"`
	PrimaryFields   []*Field `json:"primaryFields,omitempty"`
	SecondaryFields []*Field `json:"secondaryFields,omitempty"`
	TransitType     string   `json:"transitType,omitempty"`
}

// NewPassCardInfo returns a new instance of `PassCardInfo`.
func NewPassCardInfo(data *PassCard) *PassCardInfo {
	return &PassCardInfo{
		Data:      data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// PassCardInfo is a wrapper around `PassCard` with extra fields.
type PassCardInfo struct {
	ID        int64     `json:"id" db:"id"`
	Data      *PassCard `json:"data" db:"raw_data"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// IsValid checks whether input is valid or not.
func (p *PassCardInfo) IsValid() error {
	return p.Data.IsValid()
}

// String returns string representation of struct.
func (p *PassCardInfo) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(data)
}

// DefaultDataDetectorTypes is a default list for data detector.
func DefaultDataDetectorTypes() []string {
	return []string{
		PKDataDetectorTypePhoneNumber,
		PKDataDetectorTypeLink,
		PKDataDetectorTypeAddress,
		PKDataDetectorTypeCalendarEvent,
	}
}

// DefaultTextAlignment is a default test alignment.
func DefaultTextAlignment() string {
	return PKTextAlignmentNatural
}

// NewEmptyPassCard returns a new instance of `PassCard`.
func NewEmptyPassCard() *PassCard {
	return &PassCard{
		FormatVersion:       1,
		SerialNumber:        uuid.NewV4().String(),
		AuthenticationToken: secure.Token(),
	}
}

// PassCard refers `pass.json` structure.
type PassCard struct {
	// Standard Keys
	Description      string `json:"description"`
	FormatVersion    int64  `json:"formatVersion"`
	OrganizationName string `json:"organizationName"`
	PassTypeID       string `json:"passTypeIdentifier"`
	SerialNumber     string `json:"serialNumber"`
	TeamID           string `json:"teamIdentifier"`

	// Associated App Keys
	AppLaunchURL       string  `json:"appLaunchURL,omitempty"`
	AssociatedStoreIDs []int64 `json:"associatedStoreIdentifiers,omitempty"`

	// Companion App Keys
	UserInfo JSONMap `json:"userInfo,omitempty"`

	// Expiration Keys
	ExpirationDate string `json:"expirationDate,omitempty"`
	Voided         bool   `json:"voided,omitempty"`

	// Relevance Keys
	Beacons      []*Beacon   `json:"beacons,omitempty"`
	Locations    []*Location `json:"locations,omitempty"`
	MaxDistance  int64       `json:"maxDistance,omitempty"`
	RelevantDate string      `json:"relevantDate,omitempty"`

	// Style Keys
	BoardingPass *PassStructure `json:"boardingPass,omitempty"`
	Coupon       *PassStructure `json:"coupon,omitempty"`
	EventTicket  *PassStructure `json:"eventTicket,omitempty"`
	Generic      *PassStructure `json:"generic,omitempty"`
	StoreCard    *PassStructure `json:"storeCard,omitempty"`

	// Visual Appearance Keys
	Barcodes           []*Barcode `json:"barcodes,omitempty"`
	BackgroundColor    string     `json:"backgroundColor,omitempty"`
	ForegroundColor    string     `json:"foregroundColor,omitempty"`
	GroupingIdentifier string     `json:"groupingIdentifier,omitempty"`
	LabelColor         string     `json:"labelColor,omitempty"`
	LogoText           string     `json:"logoText,omitempty"`

	// Web Service Keys
	AuthenticationToken string `json:"authenticationToken"`
	WebServiceURL       string `json:"webServiceURL"`

	// NFC-Enabled Pass Keys
	NFC *NFC `json:"nfc,omitempty"`
}

// IsValid checks whether input is valid or not.
func (p *PassCard) IsValid() error {
	if p.Description == "" {
		return errors.New("standard: description is empty")
	}
	if p.FormatVersion != 1 {
		return errors.New("standard: format version must be 1")
	}
	if p.OrganizationName == "" {
		return errors.New("standard: organization name is empty")
	}
	if p.PassTypeID == "" {
		return errors.New("standard: pass type id is empty")
	}
	if p.SerialNumber == "" {
		return errors.New("standard: serial number is empty")
	}
	if p.TeamID == "" {
		return errors.New("standard: team id is empty")
	}
	if p.ExpirationDate != "" {
		if _, err := time.Parse(w3cDate, p.ExpirationDate); err != nil {
			return errors.New("expiration: date has invalid format")
		}
	}
	if p.RelevantDate != "" {
		if _, err := time.Parse(w3cDate, p.RelevantDate); err != nil {
			return errors.New("relevance: date has invalid format")
		}
	}
	if !hasOneStyle(
		p.BoardingPass,
		p.Coupon,
		p.EventTicket,
		p.Generic,
		p.StoreCard) {
		return errors.New("pass structure: only one style allowed")
	}
	if err := hasValidFields(
		p.BoardingPass,
		p.Coupon,
		p.EventTicket,
		p.Generic,
		p.StoreCard); err != nil {
		return err
	}
	if p.BoardingPass != nil && p.BoardingPass.TransitType == "" {
		return errors.New("boarding pass: transit type is empty")
	}
	if p.AuthenticationToken == "" {
		return errors.New("web service: authentication token is empty")
	}
	if p.WebServiceURL == "" {
		return errors.New("web service: url is empty")
	}
	for _, beacon := range p.Beacons {
		if beacon.ProximityUUID == "" {
			return errors.New("beacon: proximity uuid is empty")
		}
	}
	for _, location := range p.Locations {
		if location.Latitude == 0 {
			return errors.New("location: latitude must have value")
		}
		if location.Longitude == 0 {
			return errors.New("location: longitude must have value")
		}
	}
	for _, barcode := range p.Barcodes {
		if barcode.Format == "" {
			return errors.New("barcode: format is empty")
		}
		if barcode.Message == "" {
			return errors.New("barcode: message is empty")
		}
		if barcode.MessageEncoding == "" {
			return errors.New("barcode: message encoding is empty")
		}
	}
	if p.NFC != nil && p.NFC.Message == "" {
		return errors.New("nfc: message is empty")
	}
	return nil
}

// String returns string representation of struct.
func (p *PassCard) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(data)
}

// Value is a value that drivers must be able to handle.
func (p *PassCard) Value() (driver.Value, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return driver.Value(""), err
	}
	return driver.Value(string(data)), nil
}

// Scan value from database.
func (p *PassCard) Scan(src interface{}) error {
	var source []byte
	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	case sql.NullString:
		source = []byte("")
	default:
		return errors.New("invalid data type for PassCard")
	}

	if len(source) == 0 {
		source = []byte("{}")
	}
	return json.Unmarshal(source, &p)
}

func hasOneStyle(styles ...*PassStructure) bool {
	styleCnt := 0
	for _, style := range styles {
		if style != nil {
			styleCnt++
		}
	}
	return styleCnt == 1
}

func hasValidFields(styles ...*PassStructure) error {
	for _, style := range styles {
		if style != nil {
			for _, field := range style.AuxiliaryFields {
				if err := validField("auxiliary fields", field); err != nil {
					return err
				}
			}
			for _, field := range style.BackFields {
				if err := validField("back fields", field); err != nil {
					return err
				}
			}
			for _, field := range style.HeaderFields {
				if err := validField("header fields", field); err != nil {
					return err
				}
			}
			for _, field := range style.PrimaryFields {
				if err := validField("primary fields", field); err != nil {
					return err
				}
			}
			for _, field := range style.SecondaryFields {
				if err := validField("secondary fields", field); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func validField(prefix string, field *Field) error {
	if field.Key == "" {
		return fmt.Errorf("%s: key is empty", prefix)
	}
	if field.Value == nil {
		return fmt.Errorf("%s: value is nil", prefix)
	}
	return nil
}
