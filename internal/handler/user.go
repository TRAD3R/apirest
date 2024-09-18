package handler

import (
	"errors"
	"github.com/trad3r/hskills/apirest/internal/customerrors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

func (h *Handler) getUsers(c *gin.Context) {
	// Create a tracer. Usually, tracer is a global variable.
	tracer := otel.Tracer("")

	// Create a root span (a trace) to measure some operation.
	ctx, span := tracer.Start(c.Request.Context(), "getUsers")
	defer span.End()

	users, err := h.userService.UserList(ctx, c.Request)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

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
	user, err := h.userService.UserAdd(c.Request)
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

	if err := h.userService.UserUpdate(id, c.Request); err != nil {
		if errors.Is(err, customerrors.ErrUserNotFound) {
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

	if err := h.userService.UserDelete(id, c.Request); err != nil {
		if errors.Is(err, customerrors.ErrUserNotFound) {
			c.Writer.WriteHeader(http.StatusNotFound)
		} else {
			c.Writer.WriteHeader(http.StatusBadRequest) // 404
		}
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
}
