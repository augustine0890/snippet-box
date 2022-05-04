package mysql

import (
	"database/sql"
	"snippetbox/pkg/models"
)

// Define a SnippetModel type
type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use Exec() on the embedded connection pool to execute statement
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Get the ID of newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Return a specific snippet based on its id
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// Return the 10 most recently snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
