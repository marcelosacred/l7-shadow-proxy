package router

import "net/http"

type Route struct {
	Name 	 string
	Upstream string
}

type Router interface {
	Match(r *http.Request) (*Route, bool)
}