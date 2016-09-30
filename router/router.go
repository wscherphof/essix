package router

import (
	"github.com/julienschmidt/httprouter"
	"errors"
)

var (
	Router                 = httprouter.New()
	ErrInternalServerError = errors.New("ErrInternalServerError")
)
