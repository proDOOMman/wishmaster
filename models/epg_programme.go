package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type EpgProgramme struct {
	ID           int          `json:"id" db:"id"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
	UrlHash      uint32       `json:"url_hash" db:"url_hash"`
	ChannelEpgID string       `json:"channel_epg_id" db:"channel_epg_id"`
	ProgrammeID  nulls.String `json:"programme_id" db:"programme_id"`
	Start        int64        `json:"start" db:"start"`
	Stop         int64        `json:"stop" db:"stop"`
	Title        string       `json:"title" db:"title"`
	Desc         nulls.String `json:"desc" db:"desc"`
	Categories   nulls.String `json:"categories" db:"categories"`
}

// String is not required by pop and may be deleted
func (e EpgProgramme) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// EpgProgrammes is not required by pop and may be deleted
type EpgProgrammes []EpgProgramme

// String is not required by pop and may be deleted
func (e EpgProgrammes) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (e *EpgProgramme) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: e.ChannelEpgID, Name: "ChannelEpgID"},
		&validators.StringIsPresent{Field: e.Title, Name: "Title"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (e *EpgProgramme) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (e *EpgProgramme) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
