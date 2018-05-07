package actions

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/pop"
	"gitlab.com/SML-482HD/wishmaster/models"
)

func SmlApiV1Init(app *buffalo.App) {

	appNewApi := app.Group("/api/v1")

	appNewApi.GET("/device/authorize", SmlV1Authorize)

}

func SmlV1Authorize(c buffalo.Context) error {
	uuid := c.Request().URL.Query().Get("uuid")
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		smlError := models.SmlV1ErrorResult{Error: &models.SmlV1Error{Code: 1, Message: "No transaction found"}}
		return c.Render(http.StatusInternalServerError, render.JSON(smlError))
	}

	acc := &models.ServiceAccount{}
	err := tx.Where("uuid = ?", uuid).First(acc)
	if err != nil {
		acc.Uuid = uuid
		verrs, err := tx.ValidateAndCreate(acc)
		if err != nil {
			smlError := models.SmlV1ErrorResult{Error: &models.SmlV1Error{Code: 1, Message: err.Error()}}
			return c.Render(http.StatusInternalServerError, render.JSON(smlError))
		}
		if verrs.HasAny() {
			smlError := models.SmlV1ErrorResult{Error: &models.SmlV1Error{Code: 1, Message: verrs.String()}}
			return c.Render(http.StatusInternalServerError, render.JSON(smlError))
		}
	}

	result := models.SmlV1AuthorizeDevice{
		DeviceType:           "STB",
		ServiceAccountNumber: acc.ID,
		ServiceAccountState:  "ACTIVE",
		ProductOfferId:       productOfferId,
		Login:                strconv.FormatInt(int64(acc.ID), 10),
		Password:             "1111",
		Multiroom:            false,
		Unixtime:             time.Now().Unix(),
		Currency:             models.SmlV1Currency{Name: "RUB", CurrencyCode: "R00001", CurrencyStrCode: "RUB"},
		DeviceList:           []models.SmlV1Device{{DeviceType: "STB", Uuid: uuid, Id: acc.ID, IsCurrent: true, TerminalName: "STB"}},
		IsTest:               false,
	}
	return c.Render(http.StatusOK, render.JSON(models.SmlV1AuthorizeDeviceResult{Result: &result}))
}
