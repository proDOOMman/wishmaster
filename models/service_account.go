package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type ServiceAccount struct {
	ID               int          `json:"id" db:"id"`
	CreatedAt        time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at" db:"updated_at"`
	Uuid             string       `json:"uuid" db:"uuid"`
	ChannelsSort     nulls.String `json:"channels_sort" db:"channels_sort"`
	Reminders        nulls.String `json:"reminders" db:"reminders"`
	FavoriteChannels nulls.String `json:"favorite_channels" db:"favorite_channels"`
}

// String is not required by pop and may be deleted
func (s ServiceAccount) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// ServiceAccounts is not required by pop and may be deleted
type ServiceAccounts []ServiceAccount

// String is not required by pop and may be deleted
func (s ServiceAccounts) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *ServiceAccount) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.Uuid, Name: "Uuid"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *ServiceAccount) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *ServiceAccount) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
