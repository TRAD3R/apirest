package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/trad3r/hskills/apirest/internal/repository/filters"
)

type IUserPostService interface {
	AddPost(req *http.Request) error
}

type UserPostService struct {
	u IUserService
	p IPostService
}

func NewUserPostService(u IUserService, p IPostService) IUserPostService {
	return &UserPostService{u: u, p: p}
}

func (up *UserPostService) AddPost(req *http.Request) error {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second*10)
	defer cancel()

	var postAddReq filters.PostAddRequest

	if req.Body != nil {
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("Failed to read body: %v", err)
			return errors.New("Failed to read body")
		}

		err = json.Unmarshal(reqBody, &postAddReq)
		if err != nil {
			log.Printf("Failed to unmarshal body: %v", err)
			return errors.New("Failed to unmarshal body")
		}
	}

	if len(postAddReq.Subject) == 0 {
		return errors.New("subject field is required")
	}

	if postAddReq.Author < 1 {
		return errors.New("author field is required")
	}

	author, err := up.u.FindByID(ctx, postAddReq.Author)
	if err != nil {
		log.Printf("failed to get author: %v", err)
		return errors.New("failed to check author")
	}

	if author == nil {
		log.Printf("Failed to find author with id %d", postAddReq.Author)
		return errors.New("author does not exist")
	}

	return up.p.PostAdd(ctx, postAddReq.Subject, postAddReq.Body, *author)
}

func getIDFromPath(path string) (int, error) {
	pathParts := strings.Split(path, "/")
	if len(pathParts) != 3 {
		return 0, errors.New("Invalid path")
	}

	ID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		return 0, errors.New("Invalid id")
	}

	return ID, nil
}
