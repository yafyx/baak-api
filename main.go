package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Jadwal struct {
	Search string `json:"search"`
	Jadwal []Hari `json:"jadwal"`
}

type Hari struct {
	Hari string `json:"hari"`
	Data []Data `json:"data"`
}

type Data struct {
	Kelas      string `json:"kelas"`
	MataKuliah string `json:"mata_kuliah"`
	Waktu      string `json:"waktu"`
	Ruang      string `json:"ruang"`
	Dosen      string `json:"dosen"`
}

type Kegiatan struct {
	Kegiatan string `json:"kegiatan"`
	Tanggal  string `json:"tanggal"`
}

type Mahasiswa struct {
	NPM       string `json:"npm"`
	Nama      string `json:"nama"`
	KelasLama string `json:"kelas_lama"`
	KelasBaru string `json:"kelas_baru"`
}

func fetchDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error code status: %d %s", res.StatusCode, res.Status)
	}

	return goquery.NewDocumentFromReader(res.Body)
}

func getJadwal(url string, search string) (*Jadwal, error) {
	doc, err := fetchDocument(url)
	if err != nil {
		return nil, err
	}

	jadwal := &Jadwal{
		Search: search,
		Jadwal: []Hari{
			{Hari: "Senin", Data: []Data{}},
			{Hari: "Selasa", Data: []Data{}},
			{Hari: "Rabu", Data: []Data{}},
			{Hari: "Kamis", Data: []Data{}},
			{Hari: "Jum'at", Data: []Data{}},
			{Hari: "Sabtu", Data: []Data{}},
		},
	}

	doc.Find("table").First().Each(func(_ int, table *goquery.Selection) {
		table.Find("tr").Each(func(i int, row *goquery.Selection) {
			sel := row.Find("td")
			if sel.Length() == 0 {
				return
			}

			hari := strings.TrimSpace(sel.Eq(1).Text())
			data := Data{
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

			jadwal.Jadwal = append(jadwal.Jadwal, Hari{Hari: hari, Data: []Data{data}})
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

func getKegiatan(url string) ([]Kegiatan, error) {
	doc, err := fetchDocument(url)
	if err != nil {
		return nil, err
	}

	var kegiatanList []Kegiatan
	doc.Find("table").First().Find("tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() == 2 {
			kegiatan := Kegiatan{
				Kegiatan: strings.TrimSpace(cells.Eq(0).Text()),
				Tanggal:  strings.TrimSpace(cells.Eq(1).Text()),
			}
			kegiatanList = append(kegiatanList, kegiatan)
		}
	})

	return kegiatanList, nil
}

func getMahasiswa(baseURL string) ([]Mahasiswa, error) {
	var mahasiswas []Mahasiswa
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
				mahasiswa := Mahasiswa{
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

func writeJSONResponse(w http.ResponseWriter, data interface{}) {
	dataJson, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(dataJson)
}

func handlerJadwal(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 3 {
		http.Error(w, "Missing kelas in URL", http.StatusBadRequest)
		return
	}

	search := segments[2]
	url := "https://baak.gunadarma.ac.id/jadwal/cariJadKul?&teks=" + search
	jadwal, err := getJadwal(url, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, jadwal)
}

func handlerKegiatan(w http.ResponseWriter, r *http.Request) {
	kegiatanList, err := getKegiatan("https://baak.gunadarma.ac.id/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, kegiatanList)
}

func handlerMahasiswa(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 3 {
		http.Error(w, "Missing search term in URL", http.StatusBadRequest)
		return
	}

	searchTerm := segments[2]
	searchTypes := []string{"Kelas", "NPM", "Nama"}
	var mahasiswa []Mahasiswa
	var err error

	for _, searchType := range searchTypes {
		baseURL := fmt.Sprintf("https://baak.gunadarma.ac.id/cariKelasBaru?tipeKelasBaru=%s&teks=%s", searchType, searchTerm)
		mahasiswa, err = getMahasiswa(baseURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(mahasiswa) > 0 {
			break
		}
	}

	if len(mahasiswa) == 0 {
		http.Error(w, "Mahasiswa tidak ditemukan!", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, mahasiswa)
}

func main() {
	http.HandleFunc("/jadwal/", handlerJadwal)

	http.HandleFunc("/kalender", handlerKegiatan)
	http.HandleFunc("/mahasiswa/", handlerMahasiswa)

	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)

	log.Fatal(http.ListenAndServe(port, nil))
}
