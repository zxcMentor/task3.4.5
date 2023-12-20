package main

type Vacancy struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type SearchHistory struct {
	ID        int    `json:"id"`
	Query     string `json:"query"`
	Timestamp string `json:"timestamp"`
}

type Repository interface {
	SearchVacancy(query string) ([]Vacancy, error)
	GetVacancy(id string) (*Vacancy, error)
	ListVacancies() ([]Vacancy, error)
	SaveVacancy(vacancy Vacancy) error
	DeleteVacancy(id string) error

	SaveSearchHistory(query string) error
	ListSearchHistory() ([]SearchHistory, error)
	DeleteSearchHistory(id int) error
}
