package router

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/trad3r/hskills/apirest/internal/models"
	"github.com/trad3r/hskills/apirest/internal/repository/filters"
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
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	filter, err := parseUserFilters(req.URL.Query())
	if err != nil {
		r.logger.Error(err.Error())
		return nil, errors.New("invalid request params")
	}

	users, err := r.db.User.GetList(ctx, filter)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, errors.New("users not found")
	}

	return users, nil
}

func (r *Router) UserAdd(req *http.Request) (*models.User, error) {
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	var userAddReq filters.UserAddRequest

	if req.Body != nil {
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("Failed to read body: %v", err)
			return nil, errors.New("Failed to read body")
		}

		err = json.Unmarshal(reqBody, &userAddReq)
		if err != nil {
			log.Printf("Failed to unmarshal body: %v", err)
			return nil, errors.New("Failed to unmarshal body")
		}
	}

	user := &models.User{
		Name:        userAddReq.Name,
		Phonenumber: userAddReq.Phonenumber,
	}

	err := r.db.User.Add(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Router) UserUpdate(userId int, req *http.Request) error {
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	var userUpdateReq filters.UserUpdateRequest

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

	return r.db.User.Update(ctx, userId, userUpdateReq)
}

func (r *Router) UserDelete(userId int, req *http.Request) error {
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	return r.db.User.Delete(ctx, userId)
}

func parseUserFilters(query url.Values) (filters.UserFilter, error) {
	var filter filters.UserFilter

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

		filter.Offset = uint(filterOffset)
	}

	limit := query.Get("limit")
	if len(limit) > 0 {
		filterLimit, err := strconv.Atoi(limit)
		if err != nil {
			log.Printf("failed to parse limit: %v", err)
			return filter, errors.New("invalid format for limit")
		}

		filter.Limit = uint(filterLimit)
	} else {
		filter.Limit = uint(defaultLimit)
	}

	sort := query.Get("sort")
	filter.TopPostsAmount = sort

	return filter, nil
}
