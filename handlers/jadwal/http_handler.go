// handlers/http_handler.go
package handlers

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yafyx/baak-api/models"
	"github.com/yafyx/baak-api/utils"
)

func GetJadwal(url string, search string) (*models.Jadwal, error) {
	doc, err := utils.FetchDocument(url)
	if err != nil {
		return nil, err
	}

	jadwal := &models.Jadwal{
		Jadwal: []models.Hari{
			{Hari: "Senin", Kelas: search, MataKuliah: []models.MataKuliah{}},
			{Hari: "Selasa", Kelas: search, MataKuliah: []models.MataKuliah{}},
			{Hari: "Rabu", Kelas: search, MataKuliah: []models.MataKuliah{}},
			{Hari: "Kamis", Kelas: search, MataKuliah: []models.MataKuliah{}},
			{Hari: "Jum'at", Kelas: search, MataKuliah: []models.MataKuliah{}},
			{Hari: "Sabtu", Kelas: search, MataKuliah: []models.MataKuliah{}},
		},
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

			for i, h := range jadwal.Jadwal {
				if h.Hari == hari {
					jadwal.Jadwal[i].MataKuliah = append(jadwal.Jadwal[i].MataKuliah, mataKuliah)
					return
				}
			}

			jadwal.Jadwal = append(jadwal.Jadwal, models.Hari{Hari: hari, Kelas: search, MataKuliah: []models.MataKuliah{mataKuliah}})
		})
	})

	j := 0
	for _, h := range jadwal.Jadwal {
		if len(h.MataKuliah) > 0 {
			jadwal.Jadwal[j] = h
			j++
		}
	}
	jadwal.Jadwal = jadwal.Jadwal[:j]

	return jadwal, nil
}
