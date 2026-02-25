package domain

import "errors"

var (
	ErrInvalidCEPFormat    = errors.New("Invalid CEP format")
	ErrCEPNotFound         = errors.New("CEP not found")
	ErrTemperatureNotFound = errors.New("Temperature data not found")
)
