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
		Status: "success",
		Data: models.Data{
			Jadwal: []models.Hari{
				{Hari: "Senin", Kelas: search, MataKuliah: []models.MataKuliah{}},
				{Hari: "Selasa", Kelas: search, MataKuliah: []models.MataKuliah{}},
				{Hari: "Rabu", Kelas: search, MataKuliah: []models.MataKuliah{}},
				{Hari: "Kamis", Kelas: search, MataKuliah: []models.MataKuliah{}},
				{Hari: "Jum'at", Kelas: search, MataKuliah: []models.MataKuliah{}},
				{Hari: "Sabtu", Kelas: search, MataKuliah: []models.MataKuliah{}},
			},
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

			for i, h := range jadwal.Data.Jadwal {
				if h.Hari == hari {
					jadwal.Data.Jadwal[i].MataKuliah = append(jadwal.Data.Jadwal[i].MataKuliah, mataKuliah)
					return
				}
			}

			jadwal.Data.Jadwal = append(jadwal.Data.Jadwal, models.Hari{Hari: hari, Kelas: search, MataKuliah: []models.MataKuliah{mataKuliah}})
		})
	})

	j := 0
	for _, h := range jadwal.Data.Jadwal {
		if len(h.MataKuliah) > 0 {
			jadwal.Data.Jadwal[j] = h
			j++
		}
	}
	jadwal.Data.Jadwal = jadwal.Data.Jadwal[:j]

	return jadwal, nil
}
