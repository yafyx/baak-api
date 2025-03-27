package main

import (
	"fmt"
	"log"
	"net/http"

	handler "github.com/yafyx/baak-api/api"
	"github.com/yafyx/baak-api/config"
)

func main() {
	config.LoadConfig()

	// Start server (only runs locally, not on Vercel)
	port := config.AppConfig.Port
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(port, http.HandlerFunc(Handler)); err != nil {
		log.Fatal(err)
	}
}

// Handler is exported to be used by Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	handler.Handler(w, r)
}
