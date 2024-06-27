package router

import (
	"errors"
	"github.com/trad3r/hskills/apirest/storage"
	"strconv"
	"strings"
)

type Router struct {
	userStorage *storage.UserStorage
	postStorage *storage.PostStorage
}

func NewRouter(us *storage.UserStorage, ps *storage.PostStorage) *Router {
	return &Router{
		userStorage: us,
		postStorage: ps,
	}
}

func getIdFromPath(path string) (int, error) {
	pathParts := strings.Split(path, "/")
	if len(pathParts) != 3 {
		return 0, errors.New("Invalid path")
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		return 0, errors.New("Invalid id")
	}

	return id, nil
}
