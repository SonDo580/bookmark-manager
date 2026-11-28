package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type samplePathParams struct {
	Code string `uri:"code" binding:"required"`
}

type sampleQueryParams struct {
	Page int `form:"page" binding:"gte=1"`
}

type sampleReqHeaders struct {
	Authorization string `header:"Authorization" binding:"required"`
}

type sampleReqBody struct {
	Password string `json:"password" binding:"required,min=8"`
}

func newTestGinCtx(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(rec)
	return ginCtx, rec
}

func newJSONRequest(t *testing.T, method, target string, payload any) *http.Request {
	t.Helper()

	var buf bytes.Buffer
	if payload != nil {
		err := json.NewEncoder(&buf).Encode(payload)
		require.NoError(t, err)
	}

	req := httptest.NewRequest(method, target, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func assertBindFailure(
	t *testing.T,
	c *gin.Context,
	rec *httptest.ResponseRecorder,
	err error,
	result any,
) {
	t.Helper()
	require.Error(t, err)
	assert.Nil(t, result)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Message)
	assert.NotNil(t, resp.Details)
}

func assertBindSuccess[U any](
	t *testing.T,
	c *gin.Context,
	err error,
	result *U,
	expectedResult U,
) {
	t.Helper()

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, c.IsAborted())
	assert.Equal(t, expectedResult, *result)
}

func TestBindJSONWithAbort(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		reqBody    sampleReqBody
		shouldFail bool
	}{
		{
			name:       "success",
			reqBody:    sampleReqBody{Password: "secret123"},
			shouldFail: false,
		},
		{
			name:       "fail",
			reqBody:    sampleReqBody{Password: "123"}, // fail min=8
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, rec := newTestGinCtx(t)
			c.Request = newJSONRequest(t, http.MethodPost, "/", tc.reqBody)

			result, err := BindJSONWithAbort[sampleReqBody](c)

			if tc.shouldFail {
				assertBindFailure(t, c, rec, err, result)
			} else {
				assertBindSuccess(t, c, err, result, tc.reqBody)
			}
		})
	}
}

func TestBindURIWithAbort(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		code       string
		shouldFail bool
	}{
		{
			name:       "success",
			code:       "abc1234",
			shouldFail: false,
		},
		{
			name:       "fail",
			code:       "", // fail required
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, rec := newTestGinCtx(t)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.code != "" {
				c.Params = gin.Params{{Key: "code", Value: tc.code}}
			}

			result, err := BindURIWithAbort[samplePathParams](c)

			if tc.shouldFail {
				assertBindFailure(t, c, rec, err, result)
			} else {
				assertBindSuccess(t, c, err, result, samplePathParams{Code: tc.code})
			}
		})
	}
}

func TestBindQueryWithAbort(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		page       int
		shouldFail bool
	}{
		{
			name:       "success",
			page:       1,
			shouldFail: false,
		},
		{
			name:       "fail",
			page:       0, // fail gte=1
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, rec := newTestGinCtx(t)
			target := fmt.Sprintf("/?page=%d", tc.page)
			c.Request = httptest.NewRequest(http.MethodGet, target, nil)

			result, err := BindQueryWithAbort[sampleQueryParams](c)

			if tc.shouldFail {
				assertBindFailure(t, c, rec, err, result)
			} else {
				assertBindSuccess(t, c, err, result, sampleQueryParams{Page: tc.page})
			}
		})
	}
}

func TestBindHeaderWithAbort(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		authHeader string
		shouldFail bool
	}{
		{
			name:       "success",
			authHeader: "Bearer token",
			shouldFail: false,
		},
		{
			name:       "fail",
			authHeader: "", // fail required
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, rec := newTestGinCtx(t)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.authHeader != "" {
				c.Request.Header.Set("Authorization", tc.authHeader)
			}

			result, err := BindHeaderWithAbort[sampleReqHeaders](c)

			if tc.shouldFail {
				assertBindFailure(t, c, rec, err, result)
			} else {
				assertBindSuccess(t, c, err, result, sampleReqHeaders{Authorization: tc.authHeader})
			}
		})
	}
}
