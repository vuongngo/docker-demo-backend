package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSkillSet(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		description        string
		url                string
		method             string
		function           func(w http.ResponseWriter, r *http.Request)
		body               string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			description:        "test invalid body",
			url:                "/skill",
			method:             "POST",
			function:           svc.CreateSkillSet,
			body:               "{\"role\": \"Test\", \"skills\": []}",
			expectedStatusCode: 400,
			expectedBody:       "non zero value required",
		},
		{
			description:        "test json body POST",
			url:                "/skill",
			method:             "POST",
			function:           svc.CreateSkillSet,
			body:               "{\"role\": \"Test\", \"skills\": [\"Test\"]}",
			expectedStatusCode: 201,
			expectedBody:       "Test",
		},
	}

	for _, tc := range tests {
		req, err := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
		req.Header.Set("Content-type", "application/json")
		cookie := &http.Cookie{
			Name:  "token",
			Value: "123",
		}
		req.AddCookie(cookie)
		assert.NoError(err)

		w := httptest.NewRecorder()
		tc.function(w, req)

		assert.Equal(tc.expectedStatusCode, w.Code, tc.description)
		assert.Contains(w.Body.String(), tc.expectedBody, tc.description)
	}
}

func TestGetSkillSet(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		description        string
		url                string
		method             string
		function           func(w http.ResponseWriter, r *http.Request)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			description:        "test not found skillSet",
			url:                "/skill/1",
			method:             "GET",
			expectedStatusCode: 404,
			expectedBody:       "",
		},
		{
			description:        "test found skillSet",
			url:                "/skill/" + skillSetId,
			method:             "GET",
			expectedStatusCode: 200,
			expectedBody:       "Test",
		},
	}

	for _, tc := range tests {
		req, err := http.NewRequest(tc.method, tc.url, nil)
		assert.NoError(err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(tc.expectedStatusCode, w.Code, tc.description)
		assert.Contains(w.Body.String(), tc.expectedBody, tc.description)
	}
}
