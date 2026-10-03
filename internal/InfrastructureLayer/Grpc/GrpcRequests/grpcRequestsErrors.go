package Requests

import "errors"

var (
	ErrorAttempts = errors.New("attempts were expired")
)
