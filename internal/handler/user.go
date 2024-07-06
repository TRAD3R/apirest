package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// frameworks: echo, gin

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		log.Printf("failed to write body: %v\n", err)
	}
}

func (h *Handler) addUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.router.UserAdd(r)
	if err != nil {
		log.Printf("failed to add user: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write([]byte("failed to add user")); err != nil {
			log.Printf("failed to write body: %v\n", err)
		}
	} else {
		userJson, err := json.Marshal(user)
		if err != nil {
			log.Printf("failed to marshal user: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("something went wrong")); err != nil {
				log.Printf("failed to write body: %v\n", err)
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(userJson); err != nil {
				log.Printf("failed to write body: %v\n", err)
			}
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

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	err := h.router.UserDelete(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
