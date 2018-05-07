package models

import (
	"encoding/xml"
)

type SmlChannel struct {
	Id                      int      `xml:"bcid"`
	BcRId                   int      `xml:"bc_r_id"`
	BcName                  string   `xml:"bcname"`
	BcDesc                  string   `xml:"bcdesc""`
	Url                     string   `xml:"url"`
	HqUrl                   string   `xml:"hqUrl"`
	PipUrl                  []string `xml:"pipUrl>-,omitempty"`
	PlcUrl                  []string `xml:"plcUrl>-,omitempty"`
	BackupUrl1              string   `xml:"backupUrl1"`
	BackupUrl2              string   `xml:"backupUrl2"`
	Logo                    string   `xml:"logo"`
	Logo2                   string   `xml:"logo2"`
	Bcal                    int      `xml:"bcal"`
	Num                     int      `xml:"num"`
	BceId                   int      `xml:"bceid"`
	IsCrypted               int      `xml:"is_crypted"`
	IsDvrCrypted            int      `xml:"isDvrCrypted"`
	SoundVolume             []string `xml:"soundVolume>-,omitempty"`
	RaptorPort              []string `xml:"raptorPort>-,omitempty"`
	IsErotic                int      `xml:"isErotic"`
	StreamAspectRatio       int      `xml:"streamAspectRatio"`
	StartTimeRestrictUTCsec []string `xml:"startTimeRestrictUTCsec>-,omitempty"`
	EndTimeRestrictUTCsec   []string `xml:"endTimeRestrictUTCsec>-,omitempty"`
	Version                 int      `xml:"version"`
	ZoomRatio               float32  `xml:"zoomRatio"`
	OttURL                  string   `xml:"ottURL"`
	SmlOttURL               string   `xml:"smlOttURL"`
	OttDvr                  int      `xml:"ottDvr"`
	TstvOttURL              string   `xml:"tstvOttURL"`
	PlOttURL                []string `xml:"plOttURL>-,omitempty"`
	HttpUrl                 []string `xml:"httpUrl>-,omitempty"`
	EpgOffset               int      `xml:"epgOffset"`
	DvbtChannelName         []string `xml:"dvbtChannelName>-,omitempty"`
	Poster                  string   `xml:"poster"`
	IsOttEncrypted          int      `xml:"isOttEncrypted"`
	NPVRChannelID           string   `xml:"nPVRChannelID"`
	IsQualityMonitoring     int      `xml:"isQualityMonitoring"`
	IsTestStreamQuality     int      `xml:"isTestStreamQuality"`
	IsBarker                int      `xml:"isBarker"`
	PromoUrl                []string `xml:"promo_url>-,omitempty"`
	VideoServerProtocol     string   `xml:"videoServerProtocol"`
	Subjects                []string `xml:"subjects>-,omitempty"`
	PackageId               int      `xml:"packages>id"`
	LocId                   int      `xml:"loc>id"`
	Excl                    []string `xml:"excl>-,omitempty"`
	Accs                    []string `xml:"accs>-,omitempty"`
	StbFunctions            []string `xml:"stbFunctions>-,omitempty"`
	NetworkTypes            []string `xml:"networkTypes>-,omitempty"`
	AudioPIDs               []string `xml:"audioPIDs>-,omitempty"`
	SubtitlePIDs            []string `xml:"subtitlePIDs>-,omitempty"`
	Urls                    []string `xml:"urls>-,omitempty"`
	OttUrls                 []string `xml:"ott_urls>-,omitempty"`
	DvbUrls                 []string `xml:"dvbUrls>-,omitempty"`
}

type SmlTerminalService struct {
	Id                        int      `xml:"id"`
	PromoUrl                  []string `xml:"promo_url>-,omitempty"`
	Logo                      []string `xml:"logo>-,omitempty"`
	Logo2                     []string `xml:"logo2>-,omitempty"`
	RId                       int      `xml:"r_id"`
	Type                      string   `xml:"type"`
	Name                      string   `xml:"name"`
	Description               string   `xml:"description"`
	ChangeDate                int32    `xml:"changeDate"`
	Mandatory                 bool     `xml:"mandatory"`
	Price                     int      `xml:"price"`
	ExternalId                string   `xml:"externalId"`
	PackageId                 []string `xml:"packageId>-,omitempty"`
	AllowedPurchases          int      `xml:"allowedPurchases"`
	Unlimited                 bool     `xml:"unlimited"`
	CanBeSubscribedByUser     bool     `xml:"canBeSubscribedByUser"`
	IsInvisible               bool     `xml:"isInvisible"`
	ShowIfInvisibleAndDepends bool     `xml:"showIfInvisibleAndDepends"`
	QuotaTimeout              []string `xml:"quotaTimeout>-,omitempty"`
	QuotaLimit                []string `xml:"quotaLimit>-,omitempty"`
	Duration                  []string `xml:"duration>-,omitempty"`
	QuotaSizeLimit            int      `xml:"quotaSizeLimit"`
	EndDate                   int      `xml:"endDate"`
	SortOrder                 string   `xml:"sortOrder"`
	SsType                    []string `xml:"ssType>-,omitempty"`
	IsDaily                   bool     `xml:"isDaily"`
	NpvrPrice                 []string `xml:"npvr_price>-,omitempty"`
	PlType                    string   `xml:"plType"`
	TerminalType              int      `xml:"terminalType"`
	PlDiscriminator           string   `xml:"plDiscriminator"`
	PlId                      int      `xml:"plId"`
	NotRecomExclusive         bool     `xml:"notRecomExclusive"`
	UnsubscribeByLimit        bool     `xml:"unsubscribeByLimit"`
	IsAccessByLock            bool     `xml:"isAccessByLock"`
	ReplayDenied              bool     `xml:"replayDenied"`
	BundledGrp                []string `xml:"bundledGrp>-,omitempty"`
	SubscribeMode             int      `xml:"subscribeMode"`
	AppleId                   []string `xml:"appleId>-,omitempty"`
	AndroidId                 []string `xml:"androidId>-,omitempty"`
	SmarttvId                 []string `xml:"smarttvId>-,omitempty"`
	ChannelPreviewEnabled     bool     `xml:"channelPreviewEnabled"`
	IsBasic                   []string `xml:"isBasic"`
	TextOn                    []string `xml:"textOn"`
	TextOff                   []string `xml:"textOff"`
	ForAllAcc                 int      `xml:"forAllAcc"`
	PortalId                  int      `xml:"portalId"`
	NonBlockedPacksInBundle   int      `xml:"nonBlockedPacksInBundle"`
	SmExternalid              []string `xml:"sm_externalid"`
	CurrencyInternalId        int      `xml:"currencyInternalId"`
	BpriceGroup               []string `xml:"bpriceGroup"`
	UnsubscribePeriod         []string `xml:"unsubscribePeriod>-,omitempty"`
	RecurrentSubscribePeriod  []string `xml:"recurrentSubscribePeriod>-,omitempty"`
	RecommendOn               []string `xml:"recommendOn>-,omitempty"`
	RecommendOff              []string `xml:"recommendOff>-,omitempty"`
	Depends                   []int    `xml:"depends>id,omitempty"`
	Parent                    []int    `xml:"parent>id,omitempty"`
	Required                  []string `xml:"required>-,omitempty"`
	BundlePrices              []string `xml:"bundle_prices>-,omitempty"`
	PpvPrices                 []string `xml:"ppv_prices>-,omitempty"`
	StbFunctions              []string `xml:"stbFunctions>-,omitempty"`
	NetworkTypes              []string `xml:"networkTypes>-,omitempty"`
}

type SmlTerminalServices struct {
	XMLName          xml.Name             `xml:"services"`
	Code             int                  `xml:"code,attr"`
	Version          int                  `xml:"version,attr"`
	TerminalServices []SmlTerminalService `xml:"service"`
}

type SmlService struct {
	Id               int      `xml:"id"`
	Type             string   `xml:"type"`
	Name             string   `xml:"name"`
	Description      string   `xml:"description"`
	Mandatory        bool     `xml:"mandatory"`
	Price            int      `xml:"price"`
	ExternalId       string   `xml:"externalId"`
	PackageId        []string `xml:"packageId>-,omitempty"`
	EndDate          int      `xml:"endDate"`
	AllowedPurchases int      `xml:"allowedPurchases"`
	Unlimited        bool     `xml:"unlimited"`
	StartDate        int      `xml:"startDate"`
	OneTime          bool     `xml:"oneTime"`
	ServiceState     int      `xml:"serviceState"`
	IsPromoService   int      `xml:"isPromoService"`
	NotifyThreshold  []string `xml:"notifyThreshold>-,omitempty"`
	NotifyTimeout    []string `xml:"notifyTimeout>-,omitempty"`
	NotifyStartText  []string `xml:"notifyStartText>-,omitempty"`
	NotifyText       []string `xml:"notifyText>-,omitempty"`
	AppleId          []string `xml:"appleId>-,omitempty"`
	SmarttvId        []string `xml:"smarttvId>-,omitempty"`
	NotConfirmed     int      `xml:"notConfirmed"`
	PortalId         int      `xml:"portalId"`
	BonusId          []string `xml:"bonusId>-,omitempty"`
}

type SmlServices struct {

	XMLName  xml.Name     `xml:"services"`
	Code     int          `xml:"code,attr"`
	Services []SmlService `xml:"service"`
}

type SmlChannelPackage struct {

	PackageId     int          `xml:"channelPackageId,attr"`
	LocationId    int          `xml:"locationId,attr"`
	SubLocationId int          `xml:"subLocationId,attr"`
	Version       int          `xml:"version,attr"`
	Channels      []SmlChannel `xml:"channel"`
}

type SmlChannelsList struct {

	XMLName  xml.Name            `xml:"channels_list"`
	Code     int                 `xml:"code,attr"`
	Channels []SmlChannelPackage `xml:"channels"`
}

type SmlAuthorizeAnswer struct {
	ServiceAccountNumber string   `xml:"serviceAccountNumber"`
	DateTime             string   `xml:"datetime"`
	State                string   `xml:"state"`
	LocationId           int      `xml:"locationId"`
	SubLocationId        int      `xml:"subLocationId"`
	ActivationNumber     string   `xml:"activationNumber"`
	BackendVersion       string   `xml:"backend_version"`
	SessionId            []string `xml:"sessionId>-,omitempty"`
	ProductOfferId       int      `xml:"productOfferId"`
	CasId                string   `xml:"casId"`
	ShouldUseRaptor      int      `xml:"shouldUseRaptor"`
	NetLogIP             string   `xml:"netLogIP"`
	NetLogPort           int      `xml:"netLogPort"`
	ProviderId           int      `xml:"providerId"`
	ProviderExtId        int      `xml:"providerExtId"`
	TerminalType         int      `xml:"terminalType"`
	TimeZone             []string `xml:"time_zone>-,omitempty"`
	NetworkBlock         int      `xml:"networkBlock"`
	NetworkTypeId        []string `xml:"networkTypeId>-,omitempty"`
	HelpId               []string `xml:"helpId>-,omitempty"`
	SmLogin              string   `xml:"smLogin"`
	SmPassword           string   `xml:"smPassword"`
	Esamservice          []string `xml:"esamservice>-,omitempty"`
	Randomservice        []string `xml:"randomservice>-,omitempty"`
	AccountType          string   `xml:"accountType"`
	OfferId              []string `xm:"offerId"`
	IsOfferConfirmed     int      `xml:"isOfferConfirmed"`
	DefaultUrlSource     []string `xml:"defaultUrlSource>-,omitempty"`
	Reason               []string `xml:"reason>-,omitempty"`
	Multiroom            bool     `xml:"multiroom"`
	Unixtime             int64    `xml:"unixtime"`
	ResourceId           string   `xml:"resourceId"`
	BonusType            []string `xml:"bonusType>-,omitempty"`
	IsTrustedDevice      int      `xml:"isTrustedDevice"`
	IsOss                int      `xml:"isOss"`
	IsHomeGroup          int      `xml:"isHomeGroup"`
	Pin                  string   `xml:"pin"`
	ZabavaLogin          []string `xml:"zabavaLogin>-,omitempty"`
	DeviceType           string   `xml:"deviceType"`
	UseOttUrlForChannels int      `xml:"useOttUrlForChannels"`
	IP                   string   `xml:"IP"`
	StbFunctions         []int    `xml:"stbFunctions>id"`
	Multiscreen          int      `xml:"multiscreen"`
}

type SmlBillingRow struct {

	Number  int `xml:"number,attr"`
	Balance int `xml:"balance"`
}

type SmlBillingRows struct {

	XMLName    xml.Name      `xml:"rows"`
	Code       int           `xml:"code,attr"`
	BillingRow SmlBillingRow `xml:"row"`
}

type SmlIncompatiblePackages struct {
	XMLName              xml.Name `xml:"packages"`
	Code                 int      `xml:"code,attr"`
	IdArray              []int    `xml:"package>id"`
	IncompatiblePackages []int    `xml:"package>incompatible_packages>-,omitempty"`
}

type SmlProfileRow struct {
	XMLName                   xml.Name `xml:"row"`
	Id                        int      `xml:"id"`
	Name                      int      `xml:"name"`
	IsMaster                  int      `xml:"isMaster"`
	IsCurrent                 int      `xml:"isCurrent"`
	Favorite                  string   `xml:"favorite"`
	Forbidden                 []int    `xml:"forbidden>-,omitempty"`
	Sort                      string   `xml:"sort"`
	LastChannelId             []int    `xml:"lastChannelId>-,omitempty"`
	SpecChannelPIDs           []int    `xml:"specChannelPIDs>-,omitempty"`
	VodPositions              []int    `xml:"vodPositions>-,omitempty"`
	Reminders                 string   `xml:"reminders"`
	MaxAllow                  int      `xml:"maxAllow"`
	LastAccessLevelId         int      `xml:"lastAccessLevelId"`
	IsAccessLevelPersistent   int      `xml:"isAccessLevelPersistent"`
	CustomProperties          []int    `xml:"customProperties>-,omitempty"`
	SpecialData               []int    `xml:"specialData>-,omitempty"`
	AspectRatio               int      `xml:"aspectRatio"`
	OutputAspectRatio         int      `xml:"outputAspectRatio"`
	Style                     string   `xml:"style"`
	PurchaseVodAllow          int      `xml:"purchaseVodAllow"`
	PurchaseKaraokeAllow      int      `xml:"purchaseKaraokeAllow"`
	SubscribeAllow            int      `xml:"subscribeAllow"`
	PinEnable                 int      `xml:"pinEnable"`
	Pin                       string   `xml:"pin"`
	Autohide                  int      `xml:"autohide"`
	Autoopen                  int      `xml:"autoopen"`
	Lang                      []int    `xml:"lang>-,omitempty"`
	ViewLimFrom               []int    `xml:"viewLimFrom>-,omitempty"`
	ViewLimTo                 []int    `xml:"viewLimTo>-,omitempty"`
	ViewDurationHour          []int    `xml:"viewDurationHour>-,omitempty"`
	ForbiddenContent          []int    `xml:"forbiddenContent>-,omitempty"`
	NickName                  []int    `xml:"nickName>-,omitempty"`
	Privacy                   []int    `xml:"privacy>-,omitempty"`
	SpecChannelSubsPids       []int    `xml:"specChannelSubsPids>-,omitempty"`
	MessageBan                int      `xml:"messageBan"`
	IsPublic                  int      `xml:"isPublic"`
	IsPurchaseLimited         int      `xml:"isPurchaseLimited"`
	PurchaseLimitedPeriod     int      `xml:"purchaseLimitedPeriod"`
	PurchaseLimitedChangeDate []int    `xml:"purchaseLimitedChangeDate>-,omitempty"`
	PurchaseLimitedBalance    []int    `xml:"purchaseLimitedBalance>-,omitempty"`
	PurchaseLimitedQuota      []int    `xml:"purchaseLimitedQuota>-,omitempty"`
	PlIsLastModifBySa         int      `xml:"plIsLastModifBySa"`
}

type SmlProfileRows struct {

	XMLName    xml.Name      `xml:"rows"`
	Code       int           `xml:"code,attr"`
	ProfileRow SmlProfileRow `xml:"row"`
}

type SmlAccessLevelRow struct {
	XMLName                                      xml.Name `xml:"row"`
	Number                                       int      `xml:"number,attr"`
	AccessLevelDictionaryIsDefault               int      `xml:"AccessLevelDictionary_isDefault"`
	AccessLevelDictionaryRestrictTimeStart       []int    `xml:"AccessLevelDictionary_restrictTimeStart>-,omitempty"`
	AccessLevelDictionaryRestrictTimeEnd         []int    `xml:"AccessLevelDictionary_restrictTimeEnd"`
	AccessLevelDictionaryID                      int      `xml:"AccessLevelDictionary_ID"`
	AccessLevelDictionaryDISCRIMINATOR           string   `xml:"AccessLevelDictionary_DISCRIMINATOR"`
	AccessLevelDictionaryName                    string   `xml:"AccessLevelDictionary_name"`
	AccessLevelDictionaryDescription             string   `xml:"AccessLevelDictionary_description,omitempty"`
	AccessLevelDictionarySortOrder               int      `xml:"AccessLevelDictionary_sortOrder"`
	AccessLevelDictionaryLogo                    string   `xml:"AccessLevelDictionary_logo,omitempty"`
	AccessLevelDictionaryOwnerUser               string   `xml:"AccessLevelDictionary_ownerUser,omitempty"`
	OwnerUserLogin                               string   `xml:"ownerUser_login,omitempty"`
	OwnerUserID                                  string   `xml:"ownerUser_ID,omitempty"`
	OwnerUserDISCRIMINATOR                       string   `xml:"ownerUser_DISCRIMINATOR,omitempty"`
	AccessLevelDictionaryInternalId              string   `xml:"AccessLevelDictionary_internalId,omitempty"`
	AccessLevelDictionaryExternalId              string   `xml:"AccessLevelDictionary_externalId,omitempty"`
	AccessLevelDictionaryLatinExternalId         string   `xml:"AccessLevelDictionary_latinExternalId,omitempty"`
	AccessLevelDictionaryLangId                  int      `xml:"AccessLevelDictionary_langId"`
	LangIdName                                   string   `xml:"langId_name"`
	LangIdID                                     int      `xml:"langId_ID"`
	LangIdDISCRIMINATOR                          string   `xml:"langId_DISCRIMINATOR"`
	AccessLevelDictionaryStartTimeRestrictUTCsec string   `xml:"AccessLevelDictionary_startTimeRestrictUTCsec,omitempty"`
	AccessLevelDictionaryEndTimeRestrictUTCsec   string   `xml:"AccessLevelDictionary_endTimeRestrictUTCsec,omitempty"`
}

type SmlEpgList struct {
	XMLName xml.Name   `xml:"epg_list"`
	Code    int        `xml:"code,attr"`
	EpgItem SmlEpgItem `xml:"epg"`
}

type SmlEpgItem struct {

	XMLName       xml.Name      `xml:"epg"`
	ChannelId     int           `xml:"channelId,attr"`
	Day           string        `xml:"day,attr"`
	LocationId    int           `xml:"locationId,attr"`
	SubLocationId int           `xml:"subLocationId,attr"`
	Version       int64         `xml:"version,attr"`
	EpgItemPArr   []SmlEpgItemP `xml:"p"`
}

type SmlEpgItemP struct {
	XMLName     xml.Name `xml:"p"`
	Sdate       int64    `xml:"sdate"`
	Fdate       int64    `xml:"fdate"`
	C_id        int64    `xml:"c_id"`
	Id          int64    `xml:"id"`
	Name        string   `xml:"name"`
	Discr       string   `xml:"discr"`
	Desc        string   `xml:"desc"`
	Clogo       string   `xml:"clogo"`
	ClogoExtra  string   `xml:"clogoExtra>-,omitempty"`
	Al          int      `xml:"al"`
	AlSort      int      `xml:"alSort"`
	Ir          int      `xml:"ir"`
	S_id        int64    `xml:"s_id"`
	Eid         int64    `xml:"eid"`
	TstvAllowed int      `xml:"tstvAllowed"`
	PlAllowed   int      `xml:"plAllowed"`
	IsHD        int      `xml:"isHD"`
}

type SmlQuotaRow struct {
	ServiceAccountNumber string `xml:"serviceAccountNumber"`
	AllQuote             int    `xml:"allQuote"`
	UsedQuota            int    `xml:"usedQuota"`
}

type SmlQuotaRows struct {
	XMLName  xml.Name    `xml:"rows"`
	Code     int         `xml:"code,attr"`
	QuotaRow SmlQuotaRow `xml:"row"`
}

type SmlEpgDeltaRow struct {
	PoId    int    `xml:"poId,attr"`
	Version int64  `xml:"version,attr"`
	Type    string `xml:"type"`
	Add     string `xml:"add"`
	Del     string `xml:"del"`
}

type SmlEpgDeltaRows struct {
	XMLName     xml.Name       `xml:"delta"`
	Code        int            `xml:"code,attr"`
	EpgDeltaRow SmlEpgDeltaRow `xml:"row"`
}

type SmlEpgDescRow struct {
	Id          int    `xml:"id"`
	Description string `xml:"description"`
}

type SmlEpgDescRows struct {
	XMLName    xml.Name      `xml:"rows"`
	Code       int           `xml:"code,attr"`
	EpgDescRow SmlEpgDescRow `xml:"content"`
}

type SmlAuthorizeRows struct {
	XMLName         xml.Name           `xml:"rows"`
	Code            int                `xml:"code,attr"`
	AuthorizeAnswer SmlAuthorizeAnswer `xml:"row"`
}

type SmlV1ErrorResult struct {
	Error *SmlV1Error `json:"error,omitempty"`
}

type SmlV1Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SmlV1AuthorizeDeviceResult struct {
	Result *SmlV1AuthorizeDevice `json:"result,omitempty"`
}

type SmlV1AuthorizeDevice struct {
	DeviceType           string        `json:"deviceType"`
	ServiceAccountNumber int           `json:"serviceAccountNumber"`
	ServiceAccountState  string        `json:"serviceAccountState"`
	ProductOfferId       int           `json:"productOfferId"`
	Login                string        `json:"login"`
	Password             string        `json:"password"`
	Multiroom            bool          `json:"multiroom"`
	Unixtime             int64         `json:"unixtime"`
	Currency             SmlV1Currency `json:"currency"`
	DeviceList           []SmlV1Device `json:"deviceList"`
	IsTest               bool          `json:"isTest"`
}

type SmlV1Device struct {
	Id           int    `json:"id"`
	Uuid         string `json:"uuid"`
	DeviceType   string `json:"deviceType"`
	TerminalName string `json:"terminalName"`
	IsCurrent    bool   `json:"isCurrent"`
}

type SmlV1Currency struct {
	Name            string `json:"name"`
	CurrencyStrCode string `json:"currencyStrCode"`
	CurrencyCode    string `json:"currencyCode"`
}
