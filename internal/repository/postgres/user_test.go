package postgres_test

// Blackbox testing

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trad3r/hskills/apirest/internal/migrator"
	"github.com/trad3r/hskills/apirest/internal/models"
	"github.com/trad3r/hskills/apirest/internal/repository/filters"
	"github.com/trad3r/hskills/apirest/internal/repository/postgres"
	"github.com/trad3r/hskills/apirest/internal/storage"
	"github.com/trad3r/hskills/apirest/internal/testutils"
)

func TestUserAdd(t *testing.T) {
	t.Parallel()

	var err error
	ctx := context.Background()

	pgRepo := getUserRepo(t)

	testUser := &models.User{
		Name:        faker.Name(),
		Phonenumber: faker.Phonenumber(),
	}

	err = pgRepo.Add(ctx, testUser)
	require.NoError(t, err)
	require.NotEmpty(t, testUser.ID)

	dbUser, err := pgRepo.FindById(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, testUser.Name, dbUser.Name)
	assert.Equal(t, testUser.Phonenumber, dbUser.Phonenumber)
	assert.NotEmpty(t, dbUser.CreatedAt)
	assert.Empty(t, dbUser.UpdatedAt)
}

func TestUserGetList(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pgRepo := getUserRepo(t)

	testCases := []struct {
		name       string
		filter     filters.UserFilter
		count      int
		expectedId int
	}{
		{
			name:       "All asc",
			filter:     filters.UserFilter{},
			count:      5,
			expectedId: 5,
		},
		{
			name: "All desc",
			filter: filters.UserFilter{
				TopPostsAmount: "desc",
			},
			count:      5,
			expectedId: 1,
		},
		{
			name: "Limit 2",
			filter: filters.UserFilter{
				Limit: 2,
			},
			count:      2,
			expectedId: 5,
		},
		{
			name: "Offset 2",
			filter: filters.UserFilter{
				Offset: 2,
				Limit:  2,
			},
			count:      2,
			expectedId: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			list, err := pgRepo.GetList(ctx, tc.filter)
			require.NoError(t, err)
			assert.Equal(t, tc.count, len(list))
			assert.Equal(t, tc.expectedId, list[0].ID)
		})
	}
}

func TestUserUpdateName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pgRepo := getUserRepo(t)

	user, err := pgRepo.FindById(ctx, 1)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	user.Name = faker.Name()

	err = pgRepo.Update(context.Background(), user.ID, filters.UserUpdateRequest{
		Name: user.Name,
	})
	require.NoError(t, err)

	foundUser, err := pgRepo.FindById(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, user.Name, foundUser.Name)
}

func TestUserUpdatePhone(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pgRepo := getUserRepo(t)

	user, err := pgRepo.FindById(ctx, 1)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	user.Phonenumber = faker.Phonenumber()

	err = pgRepo.Update(context.Background(), user.ID, filters.UserUpdateRequest{
		Phonenumber: user.Phonenumber,
	})
	require.NoError(t, err)

	foundUser, err := pgRepo.FindById(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, user.Phonenumber, foundUser.Phonenumber)
}

func TestUserUpdateNameAndPhone(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pgRepo := getUserRepo(t)

	user, err := pgRepo.FindById(ctx, 1)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	user.Name = faker.Name()
	user.Phonenumber = faker.Phonenumber()

	err = pgRepo.Update(context.Background(), user.ID, filters.UserUpdateRequest{
		Name:        user.Name,
		Phonenumber: user.Phonenumber,
	})
	require.NoError(t, err)

	foundUser, err := pgRepo.FindById(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, user.Name, foundUser.Name)
	require.Equal(t, user.Phonenumber, foundUser.Phonenumber)
}

func TestUserDelete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pgRepo := getUserRepo(t)

	user, err := pgRepo.FindById(ctx, 1)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	err = pgRepo.Delete(context.Background(), user.ID)
	require.NoError(t, err)

	user, err = pgRepo.FindById(context.Background(), user.ID)
	require.NoError(t, err)
	require.Empty(t, user)
}

func getUserRepo(t *testing.T) postgres.IUserRepository {
	dsn := testutils.PreparePostgres(t)
	err := migrator.ApplyPostgresMigrations("../../../migrations", dsn)
	require.NoError(t, err)

	err = testutils.RunFixtures("../../../fixtures", dsn)
	require.NoError(t, err)

	db, err := storage.NewDB(context.Background(), dsn)
	require.NoError(t, err)

	return postgres.NewUserRepository(db)
}

//func BenchmarkUserAdd(b *testing.B) {
//	ctx := context.Background()
//
//	pgStorage := getUserRepo(b)
//
//	for i := 0; i < b.N; i++ {
//		b.Run("add new user", func(b *testing.B) {
//			u := models.User{
//				Name:        faker.Name(),
//				Phonenumber: faker.Phonenumber(),
//			}
//
//			user, err := su.Add(context.Background(), u)
//			require.NoError(b, err)
//			require.NotNil(b, user)
//		})
//	}
//}
