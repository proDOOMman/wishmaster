package actions

import (
	"strconv"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"gitlab.com/SML-482HD/wishmaster/models"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Channel)
// DB Table: Plural (channels)
// Resource: Plural (Channels)
// Path: Plural (/channels)
// View Template Folder: Plural (/templates/channels/)

// ChannelsResource is the resource for the Channel model
type ChannelsResource struct {
	buffalo.Resource
}

// List gets all Channels. This function is mapped to the path
// GET /channels
func (v ChannelsResource) List(c buffalo.Context) error {
	channels := &models.Channels{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := models.DB.PaginateFromParams(c.Params())
	sort := c.Request().URL.Query().Get("sort")
	if len(sort) > 0 {
		q = q.Order(sort)
	} else {
		q = q.Order("num")
	}

	filterText := c.Request().URL.Query().Get("q")
	if len(filterText) > 0 {
		q = q.Where("search_name LIKE '%' || ? || '%'", strings.ToUpper(filterText))
	}

	// Retrieve all Channels from the DB
	if err := q.All(channels); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, channels))
}

// Show gets the data for one Channel. This function is mapped to
// the path GET /channels/{channel_id}
func (v ChannelsResource) Show(c buffalo.Context) error {
	// Allocate an empty Channel
	channel := &models.Channel{}

	// To find the Channel the parameter channel_id is used.
	if err := models.DB.Eager("ChannelsPackage").Find(channel, c.Param("channel_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, channel))
}

func setChannelsPackagesSelectable(tx *pop.Connection, c *buffalo.Context) error {
	channelsPackages := &models.ChannelsPackages{}
	if err := tx.Q().All(channelsPackages); err != nil {
		return errors.WithStack(errors.New("Can't retrieve channels packages!"))
	}
	(*c).Set("channelsPackages", *channelsPackages)
	return nil
}

// New renders the form for creating a new Channel.
// This function is mapped to the path GET /channels/new
func (v ChannelsResource) New(c buffalo.Context) error {
	err := setChannelsPackagesSelectable(models.DB, &c)
	if err != nil {
		return err
	}

	return c.Render(200, r.Auto(c, &models.Channel{ZoomRatio: 0.5, StreamAspectRatio: 1, ChannelsPackageID: 1, Crypted: true}))
}

// Create adds a Channel to the DB. This function is mapped to the
// path POST /channels
func (v ChannelsResource) Create(c buffalo.Context) error {
	// Allocate an empty Channel
	channel := &models.Channel{}

	// Bind channel to the html form elements
	if err := c.Bind(channel); err != nil {
		return errors.WithStack(err)
	}

	err := models.DB.Transaction(func(tx *pop.Connection) error {
		// Validate the data from the html form
		verrs, err := tx.ValidateAndCreate(channel)
		if err != nil {
			return errors.WithStack(err)
		}

		if verrs.HasAny() {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			err := setChannelsPackagesSelectable(tx, &c)
			if err != nil {
				return err
			}
			// Render again the new.html template that the user can
			// correct the input.
			return c.Render(422, r.Auto(c, channel))
		}

		return nil
	})
	if err != nil {
		return err
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Канал успешно создан")

	// and redirect to the channels index page
	return c.Render(201, r.Auto(c, channel))
}

// Edit renders a edit form for a Channel. This function is
// mapped to the path GET /channels/{channel_id}/edit
func (v ChannelsResource) Edit(c buffalo.Context) error {
	err := setChannelsPackagesSelectable(models.DB, &c)
	if err != nil {
		return err
	}

	// Allocate an empty Channel
	channel := &models.Channel{}

	if err := models.DB.Find(channel, c.Param("channel_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, channel))
}

// Update changes a Channel in the DB. This function is mapped to
// the path PUT /channels/{channel_id}
func (v ChannelsResource) Update(c buffalo.Context) error {
	// Allocate an empty Channel
	channel := &models.Channel{}
	err := models.DB.Transaction(func(tx *pop.Connection) error {
		if err := tx.Find(channel, c.Param("channel_id")); err != nil {
			return c.Error(404, err)
		}

		// Bind Channel to the html form elements
		if err := c.Bind(channel); err != nil {
			return errors.WithStack(err)
		}

		verrs, err := tx.ValidateAndUpdate(channel)
		if err != nil {
			return errors.WithStack(err)
		}

		if verrs.HasAny() {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			return c.Render(422, r.Auto(c, channel))
		}

		return nil
	})
	if err != nil {
		return err
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Канал успешно изменен")

	// and redirect to the channels index page
	return c.Render(200, r.Auto(c, channel))
}

// Destroy deletes a Channel from the DB. This function is mapped
// to the path DELETE /channels/{channel_id}
func (v ChannelsResource) Destroy(c buffalo.Context) error {
	// Allocate an empty Channel
	channel := &models.Channel{}

	err := models.DB.Transaction(func(tx *pop.Connection) error {
		// To find the Channel the parameter channel_id is used.
		if err := tx.Find(channel, c.Param("channel_id")); err != nil {
			return c.Error(404, err)
		}

		if err := tx.Destroy(channel); err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Канал успешно удален")

	// Redirect to the channels index page
	return c.Render(200, r.Auto(c, channel))
}

func ChannelNameTypeahead(c buffalo.Context) error {
	// Allocate an empty Channel
	channelsPackage := &models.ChannelsPackage{}

	channelsPackageId, err := strconv.ParseInt(c.Param("id"), 10, 0)
	if err != nil {
		return errors.WithStack(err)
	}
	// To find the Channel the parameter channel_id is used.
	if err = models.DB.Find(channelsPackage, int(channelsPackageId)); err != nil {
		return errors.WithStack(err)
	}

	epgChannels := &models.EpgChannels{}
	models.EPG_DB_RO.Q().Limit(10).Where("url_hash = ? AND search_name LIKE ? || '%'", uint32(channelsPackage.GetEpgUrlHash()), strings.ToUpper(c.Param("search"))).Order("display_name").All(epgChannels)

	namesArray := []string{}
	for _, element := range *epgChannels {
		namesArray = append(namesArray, element.DisplayName)
	}
	c.LogField("epgChannels", len(*epgChannels))

	return c.Render(200, r.JSON(namesArray))
}
