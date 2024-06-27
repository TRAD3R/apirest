package storage

import (
	"context"
	"errors"
	"github.com/trad3r/hskills/apirest/models"
	"strings"
	"sync"
	"time"
)

type PostStorage struct {
	posts map[int]*models.Post
	mu    sync.RWMutex
}

type PostFilter struct {
	Offset        int
	Limit         int
	FromCreatedAt time.Time
	ToCreatedAt   time.Time
	Subject       string
	Author        int
	Authors       []int
}

type PostAddRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Author  int    `json:"author"`
}

type PostUpdateRequest struct {
	Subject string `json:"subject,omitempty"`
	Body    string `json:"body,omitempty"`
}

func NewPostStorage() *PostStorage {
	return &PostStorage{
		posts: make(map[int]*models.Post),
	}
}

// Add adds new post to post map
func (s *PostStorage) Add(ctx context.Context, post models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	s.mu.RLock()
	defer s.mu.RUnlock()
	id := s.getNextPostID()
	post.ID = id
	post.CreatedAt = time.Now()
	s.posts[id] = &post

	return nil
}

// GetList returns post list
func (s *PostStorage) GetList(ctx context.Context, filter PostFilter) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	var posts []models.Post

	for _, post := range s.posts {
		if filter.FromCreatedAt.Year() > 2023 && post.CreatedAt.After(filter.FromCreatedAt) {
			posts = append(posts, *post)
			continue
		}

		if filter.ToCreatedAt.Year() > 2023 && post.CreatedAt.Before(filter.ToCreatedAt) {
			posts = append(posts, *post)
			continue
		}

		if strings.Contains(post.Subject, filter.Subject) {
			posts = append(posts, *post)
			continue
		}

		for _, user := range filter.Authors {
			if post.Author == user {
				posts = append(posts, *post)
				break
			}
		}
	}

	return s.sort(posts, filter), nil
}

// Update updates post's subject or body
func (s *PostStorage) Update(ctx context.Context, id int, postReq PostUpdateRequest) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[id]
	if !ok {
		return errors.New("post does not found")
	}

	post.Subject = postReq.Subject
	post.Body = postReq.Body // typo
	post.UpdatedAt = time.Now()

	return nil
}

// Delete removes post by ID
func (s *PostStorage) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.posts[id]
	if !ok {
		return nil
	}

	delete(s.posts, id)

	return nil
}

// getNextID returns next post ID
func (s *PostStorage) getNextPostID() int {
	var maxID int

	if len(s.posts) == 0 {
		return 1
	}

	for _, post := range s.posts {
		if post.ID > maxID {
			maxID = post.ID
		}
	}

	return maxID + 1
}

// sort returns a set of posts
func (s *PostStorage) sort(posts []models.Post, filter PostFilter) []models.Post {
	if filter.Offset >= len(posts) {
		return []models.Post{}
	}

	if filter.Limit+filter.Offset >= len(posts) {
		filter.Limit = len(posts)
	}

	return posts[filter.Offset : filter.Offset+filter.Limit]
}
