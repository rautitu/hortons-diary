package attacks

type Repository interface {
	//Create(ctx context.Context, userID UserID,
	//	input CreateInput) (Attack, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

//func (s *Service) Create(
//    ctx context.Context,
//    userID UserID,
//    input CreateInput,
//) (Attack, error) {
//    return s.repo.Create(ctx, userID, input)
//}
