package handler

import (
	"github.com/trad3r/hskills/apirest/router"
	"net/http"
)

type Handler struct {
	router *router.Router
}

func NewHandler(router *router.Router) *Handler {
	return &Handler{
		router: router,
	}
}

func (h *Handler) Handlers() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("GET /users", h.getUsers())

	r.HandleFunc("POST /user", h.addUser())

	r.HandleFunc("PATCH /user/{id}", h.UpdateUser) // так проще

	r.HandleFunc("DELETE /user/{id}", h.deleteUser())

	r.HandleFunc("GET /posts", h.getPosts())

	r.HandleFunc("POST /post", h.addPost())

	r.HandleFunc("PATCH /post/{id}", h.updatePost())

	r.HandleFunc("DELETE /post/{id}", h.deletePost())

	return r
}
