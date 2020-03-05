package cors

import (
	"net/http"

	"github.com/gorilla/mux"

	"GoMailer/conf"
)

func CORS(r *mux.Router) func(http.Handler) http.Handler {
	// required so we don't get a code 405
	r.Methods(http.MethodOptions).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", conf.Cors().AllowedOrigins)
			w.Header().Set("Access-Control-Allow-Methods", conf.Cors().AllowedMethods)
			w.Header().Set("Access-Control-Allow-Headers", conf.Cors().AllowedHeaders)
			w.Header().Set("Access-Control-Max-Age", conf.Cors().MaxAge)

			if r.Method == http.MethodOptions {
				// we only need headers for OPTIONS request, no need to go down the handler chain
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
