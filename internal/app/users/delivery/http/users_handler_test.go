package users_handler

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"tp-db-project/internal/pkg/router"
)

type SuiteUsersHandler struct {
	SuiteHandler
	handler *UsersHandler
}

func (s *SuiteUsersHandler) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.handler = NewUsersHandler(router.NewRouter(), s.Logger, s.MockUsersUsecase)
}

func TestUsersHandler(t *testing.T) {
	suite.Run(t, new(SuiteUsersHandler))
}
