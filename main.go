package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	BaseURL  = "https://baak.gunadarma.ac.id"
	HTTPPort = 8080
)

type Jadwal struct {
	Senin  interface{} `json:"senin"`
	Selasa interface{} `json:"selasa"`
	Rabu   interface{} `json:"rabu"`
	Kamis  interface{} `json:"kamis"`
	Jumat  interface{} `json:"jumat"`
	Sabtu  interface{} `json:"sabtu"`
}

type MataKuliah struct {
	Nama  string `json:"nama"`
	Waktu string `json:"waktu"`
	Ruang string `json:"ruang"`
	Dosen string `json:"dosen"`
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

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type UTS struct {
	Nama  string `json:"nama"`
	Waktu string `json:"waktu"`
	Ruang string `json:"ruang"`
	Dosen string `json:"dosen"`
}

func GetJadwal(url string, search string) (map[string]interface{}, error) {
	doc, err := FetchDocument(url)
	if err != nil {
		return nil, err
	}

	jadwal := Jadwal{
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
			mataKuliah := MataKuliah{
				Nama:  strings.TrimSpace(sel.Eq(2).Text()),
				Waktu: strings.TrimSpace(sel.Eq(3).Text()),
				Ruang: strings.TrimSpace(sel.Eq(4).Text()),
				Dosen: strings.TrimSpace(sel.Eq(5).Text()),
			}

			switch hari {
			case "Senin":
				if jadwal.Senin == "" {
					jadwal.Senin = []MataKuliah{}
				}
				jadwal.Senin = append(jadwal.Senin.([]MataKuliah), mataKuliah)
			case "Selasa":
				if jadwal.Selasa == "" {
					jadwal.Selasa = []MataKuliah{}
				}
				jadwal.Selasa = append(jadwal.Selasa.([]MataKuliah), mataKuliah)
			case "Rabu":
				if jadwal.Rabu == "" {
					jadwal.Rabu = []MataKuliah{}
				}
				jadwal.Rabu = append(jadwal.Rabu.([]MataKuliah), mataKuliah)
			case "Kamis":
				if jadwal.Kamis == "" {
					jadwal.Kamis = []MataKuliah{}
				}
				jadwal.Kamis = append(jadwal.Kamis.([]MataKuliah), mataKuliah)
			case "Jum'at":
				if jadwal.Jumat == "" {
					jadwal.Jumat = []MataKuliah{}
				}
				jadwal.Jumat = append(jadwal.Jumat.([]MataKuliah), mataKuliah)
			case "Sabtu":
				if jadwal.Sabtu == "" {
					jadwal.Sabtu = []MataKuliah{}
				}
				jadwal.Sabtu = append(jadwal.Sabtu.([]MataKuliah), mataKuliah)
			}
		})
	})

	return map[string]interface{}{
		"kelas":  search,
		"jadwal": jadwal,
	}, nil
}

func GetKegiatan(url string) ([]Kegiatan, error) {
	doc, err := FetchDocument(url)
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

func GetMahasiswa(baseURL string) ([]Mahasiswa, error) {
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

func GetUTS(url string, search string) ([]UTS, error) {
	doc, err := FetchDocument(url)
	if err != nil {
		return nil, err
	}

	var utsList []UTS
	doc.Find("table").First().Find("tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() == 5 {
			uts := UTS{
				Nama:  strings.TrimSpace(cells.Eq(1).Text()),
				Waktu: strings.TrimSpace(cells.Eq(2).Text()),
				Ruang: strings.TrimSpace(cells.Eq(3).Text()),
				Dosen: strings.TrimSpace(cells.Eq(4).Text()),
			}
			utsList = append(utsList, uts)
		}
	})
	return utsList, nil
}

const BASE_URL = "https://baak.gunadarma.ac.id"

func HandlerHomepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "/jadwal/{cari berdasarkan kelas/dosen}, /kalender, /kelasbaru/{cari berdasarkan kelas/npm/nama}, /uts/{cari berdasarkan kelas/dosen}")
}

func HandlerJadwal(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 3 || segments[2] == "" {
		http.Error(w, "Missing kelas in URL", http.StatusBadRequest)
		return
	}

	search := segments[2]
	url := BASE_URL + "/jadwal/cariJadKul?&teks=" + search
	jadwal, err := GetJadwal(url, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONResponse(w, jadwal)
}

func HandlerKegiatan(w http.ResponseWriter, r *http.Request) {
	kegiatanList, err := GetKegiatan(BASE_URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONResponse(w, kegiatanList)
}

func HandlerMahasiswa(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 3 || segments[2] == "" {
		http.Error(w, "Missing search term in URL", http.StatusBadRequest)
		return
	}

	searchTerm := segments[2]
	searchTypes := []string{"Kelas", "NPM", "Nama"}
	var mahasiswa []Mahasiswa
	var err error

	for _, searchType := range searchTypes {
		baseURL := BASE_URL + "/cariKelasBaru?tipeKelasBaru=" + searchType + "&teks=" + searchTerm
		mahasiswa, err = GetMahasiswa(baseURL)
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

	WriteJSONResponse(w, mahasiswa)
}

func HandlerUTS(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 3 || segments[2] == "" {
		http.Error(w, "Missing search term in URL", http.StatusBadRequest)
		return
	}

	search := segments[2]
	url := BASE_URL + "/jadwal/cariUts?&teks=" + search
	uts, err := GetUTS(url, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONResponse(w, uts)
}

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func FetchDocument(url string) (*goquery.Document, error) {
	res, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error code status: %d %s", res.StatusCode, res.Status)
	}

	return goquery.NewDocumentFromReader(res.Body)
}

func WriteJSONResponse(w http.ResponseWriter, data interface{}) {
	response := Response{
		Status: "success",
		Data:   data,
	}

	dataJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(dataJson)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	case path == "/":
		HandlerHomepage(w, r)
	case strings.HasPrefix(path, "/jadwal/"):
		HandlerJadwal(w, r)
	case strings.HasPrefix(path, "/kalender"):
		HandlerKegiatan(w, r)
	case strings.HasPrefix(path, "/kelasbaru/"):
		HandlerMahasiswa(w, r)
	case strings.HasPrefix(path, "/uts/"):
		HandlerUTS(w, r)
	default:
		http.Error(w, "404 not found.", http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/", HandlerHomepage)
	http.HandleFunc("/jadwal/", HandlerJadwal)
	http.HandleFunc("/kalender", HandlerKegiatan)
	http.HandleFunc("/kelasbaru/", HandlerMahasiswa)
	http.HandleFunc("/uts/", HandlerUTS)

	port := HTTPPort
	log.Printf("Server running on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
