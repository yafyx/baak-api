package utils

import (
	"encoding/json"
	"net/http"

	"github.com/yafyx/baak-api/models"
)

func WriteJSONResponse(w http.ResponseWriter, data interface{}) {
	response := models.Response{
		Status: "success",
		Data:   data,
	}

	dataJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(dataJson)
}
