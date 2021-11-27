package forum_repository

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SuiteForumRepository struct {
	Suite
	repo *ForumRepository
}

func (s *SuiteForumRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewForumRepository(s.DB)
}

func (s *SuiteForumRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}

func TestForumRepository(t *testing.T) {
	suite.Run(t, new(SuiteForumRepository))
}
