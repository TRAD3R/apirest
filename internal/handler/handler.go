package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trad3r/hskills/apirest/internal/service"
)

type Handler struct {
	userService     service.IUserService
	postService     service.IPostService
	userPostService service.IUserPostService
}

func NewHandler(u service.IUserService, p service.IPostService, up service.IUserPostService) *Handler {
	return &Handler{
		userService:     u,
		postService:     p,
		userPostService: up,
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
