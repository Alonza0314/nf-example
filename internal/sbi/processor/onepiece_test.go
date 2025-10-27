package processor_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alonza0314/nf-example/internal/sbi"
	"github.com/Alonza0314/nf-example/pkg/factory"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func newTestServer(t *testing.T) *sbi.Server {
	t.Helper()
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	nfApp := sbi.NewMocknfApp(ctrl)
	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{Port: 8000},
		},
	}).AnyTimes()
	return sbi.NewServer(nfApp, "")
}

func Test_HTTPOnePieceGreeting(t *testing.T) {
	server := newTestServer(t)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req, err := http.NewRequest(http.MethodGet, "/onepiece", nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	ctx.Request = req

	server.HTTPOnePieceGreeting(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", recorder.Code, http.StatusOK)
	}

	var actual string
	if err = json.Unmarshal(recorder.Body.Bytes(), &actual); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	const expected = "Hello Straw Hat Pirates!"
	if actual != expected {
		t.Fatalf("unexpected body: got %q want %q", actual, expected)
	}
}

func Test_HTTPOnePieceRecruit(t *testing.T) {
	server := newTestServer(t)

	t.Run("Success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		body := bytes.NewBufferString(`{"name":"Jinbe"}`)
		req, err := http.NewRequest(http.MethodPost, "/onepiece/crew", body)
		if err != nil {
			t.Fatalf("create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req

		server.HTTPOnePieceRecruit(ctx)

		if recorder.Code != http.StatusCreated {
			t.Fatalf("unexpected status: got %d want %d", recorder.Code, http.StatusCreated)
		}

		var actual string
		if err = json.Unmarshal(recorder.Body.Bytes(), &actual); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		const expected = "Jinbe has joined the Straw Hat crew!"
		if actual != expected {
			t.Fatalf("unexpected body: got %q want %q", actual, expected)
		}
	})

	t.Run("MissingName", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		body := bytes.NewBufferString(`{}`)
		req, err := http.NewRequest(http.MethodPost, "/onepiece/crew", body)
		if err != nil {
			t.Fatalf("create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req

		server.HTTPOnePieceRecruit(ctx)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("unexpected status: got %d want %d", recorder.Code, http.StatusBadRequest)
		}

		var payload map[string]string
		if err = json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		const expected = "name is required"
		if payload["error"] != expected {
			t.Fatalf("unexpected error message: got %q want %q", payload["error"], expected)
		}
	})
}
