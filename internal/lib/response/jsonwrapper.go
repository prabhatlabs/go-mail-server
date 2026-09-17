package response

import (
	"encoding/json/v2"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SendJsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)

	res := data
	if res == nil {
		res = struct {
			Message string `json:"message"`
		}{Message: "No data"}
	}

	jsonData, err := json.Marshal(res)
	if err != nil {
		log.Fatalln(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
