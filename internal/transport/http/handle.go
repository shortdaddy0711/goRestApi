package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/shortdaddy0711/go-rest-api/internal/comment"
)

type Handler struct {
	Router  *mux.Router
	Service *comment.Service
}

func NewHandler(service *comment.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) SetupRoutes() {
	fmt.Println("Setting Up Routes")
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", h.DeleteComment).Methods("DELETE")
	h.Router.HandleFunc("/api/comment/{id}", h.UpdateComment).Methods("PUT")
	h.Router.HandleFunc("/api/comment", h.PostComment).Methods("POST")
	h.Router.HandleFunc("/api/comment", h.GetAllComments).Methods("GET")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "I am alive!")
	})
}

func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse UINT from ID")
	}

	comment, err := h.Service.GetComment(uint(i))
	if err != nil {
		fmt.Fprintf(w, "Error Retrieving comment by ID: %s\n", id)
	}
	fmt.Fprintf(w, "%+v", comment)
}

func (h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.Service.GetAllComment()
	if err != nil {
		fmt.Fprintf(w, "Failed to retrieve all comments")
	}
	fmt.Fprintf(w, "%+v", comments)
}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	comment, err := h.Service.PostComment(comment.Comment{
		Slug: "/",
	})
	if err != nil {
		fmt.Fprintf(w, "Failed to post new comment")
	}

	fmt.Fprintf(w, "%+v", comment)
}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse UINT from ID")
	}

	comment, err := h.Service.UpdateComment(uint(i), comment.Comment{
		Slug: "/new",
	})
	if err != nil {
		fmt.Fprintf(w, "Failed to update comment by ID: %s\n", id)
	}

	fmt.Fprintf(w, "%+v", comment)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse UINT from ID")
	}

	err = h.Service.DeleteComment(uint(i))
	if err != nil {
		fmt.Fprintf(w, "Failed to delete comment by ID: %s\n", id)
	}

	fmt.Fprintf(w, "Successfully deleted comment")
}
