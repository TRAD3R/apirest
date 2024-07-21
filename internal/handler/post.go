package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) getPosts(c *gin.Context) {
	users, err := h.router.PostList(c.Request)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.JSON(http.StatusOK, users)
}

func (h *Handler) addPost(c *gin.Context) {
	err := h.router.PostAdd(c.Request)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		if _, err = c.Writer.Write([]byte(err.Error())); err != nil {
			log.Println(err)
		}
	} else {
		c.Writer.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) updatePost(c *gin.Context) {
	err := h.router.PostUpdate(c.Request)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		if _, err = c.Writer.Write([]byte(err.Error())); err != nil {
			log.Println(err)
		}
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) deletePost(c *gin.Context) {
	err := h.router.PostDelete(c.Request)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		if _, err = c.Writer.Write([]byte(err.Error())); err != nil {
			log.Println(err)
		}
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
}
