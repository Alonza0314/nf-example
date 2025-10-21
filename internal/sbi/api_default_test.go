package sbi_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Alonza0314/nf-example/internal/sbi"
	"github.com/gin-gonic/gin"
)

func TestDefaultRoot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	s := &sbi.Server{}
	group := r.Group("/default")
	s.RegisterDefaultRoutes(group)

	req := httptest.NewRequest(http.MethodGet, "/default/", nil).WithContext(context.Background())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200", w.Code)
	}
	got := strings.TrimSpace(w.Body.String())
	if got != `"Hello free5GC!"` {
		t.Fatalf("body=%s, want %q", got, "Hello free5GC!")
	}
}

func TestDefaultExerciseGET(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	s := &sbi.Server{}
	group := r.Group("/default")
	s.RegisterDefaultRoutes(group)

	req := httptest.NewRequest(http.MethodGet, "/default/exercise", nil).WithContext(context.Background())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200", w.Code)
	}
	if strings.TrimSpace(w.Body.String()) != `"This is get"` {
		t.Fatalf("body=%s, want %q", w.Body.String(), "This is get")
	}
}

func TestDefaultExercisePOST(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	s := &sbi.Server{}
	group := r.Group("/default")
	s.RegisterDefaultRoutes(group)

	req := httptest.NewRequest(http.MethodPost, "/default/exercise", bytes.NewBufferString(`{}`)).
		WithContext(context.Background())
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200", w.Code)
	}
	if strings.TrimSpace(w.Body.String()) != `"This is post"` {
		t.Fatalf("body=%s, want %q", w.Body.String(), "This is post")
	}
}
