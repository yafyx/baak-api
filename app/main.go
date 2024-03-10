package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yafyx/baak-api/config"
	"github.com/yafyx/baak-api/handlers"
)

func main() {
	http.HandleFunc("/jadwal/", handlers.HandlerJadwal)
	http.HandleFunc("/kalender", handlers.HandlerKegiatan)
	http.HandleFunc("/kelasbaru/", handlers.HandlerMahasiswa)

	port := config.GlobalEnv.HTTPPort
	log.Printf("Server running on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
