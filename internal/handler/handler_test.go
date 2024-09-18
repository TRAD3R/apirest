package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/trad3r/hskills/apirest/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	h *Handler
)

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
