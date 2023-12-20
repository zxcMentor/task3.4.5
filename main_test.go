package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	testDB = setupTestDB()
	defer testDB.Close()

	exitCode := m.Run()

	tearDownTestDB()

	os.Exit(exitCode)
}

func setupTestDB() *sql.DB {
	db, err := sql.Open("postgres", "user=youruser_test dbname=yourdbname_test sslmode=disable password=yourpassword")
	if err != nil {
		fmt.Println("Error connecting to test database:", err)
		os.Exit(1)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS test_vacancies (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255),
			company VARCHAR(255),
			location VARCHAR(255),
			description TEXT
		);

		CREATE TABLE IF NOT EXISTS test_search_history (
			id SERIAL PRIMARY KEY,
			query VARCHAR(255),
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		fmt.Println("Error setting up test database:", err)
		os.Exit(1)
	}

	return db
}

func tearDownTestDB() {
	_, err := testDB.Exec(`
		DROP TABLE IF EXISTS test_vacancies;
		DROP TABLE IF EXISTS test_search_history;
	`)
	if err != nil {
		fmt.Println("Error tearing down test database:", err)
		os.Exit(1)
	}
}

func TestSearchVacancy(t *testing.T) {
	repo := NewSQLRepository(testDB)
	service := NewMyService(repo)
	testVacancy := Vacancy{
		Title:       "Software Engineer",
		Company:     "Test Company",
		Location:    "Test Location",
		Description: "Test Description",
	}
	err := repo.SaveVacancy(testVacancy)
	if err != nil {
		t.Fatalf("Failed to save test vacancy: %v", err)
	}

	query := "Software"
	vacancies, err := service.SearchVacancy(query)
	if err != nil {
		t.Fatalf("SearchVacancy failed: %v", err)
	}

	if len(vacancies) == 0 {
		t.Error("Expected at least one result, but got none.")
	}

	firstResult := vacancies[0]
	if firstResult.Title != testVacancy.Title {
		t.Errorf("Expected title %s, but got %s", testVacancy.Title, firstResult.Title)
	}
}

func TestGetVacancy(t *testing.T) {
	repo := NewSQLRepository(testDB)
	service := NewMyService(repo)
	testVacancy := Vacancy{
		Title:       "Software Engineer",
		Company:     "Test Company",
		Location:    "Test Location",
		Description: "Test Description",
	}
	err := repo.SaveVacancy(testVacancy)
	if err != nil {
		t.Fatalf("Failed to save test vacancy: %v", err)
	}

	vacancyID := testVacancy.ID
	resultVacancy, err := service.GetVacancy(vacancyID)
	if err != nil {
		t.Fatalf("GetVacancy failed: %v", err)
	}

	if resultVacancy == nil {
		t.Error("Expected a result, but got nil.")
	}

	if resultVacancy.Title != testVacancy.Title {
		t.Errorf("Expected title %s, but got %s", testVacancy.Title, resultVacancy.Title)
	}
}

func TestListVacancies(t *testing.T) {
	repo := NewSQLRepository(testDB)
	service := NewMyService(repo)
	testVacancies := []Vacancy{
		{
			Title:       "Software Engineer",
			Company:     "Test Company",
			Location:    "Test Location",
			Description: "Test Description",
		},
		{
			Title:       "Data Scientist",
			Company:     "Another Company",
			Location:    "Another Location",
			Description: "Another Description",
		},
	}

	for _, v := range testVacancies {
		err := repo.SaveVacancy(v)
		if err != nil {
			t.Fatalf("Failed to save test vacancy: %v", err)
		}
	}

	vacancies, err := service.ListVacancies()
	if err != nil {
		t.Fatalf("ListVacancies failed: %v", err)
	}

	if len(vacancies) != len(testVacancies) {
		t.Errorf("Expected %d vacancies, but got %d", len(testVacancies), len(vacancies))
	}

	for i, expected := range testVacancies {
		if vacancies[i].Title != expected.Title {
			t.Errorf("Expected title %s, but got %s", expected.Title, vacancies[i].Title)
		}
	}

}

func TestDeleteVacancy(t *testing.T) {
	repo := NewSQLRepository(testDB)
	service := NewMyService(repo)
	testVacancy := Vacancy{
		Title:       "Software Engineer",
		Company:     "Test Company",
		Location:    "Test Location",
		Description: "Test Description",
	}

	err := repo.SaveVacancy(testVacancy)
	if err != nil {
		t.Fatalf("SaveVacancy failed: %v", err)
	}

	err = service.DeleteVacancy(testVacancy.ID)
	if err != nil {
		t.Fatalf("DeleteVacancy failed: %v", err)
	}

	deletedVacancy, err := repo.GetVacancy(testVacancy.ID)

	if deletedVacancy != nil {
		t.Error("Expected deleted vacancy to be nil, but got a result.")
	}
}

func TestListSearchHistory(t *testing.T) {
	repo := NewSQLRepository(testDB)
	service := NewMyService(repo)
	testQueries := []string{"Software", "Data Science", "Web Developer"}
	for _, query := range testQueries {
		err := repo.SaveSearchHistory(query)
		if err != nil {
			t.Fatalf("Failed to save search history: %v", err)
		}
	}

	history, err := service.ListSearchHistory()
	if err != nil {
		t.Fatalf("ListSearchHistory failed: %v", err)
	}

	if len(history) != len(testQueries) {
		t.Errorf("Expected %d search history entries, but got %d", len(testQueries), len(history))
	}

	for i, expected := range testQueries {
		if history[i].Query != expected {
			t.Errorf("Expected query %s, but got %s", expected, history[i].Query)
		}
	}
}
