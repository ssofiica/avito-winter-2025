package request

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequestData(t *testing.T) {
	requestData := map[string]string{"name": "sofia"}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/sendCoin", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	var res map[string]string
	err = GetRequestData(req, &res)
	if err != nil {
		t.Fatalf("GetRequestData returned an error: %v", err)
	}

	assert.Equal(t, res, requestData)
	assert.NoError(t, err)
}
