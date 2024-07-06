package storage

import (
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/trad3r/hskills/apirest/internal/custom_errors"
	"github.com/trad3r/hskills/apirest/internal/models"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	su *UserStorage
)

func TestMain(m *testing.M) {
	su = NewUserStorage()

	for range 10 {
		u := models.User{
			Name:        faker.Name(),
			Phonenumber: faker.Phonenumber(),
		}

		_, _ = su.Add(context.Background(), u)
	}
}

func TestPostAdd(t *testing.T) {
	t.Parallel()

	sp := NewPostStorage()
	userLastId := su.getNextUserID() - 1
	for range 10 {
		t.Run("add new post", func(t *testing.T) {
			post := models.Post{
				Subject: faker.Sentence(),
				Body:    faker.Paragraph(),
				Author:  rand.Intn(userLastId),
			}

			err := sp.Add(context.Background(), &post)
			require.NoError(t, err)
			assert.Greater(t, post.ID, 0)
		})
	}

	require.Equal(t, 10, len(sp.posts))
}

func BenchmarkPostAdd(b *testing.B) {
	sp := NewPostStorage()
	userLastId := su.getNextUserID() - 1
	for i := 0; i < b.N; i++ {
		b.Run("add new post", func(b *testing.B) {
			p := models.Post{
				Subject: faker.Sentence(),
				Body:    faker.Paragraph(),
				Author:  rand.Intn(userLastId),
			}

			err := sp.Add(context.Background(), &p)
			assert.NoError(b, err)
		})
	}
}

func TestPostUpdate(t *testing.T) {
	t.Parallel()

	sp := NewPostStorage()
	t.Run("update post subject and body", func(t *testing.T) {
		t.Parallel()
		postID := 1
		post := sp.posts[postID]
		require.NotNil(t, post)

		req := PostUpdateRequest{
			Subject: faker.Sentence(),
			Body:    faker.Paragraph(),
		}

		err := sp.Update(context.Background(), postID, req)
		require.NoError(t, err)

		newPost := sp.posts[postID]
		require.Equal(t, req.Subject, newPost.Subject)
		require.Equal(t, post.Body, newPost.Body)
	})
}

func TestPostGetList(t *testing.T) {
	t.Parallel()

	sp := NewPostStorage()
	testCases := []struct {
		name          string
		users         []int
		subject       string
		from          time.Time
		to            time.Time
		limit         int
		offset        int
		orderBy       string
		expectedCount int
	}{
		{
			name:          "get posts by user",
			users:         []int{1},
			subject:       "",
			from:          time.Now().AddDate(-2, 0, 0),
			to:            time.Now().AddDate(-2, 0, 0),
			limit:         10,
			offset:        0,
			orderBy:       "asc",
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filter := PostFilter{
				Offset:        tc.offset,
				Limit:         tc.limit,
				FromCreatedAt: tc.from,
				ToCreatedAt:   tc.to,
				Subject:       tc.subject,
			}
			users, err := sp.GetList(context.Background(), filter)
			require.NoError(t, err)
			require.Equal(t, tc.expectedCount, len(users))
		})
	}

}

func TestPostDelete(t *testing.T) {
	t.Parallel()
	sp := NewPostStorage()

	newPost := models.Post{
		Subject: faker.Sentence(),
		Body:    faker.Paragraph(),
		Author:  1,
	}

	err := sp.Add(context.Background(), &newPost)
	require.NoError(t, err)

	post, err := sp.FindById(context.Background(), newPost.ID)
	require.NoError(t, err)
	require.Equal(t, newPost, post)

	err = sp.Delete(context.Background(), newPost.ID)
	require.NoError(t, err)

	post, err = sp.FindById(context.Background(), newPost.ID)
	require.ErrorIs(t, err, custom_errors.ErrPostNotFound)
}
