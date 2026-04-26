package domain

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidParam = errors.New("invalid parameter")
	ErrInvalidPage  = errors.New("invalid page parameter")
	ErrInvalidLimit = errors.New("invalid limit parameter")
)
