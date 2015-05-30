package router

import (
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/util2"
)

var Router = httprouter.New()

func errorHandle (method, pattern string, handle util2.ErrorHandle) {
  Router.Handle(method, pattern, util2.ErrorHandleFunc(handle))
}
func GET    (pattern string, handle util2.ErrorHandle) {errorHandle("GET",    pattern, handle)}
func PUT    (pattern string, handle util2.ErrorHandle) {errorHandle("PUT",    pattern, handle)}
func POST   (pattern string, handle util2.ErrorHandle) {errorHandle("POST",   pattern, handle)}
func DELETE (pattern string, handle util2.ErrorHandle) {errorHandle("DELETE", pattern, handle)}
