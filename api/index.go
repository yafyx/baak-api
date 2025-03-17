package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/yafyx/baak-api/config"
	"github.com/yafyx/baak-api/handlers"
	"github.com/yafyx/baak-api/middleware"
	"github.com/yafyx/baak-api/utils"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Apply middleware chain
	handler := middleware.RecoveryMiddleware(
		middleware.LoggingMiddleware(
			middleware.CORSMiddleware(
				middleware.RateLimitMiddleware(
					http.HandlerFunc(handleRoutes),
				),
			),
		),
	)
	handler.ServeHTTP(w, r)
}

func handleRoutes(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		handlers.HandlerHomepage(w, r)
	case r.URL.Path == "/health":
		handlers.HandlerHealth(w, r)
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
		utils.WriteNotFoundError(w)
	}
}

func main() {
	config.LoadConfig()

	port := config.AppConfig.Port
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(port, http.HandlerFunc(Handler)); err != nil {
		log.Fatal(err)
	}
}
