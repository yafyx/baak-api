package handlers

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yafyx/baak-api/models"
	"github.com/yafyx/baak-api/utils"
)

func GetJadwal(url string, search string) (map[string]interface{}, error) {
	doc, err := utils.FetchDocument(url)
	if err != nil {
		return nil, err
	}

	jadwal := map[string]interface{}{
		"senin":  "",
		"selasa": "",
		"rabu":   "",
		"kamis":  "",
		"jumat":  "",
		"sabtu":  "",
	}

	doc.Find("table").First().Each(func(_ int, table *goquery.Selection) {
		table.Find("tr").Each(func(i int, row *goquery.Selection) {
			sel := row.Find("td")
			if sel.Length() == 0 {
				return
			}

			hari := strings.TrimSpace(sel.Eq(1).Text())
			mataKuliah := models.MataKuliah{
				Nama:  strings.TrimSpace(sel.Eq(2).Text()),
				Waktu: strings.TrimSpace(sel.Eq(3).Text()),
				Ruang: strings.TrimSpace(sel.Eq(4).Text()),
				Dosen: strings.TrimSpace(sel.Eq(5).Text()),
			}

			switch hari {
			case "Senin":
				if jadwal["senin"] == "" {
					jadwal["senin"] = []models.MataKuliah{}
				}
				jadwal["senin"] = append(jadwal["senin"].([]models.MataKuliah), mataKuliah)
			case "Selasa":
				if jadwal["selasa"] == "" {
					jadwal["selasa"] = []models.MataKuliah{}
				}
				jadwal["selasa"] = append(jadwal["selasa"].([]models.MataKuliah), mataKuliah)
			case "Rabu":
				if jadwal["rabu"] == "" {
					jadwal["rabu"] = []models.MataKuliah{}
				}
				jadwal["rabu"] = append(jadwal["rabu"].([]models.MataKuliah), mataKuliah)
			case "Kamis":
				if jadwal["kamis"] == "" {
					jadwal["kamis"] = []models.MataKuliah{}
				}
				jadwal["kamis"] = append(jadwal["kamis"].([]models.MataKuliah), mataKuliah)
			case "Jum'at":
				if jadwal["jumat"] == "" {
					jadwal["jumat"] = []models.MataKuliah{}
				}
				jadwal["jumat"] = append(jadwal["jumat"].([]models.MataKuliah), mataKuliah)
			case "Sabtu":
				if jadwal["sabtu"] == "" {
					jadwal["sabtu"] = []models.MataKuliah{}
				}
				jadwal["sabtu"] = append(jadwal["sabtu"].([]models.MataKuliah), mataKuliah)
			}
		})
	})

	return map[string]interface{}{
		"kelas":  search,
		"jadwal": jadwal,
	}, nil
}
