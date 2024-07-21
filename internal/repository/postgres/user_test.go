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
	"github.com/trad3r/hskills/apirest/internal/storage"
	"github.com/trad3r/hskills/apirest/internal/testutils"
)

func TestUserAdd(t *testing.T) {
	t.Parallel()

	var err error
	ctx := context.Background()

	pgStorage := setup(t)

	testUser := &models.User{
		Name:        faker.Name(),
		Phonenumber: faker.Phonenumber(),
	}

	err = pgStorage.User.Add(ctx, testUser)
	require.NoError(t, err)
	require.NotEmpty(t, testUser.ID)

	dbUser, err := pgStorage.User.FindById(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, testUser.Name, dbUser.Name)
	assert.Equal(t, testUser.Phonenumber, dbUser.Phonenumber)
	assert.NotEmpty(t, dbUser.CreatedAt)
	assert.Empty(t, dbUser.UpdatedAt)
}

func TestUserGetList(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pgStorage := setup(t)

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

			list, err := pgStorage.User.GetList(ctx, tc.filter)
			require.NoError(t, err)
			assert.Equal(t, tc.count, len(list))
			assert.Equal(t, tc.expectedId, list[0].ID)
		})
	}
}

func TestUserUpdateName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pgStorage := setup(t)

	user, err := pgStorage.User.FindById(ctx, 1)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	user.Name = faker.Name()

	err = pgStorage.User.Update(context.Background(), user.ID, filters.UserUpdateRequest{
		Name: user.Name,
	})
	require.NoError(t, err)

	foundUser, err := pgStorage.User.FindById(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, user.Name, foundUser.Name)
}

func TestUserUpdatePhone(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pgStorage := setup(t)

	user, err := pgStorage.User.FindById(ctx, 1)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	user.Phonenumber = faker.Phonenumber()

	err = pgStorage.User.Update(context.Background(), user.ID, filters.UserUpdateRequest{
		Phonenumber: user.Phonenumber,
	})
	require.NoError(t, err)

	foundUser, err := pgStorage.User.FindById(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, user.Phonenumber, foundUser.Phonenumber)
}

func TestUserUpdateNameAndPhone(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pgStorage := setup(t)

	user, err := pgStorage.User.FindById(ctx, 1)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	user.Name = faker.Name()
	user.Phonenumber = faker.Phonenumber()

	err = pgStorage.User.Update(context.Background(), user.ID, filters.UserUpdateRequest{
		Name:        user.Name,
		Phonenumber: user.Phonenumber,
	})
	require.NoError(t, err)

	foundUser, err := pgStorage.User.FindById(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, user.Name, foundUser.Name)
	require.Equal(t, user.Phonenumber, foundUser.Phonenumber)
}

func TestUserDelete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pgStorage := setup(t)

	user, err := pgStorage.User.FindById(ctx, 1)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	err = pgStorage.User.Delete(context.Background(), user.ID)
	require.NoError(t, err)

	user, err = pgStorage.User.FindById(context.Background(), user.ID)
	require.NoError(t, err)
	require.Empty(t, user)
}

func setup(t *testing.T) *storage.Storage {
	pgStorage, dsnStr := testutils.PreparePostgres(t)
	err := migrator.ApplyPostgresMigrations("../../../migrations", dsnStr)
	require.NoError(t, err)

	err = testutils.RunFixtures("../../../fixtures", dsnStr)
	require.NoError(t, err)

	return pgStorage
}

//func BenchmarkUserAdd(b *testing.B) {
//	ctx := context.Background()
//
//	pgStorage := setup(b)
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
