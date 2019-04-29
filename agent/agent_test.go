package agent

import (
	"testing"

	mock_agent "px/agent/mocks"

	"github.com/golang/mock/gomock"
)

func TestAgent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storeMock := mock_agent.NewMockIStore(ctrl)
	_ = New(storeMock)
}
