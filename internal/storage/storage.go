package storage

import "errors"

//TODO: перевести на postgres

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)
