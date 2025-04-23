package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/yafyx/baak-api/config"
	"github.com/yafyx/baak-api/models"
	"github.com/yafyx/baak-api/utils"
)

func HandlerMahasiswaBaru(w http.ResponseWriter, r *http.Request) {
	searchTerm := strings.TrimPrefix(r.URL.Path, "/mahasiswabaru/")
	if searchTerm == "" {
		utils.WriteValidationError(w, "Missing search term in URL")
		return
	}

	searchTypes := []string{"Kelas", "Nama"}
	var mahasiswaBaru []models.MahasiswaBaru
	var err error

	mhsBaruBaseURL := fmt.Sprintf("%s/cariMhsBaru", config.AppConfig.BaseURL)
	token, err := utils.GetCSRFToken(mhsBaruBaseURL)
	if err != nil {
		token, err = utils.GetCSRFToken(config.AppConfig.BaseURL)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get CSRF token for MahasiswaBaru: %v", err))
			return
		}
	}

	for _, searchType := range searchTypes {
		searchURL := fmt.Sprintf("%s/cariMhsBaru?_token=%s&tipeMhsBaru=%s&teks=%s",
			config.AppConfig.BaseURL,
			url.QueryEscape(token),
			url.QueryEscape(searchType),
			url.QueryEscape(searchTerm),
		)

		mahasiswaBaru, err = utils.GetMahasiswaBaru(searchURL)
		if err != nil {
			utils.WriteHTTPError(w, err)
			return
		}
		if len(mahasiswaBaru) > 0 {
			break
		}
	}

	if len(mahasiswaBaru) == 0 {
		utils.WriteNotFoundError(w)
		return
	}

	utils.WriteJSONResponse(w, mahasiswaBaru)
}
