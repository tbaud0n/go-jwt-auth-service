package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

// TokenNewHandler generate a new token
type TokenNewHandler struct {
	JWTService *JWTService
}

func (h TokenNewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Error(data)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	ts, err := h.JWTService.NewTokenString(data)
	if err != nil {
		log.Error(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	SetJSONResponse(w, ts)
}

// TokenCheckHandler validate token
type TokenCheckHandler struct {
	JWTService *JWTService
}

func (h TokenCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	tokenString, err := readTokenString(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if tokenString == "" {
		http.Error(w, "Empty token string", http.StatusUnauthorized)
		return
	}

	u, err := h.JWTService.ValidateToken(tokenString)
	if err != nil {
		log.Error(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	SetJSONResponse(w, map[string]interface{}{
		"valid": u != nil,
	})
}

// TokenDecodeHandler validate token
type TokenDecodeHandler struct {
	JWTService *JWTService
}

func (h TokenDecodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tokenString, err := readTokenString(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if tokenString == "" {
		http.Error(w, "Empty token string", http.StatusUnauthorized)
		return
	}

	u, err := h.JWTService.ValidateToken(tokenString)
	if err != nil {
		log.Error(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if u == nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	SetJSONResponse(w, u)
}

func readTokenString(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	authHeaderParts := strings.Split(authHeader, ` `)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != `bearer` {
		err := errors.New("Bad Authorization format header")
		log.Error(err)
		return "", err
	}

	return authHeaderParts[1], nil
}

// SetJSONResponse writes the data as JSON as response of the http request
func SetJSONResponse(w http.ResponseWriter, data interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(data); err != nil {
		log.Error(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(b); err != nil {
		log.Error(err)
		return
	}

	return
}
