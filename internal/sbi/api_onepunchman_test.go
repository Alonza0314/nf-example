package sbi_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alonza0314/nf-example/internal/sbi"
	"github.com/Alonza0314/nf-example/pkg/factory"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func Test_OnePunchManAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 建立 Mock nfApp
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	nfApp := sbi.NewMocknfApp(mockCtrl)
	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{Port: 8000},
		},
	}).AnyTimes()

	// 建立整個伺服器
	server := sbi.NewServer(nfApp, "")

	// -----------------------
	// 測試 POST /onepunchman/echo
	// -----------------------
	t.Run("POST /onepunchman/echo", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusOK
		const EXPECTED_BODY = `{"data":{"name":"Saitama","power":100},"message":"Received your data!"}`

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		body := bytes.NewBuffer([]byte(`{"name":"Saitama","power":100}`))

		var err error
		ginCtx.Request, err = http.NewRequest("POST", "/onepunchman/echo", body)
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}
		server.HTTPOnePunchManEcho(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}
		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})

	// -----------------------
	// 測試 POST /onepunchman/echo 傳錯格式
	// -----------------------
	t.Run("POST /onepunchman/echo invalid json", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusBadRequest

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		body := bytes.NewBuffer([]byte(`invalid_json`))
		var err error
		ginCtx.Request, err = http.NewRequest("POST", "/onepunchman/echo", body)
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}
		server.HTTPOnePunchManEcho(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}
	})
}
