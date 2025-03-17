package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/yafyx/baak-api/handlers"
	"github.com/yafyx/baak-api/utils"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if !utils.Limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

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

func main() {
	port := ":8080"
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(port, http.HandlerFunc(Handler)); err != nil {
		log.Fatal(err)
	}
}
