package handler

import (
	"net/http"

	"github.com/yafyx/baak-api/handlers"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/jadwal/":
		handlers.HandlerJadwal(w, r)
	case "/kalender":
		handlers.HandlerKegiatan(w, r)
	case "/kelasbaru/":
		handlers.HandlerMahasiswa(w, r)
	default:
		http.Error(w, "404 not found.", http.StatusNotFound)
	}
}

// func main() {
// 	http.HandleFunc("/jadwal/", handlers.HandlerJadwal)
// 	http.HandleFunc("/kalender", handlers.HandlerKegiatan)
// 	http.HandleFunc("/kelasbaru/", handlers.HandlerMahasiswa)

// 	port := config.GlobalEnv.HTTPPort
// 	log.Printf("Server running on port %d", port)
// 	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
// }
