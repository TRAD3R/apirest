package router

import (
	"errors"
	"github.com/trad3r/hskills/apirest/internal/storage"
	"strconv"
	"strings"
)

type Router struct {
	db *storage.Storage
}

func NewRouter(db *storage.Storage) *Router {
	return &Router{
		db: db,
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
