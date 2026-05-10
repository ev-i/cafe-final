package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	cases := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, c := range cases {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", c, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
		// fmt.Println(response.Body.String())
	}
}
func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, r := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", r.request, nil)
		handler.ServeHTTP(response, req)
		// fmt.Println(response.Body.String())

		assert.Equal(t, r.status, response.Code)
		assert.Equal(t, r.message, strings.TrimSpace(response.Body.String()))
	}

}

func TestCafeCount(t *testing.T) {
	cases := []struct {
		request string
		count   int
	}{
		{"/cafe?count=0&city=moscow", 0},
		{"/cafe?count=1&city=moscow", 1},
		{"/cafe?count=2&city=moscow", 2},
		{"/cafe?count=100&city=moscow", 100},
	}
	handler := http.HandlerFunc(mainHandle)
	for _, c := range cases {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", c.request, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		raw := strings.TrimSpace(response.Body.String())
		sliced := []string{}
		if raw != "" {
			sliced = strings.Split(raw, ",")
		}

		length := len(sliced)
		minLength := min(len(cafeList["moscow"]), length)
		assert.Equal(t, length, minLength)
	}
}

func TestCafeSearch(t *testing.T) {
	cases := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}
	urlConst := "/cafe?city=moscow&search="
	handler := http.HandlerFunc(mainHandle)
	for _, c := range cases {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", urlConst+c.search, nil)

		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)

		raw := strings.TrimSpace(response.Body.String())
		sliced := []string{}
		if raw != "" {
			sliced = strings.Split(raw, ",")
		}
		count := 0

		for _, v := range sliced {
			if strings.Contains(strings.ToLower(v), strings.ToLower(c.search)) {
				count += 1
			}
		}
		assert.Equal(t, count, c.wantCount)
	}
}
