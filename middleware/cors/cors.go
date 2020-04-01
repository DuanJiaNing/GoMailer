package cors

import (
	"net/http"

	"github.com/gorilla/mux"

	"GoMailer/handler/endpoint"
	"GoMailer/handler/userapp"
	"GoMailer/log"
	"gowebsitemailer/common/key"
)

func CORS(r *mux.Router) func(http.Handler) http.Handler {
	// required so we don't get a code 405
	r.Methods(http.MethodOptions).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// No need to verify, already pass the Guard.
			ak := key.AppKeyFromRequest(r)
			ep, _ := endpoint.FindByKey(ak)
			app, err := userapp.FindById(ep.UserAppId)
			if err != nil {
				log.Error("got error when find host for CORS origin: user app not exist")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", app.Host)
			w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept-Encoding, User-Agent, Accept")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				// we only need headers for OPTIONS request, no need to go down the handler chain
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
