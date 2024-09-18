package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/trad3r/hskills/apirest/internal/customerrors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trad3r/hskills/apirest/internal/models"
	"github.com/trad3r/hskills/apirest/internal/repository/filters"
	"go.opentelemetry.io/otel"
)

type IUserRepository interface {
	Add(ctx context.Context, user *models.User) error
	GetList(ctx context.Context, filter filters.UserFilter) ([]models.User, error)
	Update(ctx context.Context, id int, userReq filters.UserUpdateRequest) error
	Delete(ctx context.Context, id int) error
	FindByID(ctx context.Context, id int) (*models.User, error)
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) IUserRepository {
	return UserRepository{
		db: db,
	}
}

// Add adds new user
func (s UserRepository) Add(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	ds := goqu.Insert("author").
		Cols("name", "phonenumber").
		Vals(
			goqu.Vals{user.Name, user.Phonenumber},
		).
		Returning("id")

	sql, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("error while creating sql: %w", err)
	}

	err = s.db.QueryRow(ctx, sql, args...).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("error while inserting user: %w", err)
	}

	return nil
}

// GetList returns user list
func (s UserRepository) GetList(ctx context.Context, filter filters.UserFilter) ([]models.User, error) {
	ctx, span := otel.Tracer("").Start(ctx, "getUserList")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	var errs error
	users := make([]models.User, 0, filter.Limit)
	var wheres []goqu.Expression

	postCountSubquery := goqu.From(goqu.T("post").As("p")).Select(goqu.COUNT("id")).Where(goqu.I("p.author_id").Eq(goqu.I("a.id")))

	ds := goqu.From(goqu.T("author").As("a")).
		Select("a.id", "a.name", "a.phonenumber", "a.created_at", "a.updated_at", postCountSubquery.As("post_count"))

	if filter.FromCreatedAt != nil {
		wheres = append(wheres, goqu.C("created_at").Gte(filter.FromCreatedAt))
	}

	if filter.ToCreatedAt != nil {
		wheres = append(wheres, goqu.C("created_at").Lte(filter.ToCreatedAt))
	}

	if len(filter.Name) > 0 {
		wheres = append(wheres, goqu.C("name").In(filter.Name))
	}

	if len(wheres) > 0 {
		ds = ds.Where(wheres...)
	}

	ds = ds.
		Offset(filter.Offset).
		Limit(filter.Limit)

	if filter.TopPostsAmount == "desc" {
		ds = ds.Order(goqu.C("post_count").Desc())
	} else {
		ds = ds.Order(goqu.C("post_count").Asc())
	}

	sql, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("error while creating sql: %w", err)
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error while querying users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Phonenumber, &user.CreatedAt, &user.UpdatedAt, &user.PostCount)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		errs = multierror.Append(errs, err)
	}

	return users, errs
}

// Update updates user's name or phone
func (s UserRepository) Update(ctx context.Context, ID int, userReq filters.UserUpdateRequest) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	var userID int

	ds := goqu.Update("author").
		Where(goqu.Ex{"ID": ID}).
		Returning("id")

	updates := make(map[string]interface{}, 3)
	if len(userReq.Name) > 0 {
		updates["name"] = userReq.Name
	}

	if len(userReq.Phonenumber) > 0 {
		updates["phonenumber"] = userReq.Phonenumber
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
	}

	if len(updates) > 0 {
		ds = ds.Set(updates)
	}

	sql, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("error while preparing update user: %w", err)
	}

	err = s.db.QueryRow(ctx, sql, args...).Scan(&userID)
	if err != nil {
		return fmt.Errorf("error while updating user: %w", err)
	}

	if userID == 0 {
		return customerrors.ErrUserNotFound
	}

	return nil
}

// Delete removes user by ID
func (s UserRepository) Delete(ctx context.Context, ID int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	ds := goqu.Delete("author").Where(goqu.Ex{"id": ID}).Returning("id")
	sql, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("error while preparing delete user: %w", err)
	}

	err = s.db.QueryRow(ctx, sql, args...).Scan(&ID)
	if err != nil {
		return fmt.Errorf("error while deleting user: %w", err)
	}

	log.Println(ID)
	return nil
}

// FindByID returns user by ID
func (s UserRepository) FindByID(ctx context.Context, ID int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	ds := goqu.From("author").
		Select("id", "name", "phonenumber", "created_at", "updated_at").
		Where(goqu.Ex{"id": ID})

	sql, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("error while preparing find user by ID %d: %w", ID, err)
	}

	var user models.User

	if err := s.db.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Name, &user.Phonenumber, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("error while find user by ID %d: %w", ID, err)
	}

	return &user, nil
}
