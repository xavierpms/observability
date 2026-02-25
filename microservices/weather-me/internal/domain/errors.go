package domain

import "errors"

var (
	ErrInvalidCEPFormat = errors.New("invalid CEP format")
	ErrForwardCEP       = errors.New("failed to forward CEP to service-b")
)
