package processor_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	nf_context "github.com/Alonza0314/nf-example/internal/context"
	"github.com/Alonza0314/nf-example/internal/sbi/processor"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_PostFortune_AddsFortune(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockNf := processor.NewMockProcessorNf(mockCtrl)
	p, err := processor.NewProcessor(mockNf)
	assert.NoError(t, err)

	mockCtx := &nf_context.NFContext{
		Fortunes: []string{},
	}
	mockNf.EXPECT().Context().Return(mockCtx).Times(1)

	rec := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(rec)

	req := processor.PostFortuneRequest{Fortune: "Lucky day"}
	p.PostFortune(ginCtx, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "Fortune added successfully", resp["message"])
	assert.Equal(t, "Lucky day", resp["fortune"])
	assert.Len(t, mockCtx.Fortunes, 1)
	assert.Equal(t, "Lucky day", mockCtx.Fortunes[0])
}

func Test_GetFortune_ReturnsOneOfFortunes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockNf := processor.NewMockProcessorNf(mockCtrl)
	p, err := processor.NewProcessor(mockNf)
	assert.NoError(t, err)

	fortunes := []string{"f1", "f2", "f3"}
	mockCtx := &nf_context.NFContext{
		Fortunes: fortunes,
	}
	mockNf.EXPECT().Context().Return(mockCtx).Times(1)

	rec := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(rec)

	p.GetFortune(ginCtx)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	f, ok := resp["fortune"].(string)
	assert.True(t, ok, "fortune should be a string")
	// check fortune is one of the provided values
	found := false
	for _, v := range fortunes {
		if v == f {
			found = true
			break
		}
	}
	assert.True(t, found, "returned fortune must be one of the provided fortunes")
	assert.Equal(t, "Here is your fortune for today!", resp["message"])
}

func Test_GetFortune_NoFortunes_ReturnsError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockNf := processor.NewMockProcessorNf(mockCtrl)
	p, err := processor.NewProcessor(mockNf)
	assert.NoError(t, err)

	mockCtx := &nf_context.NFContext{
		Fortunes: []string{},
	}
	mockNf.EXPECT().Context().Return(mockCtx).Times(1)

	rec := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(rec)

	p.GetFortune(ginCtx)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "No fortunes available.", resp["message"])
}
