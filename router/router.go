package router

import (
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/util2"
)

var Router = httprouter.New()

func Handle (method, pattern string, handle util2.ErrorHandle) {
  Router.Handle(method, pattern, util2.ErrorHandleFunc(handle))
}
func GET     (pattern string, handle util2.ErrorHandle) {Handle("GET",     pattern, handle)}
func PUT     (pattern string, handle util2.ErrorHandle) {Handle("PUT",     pattern, handle)}
func POST    (pattern string, handle util2.ErrorHandle) {Handle("POST",    pattern, handle)}
func DELETE  (pattern string, handle util2.ErrorHandle) {Handle("DELETE",  pattern, handle)}
func PATCH   (pattern string, handle util2.ErrorHandle) {Handle("PATCH",   pattern, handle)}
func OPTIONS (pattern string, handle util2.ErrorHandle) {Handle("OPTIONS", pattern, handle)}
func HEAD    (pattern string, handle util2.ErrorHandle) {Handle("HEAD",    pattern, handle)}
