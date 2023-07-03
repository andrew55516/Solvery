package api

import (
	db "Solvery/db/sqlc"
	"Solvery/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		//EmailAddress:  "solvery55516@gmail.com",
		//EmailPassword: "vxqmyxvrsdutduzk",
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
