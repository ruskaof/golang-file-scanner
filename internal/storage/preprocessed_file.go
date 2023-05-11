package storage

import "database/sql"

// PreprocessedFileDao is a dao for a table that keeps track of all files that were sent to queue
// for preprocessing
type PreprocessedFileDao interface {
	WasPreprocessed(fileName string) (bool, error)
	Add(fileName string) error
}

type PostgresPreprocessedFileDao struct {
	db *sql.DB
}

func NewPostgresPreprocessedFileDao(db *sql.DB) *PostgresPreprocessedFileDao {
	return &PostgresPreprocessedFileDao{db: db}
}

func (dao PostgresPreprocessedFileDao) Add(fileName string) error {
	query := `INSERT INTO preprocessed_file (name) VALUES ($1)`

	_, err := dao.db.Exec(query, fileName)

	if err != nil {
		return err
	}

	return nil
}

func (dao PostgresPreprocessedFileDao) WasPreprocessed(fileName string) (bool, error) {
	row := dao.db.QueryRow("SELECT EXISTS(SELECT 1 FROM preprocessed_file WHERE name = $1)", fileName)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
