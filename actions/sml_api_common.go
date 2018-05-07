package actions

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/image/font"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/render"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"golang.org/x/image/math/fixed"
)

var workingDirectory string

const locationId = 1
const currencyInternalId = 1
const productOfferId = 1
const casId = ""
const providerId = 1
const providerExtId = 1
const pin = "0000"
const startDate = 1500000000
const endDate = 1600000000

var stbFunctions = []int{8000948}

func SmlApiInit(app *buffalo.App) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	workingDirectory = filepath.Dir(ex)

	SmlApiV0Init(app)
	SmlApiSpineInit(app)
	SmlApiCacheClientInit(app)

	SmlApiV1Init(app)

	app.GET("/license.txt", SmlLicense)

	app.GET("/checknetwork", SmlCheckNetwork)

	app.GET("/extvod/qml_moyo_image.data", SmlMoyoImageData)

	app.GET("/onion/check_version_ext", SmlCheckVersion)
	app.GET("/firmware.bin", SmlFirmware)

	app.GET("/images/c{width:\\d+}x{height:\\d+}/create/{text}", SmlCreateLogo)

	app.GET("/images/c{width:\\d+}x{height:\\d+}/create/{text}/{key}/{cx}", SmlCreateLogo)

	app.GET("/images/c{width:\\d+}x{height:\\d+}/epg/{epg_id:\\d+}/images/{fname}", SmlChannelLogo)

	app.Middleware.Skip(middleware.ParameterLogger, SmlChannelLogo)

	app.ErrorHandlers[404] = func(status int, err error, c buffalo.Context) error {
		if app.Env == "development" {
			c.Logger().Info("Not found: " + c.Request().URL.String())
		}
		if strings.HasPrefix(c.Request().URL.Path, "/api/v1") {
			return c.Render(http.StatusOK, render.String("{\"result\":{\"list\":[]}}"))
		} else {
			return c.Render(http.StatusOK, render.String("<?xml version=\"1.0\" encoding=\"UTF-8\"?><rows code=\"0\"/>"))
		}
	}
}

func SmlMoyoImageData(c buffalo.Context) error {
	c.Response().WriteHeader(http.StatusNotFound)
	return nil
}

func SmlCheckNetwork(c buffalo.Context) error {
	return c.Render(http.StatusOK, render.String("smartlabs"))
}

func SmlFirmware(c buffalo.Context) error {
	c.Response().WriteHeader(http.StatusOK)
	return nil
}

func SmlCheckVersion(c buffalo.Context) error {
	currentFirmwareVersion := c.Request().URL.Query().Get("current_firmware_version")
	return c.Render(http.StatusOK, render.String(currentFirmwareVersion+", Wishmaster "+currentFirmwareVersion+" (-),http://"+c.Request().Host+"/firmware.bin"))
}

func SmlLicense(c buffalo.Context) error {
	return c.Render(http.StatusOK, render.String("AS IS"))
}

func SmlCreateLogo(c buffalo.Context) error {
	title, err := url.PathUnescape(c.Param("text"))
	if err != nil {
		return errors.WithStack(err)
	}

	width, err := strconv.ParseInt(c.Param("width"), 10, 0)
	if err != nil {
		return errors.WithStack(err)
	}
	height, err := strconv.ParseInt(c.Param("height"), 10, 0)
	if err != nil {
		return errors.WithStack(err)
	}

	googleCx := c.Param("cx")
	googleKey := c.Param("key")
	if len(googleCx) > 0 && len(googleKey) > 0 {
		c.LogField("GOOGLE_CX", googleCx)

		type GoogleSearchItem struct {
			Link string `json:"link"`
		}
		type GoogleSearchItems struct {
			Items []GoogleSearchItem `json:"items"`
		}

		u, err := url.Parse("https://www.googleapis.com/customsearch/v1?fileType=jpg%2Cpng%2Cgif&imgSize=medium&num=1&safe=medium&searchType=image&fields=items(link)")
		if err == nil {
			query := u.Query()
			query.Set("cx", googleCx)
			query.Set("key", googleKey)
			query.Set("q", title)
			u.RawQuery = query.Encode()

			request, err := http.Get(u.String())
			if err == nil && request.StatusCode == http.StatusOK {
				defer request.Body.Close()
				content, err := ioutil.ReadAll(request.Body)
				if err == nil {
					googleSearchResult := GoogleSearchItems{}
					err = json.Unmarshal(content, &googleSearchResult)
					if err == nil && len(googleSearchResult.Items) > 0 {
						imageRequest, err := http.Get(googleSearchResult.Items[0].Link)
						if err == nil && imageRequest.StatusCode == http.StatusOK {
							defer imageRequest.Body.Close()

							img, _, err := image.Decode(imageRequest.Body)
							if err == nil {
								m := resize.Thumbnail(uint(width), uint(height), img, resize.Lanczos3)

								var b bytes.Buffer
								writer := bufio.NewWriter(&b)
								err = png.Encode(writer, m)
								if err == nil {
									err = writer.Flush()
									if err == nil {
										c.Response().WriteHeader(http.StatusOK)
										c.Response().Write(b.Bytes())
									}
								}
							}
						}
					}
				}
			}
		}
		if err != nil {
			c.Logger().Debug(err)
		}
	}

	fg, bg := image.White, image.Black
	dpi := float64(72.0)
	imgW := int(width)
	imgH := int(height)
	size := float64(imgH) * 96 / dpi / 2
	size = math.Min(size, 50)

	// Read the font data.
	fontBytes := assetsBox.Bytes("fonts/BebasNeue Bold.ttf")
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return errors.WithStack(err)
	}

	rgba := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	// Draw the text.
	h := font.HintingFull
	d := &font.Drawer{
		Dst: rgba,
		Src: fg,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     dpi,
			Hinting: h,
		}),
	}
	fontSize := int(size) * int(dpi) / 96

	if d.MeasureString(title) >= fixed.I(imgW) {
		for d.MeasureString(title + "…") >= fixed.I(imgW) && len(title) > 1 {
			runes := []rune(title)
			title = string(runes[:len(runes)-1])
		}
		title = title + "…"
	}

	d.Dot = fixed.Point26_6{
		X: (fixed.I(imgW) - d.MeasureString(title)) / 2,
		Y: fixed.I(fontSize + (imgH-fontSize)/2),
	}
	d.DrawString(title)

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err = png.Encode(writer, rgba)
	if err != nil {
		return errors.WithStack(err)
	}
	err = writer.Flush()
	if err != nil {
		return errors.WithStack(err)
	}
	c.Response().WriteHeader(http.StatusOK)
	c.Response().Write(b.Bytes())
	return nil
}

func SmlChannelLogo(c buffalo.Context) error {
	width, err := strconv.ParseInt(c.Param("width"), 10, 0)
	if err != nil {
		return errors.WithStack(err)
	}
	height, err := strconv.ParseInt(c.Param("height"), 10, 0)
	if err != nil {
		return errors.WithStack(err)
	}

	redirectTo := strings.TrimPrefix(c.Request().URL.String(), "/images/c"+c.Param("width")+"x"+c.Param("height"))

	file, err := os.Open(filepath.Join(workingDirectory, filepath.FromSlash(redirectTo)))
	if err != nil {
		return errors.WithStack(err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return errors.WithStack(err)
	}
	file.Close()

	m := resize.Thumbnail(uint(width), uint(height), img, resize.Lanczos3)

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err = png.Encode(writer, m)
	if err != nil {
		return errors.WithStack(err)
	}
	err = writer.Flush()
	if err != nil {
		return errors.WithStack(err)
	}
	c.Response().WriteHeader(http.StatusOK)
	c.Response().Write(b.Bytes())
	return nil
}
