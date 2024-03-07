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
		Search: search,
		Jadwal: []models.Hari{
			{Hari: "Senin", Data: []models.Data{}},
			{Hari: "Selasa", Data: []models.Data{}},
			{Hari: "Rabu", Data: []models.Data{}},
			{Hari: "Kamis", Data: []models.Data{}},
			{Hari: "Jum'at", Data: []models.Data{}},
			{Hari: "Sabtu", Data: []models.Data{}},
		},
	}

	doc.Find("table").First().Each(func(_ int, table *goquery.Selection) {
		table.Find("tr").Each(func(i int, row *goquery.Selection) {
			sel := row.Find("td")
			if sel.Length() == 0 {
				return
			}

			hari := strings.TrimSpace(sel.Eq(1).Text())
			data := models.Data{
				Kelas:      strings.TrimSpace(sel.Eq(0).Text()),
				MataKuliah: strings.TrimSpace(sel.Eq(2).Text()),
				Waktu:      strings.TrimSpace(sel.Eq(3).Text()),
				Ruang:      strings.TrimSpace(sel.Eq(4).Text()),
				Dosen:      strings.TrimSpace(sel.Eq(5).Text()),
			}

			for i, h := range jadwal.Jadwal {
				if h.Hari == hari {
					jadwal.Jadwal[i].Data = append(jadwal.Jadwal[i].Data, data)
					return
				}
			}

			jadwal.Jadwal = append(jadwal.Jadwal, models.Hari{Hari: hari, Data: []models.Data{data}})
		})
	})

	j := 0
	for _, h := range jadwal.Jadwal {
		if len(h.Data) > 0 {
			jadwal.Jadwal[j] = h
			j++
		}
	}
	jadwal.Jadwal = jadwal.Jadwal[:j]

	return jadwal, nil
}
