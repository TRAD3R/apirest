package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trad3r/hskills/apirest/internal/custom_errors"
)

// frameworks: echo, gin

func (h *Handler) getUsers(c *gin.Context) {
	users, err := h.router.UserList(c.Request)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		if _, err := c.Writer.Write([]byte(err.Error())); err != nil {
			log.Printf("failed to write response: %v", err)
		}
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.JSON(http.StatusOK, users)
}

func (h *Handler) addUser(c *gin.Context) {
	user, err := h.router.UserAdd(c.Request)
	if err != nil {
		log.Printf("failed to add user: %v\n", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		if _, err = c.Writer.Write([]byte("failed to add user")); err != nil {
			log.Printf("failed to write body: %v\n", err)
		}
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.Atoi(idVal)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		if _, err = c.Writer.Write([]byte("invalid path param ID")); err != nil {
			log.Printf("failed to write body: %v\n", err)
		}

		return
	}

	if err := h.router.UserUpdate(id, c.Request); err != nil {
		if errors.Is(err, custom_errors.ErrUserNotFound) {
			c.Writer.WriteHeader(http.StatusNotFound)
		} else {
			c.Writer.WriteHeader(http.StatusBadRequest) // 404
		}
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) deleteUser(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.Atoi(idVal)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		if _, err = c.Writer.Write([]byte("invalid path param ID")); err != nil {
			log.Printf("failed to write body: %v\n", err)
		}

		return
	}

	if err := h.router.UserDelete(id, c.Request); err != nil {
		if errors.Is(err, custom_errors.ErrUserNotFound) {
			c.Writer.WriteHeader(http.StatusNotFound)
		} else {
			c.Writer.WriteHeader(http.StatusBadRequest) // 404
		}
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
}
