package handlers

import (
	"fmt"
	"net/http"
	"strings"

	jadwalH "baak-api/handlers/jadwal"
	kegiatanH "baak-api/handlers/kegiatan"
	mahasiswaH "baak-api/handlers/mahasiswa"
	"baak-api/models"
	"baak-api/utils"
)

const BASE_URL = "https://baak.gunadarma.ac.id"

func HandlerHomepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "/jadwal/{cari berdasarkan kelas/dosen}, /kalender, /kelasbaru/{cari berdasarkan kelas/npm/nama}")
}
func HandlerJadwal(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 3 || segments[2] == "" {
		http.Error(w, "Missing kelas in URL", http.StatusBadRequest)
		return
	}

	search := segments[2]
	url := BASE_URL + "/jadwal/cariJadKul?&teks=" + search
	jadwal, err := jadwalH.GetJadwal(url, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSONResponse(w, jadwal)
}

func HandlerKegiatan(w http.ResponseWriter, r *http.Request) {
	kegiatanList, err := kegiatanH.GetKegiatan(BASE_URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSONResponse(w, kegiatanList)
}

func HandlerMahasiswa(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 3 || segments[2] == "" {
		http.Error(w, "Missing search term in URL", http.StatusBadRequest)
		return
	}

	searchTerm := segments[2]
	searchTypes := []string{"Kelas", "NPM", "Nama"}
	var mahasiswa []models.Mahasiswa
	var err error

	for _, searchType := range searchTypes {
		baseURL := BASE_URL + "/cariKelasBaru?tipeKelasBaru=" + searchType + "&teks=" + searchTerm
		mahasiswa, err = mahasiswaH.GetMahasiswa(baseURL)
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

	utils.WriteJSONResponse(w, mahasiswa)
}
