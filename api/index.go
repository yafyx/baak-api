package handler

import (
	"net/http"
	"strings"

	"github.com/yafyx/baak-api/handlers"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		handlers.HandlerHomepage(w, r)
	case strings.HasPrefix(r.URL.Path, "/jadwal/"):
		handlers.HandlerJadwal(w, r)
	case r.URL.Path == "/kalender":
		handlers.HandlerKegiatan(w, r)
	case strings.HasPrefix(r.URL.Path, "/kelasbaru/"):
		handlers.HandlerKelasbaru(w, r)
	case strings.HasPrefix(r.URL.Path, "/uts/"):
		handlers.HandlerUTS(w, r)
	case strings.HasPrefix(r.URL.Path, "/mahasiswabaru/"):
		handlers.HandlerMahasiswaBaru(w, r)
	default:
		http.Error(w, "404 not found.", http.StatusNotFound)
	}
}
