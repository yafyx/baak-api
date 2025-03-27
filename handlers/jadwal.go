package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yafyx/baak-api/config"
	"github.com/yafyx/baak-api/models"
	"github.com/yafyx/baak-api/utils"
)

func HandlerJadwal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	search := strings.TrimPrefix(r.URL.Path, "/jadwal/")
	if search == "" {
		utils.WriteValidationError(w, "Missing kelas in URL")
		return
	}

	// Validate input
	if len(search) < 3 {
		utils.WriteValidationError(w, "Kelas must be at least 3 characters long")
		return
	}

	url := fmt.Sprintf("%s/jadwal/cariJadKul?&teks=%s", config.AppConfig.BaseURL, search)
	jadwal, err := utils.GetJadwal(url)
	if err != nil {
		utils.WriteHTTPError(w, err)
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
