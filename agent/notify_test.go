package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mock_agent "px/agent/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestNotify(t *testing.T) {
	router := gin.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storeMock := mock_agent.NewMockIStore(ctrl)
	// storeMock.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil)
	storeMock.EXPECT().Get(gomock.Any()).Return(nil, nil)
	router.GET("/notify", notify(storeMock))

	t.Run("EmptyQuery", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/notify", nil)
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

	t.Run("NoSuchJob", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/notify", nil)
		resp := httptest.NewRecorder()
		qs := req.URL.Query()
		qs.Set("pid", "123456")
		qs.Set("signal", "5000")
		req.URL.RawQuery = qs.Encode()
		router.ServeHTTP(resp, req)
		if resp.Code != 404 {
			t.Fatalf("should get %d but got %d", 200, resp.Code)
		}
	})
}
