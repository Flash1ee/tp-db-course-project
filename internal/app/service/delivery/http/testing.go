package service_handler

type TestTable struct {
	Name              string
	Data              interface{}
	ExpectedMockTimes int
	ExpectedCode      int
}

//type SuiteHandler struct {
//	suite.Suite
//	Mock             *gomock.Controller
//	MockUsersUsecase *mock_users.UsersUsecase
//	Logger           *logrus.Logger
//}
//
//func (s *SuiteHandler) SetupSuite() {
//	s.Mock = gomock.NewController(s.T())
//	s.MockUsersUsecase = mock_users.NewUsersUsecase(s.Mock)
//	s.Logger = logrus.New()
//	s.Logger.SetOutput(io.Discard)
//}
//
//func (s *SuiteHandler) TearDownSuite() {
//	s.Mock.Finish()
//}
