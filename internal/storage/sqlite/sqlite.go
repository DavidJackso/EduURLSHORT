package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "modernc.org/sqlite" // init sqlite3 driver
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	//
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url (url, alias) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		// Проверяем на ошибку уникального ограничения
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, sql.ErrTxDone) {
			return fmt.Errorf("%s: %w", op, err)
		}
		// Проверяем на ошибку уникального ограничения по тексту ошибки
		if err.Error() == "UNIQUE constraint failed: url.alias" {
			return fmt.Errorf("%s: alias already exists", op)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	_ = res
	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	res := stmt.QueryRow(alias)
	var url string
	err = res.Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}

func (s *Storage) Delete(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}
