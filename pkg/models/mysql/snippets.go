package mysql

import (
	"database/sql"
	"errors"
	"log"

	"snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (mdl *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires) 
			 VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := mdl.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (mdl *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE 
			 expires > UTC_TIMESTAMP() AND id = ?`

	row := mdl.DB.QueryRow(stmt, id)
	spt := &models.Snippet{}
	err := row.Scan(&spt.ID, &spt.Title, &spt.Content, &spt.Created, &spt.Expires)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, models.ErrNoRecord
	case err != nil:
		return nil, err
	default:
		return spt, nil
	}
}

func (mdl *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := mdl.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v", err)
		}
	}(rows)
	var snippets []*models.Snippet
	for rows.Next() {
		spt := &models.Snippet{}
		err := rows.Scan(&spt.ID, &spt.Title, &spt.Content, &spt.Created, &spt.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, spt)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
