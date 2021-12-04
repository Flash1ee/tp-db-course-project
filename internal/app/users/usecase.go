package users

//go:generate mockgen -destination=mocks/usecase.go -package=mock_users -mock_names=Usecase=UsersUsecase . Usecase

type Usecase interface {
}
