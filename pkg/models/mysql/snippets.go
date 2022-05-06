package mysql

import (
	"database/sql"
	"errors"
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
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Returns a pointer to a sql.Row object
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct
	s := &models.Snippet{}

	// Copy the values from each field in sql.Row to the corresponding field in the Snippet struct
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
	}

	return s, nil
}

// Return the 10 most recently snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
