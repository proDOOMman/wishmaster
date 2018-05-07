package actions

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

func SmlApiV0Init(app *buffalo.App) {
	app.GET("/networkConfig/config.cfg", SmlNetworkConfig)

	app.GET("/networkConfig/SML-482/config.cfg", SmlNetworkConfig)

	app.GET("/subscribe/{mac}/{localId}/{id}/{all}", SmlSubscribe)
}

func SmlNetworkConfig(c buffalo.Context) error {
	hostname := c.Request().Host
	if strings.Contains(hostname, ":") == false {
		hostname = hostname + ":80"
	}
	c.Set("hostname", hostname)
	return c.Render(http.StatusOK, render.String(`[Connection]
    insecureHttpPort = 80
    sdpAddress = http://<%= hostname %>/spine/services/
    sdpCacheAddress = http://<%= hostname %>/CacheClient/
    backupDHCPSettigsServer = http://<%= hostname %>/network.xml
[Server]
    hdImagePrefixBase = images/hd/
    sdImagePrefixBase = images/sd/
    channelsLogoPrefix = channel/
    videoServerProtocol = hls,rtsp
[Feature]
    weatherEnabled = false
    dlnaAllowed = true
[Vk]
    enabled=true
    extendedMode = true]`))
}

func SmlSubscribe(c buffalo.Context) error {
	c.Response().WriteHeader(http.StatusOK)
	select {
	case <-time.After(600e9):
		log.Println("Timeout")
	}
	return nil
}
