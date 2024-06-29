package storage

import (
	"context"
	"errors"
	"github.com/trad3r/hskills/apirest/cmd/custom_errors"
	"github.com/trad3r/hskills/apirest/models"
	"sync"
	"time"
)

type UserStorage struct {
	users map[int]*models.User
	mu    sync.RWMutex
}

type UserFilter struct {
	Offset         int
	Limit          int
	FromCreatedAt  *time.Time
	ToCreatedAt    *time.Time
	Name           []string
	TopPostsAmount string
}

type UserAddRequest struct {
	Name        string `json:"name"`
	Phonenumber string `json:"phonenumber"`
}

type UserUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Phonenumber string `json:"phonenumber,omitempty"`
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		users: make(map[int]*models.User),
	}
}

// Add adds new user to user map
func (s *UserStorage) Add(ctx context.Context, user models.User) (*models.User, error) {
	// ctx не используется
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	s.mu.RLock()
	defer s.mu.RUnlock()
	id := s.getNextUserID()
	user.ID = id
	user.CreatedAt = time.Now()
	s.users[id] = &user

	return &user, nil
}

// GetList returns user list
func (s *UserStorage) GetList(_ context.Context, filter UserFilter) ([]models.User, error) {
	var users []models.User

	// go run -race
	// mutex Lock
	// lifehack with mutex
	s.mu.RLock()
	for _, user := range s.users {
		//s.mu.RUnlock()

		// AND
		// FROM: 2022 TO: 2024

		if filter.FromCreatedAt != nil && user.CreatedAt.Before(*filter.FromCreatedAt) {
			continue
		}

		if filter.ToCreatedAt != nil && user.CreatedAt.After(*filter.ToCreatedAt) {
			continue
		}

		match := true
		if len(filter.Name) > 0 {
			match = false
			for _, name := range filter.Name {
				if name == user.Name {
					match = true
					break
				}
			}
		}

		if match {
			users = append(users, *user)
		}

		//s.mu.RLock()
	}
	s.mu.RUnlock()

	return sort(users, filter.Limit, filter.Offset, filter.TopPostsAmount), nil
}

// Update updates user's name or phone
func (s *UserStorage) Update(_ context.Context, id int, userReq UserUpdateRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[id]
	if !ok {
		return errors.New("user does not found")
	}

	if len(userReq.Name) > 0 {
		user.Name = userReq.Name
		user.UpdatedAt = time.Now()
	}

	if len(userReq.Phonenumber) > 0 {
		user.Phonenumber = userReq.Phonenumber
		user.UpdatedAt = time.Now()
	}

	return nil
}

// Delete removes user by ID
func (s *UserStorage) Delete(_ context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.users[id]
	if !ok {
		return nil // вернуть ошибку HTTP DELETE /:id
	}

	delete(s.users, id)

	return nil
}

// FindById returns user by ID
func (s *UserStorage) FindById(_ context.Context, id int) (*models.User, error) {
	user, ok := s.users[id]
	if !ok {
		return nil, custom_errors.ErrUserNotFound
	}

	return user, nil
}

// getNextID returns next user ID
func (s *UserStorage) getNextUserID() int {
	var maxID int

	if len(s.users) == 0 {
		return 1
	}

	for _, user := range s.users {
		if user.ID > maxID {
			maxID = user.ID
		}
	}

	return maxID + 1
}

// Increment user posts
func (s *UserStorage) IncrPostToUser(userId int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.ID == userId {
			user.PostCount++
			break
		}
	}

	return nil
}

// Decrement user posts
func (s *UserStorage) DecrPostToUser(userId int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.ID == userId && user.PostCount > 0 {
			user.PostCount--
			break
		}
	}

	return nil
}

// sort returns a set of users sorted by number of posts
func sort(users []models.User, limit, offset int, order string) []models.User {
	if offset >= len(users) {
		return []models.User{}
	}

	for i := 0; i < len(users)-1; i++ {
		for j := i; j < len(users); j++ {
			switch order {
			case "desc":
				if users[i].PostCount < users[j].PostCount {
					users[i], users[j] = users[j], users[i]
				}
				break
			default:
				if users[i].PostCount > users[j].PostCount {
					users[i], users[j] = users[j], users[i]
				}
			}
		}
	}

	if limit+offset >= len(users) {
		limit = len(users)
	}

	return users[offset : offset+limit]
}
