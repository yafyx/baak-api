package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yafyx/baak-api/models"
	"github.com/yafyx/baak-api/utils"
)

func HandlerMahasiswa(w http.ResponseWriter, r *http.Request) {
	searchTerm := strings.TrimPrefix(r.URL.Path, "/kelasbaru/")
	if searchTerm == "" {
		http.Error(w, "Missing search term in URL", http.StatusBadRequest)
		return
	}

	searchTypes := []string{"Kelas", "NPM", "Nama"}
	var mahasiswas []models.Mahasiswa
	var err error

	for _, searchType := range searchTypes {
		url := fmt.Sprintf("%s/cariKelasBaru?tipeKelasBaru=%s&teks=%s", utils.BaseURL, searchType, searchTerm)
		mahasiswas, err = utils.GetMahasiswa(url)
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

	utils.WriteJSONResponse(w, mahasiswas)
}
