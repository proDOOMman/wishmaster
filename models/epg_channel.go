package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type EpgChannel struct {
	ID          int          `json:"id" db:"id"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	UrlHash     uint32       `json:"url_hash" db:"url_hash"`
	EpgID       string       `json:"epg_id" db:"epg_id"`
	DisplayName string       `json:"display_name" db:"display_name"`
	SearchName  string       `json:"search_name" db:"search_name"`
	IconSrc     nulls.String `json:"icon_src" db:"icon_src"`
}

// String is not required by pop and may be deleted
func (e EpgChannel) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// EpgChannels is not required by pop and may be deleted
type EpgChannels []EpgChannel

// String is not required by pop and may be deleted
func (e EpgChannels) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (e *EpgChannel) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: e.EpgID, Name: "EpgID"},
		&validators.StringIsPresent{Field: e.DisplayName, Name: "DisplayName"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (e *EpgChannel) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (e *EpgChannel) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (e *EpgChannel) BeforeSave(tx *pop.Connection) error {
	e.SearchName = strings.ToUpper(e.DisplayName)
	return nil
}
