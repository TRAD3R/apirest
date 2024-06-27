package router

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/trad3r/hskills/apirest/models"
	"github.com/trad3r/hskills/apirest/storage"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	defaultLimit = 10
)

func (r *Router) UserList(req *http.Request) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second*10)
	defer cancel()

	filter, err := parseUserFilters(req.URL.Query())
	if err != nil {
		return nil, err
	}

	return r.userStorage.GetList(ctx, filter)
}

func (r *Router) UserAdd(req *http.Request) error {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second*10)
	defer cancel()

	var userAddReq storage.UserAddRequest

	if req.Body != nil {
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("Failed to read body: %v", err)
			return errors.New("Failed to read body")
		}

		err = json.Unmarshal(reqBody, &userAddReq)
		if err != nil {
			log.Printf("Failed to unmarshal body: %v", err)
			return errors.New("Failed to unmarshal body")
		}
	}

	user := models.User{
		Name:        userAddReq.Name,
		Phonenumber: userAddReq.Phonenumber,
	}

	return r.userStorage.Add(ctx, user)
}

func (r *Router) UserUpdate(req *http.Request) error {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second*10)
	defer cancel()

	id, err := getIdFromPath(req.URL.Path)
	if err != nil {
		return err
	}

	var userUpdateReq storage.UserUpdateRequest

	if req.Body != nil {
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("Failed to read body: %v", err)
			return errors.New("Failed to read body")
		}

		err = json.Unmarshal(reqBody, &userUpdateReq)
		if err != nil {
			log.Printf("Failed to unmarshal body: %v", err)
			return errors.New("Failed to unmarshal body")
		}
	}

	return r.userStorage.Update(ctx, id, userUpdateReq)
}

func (r *Router) UserDelete(req *http.Request) error {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second*10)
	defer cancel()

	id, err := getIdFromPath(req.URL.Path)
	if err != nil {
		return err
	}

	return r.userStorage.Delete(ctx, id)
}

func parseUserFilters(query url.Values) (storage.UserFilter, error) {
	var filter storage.UserFilter

	from := query.Get("from")
	if len(from) > 0 {
		filterFrom, err := time.Parse("2006-01-02", from)
		if err != nil {
			log.Printf("failed to parse from: %v", err)
			return filter, errors.New("invalid format for from")
		}

		filter.FromCreatedAt = &filterFrom
	}

	to := query.Get("to")
	if len(to) > 0 {
		filterTo, err := time.Parse("2006-01-02", to)
		if err != nil {
			log.Printf("failed to parse to: %v", err)
			return filter, errors.New("invalid format for to")
		}

		filter.ToCreatedAt = &filterTo
	}

	names := query.Get("name")
	if len(names) > 0 {
		filterNames := strings.Split(names, ",")
		filter.Name = filterNames
	}

	offset := query.Get("offset")
	if len(offset) > 0 {
		filterOffset, err := strconv.Atoi(offset)
		if err != nil {
			log.Printf("failed to parse offset: %v", err)
			return filter, errors.New("invalid format for offset")
		}

		filter.Offset = filterOffset
	}

	limit := query.Get("limit")
	if len(limit) > 0 {
		filterLimit, err := strconv.Atoi(limit)
		if err != nil {
			log.Printf("failed to parse limit: %v", err)
			return filter, errors.New("invalid format for limit")
		}

		filter.Limit = filterLimit
	} else {
		filter.Limit = defaultLimit
	}

	sort := query.Get("sort")
	filter.TopPostsAmount = sort

	return filter, nil
}
