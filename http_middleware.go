package middlewares

import (
	"net/http"
)

// AuthMiddleware vérifie l'authentification des requêtes entrantes
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handleRequest(w, r); err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}

		// Authentification
		if err := AuthenticateRequest(r); err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}

		// Appeler le middleware suivant
		next.ServeHTTP(w, r)
	})
}

// ValidationMiddleware valide les données des requêtes entrantes
func ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handleRequest(w, r); err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}

		// Validation (exemple simplifié)
		body, err := ReadAndResetBody(r)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		if len(body) == 0 {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		// Appeler le middleware suivant
		next.ServeHTTP(w, r)
	})
}

// handleRequest gère la lecture du corps de la requête et l'authentification
func handleRequest(w http.ResponseWriter, r *http.Request) *RequestError {
	// Lire et réinitialiser le corps de la requête
	if _, err := ReadAndResetBody(r); err != nil {
		return &RequestError{Message: "Failed to read request body", StatusCode: http.StatusInternalServerError}
	}

	return nil
}
