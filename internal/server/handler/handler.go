package handler

import (
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	appMiddleware "github.com/e1m0re/grdn/internal/server/middleware"
	"github.com/e1m0re/grdn/internal/server/service"
)

type Handler struct {
	services *service.Services
}

// NewHandler is Handler constructor.
func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

// NewRouter initializes new router.
func (h *Handler) NewRouter(signKey string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(appMiddleware.Logging())
	r.Use(appMiddleware.UnzipContent())
	if len(signKey) > 0 {
		r.Use(appMiddleware.SignChecking(signKey))
	}
	r.Use(middleware.Compress(5, "text/html", "application/json"))
	if len(signKey) > 0 {
		r.Use(appMiddleware.SignResponse(signKey))
	}

	r.Route("/", func(r chi.Router) {
		r.Get("/", h.getMainPage)
		r.Get("/ping", h.checkDBConnection)
		r.Route("/value", func(r chi.Router) {
			r.Post("/", h.getMetricValueV2)
			r.Get("/{mType}/{mName}", h.getMetricValue)
		})
		r.Route("/update", func(r chi.Router) {
			r.Post("/", h.updateMetricV2)
			r.Post("/{mType}/{mName}/{mValue}", h.updateMetric)
		})
		r.Route("/updates", func(r chi.Router) {
			r.Post("/", h.updateMetricsList)
		})

		r.Route("/debug/pprof/", func(r chi.Router) {
			r.Get("/", pprof.Index)
			r.Get("/cmdline", pprof.Cmdline)
			r.Get("/profile", pprof.Profile)
			r.Get("/symbol", pprof.Symbol)
			r.Get("/trace", pprof.Trace)
		})
	})

	return r
}
