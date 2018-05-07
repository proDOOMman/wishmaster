package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/rakanalh/scheduler"
	"github.com/rakanalh/scheduler/storage"
	"gitlab.com/SML-482HD/wishmaster/actions"
	"gitlab.com/SML-482HD/wishmaster/models"
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	app := actions.App()

	memStorage := storage.NewMemoryStorage()

	sched := scheduler.New(memStorage)

	sched.RunEvery(4*time.Hour, DownloadAndParseM3U)

	models.EpgInit(sched)

	err := sched.Start()
	if err != nil {
		log.Println(err)
	}

	if err = app.Serve(); err != nil {
		log.Fatal(err)
	}
}

func DownloadAndParseM3U() {
	log.Printf("Download and parse M3U...")
	err := models.DB.Transaction(
		func(tx *pop.Connection) error {
			channelPackages := []models.ChannelsPackage{}
			models.DB.Q().All(&channelPackages)
			for _, channelsPackage := range channelPackages {
				channelsPackage.UpdateChannelsFromUrl(tx)
			}
			return nil
		})
	if err != nil {
		log.Println(err.Error())
	}
}
