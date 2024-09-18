package domain

import (
	"context"
	"github.com/trad3r/hskills/apirest/internal/models"
	"github.com/trad3r/hskills/apirest/internal/repository/filters"
)

type UserRepository interface {
	GetList(ctx context.Context, filter filters.UserFilter) ([]models.User, error)
	Add(ctx context.Context, user *models.User) error
	Update(ctx context.Context, id int, userReq filters.UserUpdateRequest) error
	Delete(ctx context.Context, id int) error
	FindByID(ctx context.Context, id int) (*models.User, error)
}

type PostRepository interface {
	GetList(ctx context.Context, filter filters.PostFilter) ([]models.Post, error)
	Add(ctx context.Context, post *models.Post) error
	Update(ctx context.Context, id int, postReq filters.PostUpdateRequest) error
	Delete(ctx context.Context, id int) error
	FindByID(ctx context.Context, id int) (*models.Post, error)
}
