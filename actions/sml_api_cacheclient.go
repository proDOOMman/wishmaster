// sml_api_cacheclient
package actions

import (
	"encoding/xml"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/pkg/errors"
	"gitlab.com/SML-482HD/wishmaster/models"
)

const maxEpgChannelId = 10000
const epgProgrammeStartDiff = 1523000000
const tstvServiceId = 10000

func SmlApiCacheClientInit(app *buffalo.App) {

	appCacheClient := app.Group("/CacheClient")

	appCacheClient.GET("/ncdxml/ChannelPackage/list_channels", SmlListChannels)

	appCacheClient.GET("/simple/ncdxml/ProductOffer/list_services_terminal", SmlListServicesTerminal)

	appCacheClient.GET("/simple/ncdxml/ChannelPackage/list_incompatible_packages", SmlListIncompatiblePackages)

	appCacheClient.GET("/simple/ncdxml/AccessLevelDictionary/list", SmlListAccessLevelDictionary)

	appCacheClient.GET("/ncdxml/EPG/get_by_Chnnl", SmlGetEpg)

	appCacheClient.GET("/ncdxml/EPG/get_by_pkg", SmlEpgByPkg)

	appCacheClient.GET("/ncdxml/EPG/get_dsc", SmlGetDsc)

	appCacheClient.GET("/ncdxml/EPG/get_delta", SmlGetDelta)

}

func SmlGetDelta(c buffalo.Context) error {
	poId, err := strconv.ParseInt(c.Request().URL.Query().Get("poId"), 10, 0)
	if err != nil {
		return errors.WithStack(err)
	}
	// version, err := strconv.ParseInt(c.Request().URL.Query().Get("version"), 10, 0)
	// if err != nil {
	// 	return errors.WithStack(err)
	// }

	channelsPackage := models.ChannelsPackage{}
	err = models.DB.Find(&channelsPackage, int(poId))
	if err != nil {
		return errors.WithStack(err)
	}
	v := models.SmlEpgDeltaRows{}
	v.EpgDeltaRow.Type = "DIFF"
	v.EpgDeltaRow.Version = time.Now().Unix()
	v.EpgDeltaRow.PoId = int(poId)

	// t := time.Unix(version, 0)
	// epgProgrammes := &models.EpgProgrammes{}
	// err = models.EPG_DB_RO.Where("url_hash = ? AND updated_at > ?", channelsPackage.GetEpgUrlHash(), t).All(epgProgrammes)
	// if err == nil {
	// 	var buffer bytes.Buffer
	// 	for _, epgProgramme := range *epgProgrammes {
	// 		buffer.WriteString(strconv.FormatInt(int64(epgProgramme.ID), 10))
	// 		buffer.WriteString(",")
	// 	}
	// 	v.EpgDeltaRow.Add = buffer.String()
	// }

	return c.Render(http.StatusOK, render.XML(v))
}

func SmlGetDsc(c buffalo.Context) error {
	epgProgramme := &models.EpgProgramme{}
	strId := c.Request().URL.Query().Get("ID")
	id, err := strconv.ParseInt(strId, 10, 0)
	if err != nil {
		return errors.WithStack(err)
	}
	channelId := int(id % maxEpgChannelId)
	c.Logger().Debug(id/maxEpgChannelId + epgProgrammeStartDiff)
	start := int64(id/maxEpgChannelId + epgProgrammeStartDiff)
	epgChannel := &models.EpgChannel{}
	err = models.EPG_DB_RO.Find(epgChannel, int(channelId))
	if err != nil {
		return errors.WithStack(err)
	}
	err = models.EPG_DB_RO.Where("url_hash = ? AND channel_epg_id = ? AND start = ?", epgChannel.UrlHash, epgChannel.EpgID, start).First(epgProgramme)
	desc := models.SmlEpgDescRows{}
	desc.EpgDescRow.Id = int(id)
	if err != nil {
		return errors.WithStack(err)
	} else {
		desc.EpgDescRow.Description = epgProgramme.Desc.String
		return c.Render(http.StatusOK, render.XML(desc))
	}
}

func SmlEpgByPkg(c buffalo.Context) error {
	type EpgList struct {
		XMLName     xml.Name `xml:"epg_list"`
		Code        int      `xml:"code,attr"`
		EpgListItem struct {
			XMLName xml.Name `xml:"epg"`
			Version int64    `xml:"version,attr"`
		}
	}
	epgList := EpgList{}
	epgList.EpgListItem.Version = time.Now().Unix()
	return c.Render(http.StatusOK, render.XML(epgList))
}

func SmlListChannels(c buffalo.Context) error {

	channelsPackageId, err := strconv.ParseInt(c.Request().URL.Query().Get("channelPackageId"), 10, 0)
	if err != nil {
		return errors.WithStack(errors.New("Can't parse channel package ID"))
	}

	channelsPackage := &models.ChannelsPackage{}
	err = models.DB.Eager("Channels").Find(channelsPackage, int(channelsPackageId))
	if err != nil {
		return errors.WithStack(errors.New("Channels package not found!"))
	}

	smlChannels := []models.SmlChannel{}
	for _, channel := range (*channelsPackage).Channels {
		epg_channel := models.EpgChannel{}
		var logoUrl string
		if len(channel.Name) > 0 {
			logoUrl = "create/" + url.PathEscape(strconv.FormatInt(int64(channel.Num.Int), 10))
		}
		if len(channel.EpgID.String) > 0 {
			models.EPG_DB_RO.Q().Where("url_hash = ? AND epg_id = ?", channelsPackage.GetEpgUrlHash(), channel.EpgID).First(&epg_channel)
		} else {
			models.EPG_DB_RO.Q().Where("url_hash = ? AND display_name = ?", channelsPackage.GetEpgUrlHash(), channel.Name).First(&epg_channel)
		}
		if len(epg_channel.IconSrc.String) > 0 {
			logoUrl = epg_channel.IconSrc.String
		}

		smlChannel := models.SmlChannel{
			Id:                channel.ID,
			BcName:            channel.Name,
			BcDesc:            channel.Description.String,
			Num:               channel.Num.Int,
			Logo:              logoUrl,
			Logo2:             logoUrl,
			Url:               channel.Url,
			StreamAspectRatio: channel.StreamAspectRatio,
			ZoomRatio:         channel.ZoomRatio,
			BcRId:             channel.ID,
			OttURL:            channel.Url,
			SmlOttURL:         channel.Url,
			TstvOttURL:        channel.Url,
			OttDvr:            1,
			NPVRChannelID:     "npvr",
			LocId:             locationId,
			PackageId:         channel.ChannelsPackageID,
			EpgOffset:         channel.EpgOffset,
		}
		if channel.Crypted {
			smlChannel.IsCrypted = 1
		}
		if channel.Erotic {
			smlChannel.IsErotic = 1
		}
		smlChannel.IsDvrCrypted = smlChannel.IsCrypted
		smlChannel.IsOttEncrypted = smlChannel.IsCrypted
		smlChannels = append(smlChannels, smlChannel)
	}
	channelPackage := models.SmlChannelPackage{
		PackageId:     int(channelsPackageId),
		LocationId:    locationId,
		SubLocationId: locationId,
		Channels:      smlChannels,
	}

	v := &models.SmlChannelsList{Channels: []models.SmlChannelPackage{channelPackage}}

	return c.Render(http.StatusOK, render.XML(v))
}

func SmlListServicesTerminal(c buffalo.Context) error {

	channelsPackages := &models.ChannelsPackages{}

	// Retrieve all ChannelsPackages from the DB
	if err := models.DB.Q().Where("Active = ?", true).All(channelsPackages); err != nil {
		return errors.WithStack(errors.New("Can't retrieve all ChannelsPackages from the DB"))
	}

	var terminalServices []models.SmlTerminalService

	tstvService := models.SmlTerminalService{
		Id:                        tstvServiceId,
		RId:                       tstvServiceId,
		Type:                      "TSTVCHANNELPACKAGE",
		Name:                      "TSTV_All",
		Description:               "",
		ChangeDate:                int32(time.Now().Unix()),
		Mandatory:                 false,
		Price:                     -1,
		ExternalId:                "TSTV",
		AllowedPurchases:          0,
		Unlimited:                 true,
		CanBeSubscribedByUser:     false,
		IsInvisible:               true,
		ShowIfInvisibleAndDepends: false,
		QuotaSizeLimit:            -1,
		EndDate:                   endDate,
		IsDaily:                   false,
		PlType:                    "subscribe",
		TerminalType:              0,
		PlDiscriminator:           "TSTVCHANNELPACKAGE_S_PRICELIST",
		PlId:                      0,
		NotRecomExclusive:         false,
		UnsubscribeByLimit:        false,
		IsAccessByLock:            true,
		ReplayDenied:              false,
		SubscribeMode:             1,
		ChannelPreviewEnabled:     true,
		ForAllAcc:                 1,
		PortalId:                  1,
		NonBlockedPacksInBundle:   0,
		CurrencyInternalId:        currencyInternalId,
	}

	for _, element := range *channelsPackages {
		service := models.SmlTerminalService{
			Id:                        element.ID,
			RId:                       element.ID,
			Type:                      "CHANNELPACKAGE",
			Name:                      element.Name,
			Description:               element.Description.String,
			ChangeDate:                int32(time.Now().Unix()),
			Mandatory:                 false,
			Price:                     0,
			ExternalId:                "CUSTOM",
			AllowedPurchases:          0,
			Unlimited:                 true,
			CanBeSubscribedByUser:     true,
			IsInvisible:               false,
			ShowIfInvisibleAndDepends: false,
			QuotaSizeLimit:            -1,
			EndDate:                   endDate,
			IsDaily:                   false,
			PlType:                    "subscribe",
			TerminalType:              0,
			PlDiscriminator:           "CHANNELPACKAGE_S_PRICELIST",
			PlId:                      0,
			NotRecomExclusive:         false,
			UnsubscribeByLimit:        false,
			IsAccessByLock:            true,
			ReplayDenied:              false,
			SubscribeMode:             1,
			ChannelPreviewEnabled:     true,
			ForAllAcc:                 1,
			PortalId:                  1,
			NonBlockedPacksInBundle:   0,
			CurrencyInternalId:        currencyInternalId,
			Depends:                   []int{tstvServiceId},
		}
		terminalServices = append(terminalServices, service)
		tstvService.Parent = append(tstvService.Parent, service.Id)
	}
	terminalServices = append(terminalServices, tstvService)

	v := &models.SmlTerminalServices{Code: 0, Version: 0, TerminalServices: terminalServices}

	return c.Render(http.StatusOK, render.XML(v))
}

func SmlListIncompatiblePackages(c buffalo.Context) error {

	v := models.SmlIncompatiblePackages{IdArray: []int{}}

	return c.Render(http.StatusOK, render.XML(v))
}

func SmlListAccessLevelDictionary(c buffalo.Context) error {

	type Rows struct {
		// <rows code="0">
		XMLName         xml.Name                   `xml:"rows"`
		Number          int                        `xml:"code,attr"`
		AccessLevelRows []models.SmlAccessLevelRow `xml:"row"`
	}

	accessLevelRowArray := []models.SmlAccessLevelRow{
		{
			Number:                             0,
			AccessLevelDictionaryID:            0,
			AccessLevelDictionaryDISCRIMINATOR: "ACCESSLEVELDICTIONARY",
			AccessLevelDictionaryName:          "0+",
			AccessLevelDictionarySortOrder:     1,
			AccessLevelDictionaryLangId:        11544123,
			LangIdName:                         "Русский",
			LangIdID:                           11544123,
			LangIdDISCRIMINATOR:                "LANGUAGE"},
		{
			Number:                             1,
			AccessLevelDictionaryID:            1,
			AccessLevelDictionaryDISCRIMINATOR: "ACCESSLEVELDICTIONARY",
			AccessLevelDictionaryName:          "3+",
			AccessLevelDictionarySortOrder:     2,
			AccessLevelDictionaryLangId:        11544123,
			LangIdName:                         "Русский",
			LangIdID:                           11544123,
			LangIdDISCRIMINATOR:                "LANGUAGE"},
		{
			Number:                             2,
			AccessLevelDictionaryID:            2,
			AccessLevelDictionaryDISCRIMINATOR: "ACCESSLEVELDICTIONARY",
			AccessLevelDictionaryName:          "6+",
			AccessLevelDictionarySortOrder:     3,
			AccessLevelDictionaryLangId:        11544123,
			LangIdName:                         "Русский",
			LangIdID:                           11544123,
			LangIdDISCRIMINATOR:                "LANGUAGE"},
		{
			Number:                             3,
			AccessLevelDictionaryID:            3,
			AccessLevelDictionaryDISCRIMINATOR: "ACCESSLEVELDICTIONARY",
			AccessLevelDictionaryName:          "12+",
			AccessLevelDictionarySortOrder:     4,
			AccessLevelDictionaryLangId:        11544123,
			LangIdName:                         "Русский",
			LangIdID:                           11544123,
			LangIdDISCRIMINATOR:                "LANGUAGE"},
		{
			Number:                             4,
			AccessLevelDictionaryID:            4,
			AccessLevelDictionaryDISCRIMINATOR: "ACCESSLEVELDICTIONARY",
			AccessLevelDictionaryName:          "14+",
			AccessLevelDictionarySortOrder:     5,
			AccessLevelDictionaryLangId:        11544123,
			LangIdName:                         "Русский",
			LangIdID:                           11544123,
			LangIdDISCRIMINATOR:                "LANGUAGE"},
		{
			Number:                             5,
			AccessLevelDictionaryID:            5,
			AccessLevelDictionaryDISCRIMINATOR: "ACCESSLEVELDICTIONARY",
			AccessLevelDictionaryName:          "16+",
			AccessLevelDictionarySortOrder:     6,
			AccessLevelDictionaryLangId:        11544123,
			LangIdName:                         "Русский",
			LangIdID:                           11544123,
			LangIdDISCRIMINATOR:                "LANGUAGE"},
		{
			Number:                             6,
			AccessLevelDictionaryID:            6,
			AccessLevelDictionaryDISCRIMINATOR: "ACCESSLEVELDICTIONARY",
			AccessLevelDictionaryName:          "18+",
			AccessLevelDictionarySortOrder:     7,
			AccessLevelDictionaryLangId:        11544123,
			LangIdName:                         "Русский",
			LangIdID:                           11544123,
			LangIdDISCRIMINATOR:                "LANGUAGE"},
	}

	v := &Rows{AccessLevelRows: accessLevelRowArray}

	return c.Render(http.StatusOK, render.XML(v))
}

func SmlGetEpg(c buffalo.Context) error {
	startUnixDt, err := strconv.ParseInt(c.Request().URL.Query().Get("startUnixDt"), 10, 0)
	if err != nil {
		startUnixDt = 0
	}

	//from, err := strconv.ParseInt(c.Request().URL.Query().Get("from"), 10, 0);
	//if err != nil {
	//	from = 0
	//}
	//
	//to, err := strconv.ParseInt(c.Request().URL.Query().Get("to"), 10, 0);
	//if err != nil {
	//	to = 99999
	//}

	v := models.SmlEpgList{Code: 0}

	channelId, _ := strconv.ParseInt(c.Request().URL.Query().Get("channelId"), 10, 0)
	channel := &models.Channel{}
	err = models.DB.Eager("ChannelsPackage").Find(channel, int(channelId))
	if err != nil {
		return errors.WithStack(err)
	}

	epg_channel := models.EpgChannel{}
	urlHash := channel.ChannelsPackage.GetEpgUrlHash()
	if len(channel.EpgID.String) > 0 {
		err = models.EPG_DB_RO.Q().Where("url_hash = ? AND epg_id = ?", uint32(urlHash), channel.EpgID).First(&epg_channel)
	} else {
		err = models.EPG_DB_RO.Q().Where("url_hash = ? AND display_name = ?", uint32(urlHash), channel.Name).First(&epg_channel)
	}

	epgProgrammes := &models.EpgProgrammes{}
	if err == nil {
		err = models.EPG_DB_RO.Q().Where("url_hash = ? AND channel_epg_id = ? AND start >= ?", urlHash, epg_channel.EpgID, startUnixDt).Order("start asc").All(epgProgrammes)
		if err != nil {
			return errors.WithStack(errors.New("Can't load epg programmes!"))
		}
		v.EpgItem.Version = time.Now().Unix()
		v.EpgItem.ChannelId = channel.ID
		v.EpgItem.LocationId = locationId
		v.EpgItem.SubLocationId = locationId
		for _, element := range *epgProgrammes {
			elementId := (element.Start-epgProgrammeStartDiff)*maxEpgChannelId + int64(epg_channel.ID)
			epgItemP := models.SmlEpgItemP{
				Sdate:       element.Start,
				Fdate:       element.Stop,
				C_id:        channelId,
				Id:          elementId,
				Name:        element.Title,
				Discr:       "EPG",
				Desc:        element.Desc.String,
				Clogo:       "create/" + url.PathEscape(element.Title),
				Al:          0,
				Ir:          0,
				S_id:        elementId,
				Eid:         elementId,
				TstvAllowed: 1,
				PlAllowed:   1,
			}
			if len(channel.ChannelsPackage.GoogleKey.String) > 0 && len(channel.ChannelsPackage.GoogleCx.String) > 0 {
				epgItemP.Clogo = epgItemP.Clogo + "/" + url.PathEscape(channel.ChannelsPackage.GoogleKey.String) + "/" + url.PathEscape(channel.ChannelsPackage.GoogleCx.String)
			}
			v.EpgItem.EpgItemPArr = append(v.EpgItem.EpgItemPArr, epgItemP)
		}
	}

	return c.Render(http.StatusOK, render.XML(v))
}
