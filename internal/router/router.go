package router

import (
	"errors"
	"github.com/TRAD3R/tlog"
	"github.com/trad3r/hskills/apirest/internal/storage"
	"strconv"
	"strings"
)

type Router struct {
	db     *storage.Storage
	logger *tlog.Logger
}

func NewRouter(logger *tlog.Logger, db *storage.Storage) *Router {
	return &Router{
		db:     db,
		logger: logger,
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
