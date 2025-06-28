package router

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dzhisl/license-api/internal/storage"
	"github.com/dzhisl/license-api/pkg/config"
	"github.com/dzhisl/license-api/pkg/logger"
	"github.com/spf13/viper"
	"github.com/test-go/testify/assert"
)

var (
	ctx context.Context
	r   http.Handler
)

func TestMain(m *testing.M) {

	ctx = context.TODO()
	config.InitConfig()
	logger.InitLogger()
	storage.InitStorage(ctx)
	r = InitRouter()

	code := m.Run()
	os.Exit(code)
}

func TestAuthorization(t *testing.T) {
	w := httptest.NewRecorder()

	exampleUser := map[string]interface{}{
		"max_activations": 3,
		"expires_at":      1750721178,
		"discord_id":      809123,
	}

	userJson, err := json.Marshal(exampleUser)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/api/user/create", strings.NewReader(string(userJson)))
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestCreateUser(t *testing.T) {
	w := httptest.NewRecorder()

	exampleUser := map[string]interface{}{
		"max_activations": 3,
		"expires_at":      1750721178,
		"discord_id":      809123,
	}

	body, err := json.Marshal(exampleUser)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/user/create", bytes.NewBuffer(body))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", viper.GetString("ADMIN_SECRET_KEY"))

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
