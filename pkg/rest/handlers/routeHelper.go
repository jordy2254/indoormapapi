package handlers

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/gorilla/mux"
	"net/http"
)

type RouteHelper struct{
	router *mux.Router
	auth *jwtmiddleware.JWTMiddleware
}

func NewRouteHelper(router *mux.Router, auth *jwtmiddleware.JWTMiddleware) *RouteHelper {
	return &RouteHelper{
		router: router,
		auth:   auth,
	}
}

func (rb *RouteHelper) protectedRoute(path string, f http.HandlerFunc) *mux.Route {
	return rb.router.Handle(path, rb.auth.Handler(f))
}

func (rb *RouteHelper) openRoute(path string, f http.HandlerFunc) *mux.Route {
	return rb.router.HandleFunc(path,f)
}