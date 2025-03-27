package main

import (
	"net/http"

	handler "github.com/yafyx/baak-api/api"
)

func main() {
	// The actual server is started in the handler package when deployed to Vercel
}

// Handler is exported to be used by Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	handler.Handler(w, r)
}
