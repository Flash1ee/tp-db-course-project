package vote

//go:generate mockgen -destination=usecase/mocks/usecase.go -package=mock_vote -mock_names=Usecase=VoteUsecase . Usecase

type Usecase interface {
}
