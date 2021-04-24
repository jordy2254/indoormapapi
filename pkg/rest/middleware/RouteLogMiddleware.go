package middleware

import (
	"github.com/op/go-logging"
	"net/http"
)

type RouteLogMiddleware struct{
	logger *logging.Logger
}

func NewRouteLogger(logger *logging.Logger) *RouteLogMiddleware{
	return &RouteLogMiddleware{logger: logger}
}

func (m *RouteLogMiddleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Infof("Endpoint hit %s, %s", r.RequestURI, r.Method)
		h.ServeHTTP(w, r)
	})
}