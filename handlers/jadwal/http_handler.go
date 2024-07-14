package handlers

import (
	"strings"

	"baak-api/models"
	"baak-api/utils"

	"github.com/PuerkitoBio/goquery"
)

func GetJadwal(url string, search string) (map[string]interface{}, error) {
	doc, err := utils.FetchDocument(url)
	if err != nil {
		return nil, err
	}

	jadwal := &models.Jadwal{
		Senin:  "",
		Selasa: "",
		Rabu:   "",
		Kamis:  "",
		Jumat:  "",
		Sabtu:  "",
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
				if jadwal.Senin == "" {
					jadwal.Senin = []models.MataKuliah{}
				}
				jadwal.Senin = append(jadwal.Senin.([]models.MataKuliah), mataKuliah)
			case "Selasa":
				if jadwal.Selasa == "" {
					jadwal.Selasa = []models.MataKuliah{}
				}
				jadwal.Selasa = append(jadwal.Selasa.([]models.MataKuliah), mataKuliah)
			case "Rabu":
				if jadwal.Rabu == "" {
					jadwal.Rabu = []models.MataKuliah{}
				}
				jadwal.Rabu = append(jadwal.Rabu.([]models.MataKuliah), mataKuliah)
			case "Kamis":
				if jadwal.Kamis == "" {
					jadwal.Kamis = []models.MataKuliah{}
				}
				jadwal.Kamis = append(jadwal.Kamis.([]models.MataKuliah), mataKuliah)
			case "Jum'at":
				if jadwal.Jumat == "" {
					jadwal.Jumat = []models.MataKuliah{}
				}
				jadwal.Jumat = append(jadwal.Jumat.([]models.MataKuliah), mataKuliah)
			case "Sabtu":
				if jadwal.Sabtu == "" {
					jadwal.Sabtu = []models.MataKuliah{}
				}
				jadwal.Sabtu = append(jadwal.Sabtu.([]models.MataKuliah), mataKuliah)
			}
		})
	})

	return map[string]interface{}{
		"kelas":  search,
		"jadwal": jadwal,
	}, nil
}
