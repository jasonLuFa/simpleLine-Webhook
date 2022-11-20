package api

import (
	"jasonLuFa/simpleLine-Webhook/save/query"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, iUserMessageRepository query.IUserMessageRepository) *Server {
	server, err := NewServer(iUserMessageRepository)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
