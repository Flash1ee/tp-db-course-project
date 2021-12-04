// Code generated by MockGen. DO NOT EDIT.
// Source: tp-db-project/internal/app/users (interfaces: Usecase)

// Package mock_users is a generated GoMock package.
package mock_users

import (
	gomock "github.com/golang/mock/gomock"
)

// UsersUsecase is a mock of Usecase interface.
type UsersUsecase struct {
	ctrl     *gomock.Controller
	recorder *UsersUsecaseMockRecorder
}

// UsersUsecaseMockRecorder is the mock recorder for UsersUsecase.
type UsersUsecaseMockRecorder struct {
	mock *UsersUsecase
}

// NewUsersUsecase creates a new mock instance.
func NewUsersUsecase(ctrl *gomock.Controller) *UsersUsecase {
	mock := &UsersUsecase{ctrl: ctrl}
	mock.recorder = &UsersUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *UsersUsecase) EXPECT() *UsersUsecaseMockRecorder {
	return m.recorder
}
