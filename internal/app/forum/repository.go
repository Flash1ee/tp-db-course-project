package forum

//go:generate mockgen -destination=mocks/repository.go -package=mock_forum -mock_names=Repository=ForumRepository . Repository

type Repository interface {
}
