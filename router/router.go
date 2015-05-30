package router

import (
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/util2"
)

var Router = httprouter.New()

func Handle (method, path string, handle util2.ErrorHandle) {
  Router.Handle(method, path, util2.ErrorHandleFunc(handle))
}
func GET     (path string, handle util2.ErrorHandle) {Handle("GET",     path, handle)}
func PUT     (path string, handle util2.ErrorHandle) {Handle("PUT",     path, handle)}
func POST    (path string, handle util2.ErrorHandle) {Handle("POST",    path, handle)}
func DELETE  (path string, handle util2.ErrorHandle) {Handle("DELETE",  path, handle)}
func PATCH   (path string, handle util2.ErrorHandle) {Handle("PATCH",   path, handle)}
func OPTIONS (path string, handle util2.ErrorHandle) {Handle("OPTIONS", path, handle)}
func HEAD    (path string, handle util2.ErrorHandle) {Handle("HEAD",    path, handle)}
