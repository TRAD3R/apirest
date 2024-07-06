package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) getPosts(w http.ResponseWriter, r *http.Request) {
	users, err := h.router.PostList(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		log.Printf("failed to write body: %v\n", err)
	}
}

func (h *Handler) addPost(w http.ResponseWriter, r *http.Request) {
	err := h.router.PostAdd(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request) {
	err := h.router.PostUpdate(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	err := h.router.PostDelete(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
