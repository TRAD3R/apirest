package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"github.com/trad3r/hskills/apirest/internal/migrator"
	"github.com/trad3r/hskills/apirest/internal/repository/filters"
	"github.com/trad3r/hskills/apirest/internal/repository/postgres"
	"github.com/trad3r/hskills/apirest/internal/storage"
	"github.com/trad3r/hskills/apirest/internal/testutils"
)

//func TestPostAdd(t *testing.T) {
//
//	var err error
//	ctx := context.Background()
//
//	pgStorage := getPostRepo(t)
//
//	testUser, err := pgStorage.User.FindById(ctx, 1)
//	require.NoError(t, err)
//
//	post := &models.Post{
//		Subject: faker.Sentence(),
//		Body:    faker.Paragraph(),
//		Author:  *testUser,
//	}
//
//	err = pgRepo.Add(ctx, post)
//	require.NoError(t, err)
//	require.NotEmpty(t, post.ID)
//
//	dbPost, err := pgRepo.FindById(ctx, post.ID)
//	require.NoError(t, err)
//	assert.Equal(t, post.Subject, dbPost.Subject)
//	assert.Equal(t, post.Body, dbPost.Body)
//	assert.NotEmpty(t, dbPost.CreatedAt)
//	assert.Empty(t, dbPost.UpdatedAt)
//}

func TestPostUpdate(t *testing.T) {
	t.Parallel()

	var err error
	ctx := context.Background()

	pgRepo := getPostRepo(t)

	postID := 1
	post, err := pgRepo.FindById(ctx, postID)
	require.NoError(t, err)
	require.NotNil(t, post)

	req := filters.PostUpdateRequest{
		Subject: faker.Sentence(),
		Body:    faker.Paragraph(),
	}

	err = pgRepo.Update(ctx, postID, req)
	require.NoError(t, err)

	newPost, err := pgRepo.FindById(ctx, postID)
	require.NoError(t, err)
	require.Equal(t, req.Subject, newPost.Subject)
	require.Equal(t, req.Body, newPost.Body)
	require.NotEmpty(t, newPost.UpdatedAt)
}

func TestPostGetList(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pgRepo := getPostRepo(t)

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
			from:          time.Time{},
			to:            time.Time{},
			limit:         10,
			offset:        0,
			orderBy:       "asc",
			expectedCount: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filter := filters.PostFilter{
				Offset:        tc.offset,
				Limit:         tc.limit,
				FromCreatedAt: tc.from,
				ToCreatedAt:   tc.to,
				Subject:       tc.subject,
				Authors:       tc.users,
			}
			posts, err := pgRepo.GetList(ctx, filter)
			require.NoError(t, err)
			require.Equal(t, tc.expectedCount, len(posts))
		})
	}

}

func TestPostDelete(t *testing.T) {
	t.Parallel()

	var err error
	ctx := context.Background()

	pgRepo := getPostRepo(t)

	post, err := pgRepo.FindById(ctx, 1)
	require.NoError(t, err)

	err = pgRepo.Delete(ctx, post.ID)
	require.NoError(t, err)

	post, err = pgRepo.FindById(ctx, post.ID)
	require.NoError(t, err)
	require.Nil(t, post)
}

func getPostRepo(t *testing.T) postgres.IPostRepository {
	dsn := testutils.PreparePostgres(t)
	err := migrator.ApplyPostgresMigrations("../../../migrations", dsn)
	require.NoError(t, err)

	err = testutils.RunFixtures("../../../fixtures", dsn)
	require.NoError(t, err)

	db, err := storage.NewDB(context.Background(), dsn)
	require.NoError(t, err)

	return postgres.NewPostRepository(db)
}

//func BenchmarkPostAdd(b *testing.B) {
//	sp := NewPostStorage()
//	userLastId := su.getNextUserID() - 1
//	for i := 0; i < b.N; i++ {
//		b.Run("add new post", func(b *testing.B) {
//			p := models.Post{
//				Subject: faker.Sentence(),
//				Body:    faker.Paragraph(),
//				Author:  rand.Intn(userLastId),
//			}
//
//			err := sp.Add(context.Background(), &p)
//			assert.NoError(b, err)
//		})
//	}
//}
//
