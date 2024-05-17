package handler

import (
	"net/http"

	"github.com/yafyx/baak-api/handlers"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	handlers.HandlerJadwal(w, r)
}
