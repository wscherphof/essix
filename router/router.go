package router

import (
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/util"
)

var Router = httprouter.New()

func Handle (method, path string, handle util.ErrorHandle) {
  Router.Handle(method, path, util.ErrorHandleFunc(handle))
}
func GET     (path string, handle util.ErrorHandle) {Handle("GET",     path, handle)}
func PUT     (path string, handle util.ErrorHandle) {Handle("PUT",     path, handle)}
func POST    (path string, handle util.ErrorHandle) {Handle("POST",    path, handle)}
func DELETE  (path string, handle util.ErrorHandle) {Handle("DELETE",  path, handle)}
func PATCH   (path string, handle util.ErrorHandle) {Handle("PATCH",   path, handle)}
func OPTIONS (path string, handle util.ErrorHandle) {Handle("OPTIONS", path, handle)}
func HEAD    (path string, handle util.ErrorHandle) {Handle("HEAD",    path, handle)}
