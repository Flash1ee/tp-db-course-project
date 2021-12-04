package users_usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"io"
	mock_users "tp-db-project/internal/app/users/mocks"
)

type SuiteUsecase struct {
	suite.Suite
	Mock                *gomock.Controller
	MockUsersRepository *mock_users.UsersRepository

	Logger *logrus.Logger
}

func (s *SuiteUsecase) SetupSuite() {
	s.Mock = gomock.NewController(s.T())
	s.MockUsersRepository = mock_users.NewUsersRepository(s.Mock)
	s.Logger = logrus.New()
	s.Logger.SetOutput(io.Discard)
}

func (s *SuiteUsecase) TearDownSuite() {
	s.Mock.Finish()
}
