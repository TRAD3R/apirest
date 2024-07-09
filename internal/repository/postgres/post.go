package postgres

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trad3r/hskills/apirest/internal/models"
	"github.com/trad3r/hskills/apirest/internal/repository/filters"
	"time"
)

type PostRepository struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

// Add adds new post
func (s PostRepository) Add(ctx context.Context, post *models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	ds := goqu.Insert("post").
		Cols("subject", "body", "author_id").
		Vals(
			goqu.Vals{post.Subject, post.Body, post.Author.ID},
		)

	sql, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("error while creating sql: %w", err)
	}

	_, err = s.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error while inserting post: %w", err)
	}

	return nil
}

// GetList returns user list
func (s PostRepository) GetList(ctx context.Context, filter filters.PostFilter) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	var errs error
	var posts []models.Post
	var wheres []goqu.Expression

	ds := goqu.From(goqu.T("author").As("a")).
		Select("p.id", "p.subject", "p.body", "p.created_at", "p.updated_at", "a.id", "a.name", "a.phonenumber", "a.created_at", "a.updated_at").
		Join(goqu.T("post").As("p"), goqu.On(goqu.Ex{"a.id": goqu.I("p.author_id")}))

	if !filter.FromCreatedAt.IsZero() {
		wheres = append(wheres, goqu.C("p.created_at").Gte(filter.FromCreatedAt))
	}

	if !filter.ToCreatedAt.IsZero() {
		wheres = append(wheres, goqu.C("p.created_at").Lte(filter.ToCreatedAt))
	}

	if len(filter.Authors) > 0 {
		wheres = append(wheres, goqu.C("p.author_id").In(filter.Authors))
	}

	if len(wheres) > 0 {
		ds = ds.Where(wheres...)
	}

	ds.
		Offset(uint(filter.Offset)).
		Limit(uint(filter.Limit))

	sql, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("error while creating sql: %w", err)
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error while querying posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		var post models.Post

		err := rows.Scan(&post.ID, &post.Subject, &post.Body, &post.CreatedAt, &post.UpdatedAt, &user.ID,
			&user.Name, &user.Phonenumber, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		post.Author = user
		posts = append(posts, post)
	}

	return posts, errs
}

// Example of buildUpdatemap
//func buildUpdateMap(postReq filters.PostUpdateRequest) (m map[string]any, err error) {
//	// make map
//	if len(postReq.Body) > 0 {
//		m =
//	}
//
//	if len(m) == 0 {
//		// return err
//	}
//}

// Update updates user's name or phone
func (s PostRepository) Update(ctx context.Context, id int, postReq filters.PostUpdateRequest) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	ds := goqu.Update("post").
		Where(goqu.Ex{"id": id})

	var sets []goqu.Record
	if len(postReq.Subject) > 0 {
		sets = append(sets, goqu.Record{"subject": postReq.Subject})
	}

	if len(postReq.Body) > 0 {
		sets = append(sets, goqu.Record{"body": postReq.Body})
	}

	if len(sets) > 0 {
		ds = ds.Set(sets)
	}

	sql, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("error while preparing update post: %w", err)
	}

	_, err = s.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error while updating post: %w", err)
	}

	return nil
}

// Delete removes user by ID
func (s PostRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	ds := goqu.Delete("post").Where(goqu.Ex{"id": id})
	sql, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("error while preparing delete post: %w", err)
	}

	_, err = s.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error while deleting post: %w", err)
	}

	return nil
}

// FindById returns post by ID
func (s PostRepository) FindById(ctx context.Context, id int) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	ds := goqu.From(goqu.T("post").As("p")).
		Select("p.id", "p.subject", "p.body", "a.id", "a.name", "a.phonenumber", "a.created_at").
		Join(goqu.T("author").As("a"), goqu.On(goqu.Ex{"a.id": goqu.I("p.author_id")})).
		Where(goqu.Ex{"p.id": id})

	sql, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("error while preparing find post by ID %d: %w", id, err)
	}

	var author models.User
	var post models.Post

	if err := s.db.QueryRow(ctx, sql, args...).
		Scan(&post.ID, &post.Subject, &post.Body, &author.ID, &author.Name, &author.Phonenumber, &author.CreatedAt, &author.UpdatedAt); err != nil {
		return nil, fmt.Errorf("error while find post by ID %d: %w", id, err)
	}

	post.Author = author

	return &post, nil
}
