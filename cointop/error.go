package cointop

import "errors"

// ErrInvalidAPIChoice is error for invalid API choice
var ErrInvalidAPIChoice = errors.New("invalid API choice")

// ErrCoinNameOrSymbolRequired is error for when coin name or symbol is required
var ErrCoinNameOrSymbolRequired = errors.New("coin name or symbol is required")
