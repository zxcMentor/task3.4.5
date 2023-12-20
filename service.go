package main

type Service interface {
	SearchVacancy(query string) ([]Vacancy, error)
	GetVacancy(id string) (*Vacancy, error)
	ListVacancies() ([]Vacancy, error)
	DeleteVacancy(id string) error

	ListSearchHistory() ([]SearchHistory, error)
	DeleteSearchHistory(id int) error
}

type MyService struct {
	repo Repository
}

func NewMyService(repo Repository) *MyService {
	return &MyService{repo: repo}
}

func (s *MyService) SearchVacancy(query string) ([]Vacancy, error) {
	return s.repo.SearchVacancy(query)
}

func (s *MyService) GetVacancy(id string) (*Vacancy, error) {
	return s.repo.GetVacancy(id)
}

func (s *MyService) ListVacancies() ([]Vacancy, error) {
	return s.repo.ListVacancies()
}

func (s *MyService) DeleteVacancy(id string) error {
	return s.repo.DeleteVacancy(id)
}

func (s *MyService) ListSearchHistory() ([]SearchHistory, error) {
	return s.repo.ListSearchHistory()
}

func (s *MyService) DeleteSearchHistory(id int) error {
	return s.repo.DeleteSearchHistory(id)
}
