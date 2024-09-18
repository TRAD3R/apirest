package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/TRAD3R/tlog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trad3r/hskills/apirest/internal/models"
	"github.com/trad3r/hskills/apirest/internal/repository/filters"
	"github.com/trad3r/hskills/apirest/internal/repository/postgres"
	"go.opentelemetry.io/otel"
)

var (
	defaultLimit = 10
)

type IUserService interface {
	UserList(ctx context.Context, req *http.Request) ([]models.User, error)
	UserAdd(req *http.Request) (*models.User, error)
	UserUpdate(userID int, req *http.Request) error
	UserDelete(userID int, req *http.Request) error
	FindByID(ctx context.Context, userID int) (*models.User, error)
}

type UserService struct {
	repo   postgres.IUserRepository
	logger *tlog.Logger
}

func NewUserService(logger *tlog.Logger, db *pgxpool.Pool) IUserService {
	return &UserService{
		repo:   postgres.NewUserRepository(db),
		logger: logger,
	}
}

func (s *UserService) UserList(ctx context.Context, req *http.Request) ([]models.User, error) {
	ctx, span := otel.Tracer("").Start(ctx, "userlist")
	defer span.End()

	filter, err := parseUserFilters(req.URL.Query())
	if err != nil {
		s.logger.Error(err.Error())
		return nil, errors.New("invalid request params")
	}
	span.AddEvent("filter")

	users, err := s.repo.GetList(ctx, filter)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, errors.New("users not found")
	}

	return users, nil
}

func (s *UserService) UserAdd(req *http.Request) (*models.User, error) {
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

	err := s.repo.Add(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UserUpdate(userID int, req *http.Request) error {
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

	return s.repo.Update(ctx, userID, userUpdateReq)
}

func (s *UserService) UserDelete(userID int, req *http.Request) error {
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	return s.repo.Delete(ctx, userID)
}

func (s *UserService) FindByID(ctx context.Context, userID int) (*models.User, error) {
	return s.repo.FindByID(ctx, userID)
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
