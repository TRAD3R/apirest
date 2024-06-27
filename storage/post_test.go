package storage

import (
	"context"
	"github.com/trad3r/hskills/apirest/models"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	sp *PostStorage
	su *UserStorage
)

func TestMain(m *testing.M) {
	sp = NewPostStorage()
	su = NewUserStorage()
}

func TestPostAdd(t *testing.T) {
	userLastId := sp.getNextPostID() - 1
	for range 10 {
		t.Run("add new post", func(t *testing.T) {
			post := models.Post{
				Subject: faker.Sentence(),
				Body:    faker.Paragraph(),
				Author:  rand.Intn(userLastId),
			}

			err := sp.Add(context.Background(), post)
			require.NoError(t, err)
		})
	}

	require.Equal(t, 10, len(sp.posts))
}

func BenchmarkPostAdd(b *testing.B) {
	userLastId := su.getNextUserID() - 1
	for i := 0; i < b.N; i++ {
		b.Run("add new post", func(b *testing.B) {
			p := models.Post{
				Subject: faker.Sentence(),
				Body:    faker.Paragraph(),
				Author:  rand.Intn(userLastId),
			}

			err := sp.Add(context.Background(), p)
			assert.NoError(b, err)
		})
	}
}

func TestPostUpdate(t *testing.T) {
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
	//t.Parallel()
	//userID := 1
	//user := sp.posts[userID]
	//require.NotNil(t, user)
	//
	//err := sp.Delete(context.Background(), userID)
	//require.NoError(t, err)
	//require.Equal(t, 9, len(s.users))
}
