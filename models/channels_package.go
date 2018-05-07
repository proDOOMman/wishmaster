package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"hash/fnv"
	"strconv"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"gitlab.com/SML-482HD/wishmaster/tvg"
)

type ChannelsPackage struct {
	ID          int          `json:"id" db:"id"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	Active      bool         `json:"active" db:"active"`
	Name        string       `json:"name" db:"name"`
	Description nulls.String `json:"description" db:"description"`
	M3uUrl      nulls.String `json:"m3u_url" db:"m3u_url"`
	XmltvUrl    nulls.String `json:"xmltv_url" db:"xmltv_url"`
	GoogleKey   nulls.String `json:"google_key" db:"google_key"`
	GoogleCx    nulls.String `json:"google_cx" db:"google_cx"`
	Channels    Channels     `has_many:"channels" order_by:"num asc"`
}

// String is not required by pop and may be deleted
func (c ChannelsPackage) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// ChannelsPackages is not required by pop and may be deleted
type ChannelsPackages []ChannelsPackage

// String is not required by pop and may be deleted
func (c ChannelsPackages) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *ChannelsPackage) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Name, Name: "Name"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *ChannelsPackage) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *ChannelsPackage) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (c *ChannelsPackage) GetEpgUrlHash() uint32 {
	h := fnv.New32a()
	h.Write([]byte(c.XmltvUrl.String))
	return h.Sum32()
}

func (c *ChannelsPackage) GetEpgUrlHashString() string {
	return strconv.FormatUint(uint64(c.GetEpgUrlHash()), 10)
}

func (c ChannelsPackage) SelectValue() interface{} {
	return c.ID
}

func (c ChannelsPackage) SelectLabel() string {
	return c.Name
}

func (c *ChannelsPackage) AfterSave(tx *pop.Connection) error {
	err := c.UpdateChannelsFromUrl(tx)
	if err != nil {
		return err
	}
	go UpdateChannelPackageEpg(c, false)
	return nil
}

func (c *ChannelsPackage) CreateChannelsFromM3U(tx *pop.Connection, fileData []byte) error {
	m3uData, err := tvg.Parse(fileData)
	if err != nil {
		return errors.WithStack(err)
	}
	currentTime := time.Now()
	for _, inf := range m3uData.List {
		if inf == nil || len(inf.Title) == 0 || len(inf.Url) == 0 {
			continue
		}
		channel := Channel{}
		err := tx.Where("channels_package_id = ? AND url = ?", c.ID, inf.Url).First(&channel)
		if err != nil {
			err := tx.Where("channels_package_id = ? AND Name = ?", c.ID, inf.Title).First(&channel)
			if err != nil {
				channel.ChannelsPackageID = c.ID
				channel.Crypted = true
				channel.StreamAspectRatio = 1
				channel.ZoomRatio = 0.5
			}
		}
		channel.Name = inf.Title
		channel.EpgID = nulls.NewString(inf.Id)
		channel.Url = inf.Url
		channel.Description = nulls.NewString(inf.Name)

		verrs := &validate.Errors{}
		if channel.ID == 0 {
			verrs, err = tx.ValidateAndCreate(&channel)
		} else {
			verrs, err = tx.ValidateAndUpdate(&channel)
		}

		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return errors.New(verrs.String())
		}
	}

	err = tx.RawQuery("DELETE FROM channels WHERE channels_package_id = ? AND updated_at < ?", c.ID, currentTime).Exec()

	return err
}

func (c *ChannelsPackage) UpdateChannelsFromUrl(tx *pop.Connection) error {
	if len(c.M3uUrl.String) == 0 {
		return nil
	}

	log.Println("Updating playlist from " + c.M3uUrl.String)
	// Get the data
	resp, err := http.Get(c.M3uUrl.String)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return errors.New("Download m3u failed: bad status code")
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	// Writer the body to file
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err
	}

	return c.CreateChannelsFromM3U(tx, b.Bytes())
}
