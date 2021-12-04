package users

//go:generate mockgen -destination=mocks/repository.go -package=mock_users -mock_names=Repository=UsersRepository . Repository

type Repository interface {
}
