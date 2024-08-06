package api

import (
	"net/http/pprof"
	"net/netip"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	appMiddleware "github.com/e1m0re/grdn/internal/middleware"
	"github.com/e1m0re/grdn/internal/service"
)

type Config struct {
	TrustedSubnet  *netip.Prefix
	SignKey        string
	PrivateKeyFile string
}

type Handler struct {
	services *service.ServerServices
	config   Config
}

// NewHandler is Handler constructor.
func NewHandler(services *service.ServerServices, cfg Config) *Handler {
	return &Handler{
		services: services,
		config:   cfg,
	}
}

// NewRouter initializes new router.
func (h *Handler) NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(appMiddleware.Logging())
	if h.config.TrustedSubnet != nil {
		r.Use(appMiddleware.SubnetChecking(h.config.TrustedSubnet))
	}
	r.Use(appMiddleware.UnzipContent())
	if len(h.config.SignKey) > 0 {
		r.Use(appMiddleware.SignChecking(h.config.SignKey))
	}
	r.Use(middleware.Compress(5, "text/html", "application/json"))
	if len(h.config.PrivateKeyFile) > 0 {
		r.Use(appMiddleware.DecryptContent(h.config.PrivateKeyFile))
	}
	if len(h.config.SignKey) > 0 {
		r.Use(appMiddleware.SignResponse(h.config.SignKey))
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
