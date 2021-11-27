package forum

//go:generate mockgen -destination=mocks/usecase.go -package=mock_forum -mock_names=Usecase=ForumUsecase . Usecase

type Usecase interface {
}
