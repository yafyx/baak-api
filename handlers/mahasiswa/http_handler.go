package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"baak-api/models"

	"github.com/PuerkitoBio/goquery"
)

func GetMahasiswa(baseURL string) ([]models.Mahasiswa, error) {
	var mahasiswas []models.Mahasiswa
	page := 1

	for {
		url := fmt.Sprintf("%s&page=%d", baseURL, page)
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return nil, fmt.Errorf("error code status: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, err
		}

		doc.Find("table").First().Find("tr").Each(func(_ int, row *goquery.Selection) {
			cells := row.Find("td")
			if cells.Length() == 5 {
				mahasiswa := models.Mahasiswa{
					NPM:       strings.TrimSpace(cells.Eq(1).Text()),
					Nama:      strings.TrimSpace(cells.Eq(2).Text()),
					KelasLama: strings.TrimSpace(cells.Eq(3).Text()),
					KelasBaru: strings.TrimSpace(cells.Eq(4).Text()),
				}
				mahasiswas = append(mahasiswas, mahasiswa)
			}
		})

		if doc.Find(`a[rel="next"]`).Length() == 0 {
			break
		}

		page++
	}

	return mahasiswas, nil
}
