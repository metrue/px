package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock_agent "px/agent/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestFibonacci(t *testing.T) {
	router := gin.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storeMock := mock_agent.NewMockIStore(ctrl)
	pid := "5000"
	process := Process{
		Pid: 5000,
	}
	data, _ := json.Marshal(process)
	storeMock.EXPECT().Get(pid).Return(data, nil)
	router.GET("/inspect", inspect(storeMock))

	t.Run("EmptyQuery", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/inspect", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		if resp.Code != 400 {
			t.Fatalf("should get %d but got %d", 400, resp.Code)
		}

		msg := `{"message":"pid is required"}`
		if resp.Body.String() != msg {
			t.Fatalf("should get %s but got %s", msg, resp.Body.String())
		}
	})
	t.Run("OK", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/inspect", nil)
		resp := httptest.NewRecorder()
		qs := req.URL.Query()
		qs.Set("pid", pid)
		req.URL.RawQuery = qs.Encode()
		router.ServeHTTP(resp, req)
		if resp.Code != 200 {
			t.Fatalf("should get %d but got %d", 200, resp.Code)
		}
		body := resp.Body.String()
		if !strings.Contains(body, `\"pid\":5000`) {
			t.Fatalf("incorrect response")
		}
	})
}
