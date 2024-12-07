package db

import "errors"

var (
	ErrDBURIEmpty         = errors.New("MONGO_URI is not set")
	ErrDBConnectionFailed = errors.New("Failed to connect to database")
)
