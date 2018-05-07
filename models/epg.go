package models

import (
	"compress/gzip"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/AlekSi/xmltv"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/rakanalh/scheduler"
	"github.com/pkg/errors"
)

var workingDirectory string

func EpgInit(sched scheduler.Scheduler) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	workingDirectory = filepath.Dir(ex)

	sched.RunEvery(4*time.Hour, DownloadAndParseEpg)

	sched.RunAfter(10*time.Minute, DownloadAndParseEpg)
}

func downloadFile(from string, to string) (err error) {
	err = os.MkdirAll(filepath.Dir(to), os.ModePerm)
	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(to)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(from)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return errors.New("Bad http status")
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func readXmltv(channelsPackage *ChannelsPackage, sourceFile string, forceDownload bool) {
	f, err := os.Open(sourceFile)
	if err != nil {
		log.Println("Can't open source file!")
		return
	}
	defer f.Close()

	uncompressedStream, err := gzip.NewReader(f)
	if err != nil {
		log.Println(err.Error())
		return
	}

	decoder := xml.NewDecoder(uncompressedStream)

	urlHash := channelsPackage.GetEpgUrlHash()

	imagesDirectory := filepath.Dir(sourceFile)

	oldDebug := pop.Debug
	pop.Debug = false

	err = EPG_DB.Transaction(
		func(tx *pop.Connection) error {
			minimumTime := time.Now().Add(-12 * time.Hour)
			maximumTime := time.Now().Add(48 * time.Hour)

			tx.RawQuery("DELETE FROM epg_channels WHERE url_hash = ?", urlHash).Exec()
			tx.RawQuery("DELETE FROM epg_programmes WHERE url_hash = ?", urlHash).Exec()

			for {
				t, _ := decoder.Token()
				if t == nil {
					break
				}
				switch se := t.(type) {
				case xml.StartElement:
					if se.Name.Local == "channel" {
						var c xmltv.Channel
						decoder.DecodeElement(&c, &se)

						var logo string
						if len(c.Icon.Src) > 0 {
							var imageFilePath string
							logo = filepath.ToSlash("images/" + c.Id + "_" + path.Base(c.Icon.Src))
							imageFilePath = filepath.Join(imagesDirectory, filepath.FromSlash(logo))
							_, err = os.Stat(imageFilePath)
							if forceDownload || os.IsNotExist(err) {
								downloadFile(c.Icon.Src, imageFilePath)
							}
						}

						for _, channelName := range c.DisplayNames {
							epgChannel := EpgChannel{}
							tx.Where("display_name = ? AND epg_id = ? and url_hash = ?", channelName, c.Id, urlHash).First(&epgChannel)
							epgChannel.DisplayName = channelName
							epgChannel.EpgID =      c.Id
							epgChannel.UrlHash =     urlHash
							if len(logo) > 0 {
								epgChannel.IconSrc = nulls.NewString("epg/" + channelsPackage.GetEpgUrlHashString() + "/" + logo)
							} else {
								epgChannel.IconSrc = nulls.NewString("")
							}
							verrs, err := tx.ValidateAndSave(&epgChannel)
							if err != nil || verrs.HasAny() {
								return err
							}
						}
					} else if se.Name.Local == "programme" {
						var p xmltv.Programme
						decoder.DecodeElement(&p, &se)
						epgProgramme := EpgProgramme{
							UrlHash:      urlHash,
							ProgrammeID:  nulls.NewString(p.Id),
							Desc:         nulls.NewString(p.Desc),
							ChannelEpgID: p.ChannelId,
							Title:        p.Title,
							Start:        p.Start.Time.Unix(),
							Stop:         p.Stop.Time.Unix(),
							Categories:   nulls.NewString(""),
						}
						if p.Stop.Time.After(minimumTime) && p.Start.Time.Before(maximumTime) {
							verrs, err := tx.ValidateAndCreate(&epgProgramme)
							if err != nil {
								return err
							}
							if verrs.HasAny() {
								return errors.New(verrs.Error())
							}
						}
					}
				}
			}

//			tx.RawQuery("VACUUM").Exec()
			return nil
		})

	pop.Debug = oldDebug

	if err != nil {
		log.Println(err.Error())
	}

	log.Println("EPG parse finished")
	return
}

func DownloadAndParseEpg() {
	downloadAndParseEpg(true)
}

func downloadAndParseEpg(forceDownload bool) {
	log.Println("Downloading EPG...")

	EPG_DB.RawQuery("DELETE FROM epg_programmes WHERE stop < ?", time.Now().Add(-12*time.Hour).Unix()).Exec()

	channelsPackages := &ChannelsPackages{}
	if err := DB.Q().All(channelsPackages); err != nil {
		log.Println("Can't get channels packages from database!")
		return
	}

	var downloadedUrls []string
	for _, element := range *channelsPackages {
		isUrlDownloaded := false
		for _, downloadedUrl := range downloadedUrls {
			if downloadedUrl == element.XmltvUrl.String {
				isUrlDownloaded = false
			}
		}
		if isUrlDownloaded {
			continue
		}
		UpdateChannelPackageEpg(&element, forceDownload)
		downloadedUrls = append(downloadedUrls, element.XmltvUrl.String)
	}
}

func UpdateChannelPackageEpg(element *ChannelsPackage, forceDownload bool) {
	epgUrl := element.XmltvUrl.String
	if len(epgUrl) == 0 {
		return
	}
	log.Println("Updating EPG: " + epgUrl)

	sourceFile := filepath.Join(workingDirectory, "epg", element.GetEpgUrlHashString(), "xmltv.xml.gz")

	fileInfo, fileErr := os.Stat(sourceFile)
	if fileErr != nil && !os.IsNotExist(fileErr) {
		log.Println(fileErr)
		return
	}
	if forceDownload || (fileErr != nil && os.IsNotExist(fileErr)) || time.Since(fileInfo.ModTime()).Hours() > 12 {
		err := downloadFile(epgUrl, sourceFile)
		if err != nil {
			log.Println("Can't download file from url!")
			log.Println(err)
			return
		}
	}
	readXmltv(element, sourceFile, forceDownload)
}
