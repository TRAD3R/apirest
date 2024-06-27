package storage_test

// Blackbox testing

import (
	"context"
	"github.com/trad3r/hskills/apirest/models"
	"github.com/trad3r/hskills/apirest/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var su *storage.UserStorage

func init() {
	su = storage.NewUserStorage()
}

// Пример интеграционного теста
func TestUserAddAndGetAfter(t *testing.T) {
	t.Parallel()

	for range 10 {
		t.Run("add new user", func(t *testing.T) {
			t.Parallel()

			u := models.User{
				Name:        faker.Name(),
				Phonenumber: faker.Phonenumber(),
			}

			err := su.Add(context.Background(), u)
			require.NoError(t, err)
		})
	}
}

func TestUserAdd(t *testing.T) {
	t.Parallel()
	// Тесткейсы отмеченные как параллельные откладываются в отдельный стек
	// и после выполнения всех синхронных тестов запускаются в общем пулле паралелльно

	for range 10 {
		t.Run("add new user", func(t *testing.T) {
			t.Parallel()

			// Создавать новый сторадж под каждый тестовый случай
			su = storage.NewUserStorage()

			u := models.User{
				Name:        faker.Name(),
				Phonenumber: faker.Phonenumber(),
			}

			err := su.Add(context.Background(), u)
			require.NoError(t, err)
		})
	}
}

func BenchmarkUserAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.Run("add new user", func(b *testing.B) {
			u := models.User{
				Name:        faker.Name(),
				Phonenumber: faker.Phonenumber(),
			}

			su.Add(context.Background(), u)
		})
	}
}

func TestUserUpdate(t *testing.T) {
	t.Run("update user name", func(t *testing.T) {
		t.Parallel()
		userID := 1
		user := su.users[userID]
		require.NotNil(t, user)

		// Лучше так
		userToAdd := models.User{}
		su.Add()

		su.Update()

		expectedUser := models.User{}

		su.GetList()

		req := storage.UserUpdateRequest{
			Name:        faker.Name(),
			Phonenumber: "",
		}

		err := su.Update(context.Background(), userID, req)
		require.NoError(t, err)

		newUser := su.users[userID]
		require.Equal(t, userToAdd, expectedUser) // Лучше так

		require.Equal(t, req.Name, newUser.Name)
		require.Equal(t, user.Phonenumber, newUser.Phonenumber)

	})

	t.Run("update user phonenumber", func(t *testing.T) {
		t.Parallel()
		userID := 2
		user := su.users[userID]
		require.NotNil(t, user)

		req := UserUpdateRequest{
			Name:        "",
			Phonenumber: faker.Phonenumber(),
		}

		err := su.Update(context.Background(), userID, req)
		require.NoError(t, err)

		newUser := su.users[userID]
		require.Equal(t, req.Phonenumber, newUser.Phonenumber)
		require.Equal(t, user.Name, newUser.Name)
	})

	t.Run("update user name and phonenumber", func(t *testing.T) {
		t.Parallel()
		userID := 3
		user := su.users[userID]
		require.NotNil(t, user)

		req := UserUpdateRequest{
			Name:        faker.Name(),
			Phonenumber: faker.Phonenumber(),
		}

		err := su.Update(context.Background(), userID, req)
		require.NoError(t, err)

		newUser := su.users[userID]
		require.Equal(t, req.Phonenumber, newUser.Phonenumber)
		require.Equal(t, req.Name, newUser.Name)
	})
}

func TestUserGetList(t *testing.T) {
	testCases := []struct {
		name          string
		names         []string
		from          time.Time
		to            time.Time
		limit         int
		offset        int
		orderBy       string
		expectedCount int
	}{
		{
			name:          "get users by names",
			names:         []string{su.users[1].Name, su.users[2].Name},
			from:          time.Now().AddDate(-2, 0, 0),
			to:            time.Now().AddDate(-2, 0, 0),
			limit:         10,
			offset:        0,
			orderBy:       "asc",
			expectedCount: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filter := UserFilter{
				Offset:         tc.offset,
				Limit:          tc.limit,
				FromCreatedAt:  tc.from,
				ToCreatedAt:    tc.to,
				Name:           tc.names,
				TopPostsAmount: tc.orderBy,
			}
			users, err := su.GetList(context.Background(), filter)
			require.NoError(t, err)
			require.Equal(t, tc.expectedCount, len(users))
		})
	}

}

func TestUserDelete(t *testing.T) {
	t.Parallel()
	userID := 1
	user := su.users[userID]
	require.NotNil(t, user)

	err := su.Delete(context.Background(), userID)
	require.NoError(t, err)
	require.Equal(t, 9, len(su.users))
}
