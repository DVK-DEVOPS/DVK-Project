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

// UpsertPage saves a new page or updates it if the URL already exists
func (r *PageRepository) UpsertPage(title, url, content, language string) error {
	_, err := r.DB.Exec(`
		INSERT INTO pages (title, url, content, language, createdAt, updatedAt)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (url)
		DO UPDATE SET
			title = EXCLUDED.title,
			content = EXCLUDED.content,
			language = EXCLUDED.language,
			updatedAt = NOW()
	`, title, url, content, language)
	return err
}

// For page templating
func (r *PageRepository) FindSearchResults(searchStr string, language string) ([]Result, error) {
	query := `SELECT title, url, content, language, createdat, updatedat
			  FROM pages
			  WHERE search_vector @@ plainto_tsquery('english', $1)`
	args := []interface{}{searchStr}

	if language != "" {
		query += " AND LOWER(language) = LOWER($2)"
		args = append(args, language)
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Result
	for rows.Next() {
		var res Result
		var lang sql.NullString
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&res.Title, &res.Url, &res.Content, &lang, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		if lang.Valid {
			res.Language = lang.String
		} else {
			res.Language = ""
		}
		res.CreatedAt = createdAt
		res.UpdatedAt = updatedAt
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
