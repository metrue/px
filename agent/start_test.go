package agent

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock_agent "px/agent/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestStart(t *testing.T) {
	router := gin.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storeMock := mock_agent.NewMockIStore(ctrl)
	storeMock.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil)
	router.GET("/start", start(storeMock))

	t.Run("EmptyQuery", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/start", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		if resp.Code != 400 {
			t.Fatalf("should get %d but got %d", 400, resp.Code)
		}

		msg := `{"message":"cmd is required"}`
		if resp.Body.String() != msg {
			t.Fatalf("should get %s but got %s", msg, resp.Body.String())
		}
	})
	t.Run("OK", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/start", nil)
		resp := httptest.NewRecorder()
		qs := req.URL.Query()
		qs.Set("cmd", "sleep 5000")
		req.URL.RawQuery = qs.Encode()
		router.ServeHTTP(resp, req)
		if resp.Code != 200 {
			t.Fatalf("should get %d but got %d", 200, resp.Code)
		}
		body := resp.Body.String()
		if !strings.Contains(body, `{"message":"job started with`) {
			t.Fatalf("should get %s but got %s", `{"message":"job started with`, body)
		}
	})
}
