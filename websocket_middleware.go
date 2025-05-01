package middlewares

import (
	"encoding/json"
	"log"
	"net/http"

	logger "github.com/crettien/logger/models"
	"github.com/gorilla/websocket"
)

// Type définissant la signature de la fonction de rappel pour envoyer les logs
type SendLogFunc func(logEntry logger.LogEntry) error

// AuthWebSocketMiddleware pour WebSocket
func AuthWebSocketMiddleware(conn *websocket.Conn, r *http.Request, sendLog SendLogFunc) {
	if err := AuthenticateRequest(r); err != nil {
		log.Printf("Authentication failed: %v", err)
		conn.Close()
		return
	}
	HandleWebSocket(conn, sendLog)
}

// ValidationWebSocketMiddleware pour WebSocket
func ValidationWebSocketMiddleware(conn *websocket.Conn, r *http.Request, sendLog SendLogFunc) {
	// Utiliser un http.ResponseWriter fictif
	var w http.ResponseWriter = nil
	if err := ValidateRequestBody(w, r); err != nil {
		log.Printf("Validation failed: %v", err)
		conn.Close()
		return
	}
	HandleWebSocket(conn, sendLog)
}

// AuthAndValidationWebSocketMiddleware pour WebSocket
func AuthAndValidationWebSocketMiddleware(conn *websocket.Conn, r *http.Request, sendLog SendLogFunc) {
	if err := AuthenticateRequest(r); err != nil {
		log.Printf("Authentication failed: %v", err)
		conn.Close()
		return
	}
	// Utiliser un http.ResponseWriter fictif
	var w http.ResponseWriter = nil
	if err := ValidateRequestBody(w, r); err != nil {
		log.Printf("Validation failed: %v", err)
		conn.Close()
		return
	}
	HandleWebSocket(conn, sendLog)
}

// HandleWebSocket gère la logique de gestion des messages WebSocket
func HandleWebSocket(conn *websocket.Conn, sendLog SendLogFunc) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		log.Printf("Received: %s", message)

		var logEntry logger.LogEntry
		err = json.Unmarshal(message, &logEntry)
		if err != nil {
			log.Printf("Invalid JSON payload, err=%d", http.StatusBadRequest)
			return
		}

		// Appeler la fonction de rappel pour envoyer le log
		if err := sendLog(logEntry); err != nil {
			log.Printf("Failed to send log, err=%v", err)
			return
		}

		// Envoyer un message de réponse
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}
