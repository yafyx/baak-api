package handlers

import (
	"strings"

	"baak-api/models"
	"baak-api/utils"

	"github.com/PuerkitoBio/goquery"
)

func GetKegiatan(url string) ([]models.Kegiatan, error) {
	doc, err := utils.FetchDocument(url)
	if err != nil {
		return nil, err
	}

	var kegiatanList []models.Kegiatan
	doc.Find("table").First().Find("tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() == 2 {
			kegiatan := models.Kegiatan{
				Kegiatan: strings.TrimSpace(cells.Eq(0).Text()),
				Tanggal:  strings.TrimSpace(cells.Eq(1).Text()),
			}
			kegiatanList = append(kegiatanList, kegiatan)
		}
	})
	return kegiatanList, nil
}
