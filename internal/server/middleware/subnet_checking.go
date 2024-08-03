package middleware

import (
	"net/http"
	"net/netip"
)

// SubnetChecking executes check of requests sign.
func SubnetChecking(network *netip.Prefix) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				next.ServeHTTP(w, r)
				return
			}

			ip, err := netip.ParseAddr(r.Header.Get("X-Real-IP"))
			switch true {
			case err != nil:
				fallthrough
			case !network.Contains(ip):
				w.WriteHeader(http.StatusForbidden)
				return
			default:
				next.ServeHTTP(w, r)
			}
		})
	}
}
