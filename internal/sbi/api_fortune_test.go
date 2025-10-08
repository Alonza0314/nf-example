package sbi_test

import (
	"net/http"
	"testing"

	"github.com/Alonza0314/nf-example/internal/sbi"
	"github.com/Alonza0314/nf-example/pkg/factory"
	"go.uber.org/mock/gomock"
)

func Test_GetFortuneRoute(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	nfApp := sbi.NewMocknfApp(mockCtrl)
	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{
				Port: 8000,
			},
		},
	}).AnyTimes()

	server := sbi.NewServer(nfApp, "")

	routes := server.GetFortuneRoute()
	if routes == nil {
		t.Fatalf("Expected routes slice, got nil")
	}

	if len(routes) != 2 {
		t.Fatalf("Expected 2 routes, got %d", len(routes))
	}

	// Validate first route (GET)
	r0 := routes[0]
	if r0.Name != "Get Today's Fortune" {
		t.Errorf("Route 0 Name: expected %q, got %q", "Get Today's Fortune", r0.Name)
	}
	if r0.Method != http.MethodGet {
		t.Errorf("Route 0 Method: expected %q, got %q", http.MethodGet, r0.Method)
	}
	if r0.Pattern != "/" {
		t.Errorf("Route 0 Pattern: expected %q, got %q", "/", r0.Pattern)
	}
	if r0.APIFunc == nil {
		t.Errorf("Route 0 APIFunc should not be nil")
	}

	// Validate second route (POST)
	r1 := routes[1]
	if r1.Name != "Add a new Fortune" {
		t.Errorf("Route 1 Name: expected %q, got %q", "Add a new Fortune", r1.Name)
	}
	if r1.Method != http.MethodPost {
		t.Errorf("Route 1 Method: expected %q, got %q", http.MethodPost, r1.Method)
	}
	if r1.Pattern != "/" {
		t.Errorf("Route 1 Pattern: expected %q, got %q", "/", r1.Pattern)
	}
	if r1.APIFunc == nil {
		t.Errorf("Route 1 APIFunc should not be nil")
	}
}
