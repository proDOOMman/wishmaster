package actions

import (
	"io/ioutil"
	"strconv"

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
// Model: Singular (ChannelsPackage)
// DB Table: Plural (channels_packages)
// Resource: Plural (ChannelsPackages)
// Path: Plural (/channels_packages)
// View Template Folder: Plural (/templates/channels_packages/)

// ChannelsPackagesResource is the resource for the ChannelsPackage model
type ChannelsPackagesResource struct {
	buffalo.Resource
}

// List gets all ChannelsPackages. This function is mapped to the path
// GET /channels_packages
func (v ChannelsPackagesResource) List(c buffalo.Context) error {
	channelsPackages := &models.ChannelsPackages{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := models.DB.PaginateFromParams(c.Params())

	filterText := c.Request().URL.Query().Get("q")
	if len(filterText) > 0 {
		q = q.Where("name LIKE '%' || ? || '%'", filterText)
	}

	// Retrieve all ChannelsPackages from the DB
	if err := q.All(channelsPackages); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, channelsPackages))
}

// Show gets the data for one ChannelsPackage. This function is mapped to
// the path GET /channels_packages/{channels_package_id}
func (v ChannelsPackagesResource) Show(c buffalo.Context) error {
	// Allocate an empty ChannelsPackage
	channelsPackage := &models.ChannelsPackage{}

	// To find the ChannelsPackage the parameter channels_package_id is used.
	if err := models.DB.Find(channelsPackage, c.Param("channels_package_id")); err != nil {
		return c.Error(404, err)
	}

	if len(channelsPackage.XmltvUrl.String) > 0 {
		count, err := models.EPG_DB_RO.Where("url_hash = ?", channelsPackage.GetEpgUrlHash()).Count(models.EpgChannel{})
		if err == nil && count > 0 {
			c.Flash().Add("info", "Каналов EPG: "+strconv.FormatInt(int64(count), 10))
		}
		epgProgramme := models.EpgProgramme{}
		err = models.EPG_DB_RO.Where("url_hash = ?", channelsPackage.GetEpgUrlHash()).First(&epgProgramme)
		if err == nil {
			c.Flash().Add("info", "Время обновления EPG: "+epgProgramme.UpdatedAt.Format("02.01.2006 15:04"))
		} else {
			c.Flash().Add("warning", "Программа передач ещё не загружена")
		}
	}

	return c.Render(200, r.Auto(c, channelsPackage))
}

// New renders the form for creating a new ChannelsPackage.
// This function is mapped to the path GET /channels_packages/new
func (v ChannelsPackagesResource) New(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.ChannelsPackage{Active: true}))
}

// Create adds a ChannelsPackage to the DB. This function is mapped to the
// path POST /channels_packages
func (v ChannelsPackagesResource) Create(c buffalo.Context) error {
	// Allocate an empty ChannelsPackage
	channelsPackage := &models.ChannelsPackage{}

	// Bind channelsPackage to the html form elements
	if err := c.Bind(channelsPackage); err != nil {
		return errors.WithStack(err)
	}

	err := models.DB.Transaction(func(tx *pop.Connection) error {
		// Validate the data from the html form
		verrs, err := tx.ValidateAndCreate(channelsPackage)
		if err != nil {
			return errors.WithStack(err)
		}

		if verrs.HasAny() {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the new.html template that the user can
			// correct the input.
			return c.Render(422, r.Auto(c, channelsPackage))
		}

		err = saveM3uFileChannels(&c, tx, channelsPackage)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Пакет каналов успешно создан")

	// and redirect to the channels_packages index page
	return c.Render(201, r.Auto(c, channelsPackage))
}

func saveM3uFileChannels(c *buffalo.Context, tx *pop.Connection, channelsPackage *models.ChannelsPackage) error {
	f, err := (*c).File("M3UFile")
	if err != nil {
		return errors.WithStack(err)
	}
	if f.Valid() {
		// parse m3u
		fileData, err := ioutil.ReadAll(f)
		if err != nil {
			return errors.WithStack(err)
		}
		if len(fileData) > 0 {
			err = channelsPackage.CreateChannelsFromM3U(tx, fileData)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}
	return nil
}

// Edit renders a edit form for a ChannelsPackage. This function is
// mapped to the path GET /channels_packages/{channels_package_id}/edit
func (v ChannelsPackagesResource) Edit(c buffalo.Context) error {
	// Allocate an empty ChannelsPackage
	channelsPackage := &models.ChannelsPackage{}

	if err := models.DB.Find(channelsPackage, c.Param("channels_package_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, channelsPackage))
}

// Update changes a ChannelsPackage in the DB. This function is mapped to
// the path PUT /channels_packages/{channels_package_id}
func (v ChannelsPackagesResource) Update(c buffalo.Context) error {
	// Allocate an empty ChannelsPackage
	channelsPackage := &models.ChannelsPackage{}

	err := models.DB.Transaction(func(tx *pop.Connection) error {
		if err := tx.Find(channelsPackage, c.Param("channels_package_id")); err != nil {
			return c.Error(404, err)
		}

		// Bind ChannelsPackage to the html form elements
		if err := c.Bind(channelsPackage); err != nil {
			return errors.WithStack(err)
		}

		verrs, err := tx.ValidateAndUpdate(channelsPackage)
		if err != nil {
			return errors.WithStack(err)
		}

		if verrs.HasAny() {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			return c.Render(422, r.Auto(c, channelsPackage))
		}

		err = saveM3uFileChannels(&c, tx, channelsPackage)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Пакет каналов успешно сохранён")

	// and redirect to the channels_packages index page
	return c.Render(200, r.Auto(c, channelsPackage))
}

// Destroy deletes a ChannelsPackage from the DB. This function is mapped
// to the path DELETE /channels_packages/{channels_package_id}
func (v ChannelsPackagesResource) Destroy(c buffalo.Context) error {
	// Allocate an empty ChannelsPackage
	channelsPackage := &models.ChannelsPackage{}

	err := models.DB.Transaction(func(tx *pop.Connection) error {
		// To find the ChannelsPackage the parameter channels_package_id is used.
		if err := tx.Find(channelsPackage, c.Param("channels_package_id")); err != nil {
			return c.Error(404, err)
		}

		if err := tx.Destroy(channelsPackage); err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Пакет каналов успешно удален")

	// Redirect to the channels_packages index page
	return c.Render(200, r.Auto(c, channelsPackage))
}
