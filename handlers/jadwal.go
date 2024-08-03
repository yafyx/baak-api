package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yafyx/baak-api/models"
	"github.com/yafyx/baak-api/utils"
)

func HandlerJadwal(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimPrefix(r.URL.Path, "/jadwal/")
	if search == "" {
		http.Error(w, "Missing kelas in URL", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%s/jadwal/cariJadKul?&teks=%s", utils.BaseURL, search)
	jadwal, err := utils.GetJadwal(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Kelas  string        `json:"kelas"`
		Jadwal models.Jadwal `json:"jadwal"`
	}{
		Kelas:  search,
		Jadwal: jadwal,
	}

	utils.WriteJSONResponse(w, response)
}
