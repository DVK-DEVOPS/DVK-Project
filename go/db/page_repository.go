package db

import (
	"database/sql"
	"time"
)

type PageRepository struct {
	DB *sql.DB
}

func NewPageRepository(db *sql.DB) *PageRepository {
	return &PageRepository{DB: db}
}

// For page templating
func (r *PageRepository) FindSearchResults(searchStr string) ([]Result, error) {
	if searchStr == "" {
		return nil, nil
	}

	rows, err := r.DB.Query("SELECT Title, Url, Content, Language, CreatedAt, UpdatedAt FROM pages WHERE LOWER(Title) LIKE LOWER(?)", "%"+searchStr+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Result
	for rows.Next() {
		var res Result
		if err := rows.Scan(&res.Title, &res.Url, &res.Content, &res.Language, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

type Result struct {
	Title     string
	Url       string
	Content   string
	Language  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
