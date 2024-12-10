package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yafyx/baak-api/models"
	"github.com/yafyx/baak-api/utils"
)

func HandlerMahasiswaBaru(w http.ResponseWriter, r *http.Request) {
	searchTerm := strings.TrimPrefix(r.URL.Path, "/mahasiswabaru/")
	if searchTerm == "" {
		http.Error(w, "Missing search term in URL", http.StatusBadRequest)
		return
	}

	searchTypes := []string{"Kelas", "Nama"}
	var mahasiswaBaru []models.MahasiswaBaru
	var err error

	for _, searchType := range searchTypes {
		url := fmt.Sprintf("%s/cariMhsBaru?tipeMhsBaru=%s&teks=%s", utils.BaseURL, searchType, searchTerm)
		mahasiswaBaru, err = utils.GetMahasiswaBaru(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(mahasiswaBaru) > 0 {
			break
		}
	}

	if len(mahasiswaBaru) == 0 {
		http.Error(w, "Mahasiswa baru tidak ditemukan!", http.StatusNotFound)
		return
	}

	utils.WriteJSONResponse(w, mahasiswaBaru)
}