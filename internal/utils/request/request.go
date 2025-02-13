package request

import (
	"encoding/json"
	"io"
	"net/http"
)

func GetRequestData(r *http.Request, requestData interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if err := json.Unmarshal(body, &requestData); err != nil {
		return err
	}
	return nil
}
