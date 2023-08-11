package utils

import "github.com/pkg/errors"

var (
	ErrFileIDAlreadyExists = errors.New("file_id already exists")
	ErrUsageLimitExceeded  = errors.New("user usage limit exceeded")
	ErrUserRateLimitNotSet = errors.New("user rate limit not set")
	ErrTooManyRequests     = errors.New("too many requests. try again later")
)
