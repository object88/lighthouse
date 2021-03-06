package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/object88/lighthouse/pkg/http/router/route"
)

type DefaultRoute func(rtr *Router) http.HandlerFunc

func NoopDefaultRoute(rtr *Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}
}

func LoggingDefaultRoute(rtr *Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rtr.Logger.Info("Unhandled route", "URL", r.URL)
		w.WriteHeader(404)
	}
}

type Router struct {
	m      *mux.Router
	Logger logr.Logger
}

// New creates a new Router
func New(logger logr.Logger) *Router {
	rtr := &Router{
		m:      mux.NewRouter(),
		Logger: logger,
	}
	return rtr
}

func (rtr *Router) Route(defaultRoute DefaultRoute, routes []*route.Route) (*mux.Router, error) {
	if err := rtr.configureRoutes(rtr.m, routes); err != nil {
		return nil, err
	}

	rtr.m.PathPrefix("/").HandlerFunc(defaultRoute(rtr))

	rtr.reportRoutes()

	return rtr.m, nil
}

func (rtr *Router) reportRoutes() {
	report := func(r *mux.Route) {
		pathTemplate, _ := r.GetPathTemplate()
		pathRegexp, _ := r.GetPathRegexp()
		queriesTemplates, _ := r.GetQueriesTemplates()
		queriesRegexps, _ := r.GetQueriesRegexp()

		rtr.Logger.Info("Route",
			"name", r.GetName(),
			"path", pathTemplate,
			"path-regexp", pathRegexp,
			"queries-templates", strings.Join(queriesTemplates, ","),
			"queries-regexps", strings.Join(queriesRegexps, ","),
			"has-handler", r.GetHandler() != nil)
	}

	rtr.m.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		report(route)
		return nil
	})
}

func (rtr *Router) configureRoutes(base *mux.Router, routes []*route.Route) error {
	for _, r := range routes {
		if r.Handler != nil {
			s := base.NewRoute().Subrouter()

			if len(r.Middleware) != 0 {
				s.Use(r.Middleware...)
			}

			rt := s.Path(r.Path)
			rt = rt.Handler(r.Handler)

			if len(r.Methods) != 0 {
				rt = rt.Methods(r.Methods...)
			}
			if len(r.Queries) != 0 {
				qs := make([]string, len(r.Queries)*2)
				for k, v := range r.Queries {
					qs[k*2] = v
					qs[k*2+1] = fmt.Sprintf("{%s}", v)
				}
				rt = rt.Queries(qs...)
			}
			err := rt.GetError()
			if err != nil {
				return fmt.Errorf("failed to create route for path '%s': %w", r.Path, err)
			}
		}

		if len(r.Subroutes) != 0 {
			sub := base.PathPrefix(r.Path).Subrouter()

			if len(r.Middleware) != 0 {
				for _, mfunc := range r.Middleware {
					sub.Use(mfunc)
				}
			}

			rtr.configureRoutes(sub, r.Subroutes)
		}
	}

	return nil
}
