package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// frameworks: echo, gin

func (h *Handler) getUsers() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.router.UserList(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := json.Marshal(users)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			log.Printf("failed to write body: %v\n", err)
		}
	}
}

func (h *Handler) addUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.router.UserAdd(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	err := h.router.UserUpdate(r)
	if err != nil {
		// тут надо проверить на ошибку errs.ErrUserNotFound
		w.WriteHeader(http.StatusBadRequest) // 404
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) updateUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.router.UserUpdate(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (h *Handler) deleteUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.router.UserDelete(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
