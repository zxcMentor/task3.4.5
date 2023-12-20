package main

import "database/sql"

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS vacancies (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255),
			company VARCHAR(255),
			location VARCHAR(255),
			description TEXT
		);

		CREATE TABLE IF NOT EXISTS search_history (
			id SERIAL PRIMARY KEY,
			query VARCHAR(255),
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		panic(err)
	}

	return &SQLRepository{db: db}
}

func (r *SQLRepository) SearchVacancy(query string) ([]Vacancy, error) {
	queryBuilder := squirrel.Select("id", "title", "company", "location", "description").
		From("vacancies").
		Where("title ILIKE ?", "%"+query+"%")

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vacancies []Vacancy
	for rows.Next() {
		var v Vacancy
		if err := rows.Scan(&v.ID, &v.Title, &v.Company, &v.Location, &v.Description); err != nil {
			return nil, err
		}
		vacancies = append(vacancies, v)
	}

	return vacancies, nil
}

func (r *SQLRepository) GetVacancy(id string) (*Vacancy, error) {
	queryBuilder := squirrel.Select("id", "title", "company", "location", "description").
		From("vacancies").
		Where("id = ?", id)

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var v Vacancy
	if err := r.db.QueryRow(sql, args...).Scan(&v.ID, &v.Title, &v.Company, &v.Location, &v.Description); err != nil {
		return nil, err
	}

	return &v, nil
}

func (r *SQLRepository) ListVacancies() ([]Vacancy, error) {
	queryBuilder := squirrel.Select("id", "title", "company", "location", "description").
		From("vacancies")

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vacancies []Vacancy
	for rows.Next() {
		var v Vacancy
		if err := rows.Scan(&v.ID, &v.Title, &v.Company, &v.Location, &v.Description); err != nil {
			return nil, err
		}
		vacancies = append(vacancies, v)
	}

	return vacancies, nil
}
func (r *SQLRepository) SaveVacancy(vacancy Vacancy) error {
	_, err := squirrel.Insert("vacancies").
		Columns("title", "company", "location", "description").
		Values(vacancy.Title, vacancy.Company, vacancy.Location, vacancy.Description).
		RunWith(r.db).
		Exec()
	return err
}

func (r *SQLRepository) DeleteVacancy(id string) error {
	_, err := squirrel.Delete("vacancies").
		Where(squirrel.Eq{"id": id}).
		RunWith(r.db).
		Exec()
	return err
}

func (r *SQLRepository) SaveSearchHistory(query string) error {
	_, err := squirrel.Insert("search_history").
		Columns("query").
		Values(query).
		RunWith(r.db).
		Exec()
	return err
}
func (r *SQLRepository) ListSearchHistory() ([]SearchHistory, error) {
	queryBuilder := squirrel.Select("id", "query", "timestamp").
		From("search_history")

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []SearchHistory
	for rows.Next() {
		var h SearchHistory
		if err := rows.Scan(&h.ID, &h.Query, &h.Timestamp); err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	return history, nil
}

func (r *SQLRepository) DeleteSearchHistory(id int) error {
	_, err := squirrel.Delete("search_history").
		Where(squirrel.Eq{"id": id}).
		RunWith(r.db).
		Exec()
	return err
}
