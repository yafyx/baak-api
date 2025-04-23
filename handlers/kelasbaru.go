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

func HandlerKelasbaru(w http.ResponseWriter, r *http.Request) {
	searchTerm := strings.TrimPrefix(r.URL.Path, "/kelasbaru/")
	if searchTerm == "" {
		utils.WriteValidationError(w, "Missing search term in URL")
		return
	}

	searchTypes := []string{"Kelas", "NPM", "Nama"}
	var kelasBaru []models.KelasBaru
	var err error

	kelasBaruBaseURL := fmt.Sprintf("%s/cariKelasBaru", config.AppConfig.BaseURL)
	token, err := utils.GetCSRFToken(kelasBaruBaseURL)
	if err != nil {
		token, err = utils.GetCSRFToken(config.AppConfig.BaseURL)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get CSRF token for KelasBaru: %v", err))
			return
		}
	}

	for _, searchType := range searchTypes {
		searchURL := fmt.Sprintf("%s/cariKelasBaru?_token=%s&tipeKelasBaru=%s&teks=%s",
			config.AppConfig.BaseURL,
			url.QueryEscape(token),
			url.QueryEscape(searchType),
			url.QueryEscape(searchTerm),
		)
		kelasBaru, err = utils.GetKelasbaru(searchURL)
		if err != nil {
			utils.WriteHTTPError(w, err)
			return
		}
		if len(kelasBaru) > 0 {
			break
		}
	}

	if len(kelasBaru) == 0 {
		utils.WriteNotFoundError(w)
		return
	}

	utils.WriteJSONResponse(w, kelasBaru)
}
