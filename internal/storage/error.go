package storage

import "database/sql"

// ErrorDao is a dao for a table of errors that occurred during file preprocessing or parsing.
type ErrorDao interface {
	Add(message string) (int64, error)
}

type PostgresErrorDao struct {
	db *sql.DB
}

func NewPostgresErrorDao(db *sql.DB) *PostgresErrorDao {
	return &PostgresErrorDao{db: db}
}

func (dao PostgresErrorDao) Add(message string) (int64, error) {
	query := `INSERT INTO error (message) VALUES ($1) RETURNING id`

	var id int64 = 0

	err := dao.db.QueryRow(query, message).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
