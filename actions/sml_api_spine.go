package actions

import (
	"net/http"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/pkg/errors"
	"gitlab.com/SML-482HD/wishmaster/models"
)

func SmlApiSpineInit(app *buffalo.App) {

	appSpine := app.Group("/spine")

	appSpine.GET("/services/ncdxml/ServiceAccountExtended/authorize_stb", SmlAuthorizeStb)

	appSpine.GET("/services/ncdxml/ServiceAccount/get_quotas", SmlGetQuotas)

	appSpine.GET("/services/ncdxml/ServiceAccount/get_billing_info", SmlGetBillingInfo)

	appSpine.GET("/services/ncdxml/ServiceAccount/list_services", SmlListServices)

	appSpine.GET("/services/ncdxml/ServiceAccount/list_profiles", SmlListProfiles)

	appSpine.POST("/services/ncdxml/Profile/update_new", SmlUpdateNewPost)

	appSpine.GET("/services/ncdxml/Profile/update_new", SmlUpdateNewGet)

	appSpine.Middleware.Skip(csrf.New, SmlUpdateNewPost)
}

func SmlAuthorizeStb(c buffalo.Context) error {

	serial := c.Request().URL.Query().Get("sn")
	mac := c.Request().URL.Query().Get("mac")

	row := &models.SmlAuthorizeAnswer{
		LocationId:           locationId,
		SubLocationId:        locationId,
		SmPassword:           serial,
		ServiceAccountNumber: serial,
		DateTime:             time.Now().Format("02/01/2016 15:04:05"),
		State:                "ACTIVE",
		ActivationNumber:     serial,
		BackendVersion:       "0.0.1",
		ProductOfferId:       productOfferId,
		CasId:                casId,
		ShouldUseRaptor:      0,
		NetLogIP:             "127.0.0.1",
		NetLogPort:           5040,
		ProviderId:           providerId,
		ProviderExtId:        providerExtId,
		TerminalType:         1,
		NetworkBlock:         1,
		SmLogin:              mac,
		AccountType:          "free",
		IsOfferConfirmed:     0,
		Multiroom:            true,
		Unixtime:             time.Now().Unix(),
		ResourceId:           mac,
		IsTrustedDevice:      0,
		IsOss:                1,
		IsHomeGroup:          1,
		Pin:                  pin,
		DeviceType:           "STB",
		UseOttUrlForChannels: 0,
		IP:                   c.Request().Host,
		StbFunctions:         stbFunctions,
		Multiscreen:          0,
	}

	v := &models.SmlAuthorizeRows{Code: 0, AuthorizeAnswer: *row}

	return c.Render(http.StatusOK, render.XML(v))
}

func SmlGetQuotas(c buffalo.Context) error {

	serviceAccountNumber := c.Request().URL.Query().Get("serviceAccountNumber")

	row := &models.SmlQuotaRow{ServiceAccountNumber: serviceAccountNumber, AllQuote: 1024}

	v := &models.SmlQuotaRows{Code: 0, QuotaRow: *row}

	return c.Render(http.StatusOK, render.XML(v))
}

func SmlGetBillingInfo(c buffalo.Context) error {

	v := &models.SmlBillingRows{}

	return c.Render(http.StatusOK, render.XML(v))
}

func SmlListServices(c buffalo.Context) error {

	channelsPackages := &models.ChannelsPackages{}

	// Retrieve all ChannelsPackages from the DB
	if err := models.DB.Q().Where("Active = ?", true).All(channelsPackages); err != nil {
		return errors.WithStack(errors.New("Can't retrieve all ChannelsPackages from the DB"))
	}

	var services []models.SmlService

	tstvService := models.SmlService{
		Id:               tstvServiceId,
		Type:             "TSTVCHANNELPACKAGE",
		Name:             "TSTV_All",
		Description:      "",
		Mandatory:        false,
		Price:            0,
		ExternalId:       "TSTV",
		EndDate:          endDate,
		AllowedPurchases: 1000,
		Unlimited:        true,
		StartDate:        startDate,
		OneTime:          false,
		ServiceState:     1,
		IsPromoService:   0,
		NotConfirmed:     0,
		PortalId:         1,
	}
	services = append(services, tstvService)
	for _, element := range *channelsPackages {
		service := models.SmlService{
			Id:               element.ID,
			Type:             "CHANNELPACKAGE",
			Name:             element.Name,
			Description:      element.Description.String,
			Mandatory:        false,
			Price:            0,
			ExternalId:       "CUSTOM",
			EndDate:          endDate,
			AllowedPurchases: 1000,
			Unlimited:        true,
			StartDate:        startDate,
			OneTime:          false,
			ServiceState:     1,
			IsPromoService:   0,
			NotConfirmed:     0,
			PortalId:         1,
		}
		services = append(services, service)
	}

	v := &models.SmlServices{Code: 0, Services: services}

	return c.Render(http.StatusOK, render.XML(v))
}

func SmlListProfiles(c buffalo.Context) error {

	v := &models.SmlProfileRows{}

	serviceAccount := models.ServiceAccount{}

	mac := c.Request().Header.Get("X-Smartlabs-Mac-Address")
	err := models.DB.Where("uuid = ?", mac).First(&serviceAccount)
	if err == nil {
		v.ProfileRow.Sort = serviceAccount.ChannelsSort.String
		v.ProfileRow.Favorite = serviceAccount.FavoriteChannels.String
		v.ProfileRow.Reminders = serviceAccount.Reminders.String
	} else {
		c.Logger().Debug(err)
	}

	return c.Render(http.StatusOK, render.XML(v))
}

func SmlUpdateNewGet(c buffalo.Context) error {

	mac := c.Request().Header.Get("X-Smartlabs-Mac-Address")

	serviceAccount := models.ServiceAccount{}
	err := models.DB.Where("uuid = ?", mac).First(&serviceAccount)
	if err != nil {
		serviceAccount.Uuid = mac
		serviceAccount.ChannelsSort = nulls.NewString("")
		serviceAccount.Reminders = nulls.NewString("")
		serviceAccount.FavoriteChannels = nulls.NewString("")
	}

	channelsSortOrder := c.Request().URL.Query().Get("channelsSortOrder")
	if len(channelsSortOrder) > 0 {
		serviceAccount.ChannelsSort = nulls.NewString(channelsSortOrder)
	}

	reminders := c.Request().URL.Query().Get("reminders")
	if len(reminders) > 0 {
		serviceAccount.Reminders = nulls.NewString(reminders)
	}

	favorite := c.Request().URL.Query().Get("favorite")
	if len(favorite) > 0 {
		serviceAccount.FavoriteChannels = nulls.NewString(favorite)
	}

	err = models.DB.Transaction(func(tx *pop.Connection) error {
		verrs, err := tx.ValidateAndSave(&serviceAccount)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return errors.New(verrs.String())
		}
		return nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(http.StatusOK, render.String(""))
}

func SmlUpdateNewPost(c buffalo.Context) error {
	mac := c.Request().Header.Get("X-Smartlabs-Mac-Address")

	serviceAccount := models.ServiceAccount{}
	err := models.DB.Where("uuid = ?", mac).First(&serviceAccount)
	if err != nil {
		serviceAccount.Uuid = mac
		serviceAccount.ChannelsSort = nulls.NewString("")
		serviceAccount.Reminders = nulls.NewString("")
		serviceAccount.FavoriteChannels = nulls.NewString("")
	}

	channelsSortOrder := c.Request().FormValue("channelsSortOrder")
	if len(channelsSortOrder) > 0 {
		serviceAccount.ChannelsSort = nulls.NewString(channelsSortOrder)
	}

	reminders := c.Request().FormValue("reminders")
	if len(reminders) > 0 {
		serviceAccount.Reminders = nulls.NewString(reminders)
	}

	favorite := c.Request().FormValue("favorite")
	if len(favorite) > 0 {
		serviceAccount.FavoriteChannels = nulls.NewString(favorite)
	}

	err = models.DB.Transaction(func(tx *pop.Connection) error {
		verrs, err := tx.ValidateAndSave(&serviceAccount)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return errors.New(verrs.String())
		}
		return nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(http.StatusOK, render.String(""))
}
