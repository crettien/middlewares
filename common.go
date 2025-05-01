package middlewares

import (
	"bytes"
	"io"
	"net/http"
)

// RequestError représente une erreur de requête avec un message et un code d'état
type RequestError struct {
	Message    string
	StatusCode int
}

// Error implémente la méthode Error pour RequestError
func (e *RequestError) Error() string {
	return e.Message
}

// ReadAndResetBody lit et réinitialise le corps de la requête
func ReadAndResetBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

// CheckAuth vérifie l'authentification de la requête
func CheckAuth(r *http.Request) bool {
	return r.Header.Get("Authorization") == "Bearer valid_token"
}

// AuthenticateRequest vérifie l'authentification de la requête
func AuthenticateRequest(r *http.Request) *RequestError {
	if !CheckAuth(r) {
		return &RequestError{Message: "Unauthorized", StatusCode: http.StatusUnauthorized}
	}
	return nil
}

// ValidateRequestBody valide le corps de la requête
func ValidateRequestBody(w http.ResponseWriter, r *http.Request) *RequestError {
	body, err := ReadAndResetBody(r)
	if err != nil {
		return &RequestError{Message: "Failed to read request body", StatusCode: http.StatusInternalServerError}
	}
	if len(body) == 0 {
		return &RequestError{Message: "Request body is empty", StatusCode: http.StatusBadRequest}
	}
	return nil
}
