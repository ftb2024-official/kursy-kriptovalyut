package entities

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidParam = errors.New("invalid parameter")
	ErrInternal     = errors.New("internal error")
)
