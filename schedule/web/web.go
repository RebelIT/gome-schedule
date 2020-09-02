package web

import (
	"github.com/gorilla/mux"
	"github.com/rebelit/gome-schedule/common/config"
	"github.com/rebelit/gome-schedule/common/stat"
	"github.com/rebelit/gome-schedule/schedule"
	"log"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
		router.Use(authMiddleware)
	}

	return router
}

var routes = Routes{
	Route{"status", "GET", "/api/status", status},
	Route{"state", "GET", "/api/schedule/state", schedule.StateGetAll},
	Route{"state", "GET", "/api/schedule/state/{friendlyName}", schedule.StateGet},
	Route{"state", "POST", "/api/schedule/state/{friendlyName}", schedule.StateNew},
	Route{"state", "DELETE", "/api/schedule/state/{friendlyName}", schedule.StateDelete},
	Route{"toggle", "GET", "/api/schedule/toggle", schedule.ToggleGetAll},
	Route{"toggle", "GET", "/api/schedule/toggle/{friendlyName}", schedule.ToggleGet},
	Route{"toggle", "POST", "/api/schedule/toggle/{friendlyName}", schedule.ToggleNew},
	Route{"toggle", "DELETE", "/api/schedule/toggle/{friendlyName}", schedule.ToggleDelete},
}

func status(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	stat.Http(r.Method, "inbound", r.URL.String(), http.StatusOK)

	return
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/status" {
			//skip auth for health check
			next.ServeHTTP(w, r)

		} else {
			authorization := r.Header.Get("Authorization")
			if validateAuth(authorization) {
				// Pass down the request to the next handler
				next.ServeHTTP(w, r)

			} else {
				log.Printf("WARN: http unauthorized %s:%s from %s", r.Method, r.URL.String(), r.RemoteAddr)
				w.WriteHeader(http.StatusUnauthorized)
				stat.Http(r.Method, stat.HTTPIN, r.URL.String(), http.StatusUnauthorized)

				return
			}
		}
	})
}

func validateAuth(authorization string) bool {
	if authorization == "Bearer "+config.App.AuthToken {
		return true
	}

	return false
}

func GenerateSpec() {
	router := NewRouter()
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tmpl, _ := route.GetPathTemplate()
		log.Printf("route: %s", tmpl)
		return nil
	})
	if err != nil {

	}
	return
}
