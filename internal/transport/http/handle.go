package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/shortdaddy0711/go-rest-api/internal/comment"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Router  *mux.Router
	Service *comment.Service
}

type Response struct {
	Message string
	Error   string
}

func NewHandler(service *comment.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"Method": r.Method,
			"Path":   r.URL.Path,
		}).Info("handled request")
		next.ServeHTTP(w, r)
	})
}

func BasicAuth(callback func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Basic Auth Endpoint hit")
		user, pw, ok := r.BasicAuth()
		if user == "admin" && pw == "admin" && ok {
			callback(w, r)
		} else {
			sendErrorResponse(w, "You are not authorized", errors.New("authorization failed"))
			return
		}
	}
}

func validateToken(accessToken string) bool {
	var myKey = []byte("fruitAndVegetables")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("access token not validated")
		}
		return myKey, nil
	})
	if err != nil {
		return false
	}
	return token.Valid
}

func JWTAuth(callback func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("JWT Auth Endpoint hit")
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			sendErrorResponse(w, "You are not authorized for this", errors.New("authorization failed"))
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			sendErrorResponse(w, "You are not authorized for this", errors.New("authorization failed"))
			return
		}

		if validateToken(authHeaderParts[1]) {
			callback(w, r)
		} else {
			sendErrorResponse(w, "You are not authorized for this", errors.New("authorization failed"))
			return
		}
	}
}

func (h *Handler) SetupRoutes() {
	log.Info("Setting Up Routes")
	h.Router = mux.NewRouter()

	h.Router.Use(LoggingMiddleware)

	h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", JWTAuth(h.DeleteComment)).Methods("DELETE")
	h.Router.HandleFunc("/api/comment/{id}", JWTAuth(h.UpdateComment)).Methods("PUT")
	h.Router.HandleFunc("/api/comment", JWTAuth(h.PostComment)).Methods("POST")
	h.Router.HandleFunc("/api/comment", h.GetAllComments).Methods("GET")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Response{Message: "I am Alive"}); err != nil {
			panic(err)
		}
	})
}

func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(Response{Message: message, Error: err.Error()}); err != nil {
		log.Error(err)
	}
}
