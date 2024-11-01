// utils/response.go
package utils

import (
    "encoding/json"
    "net/http"
    "tech-test/backend/internal/domain"
)

func RespondWithError(w http.ResponseWriter, apiErr *domain.APIError) {
    RespondWithJSON(w, apiErr.StatusCode, apiErr)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}
