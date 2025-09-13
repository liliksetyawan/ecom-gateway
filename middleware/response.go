package middleware

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Timestamp string      `json:"timestamp"`
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, rawErr string) {
	resp := map[string]interface{}{
		"timestamp":  time.Now().UTC().Format(time.RFC3339), // default kalau gagal parse
		"request_id": uuid.New().String(),                   // default kalau gagal parse
		"data": map[string]interface{}{
			"error": rawErr,
		},
	}

	// Coba parse error JSON dari backend
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(rawErr), &parsed); err == nil {
		// ambil timestamp & request_id dari parsed jika ada
		if ts, ok := parsed["timestamp"].(string); ok {
			resp["timestamp"] = ts
		}
		if rid, ok := parsed["request_id"].(string); ok {
			resp["request_id"] = rid
		}

		// ambil data.error
		if d, ok := parsed["data"].(map[string]interface{}); ok {
			if e, ok := d["error"]; ok {
				resp["data"] = map[string]interface{}{
					"error": e,
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}
