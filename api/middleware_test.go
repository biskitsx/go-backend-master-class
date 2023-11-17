package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/biskitsx/go-backend-master-class/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthorization(t *testing.T, request *http.Request, tokenMaker token.Maker, authorizationType string, username string, duration time.Duration) {
	token, _, err := tokenMaker.CreateToken(username, "role", duration)
	require.NoError(t, err)

	authHeader := authorizationType + " " + token
	request.Header.Set(authorizationHeaderKey, authHeader)
}
func TestAuthMiddleware(t *testing.T) {
	teseCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "username", time.Minute)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name: "Unauthorized/NoToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// blank
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, resp.Code)
			},
		},
		{
			name: "UnsupportedAuthorizationType",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "whichtype", "username", time.Minute)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, resp.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "", "username", time.Minute)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, resp.Code)
			},
		},
		// {
		// 	name: "ExpiredToken",
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "username", -time.Minute)
		// 	},
		// 	checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusUnauthorized, resp.Code)
		// 	},
		// },
	}

	for _, tc := range teseCases {
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(
				authPath,
				authMiddleware(server.tokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(200, gin.H{})
				})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
