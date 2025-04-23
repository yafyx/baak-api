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

	// Fetch CSRF token from the base jadwal page
	jadwalBaseURL := fmt.Sprintf("%s/jadwal", config.AppConfig.BaseURL)
	token, err := utils.GetCSRFToken(jadwalBaseURL)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get CSRF token: %v", err))
		return
	}

	// Construct the search URL with the token
	searchURL := fmt.Sprintf("%s/jadwal/cariJadKul?_token=%s&teks=%s",
		config.AppConfig.BaseURL,
		url.QueryEscape(token),
		url.QueryEscape(search),
	)

	jadwal, err := utils.GetJadwal(searchURL)
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

func HandlerJadwalSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	search := r.URL.Query().Get("q")
	if search == "" {
		utils.WriteValidationError(w, "Missing search query parameter 'q'")
		return
	}

	// Validate input
	if len(search) < 3 {
		utils.WriteValidationError(w, "Search query must be at least 3 characters long")
		return
	}

	// Fetch CSRF token from the base jadwal page
	jadwalBaseURL := fmt.Sprintf("%s/jadwal", config.AppConfig.BaseURL)
	token, err := utils.GetCSRFToken(jadwalBaseURL)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get CSRF token: %v", err))
		return
	}

	// Construct the search URL with the token
	searchURL := fmt.Sprintf("%s/jadwal/cariJadKul?_token=%s&teks=%s",
		config.AppConfig.BaseURL,
		url.QueryEscape(token),
		url.QueryEscape(search),
	)

	jadwal, err := utils.GetJadwal(searchURL)
	if err != nil {
		utils.WriteHTTPError(w, err)
		return
	}

	response := struct {
		Query  string        `json:"query"`
		Jadwal models.Jadwal `json:"jadwal"`
	}{
		Query:  search,
		Jadwal: jadwal,
	}

	utils.WriteJSONResponse(w, response)
}
