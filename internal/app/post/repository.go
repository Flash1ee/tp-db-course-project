package repository

//go:generate mockgen -destination=mocks/repository.go -package=mock_post -mock_names=Repository=PostRepository . Repository

type Repository interface {
}
