package storage

import "errors"

//TODO: перевести на postgres

var (
	ErrUrlNotFound = errors.New("url not found")
	ErrUrlExists   = errors.New("url already exists")
)
