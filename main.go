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

// Constants
const (
	BaseURL  = "https://baak.gunadarma.ac.id"
	HTTPPort = 8080
)

// Structs
type (
	Jadwal struct {
		Senin  []MataKuliah `json:"senin"`
		Selasa []MataKuliah `json:"selasa"`
		Rabu   []MataKuliah `json:"rabu"`
		Kamis  []MataKuliah `json:"kamis"`
		Jumat  []MataKuliah `json:"jumat"`
		Sabtu  []MataKuliah `json:"sabtu"`
	}

	MataKuliah struct {
		Nama  string `json:"nama"`
		Waktu string `json:"waktu"`
		Ruang string `json:"ruang"`
		Dosen string `json:"dosen"`
	}

	Kegiatan struct {
		Kegiatan string `json:"kegiatan"`
		Tanggal  string `json:"tanggal"`
	}

	Mahasiswa struct {
		NPM       string `json:"npm"`
		Nama      string `json:"nama"`
		KelasLama string `json:"kelas_lama"`
		KelasBaru string `json:"kelas_baru"`
	}

	Response struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}

	UTS struct {
		Nama  string `json:"nama"`
		Waktu string `json:"waktu"`
		Ruang string `json:"ruang"`
		Dosen string `json:"dosen"`
	}
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// Handlers
func HandlerHomepage(w http.ResponseWriter, r *http.Request) {
	endpoints := []string{
		"/jadwal/{kelas}",
		"/kalender",
		"/kelasbaru/{kelas/npm/nama}",
		"/uts/{kelas/dosen}",
	}
	WriteJSONResponse(w, endpoints)
}

func HandlerJadwal(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimPrefix(r.URL.Path, "/jadwal/")
	if search == "" {
		http.Error(w, "Missing kelas in URL", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%s/jadwal/cariJadKul?&teks=%s", BaseURL, search)
	jadwal, err := GetJadwal(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONResponse(w, map[string]interface{}{"kelas": search, "jadwal": jadwal})
}

func HandlerKegiatan(w http.ResponseWriter, r *http.Request) {
	kegiatanList, err := GetKegiatan(BaseURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONResponse(w, kegiatanList)
}

func HandlerMahasiswa(w http.ResponseWriter, r *http.Request) {
	searchTerm := strings.TrimPrefix(r.URL.Path, "/kelasbaru/")
	if searchTerm == "" {
		http.Error(w, "Missing search term in URL", http.StatusBadRequest)
		return
	}

	searchTypes := []string{"Kelas", "NPM", "Nama"}
	var mahasiswas []Mahasiswa
	var err error

	for _, searchType := range searchTypes {
		url := fmt.Sprintf("%s/cariKelasBaru?tipeKelasBaru=%s&teks=%s", BaseURL, searchType, searchTerm)
		mahasiswas, err = GetMahasiswa(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(mahasiswas) > 0 {
			break
		}
	}

	if len(mahasiswas) == 0 {
		http.Error(w, "Mahasiswa tidak ditemukan!", http.StatusNotFound)
		return
	}

	WriteJSONResponse(w, mahasiswas)
}

func HandlerUTS(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimPrefix(r.URL.Path, "/uts/")
	if search == "" {
		http.Error(w, "Missing search term in URL", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%s/jadwal/cariUts?&teks=%s", BaseURL, search)
	uts, err := GetUTS(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONResponse(w, uts)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		HandlerHomepage(w, r)
	case strings.HasPrefix(r.URL.Path, "/jadwal/"):
		HandlerJadwal(w, r)
	case r.URL.Path == "/kalender":
		HandlerKegiatan(w, r)
	case strings.HasPrefix(r.URL.Path, "/kelasbaru/"):
		HandlerMahasiswa(w, r)
	case strings.HasPrefix(r.URL.Path, "/uts/"):
		HandlerUTS(w, r)
	default:
		http.Error(w, "404 not found.", http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/", Handler)
	log.Printf("Server running on port %d", HTTPPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", HTTPPort), nil))
}

// Helper Functions
func FetchDocument(url string) (*goquery.Document, error) {
	res, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	return doc, nil
}

func WriteJSONResponse(w http.ResponseWriter, data interface{}) {
	response := Response{
		Status: "success",
		Data:   data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON: %v", err), http.StatusInternalServerError)
	}
}

func GetJadwal(url string) (Jadwal, error) {
	doc, err := FetchDocument(url)
	if err != nil {
		return Jadwal{}, err
	}

	jadwal := Jadwal{}
	hariMap := map[string]*[]MataKuliah{
		"Senin":  &jadwal.Senin,
		"Selasa": &jadwal.Selasa,
		"Rabu":   &jadwal.Rabu,
		"Kamis":  &jadwal.Kamis,
		"Jum'at": &jadwal.Jumat,
		"Sabtu":  &jadwal.Sabtu,
	}

	doc.Find("table").First().Find("tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 5 {
			return
		}

		hari := strings.TrimSpace(cells.Eq(1).Text())
		mataKuliah := MataKuliah{
			Nama:  strings.TrimSpace(cells.Eq(2).Text()),
			Waktu: strings.TrimSpace(cells.Eq(3).Text()),
			Ruang: strings.TrimSpace(cells.Eq(4).Text()),
			Dosen: strings.TrimSpace(cells.Eq(5).Text()),
		}

		if hariSlice, ok := hariMap[hari]; ok {
			*hariSlice = append(*hariSlice, mataKuliah)
		}
	})

	return jadwal, nil
}

func GetKegiatan(url string) ([]Kegiatan, error) {
	doc, err := FetchDocument(url)
	if err != nil {
		return nil, err
	}

	var kegiatanList []Kegiatan
	doc.Find("table").First().Find("tr").Each(func(i int, row *goquery.Selection) {
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
		doc, err := FetchDocument(url)
		if err != nil {
			return nil, err
		}

		doc.Find("table").First().Find("tr").Each(func(i int, row *goquery.Selection) {
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

func GetUTS(url string) ([]UTS, error) {
	doc, err := FetchDocument(url)
	if err != nil {
		return nil, err
	}

	var utsList []UTS
	doc.Find("table").First().Find("tr").Each(func(i int, row *goquery.Selection) {
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
