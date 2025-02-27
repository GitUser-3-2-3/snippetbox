package mysql

import (
	"database/sql"
	"errors"

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
		return nil, err
	case err != nil:
		return nil, err
	default:
		return spt, nil
	}
}

func (mdl *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
