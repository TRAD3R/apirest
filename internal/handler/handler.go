package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/trad3r/hskills/apirest/internal/router"
	"net/http"
)

type Handler struct {
	router *router.Router
}

func NewHandler(r *router.Router) *Handler {
	return &Handler{
		router: r,
	}
}

func (h *Handler) Handlers() http.Handler {
	r := gin.Default()

	r.GET("/users", h.getUsers)
	r.POST("/user", h.addUser)
	r.PATCH("/user/:id", h.UpdateUser) // так проще
	r.DELETE("/user/:id", h.deleteUser)

	r.GET("/posts", h.getPosts)
	r.POST("/post", h.addPost)
	r.PATCH("/post/:id", h.updatePost)
	r.DELETE("/post/:id", h.deletePost)

	//r.HandleFunc("/debug/pprof/", pprof.Index)
	//r.HandleFunc("debug/pprof/cmdline", pprof.Cmdline)
	//r.HandleFunc("debug/pprof/profile", pprof.Profile)
	//r.HandleFunc("debug/pprof/symbol", pprof.Symbol)
	//r.HandleFunc("debug/pprof/trace", pprof.Trace)

	return r
}
