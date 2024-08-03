package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yafyx/baak-api/utils"
)

func HandlerUTS(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimPrefix(r.URL.Path, "/uts/")
	if search == "" {
		http.Error(w, "Missing search term in URL", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%s/jadwal/cariUts?&teks=%s", utils.BaseURL, search)
	uts, err := utils.GetUTS(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSONResponse(w, uts)
}