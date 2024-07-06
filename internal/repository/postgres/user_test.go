package postgres

// Blackbox testing

//func TestUserAdd(t *testing.T) {
//	t.Parallel()
//	// Тесткейсы отмеченные как параллельные откладываются в отдельный стек
//	// и после выполнения всех синхронных тестов запускаются в общем пулле паралелльно
//
//	for range 10 {
//		t.Run("add new user", func(t *testing.T) {
//			t.Parallel()
//
//			// Создавать новый сторадж под каждый тестовый случай
//			su := NewUserStorage()
//
//			u := models.User{
//				Name:        faker.Name(),
//				Phonenumber: faker.Phonenumber(),
//			}
//
//			user, err := su.Add(context.Background(), u)
//			require.NoError(t, err)
//			require.NotNil(t, user)
//		})
//	}
//}
//
//func BenchmarkUserAdd(b *testing.B) {
//	su := NewUserStorage()
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
//
//func TestUserUpdateName(t *testing.T) {
//	su := NewUserStorage()
//
//	t.Run("update user name", func(t *testing.T) {
//		t.Parallel()
//		userToAdd := models.User{
//			Name:        faker.Name(),
//			Phonenumber: faker.Phonenumber(),
//		}
//
//		addedUser, err := su.Add(context.Background(), userToAdd)
//		require.NoError(t, err)
//		require.Greater(t, addedUser.ID, 0)
//
//		addedUser.Name = faker.Name()
//		updateNameReq := UserUpdateRequest{
//			Name: addedUser.Name,
//		}
//
//		err = su.Update(context.Background(), addedUser.ID, updateNameReq)
//		require.NoError(t, err)
//
//		foundUser, err := su.FindById(context.Background(), addedUser.ID)
//		require.NoError(t, err)
//		require.Equal(t, addedUser, foundUser)
//	})
//}
//
//func TestUserUpdatePhone(t *testing.T) {
//	su := NewUserStorage()
//
//	t.Run("update user phone", func(t *testing.T) {
//		t.Parallel()
//		userToAdd := models.User{
//			Name:        faker.Name(),
//			Phonenumber: faker.Phonenumber(),
//		}
//
//		addedUser, err := su.Add(context.Background(), userToAdd)
//		require.NoError(t, err)
//		require.Greater(t, addedUser.ID, 0)
//
//		addedUser.Phonenumber = faker.Phonenumber()
//		updatePhoneReq := UserUpdateRequest{
//			Phonenumber: addedUser.Phonenumber,
//		}
//
//		err = su.Update(context.Background(), addedUser.ID, updatePhoneReq)
//		require.NoError(t, err)
//
//		foundUser, err := su.FindById(context.Background(), addedUser.ID)
//		require.NoError(t, err)
//		require.Equal(t, addedUser, foundUser)
//	})
//}
//
//func TestUserUpdateNameAndPhone(t *testing.T) {
//	su := NewUserStorage()
//
//	t.Run("update user phone", func(t *testing.T) {
//		t.Parallel()
//		userToAdd := models.User{
//			Name:        faker.Name(),
//			Phonenumber: faker.Phonenumber(),
//		}
//
//		addedUser, err := su.Add(context.Background(), userToAdd)
//		require.NoError(t, err)
//		require.Greater(t, addedUser.ID, 0)
//
//		addedUser.Name = faker.Name()
//		addedUser.Phonenumber = faker.Phonenumber()
//		updatePhoneReq := UserUpdateRequest{
//			Name:        addedUser.Name,
//			Phonenumber: addedUser.Phonenumber,
//		}
//
//		err = su.Update(context.Background(), addedUser.ID, updatePhoneReq)
//		require.NoError(t, err)
//
//		foundUser, err := su.FindById(context.Background(), addedUser.ID)
//		require.NoError(t, err)
//		require.Equal(t, addedUser, foundUser)
//	})
//}
//
//func TestUserGetList(t *testing.T) {
//	su := NewUserStorage()
//	for i := 0; i < 10; i++ {
//		userToAdd := models.User{
//			Name:        faker.Name(),
//			Phonenumber: faker.Phonenumber(),
//		}
//
//		_, err := su.Add(context.Background(), userToAdd)
//		require.NoError(t, err)
//	}
//
//	testCases := []struct {
//		name          string
//		names         []string
//		from          time.Time
//		to            time.Time
//		limit         int
//		offset        int
//		orderBy       string
//		expectedCount int
//	}{
//		{
//			name:          "get users by names",
//			names:         []string{su.users[1].Name, su.users[2].Name},
//			from:          time.Time{},
//			to:            time.Time{},
//			limit:         10,
//			offset:        0,
//			orderBy:       "asc",
//			expectedCount: 2,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			t.Parallel()
//
//			filter := UserFilter{
//				Offset:         tc.offset,
//				Limit:          tc.limit,
//				FromCreatedAt:  &tc.from,
//				ToCreatedAt:    &tc.to,
//				Name:           tc.names,
//				TopPostsAmount: tc.orderBy,
//			}
//			users, err := su.GetList(context.Background(), filter)
//			require.NoError(t, err)
//			require.Equal(t, tc.expectedCount, len(users))
//		})
//	}
//
//}
//
//func TestUserDelete(t *testing.T) {
//	t.Parallel()
//	su := NewUserStorage()
//
//	userToAdd := models.User{
//		Name:        faker.Name(),
//		Phonenumber: faker.Phonenumber(),
//	}
//
//	addedUser, err := su.Add(context.Background(), userToAdd)
//	require.NoError(t, err)
//
//	user, err := su.FindById(context.Background(), addedUser.ID)
//	require.NoError(t, err)
//	require.Equal(t, addedUser.ID, user.ID)
//
//	err = su.Delete(context.Background(), user.ID)
//	require.NoError(t, err)
//
//	user, err = su.FindById(context.Background(), addedUser.ID)
//	require.ErrorIs(t, err, custom_errors.ErrUserNotFound)
//}
