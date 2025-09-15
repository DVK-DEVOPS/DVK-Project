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
func (r *PageRepository) FindSearchResults(searchStr string, language string) ([]Result, error) {
	query := `SELECT Title, Url, Content, Language, CreatedAt, UpdatedAt 
	FROM pages WHERE LOWER(Title) LIKE LOWER(?)`

	args := []interface{}{"%" + searchStr + "%"}

	if language != "" {
		query += " AND LOWER(Language) = LOWER(?)"
		args = append(args, language)
	}
	rows, err := r.DB.Query(query, args...)

	//rows, err := r.DB.Query("SELECT Title, Url, Content, Language, CreatedAt, UpdatedAt FROM pages WHERE LOWER(Title) LIKE LOWER(?) AND LOWER(Language) = LOWER(?)", "%"+searchStr+"%", language)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Result
	for rows.Next() {
		var res Result
		var lang sql.NullString
		if err := rows.Scan(&res.Title, &res.Url, &res.Content, &lang, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, err
		}

		if lang.Valid {
			res.Language = lang.String
		} else {
			res.Language = ""
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
