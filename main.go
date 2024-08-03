package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
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
		Jam   string `json:"jam"`
		Ruang string `json:"ruang"`
		Dosen string `json:"dosen"`
	}

	Kegiatan struct {
		Kegiatan string `json:"kegiatan"`
		Tanggal  string `json:"tanggal"`
		Start    string `json:"start"`
		End      string `json:"end"`
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

	response := struct {
		Kelas  string `json:"kelas"`
		Jadwal Jadwal `json:"jadwal"`
	}{
		Kelas:  search,
		Jadwal: jadwal,
	}

	WriteJSONResponse(w, response)
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

	timeStampLUT, err := GetTimeStampLUT()
	if err != nil {
		return Jadwal{}, err
	}

	doc.Find("table").First().Find("tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 5 {
			return
		}

		hari := strings.TrimSpace(cells.Eq(1).Text())
		waktu := strings.TrimSpace(cells.Eq(3).Text())
		jam := convertWaktuToJam(waktu, timeStampLUT)

		mataKuliah := MataKuliah{
			Nama:  strings.TrimSpace(cells.Eq(2).Text()),
			Waktu: waktu,
			Jam:   jam,
			Ruang: strings.TrimSpace(cells.Eq(4).Text()),
			Dosen: strings.TrimSpace(cells.Eq(5).Text()),
		}

		if hariSlice, ok := hariMap[hari]; ok {
			*hariSlice = append(*hariSlice, mataKuliah)
		}
	})

	return jadwal, nil
}

func GetTimeStampLUT() ([][]string, error) {
	doc, err := FetchDocument(BaseURL + "/kuliahUjian/6")
	if err != nil {
		return nil, err
	}

	var result [][]string
	doc.Find("table.cell-xs-6 tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() >= 2 {
			timeRange := strings.TrimSpace(cells.Eq(1).Text())
			timeRange = strings.ReplaceAll(timeRange, " ", "")
			timeRange = strings.ReplaceAll(timeRange, ".", ":")
			times := strings.Split(timeRange, "-")
			if len(times) == 2 {
				result = append(result, times)
			}
		}
	})

	return result, nil
}

func convertWaktuToJam(waktu string, timeStampLUT [][]string) string {
	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindAllString(waktu, -1)

	if len(matches) < 2 || len(matches) > 3 || len(timeStampLUT) == 0 {
		return ""
	}

	start, _ := strconv.Atoi(matches[0])
	end, _ := strconv.Atoi(matches[len(matches)-1])

	if start < 1 || start > len(timeStampLUT) || end < 1 || end > len(timeStampLUT) {
		return ""
	}

	return timeStampLUT[start-1][0] + " - " + timeStampLUT[end-1][1]
}

func GetKegiatan(url string) ([]Kegiatan, error) {
	doc, err := FetchDocument(url)
	if err != nil {
		return nil, err
	}

	var kegiatanList []Kegiatan
	var parentKegiatan string

	doc.Find("table").First().Find("tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() == 2 {
			kegiatanText := strings.TrimSpace(cells.Eq(0).Text())
			tanggalText := strings.TrimSpace(cells.Eq(1).Text())

			if tanggalText == "" {
				parentKegiatan = kegiatanText
				return
			}

			start, end := parseTanggal(tanggalText)

			fullKegiatan := kegiatanText
			if parentKegiatan != "" && isSubItem(kegiatanText) {
				fullKegiatan = parentKegiatan + " " + kegiatanText
			} else {
				parentKegiatan = ""
			}

			kegiatan := Kegiatan{
				Kegiatan: fullKegiatan,
				Tanggal:  tanggalText,
				Start:    start,
				End:      end,
			}
			kegiatanList = append(kegiatanList, kegiatan)
		} else {
			parentKegiatan = ""
		}
	})

	return kegiatanList, nil
}

func isSubItem(text string) bool {
	return regexp.MustCompile(`^[a-z]\..+`).MatchString(text)
}

func parseTanggal(tanggal string) (start, end string) {
	parts := strings.Split(tanggal, "-")
	if len(parts) == 2 {
		start = strings.TrimSpace(parts[0])
		end = strings.TrimSpace(parts[1])
	} else if len(parts) == 1 {
		start = strings.TrimSpace(parts[0])
		end = start
	}

	start = addYearIfMissing(start)
	end = addYearIfMissing(end)

	return start, end
}

func addYearIfMissing(date string) string {
	if !strings.Contains(date, "20") {
		return date + " 2024"
	}
	return date
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
