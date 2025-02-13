package response

import (
	"encoding/json"
	"net/http"
)

func WithError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(`{"error":"` + err.Error() + `"}`))
}

func WriteData(w http.ResponseWriter, data interface{}, statusCode int) {
	if data == nil {
		data = "Успешный ответ"
	}
	body, err := json.Marshal(data)
	if err != nil {
		WithError(w, 500, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)
}
