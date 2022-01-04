package forum_usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"io"
	mock_forum "tp-db-project/internal/app/forum/mocks"
)

type SuiteUsecase struct {
	suite.Suite
	Mock                *gomock.Controller
	MockForumRepository *mock_forum.ForumRepository

	Logger *logrus.Logger
}

func (s *SuiteUsecase) SetupSuite() {
	s.Mock = gomock.NewController(s.T())
	s.MockForumRepository = mock_forum.NewForumRepository(s.Mock)
	s.Logger = logrus.New()
	s.Logger.SetOutput(io.Discard)
}

func (s *SuiteUsecase) TearDownSuite() {
	s.Mock.Finish()
}
