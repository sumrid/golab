package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/labstack/echo"
	"github.com/magiconair/properties/assert"
)

func TestIndex(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/")

	assert.Equal(t, rec.Code, http.StatusOK)
}

func TestGradeEndpoint(t *testing.T) {
	q := make(url.Values)
	q.Set("point", "60")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/grade?" + q.Encode())

	assert.Equal(t, rec.Code, http.StatusOK)
}
