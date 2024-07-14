package handler

import (
	"net/http"

	"baak-api/handlers"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	handlers.HandlerJadwal(w, r)
}
