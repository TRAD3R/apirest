package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getPosts(c *gin.Context) {
	users, err := h.postService.PostList(c.Request)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.JSON(http.StatusOK, users)
}

func (h *Handler) addPost(c *gin.Context) {
	err := h.userPostService.AddPost(c.Request)
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
	err := h.postService.PostUpdate(c.Request)
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
	err := h.postService.PostDelete(c.Request)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		if _, err = c.Writer.Write([]byte(err.Error())); err != nil {
			log.Println(err)
		}
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
}
