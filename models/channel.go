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

type Channel struct {
	ID                int             `json:"id" db:"id"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
	Name              string          `json:"name" db:"name"`
	SearchName        string          `json:"-" db:"search_name"`
	Description       nulls.String    `json:"description" db:"description"`
	Url               string          `json:"url" db:"url"`
	Num               nulls.Int       `json:"num" db:"num"`
	Crypted           bool            `json:"crypted" db:"crypted"`
	Erotic            bool            `json:"erotic" db:"erotic"`
	StreamAspectRatio int             `json:"stream_aspect_ratio" db:"stream_aspect_ratio"`
	ZoomRatio         float32         `json:"zoom_ratio" db:"zoom_ratio"`
	EpgID             nulls.String    `json:"epg_id" db:"epg_id"`
	ChannelsPackage   ChannelsPackage `json:"-" belongs_to:"channels_package"`
	ChannelsPackageID int             `json:"channels_package_id" db:"channels_package_id"`
	EpgOffset         int             `json:"epg_offset" db:"epg_offset"`
}

// String is not required by pop and may be deleted
func (c Channel) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Channels is not required by pop and may be deleted
type Channels []Channel

// String is not required by pop and may be deleted
func (c Channels) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Channel) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Name, Name: "Name"},
		&validators.StringIsPresent{Field: c.Url, Name: "Url"},
		&validators.IntIsPresent{Field: c.StreamAspectRatio, Name: "StreamAspectRatio"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Channel) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Channel) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (c *Channel) BeforeSave(tx *pop.Connection) error {
	c.SearchName = strings.ToUpper(c.Name)
	return nil
}

func (c *Channel) BeforeCreate(tx *pop.Connection) error {
	channelsNum1 := []Channel{}
	err := tx.Where("num = 1").Limit(1).All(&channelsNum1)
	if len(channelsNum1) == 0 {
		c.Num = nulls.NewInt(1)
	} else {
		channelNum := ChannelNum{}
		err = tx.RawQuery("select num +1 as id from channels t1 where not exists (select * from channels t2 where t1.num +1 = t2.num) order by num limit 1").First(&channelNum)
		if err == nil {
			c.Num = nulls.NewInt(channelNum.ID)
		} else {
			c.Num = nulls.NewInt(1)
		}
	}
	return nil
}
