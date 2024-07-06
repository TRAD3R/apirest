package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/trad3r/hskills/apirest/internal/models"
	"github.com/trad3r/hskills/apirest/internal/router"
	"github.com/trad3r/hskills/apirest/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	s *storage.UserStorage
	p *storage.PostStorage
	r *router.Router
	h *Handler
)

func TestMain(m *testing.M) {
	s = storage.NewUserStorage()
	p = storage.NewPostStorage()
	r = router.NewRouter(s, p)
	h = NewHandler(r)
}

func TestRouting_AddUser(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(h.Handlers())
	defer srv.Close()

	user := models.User{
		ID:          0,
		Name:        faker.Name(),
		Phonenumber: faker.Phonenumber(),
	}

	body, err := json.Marshal(user)
	assert.NoError(t, err)

	res, err := http.Post(fmt.Sprintf("%s/user", srv.URL), "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestRouting_GetUserList(t *testing.T) {
	srv := httptest.NewServer(h.Handlers())
	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s/user?limit=2", srv.URL))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Len(t, res, 1)
}
