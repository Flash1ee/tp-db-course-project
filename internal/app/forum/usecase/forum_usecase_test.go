package forum_usecase

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"tp-db-project/internal/app/forum"
)

type SuiteForumUsecase struct {
	SuiteUsecase
	uc forum.Usecase
}

func (s *SuiteForumUsecase) SetupSuite() {
	s.SuiteUsecase.SetupSuite()
	//s.uc = NewForumUsecase(s.MockForumRepository)
}

func TestUsecaseForum(t *testing.T) {
	suite.Run(t, new(SuiteForumUsecase))
}
