package services

type RepositoryService struct {
}

func (s *RepositoryService) Add(a, b int) int {
	return a + b
}

func (s *RepositoryService) Minus(a, b int) int {
	return a - b
}
