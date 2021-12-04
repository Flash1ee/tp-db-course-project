package users_handler

import (
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"io"
	mock_users "tp-db-project/internal/app/users/mocks"
)

type TestTable struct {
	Name              string
	Data              interface{}
	ExpectedMockTimes int
	ExpectedCode      int
}

type SuiteHandler struct {
	suite.Suite
	Mock             *gomock.Controller
	MockUsersUsecase *mock_users.UsersUsecase
	Logger           *logrus.Logger
}

func (s *SuiteHandler) SetupSuite() {
	s.Mock = gomock.NewController(s.T())
	s.MockUsersUsecase = mock_users.NewUsersUsecase(s.Mock)
	s.Logger = logrus.New()
	s.Logger.SetOutput(io.Discard)
}

func (s *SuiteHandler) TearDownSuite() {
	s.Mock.Finish()
}
