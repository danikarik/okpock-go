package api

import (
	"encoding/json"
	"errors"
	"time"
)

// NewUpload creates a new instance of `Upload`.
func NewUpload(uuid, fname, hash string) *Upload {
	return &Upload{
		UUID:      uuid,
		Filename:  fname,
		Hash:      hash,
		CreatedAt: time.Now(),
	}
}

// Upload holds upload filename, key and checksum.
type Upload struct {
	ID int64 `json:"id" db:"id"`

	UUID     string `json:"uuid" db:"uuid"`
	Filename string `json:"filename" db:"filename"`
	Hash     string `json:"hash" db:"hash"`

	Body        []byte `json:"-" db:"-"`
	ContentType string `json:"-" db:"-"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// IsValid checks whether input is valid or not.
func (u *Upload) IsValid() error {
	if u.UUID == "" {
		return errors.New("uuid is empty")
	}
	if u.Filename == "" {
		return errors.New("filename is empty")
	}
	if u.Hash == "" {
		return errors.New("hash is empty")
	}
	return nil
}

// String returns string representation of struct.
func (u *Upload) String() string {
	data, err := json.Marshal(u)
	if err != nil {
		return ""
	}
	return string(data)
}

// Uploads holds next page token and items.
type Uploads struct {
	Opts *PagingOptions
	Data []*Upload
}
