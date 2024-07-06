package router

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/trad3r/hskills/apirest/internal/models"
	"github.com/trad3r/hskills/apirest/internal/storage"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (r *Router) PostList(req *http.Request) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second*10)
	defer cancel()

	filter, err := parsePostFilters(req.URL.Query())
	if err != nil {
		return nil, err
	}

	return r.postStorage.GetList(ctx, filter)
}

func (r *Router) PostAdd(req *http.Request) error {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second*10)
	defer cancel()

	var postAddReq storage.PostAddRequest

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

	post := models.Post{
		Subject: postAddReq.Subject,
		Body:    postAddReq.Body,
	}

	author, err := r.userStorage.FindById(ctx, postAddReq.Author)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return errors.New("failed to check author")
	}

	if author == nil {
		log.Printf("Failed to find user with id %d", postAddReq.Author)
		return errors.New("author does not exist")
	}

	post.Author = author.ID

	if err := r.postStorage.Add(ctx, &post); err != nil {
		log.Printf("failed to add post: %v", err)
		return errors.New("failed to add post")
	}

	return r.userStorage.IncrPostToUser(post.Author)
}

func (r *Router) PostUpdate(req *http.Request) error {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second*10)
	defer cancel()

	id, err := getIdFromPath(req.URL.Path)
	if err != nil {
		return err
	}

	var postUpdateReq storage.PostUpdateRequest

	if req.Body != nil {
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("Failed to read body: %v", err)
			return errors.New("Failed to read body")
		}

		err = json.Unmarshal(reqBody, &postUpdateReq)
		if err != nil {
			log.Printf("Failed to unmarshal body: %v", err)
			return errors.New("Failed to unmarshal body")
		}
	}

	return r.postStorage.Update(ctx, id, postUpdateReq)
}

func (r *Router) PostDelete(req *http.Request) error {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second*10)
	defer cancel()

	id, err := getIdFromPath(req.URL.Path)
	if err != nil {
		return err
	}

	if err := r.userStorage.DecrPostToUser(id); err != nil {
		log.Printf("failed to delete post: %v", err)
		return errors.New("failed to delete post")
	}

	return r.postStorage.Delete(ctx, id)
}

func parsePostFilters(query url.Values) (storage.PostFilter, error) {
	var filter storage.PostFilter
	var err error

	from := query.Get("from")
	if len(from) > 0 {
		filter.FromCreatedAt, err = time.Parse("2006-01-02", from)
		if err != nil {
			log.Printf("failed to parse from: %v", err)
			return filter, errors.New("invalid format for from")
		}
	}

	to := query.Get("to")
	if len(to) > 0 {
		filter.ToCreatedAt, err = time.Parse("2006-01-02", from)
		if err != nil {
			log.Printf("failed to parse to: %v", err)
			return filter, errors.New("invalid format for to")
		}
	}

	authors := query.Get("author")
	if len(authors) > 0 {
		for _, authorId := range strings.Split(authors, ",") {
			author, err := strconv.Atoi(authorId)
			if err != nil {
				log.Printf("failed to parse author: %v", err)
				continue
			}

			filter.Authors = append(filter.Authors, author)
		}
	}

	offset := query.Get("offset")
	if len(offset) > 0 {
		filter.Offset, err = strconv.Atoi(offset)
		if err != nil {
			log.Printf("failed to parse offset: %v", err)
			return filter, errors.New("invalid format for offset")
		}
	}

	limit := query.Get("limit")
	if len(limit) > 0 {
		filter.Limit, err = strconv.Atoi(limit)
		if err != nil {
			log.Printf("failed to parse limit: %v", err)
			return filter, errors.New("invalid format for limit")
		}
	} else {
		filter.Limit = 10
	}

	filter.Subject = query.Get("subject")

	return filter, nil
}
