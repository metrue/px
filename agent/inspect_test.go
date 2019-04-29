package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mock_agent "px/agent/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestInspect(t *testing.T) {
	router := gin.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storeMock := mock_agent.NewMockIStore(ctrl)
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
}
