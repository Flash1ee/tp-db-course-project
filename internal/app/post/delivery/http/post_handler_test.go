package post_handler

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type SuiteForumHandler struct {
	SuiteHandler
	handler *PostHandler
}

func (s *SuiteForumHandler) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	//s.handler = NewPostHandler(router.NewRouter(), s.Logger, s.MockForumUsecase)
}

func TestForumHandler(t *testing.T) {
	suite.Run(t, new(SuiteForumHandler))
}
