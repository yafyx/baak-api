package handlers

import (
	"net/http"

	"github.com/yafyx/baak-api/utils"
)

func HandlerHomepage(w http.ResponseWriter, r *http.Request) {
	endpoints := []string{
		"/jadwal/{kelas}",
		"/kalender",
		"/kelasbaru/{kelas/npm/nama}",
		"/uts/{kelas/dosen}",
		"/mahasiswabaru/{kelas/nama}",
	}
	utils.WriteJSONResponse(w, endpoints)
}
