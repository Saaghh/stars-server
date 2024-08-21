package models

import (
	"errors"
	"fmt"
)

var (
	ErrNoTx        = errors.New("no tx provided in context")
	ErrNoDataInCtx = errors.New("data not found in ctx")
)

type ErrDictionaryNotFound struct {
	DictionaryName string
}

func (e *ErrDictionaryNotFound) Error() string {
	return fmt.Sprintf("dictionary %s not found", e.DictionaryName)
}
