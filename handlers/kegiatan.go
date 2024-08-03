package handlers

import (
	"net/http"

	"github.com/yafyx/baak-api/utils"
)

func HandlerKegiatan(w http.ResponseWriter, r *http.Request) {
	kegiatanList, err := utils.GetKegiatan(utils.BaseURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSONResponse(w, kegiatanList)
}
