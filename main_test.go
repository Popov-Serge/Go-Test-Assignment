package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeWhenNotOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	type TestCase struct {
		url        string
		status     int
		bodyString string
	}

	requests := []TestCase{
		{"/cafe", http.StatusBadRequest, "unknown"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.url, nil)
		handler.ServeHTTP(response, req)
		assert.Equal(t, v.status, response.Code)
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	type testData struct {
		city  string
		count int
	}

	requests := []testData{
		{"moscow", 0},
		{"moscow", 1},
		{"moscow", 2},
		{"moscow", 100},
		{"tula", 0},
		{"tula", 1},
		{"tula", 2},
		{"tula", 100},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()

		req := httptest.NewRequest(
			"GET",
			"/cafe?count="+strconv.Itoa(v.count)+"&city="+v.city,
			nil,
		)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		body := strings.TrimSpace(response.Body.String())

		respBody := []string{}

		if body != "" {
			respBody = strings.Split(body, ",")
		}

		want := v.count
		if want > len(cafeList[v.city]) {
			want = len(cafeList[v.city])
		}

		assert.Equal(t, want, len(respBody))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	type testData struct {
		search string
		count  int
	}

	requests := []testData{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&search="+v.search, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)
		body := strings.TrimSpace(response.Body.String())

		respBody := []string{}

		if body != "" {
			respBody = strings.Split(body, ",")
			for _, name := range respBody {
				assert.True(t, strings.Contains(strings.ToLower(name), strings.ToLower(v.search)))
			}
		}

		assert.Equal(t, v.count, len(respBody))

	}
}
