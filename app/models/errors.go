package models

import "errors"

var (
	ErrNoTx        = errors.New("no tx provided in context")
	ErrNoDataInCtx = errors.New("data not found in ctx")
)
