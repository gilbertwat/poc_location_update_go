package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/locations", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

func BenchmarkLocationRoute(b *testing.B) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/locations", nil)
	for n := 0; n < b.N; n++ {
		router.ServeHTTP(w, req)
	}
}
