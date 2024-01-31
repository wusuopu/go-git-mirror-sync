package routes

import (
	"app/initialize"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var engine *gin.Engine

func TestMain(m *testing.M) {
	engine = initialize.InitTest(nil)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestHealthCheck(t *testing.T) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/_health", nil)

	engine.ServeHTTP(recorder, req)

	assert.Equal(t, recorder.Code, 200)
}
