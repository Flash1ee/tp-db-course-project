package repository

//go:generate mockgen -destination=mocks/repository.go -package=mock_thread -mock_names=Repository=ThreadRepository . Repository

type Repository interface {
}
