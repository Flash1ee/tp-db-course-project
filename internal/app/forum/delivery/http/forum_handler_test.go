package forum_handler

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"tp-db-project/internal/pkg/router"
)

type SuiteForumHandler struct {
	SuiteHandler
	handler *ForumHandler
}

func (s *SuiteForumHandler) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.handler = NewForumHandler(router.NewRouter(), s.Logger, s.MockForumUsecase)
}

func TestForumHandler(t *testing.T) {
	suite.Run(t, new(SuiteForumHandler))
}
