package actions

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"github.com/gobuffalo/buffalo/middleware/csrf"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

type logrusWrapper struct {
	logrus.FieldLogger
}

func (l logrusWrapper) WithField(s string, i interface{}) buffalo.Logger {
	return logrusWrapper{l.FieldLogger.WithField(s, i)}
}

func (l logrusWrapper) WithFields(m map[string]interface{}) buffalo.Logger {
	return logrusWrapper{l.FieldLogger.WithFields(m)}
}

func myLogger(logLevel string) *logrus.Logger {
	l := logrus.New()
	l.Level, _ = logrus.ParseLevel(logLevel)
	l.Formatter = &logrus.TextFormatter{
		ForceColors: strings.HasPrefix(ENV, "development"),
	}
	if strings.HasSuffix(ENV, "sml482") {
		path := "/mnt/persistent/wishmaster/wishmaster.log"
		writer, err := rotatelogs.New(
			path+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithRotationTime(time.Duration(30)*time.Minute),
			rotatelogs.WithRotationCount(3),
		)
		if err != nil {
			log.Println("Can't create log file!")
		} else {
			l.Hooks.Add(lfshook.NewHook(
				lfshook.WriterMap{
					logrus.InfoLevel:  writer,
					logrus.ErrorLevel: writer,
				},
				&logrus.TextFormatter{
					ForceColors: false,
				},
			))
		}
	}
	return l
}

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		logLevel := "debug"
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_wishmaster_session",
			LooseSlash:  true,
			WorkerOff:   true,
			LogLevel:    logLevel,
			Logger:      logrusWrapper{myLogger(logLevel)},
		})

		if ENV == "development" || ENV == "development_sml482" {
			app.Use(middleware.ParameterLogger)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		app.GET("/", HomeHandler)

		app.POST("/reboot", RebootHandler)

		app.POST("/execute", ExecuteHandler)

		app.Resource("/channels_packages", ChannelsPackagesResource{})
		app.Resource("/channels", ChannelsResource{})

		app.GET("/typeahead", ChannelNameTypeahead)

		SmlApiInit(app)

		app.Mount("/d", http.DefaultServeMux)

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	if ENV == "development" || ENV == "development_sml482" {
		pop.Debug = true
	}

	return app
}
