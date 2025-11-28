package main

import (
    "net/http"
)

type Route struct {
    method string
    path string
    handler http.HandlerFunc
}

type Router struct {
    routes []Route
}

func NewRouter() Router {
    return Router{}
}

func (r *Router) defineRoute(method string, path string, handler http.HandlerFunc) {
    r.routes = append(r.routes, Route{
        method: method,
        path: path,
        handler: handler,
    })
}

func (r *Router) Get(path string, handler http.HandlerFunc) {
    r.defineRoute("GET", path, handler)
}

func (r *Router) Post(path string, handler http.HandlerFunc) {
    r.defineRoute("POST", path, handler)
}

func (router *Router) Serve(w http.ResponseWriter, r *http.Request) bool {
    found := false
    path := r.URL.Path

    for _, route := range router.routes {
        if r.Method != route.method || path != route.path {
            continue
        }

        found = true;
        route.handler(w, r)
        break
    }

    return found
}
