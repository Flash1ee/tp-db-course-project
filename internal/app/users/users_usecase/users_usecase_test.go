package users_usecase

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"tp-db-project/internal/app/users"
)

type SuiteUsersUsecase struct {
	SuiteUsecase
	uc users.Usecase
}

func (s *SuiteUsersUsecase) SetupSuite() {
	s.SuiteUsecase.SetupSuite()
	s.uc = NewUsersUsecase(s.MockUsersRepository)
}

func TestUsecaseUsers(t *testing.T) {
	suite.Run(t, new(SuiteUsersUsecase))
}
