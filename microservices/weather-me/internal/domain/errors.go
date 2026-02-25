package domain

import "errors"

var (
	ErrInvalidCEPFormat = errors.New("Invalid CEP format")
	ErrForwardCEP       = errors.New("Failed to forward CEP to weather-by-city")
)
