package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yafyx/baak-api/models"
	"github.com/yafyx/baak-api/utils"
)

func HandlerKelasbaru(w http.ResponseWriter, r *http.Request) {
	searchTerm := strings.TrimPrefix(r.URL.Path, "/kelasbaru/")
	if searchTerm == "" {
		http.Error(w, "Missing search term in URL", http.StatusBadRequest)
		return
	}

	searchTypes := []string{"Kelas", "NPM", "Nama"}
	var kelasBaru []models.KelasBaru
	var err error

	for _, searchType := range searchTypes {
		url := fmt.Sprintf("%s/cariKelasBaru?tipeKelasBaru=%s&teks=%s", utils.BaseURL, searchType, searchTerm)
		kelasBaru, err = utils.GetKelasbaru(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(kelasBaru) > 0 {
			break
		}
	}

	if len(kelasBaru) == 0 {
		http.Error(w, "Mahasiswa tidak ditemukan!", http.StatusNotFound)
		return
	}

	utils.WriteJSONResponse(w, kelasBaru)
}
