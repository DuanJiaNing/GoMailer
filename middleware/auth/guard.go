package auth

import (
	"net/http"

	"GoMailer/common/key"
	"GoMailer/handler/endpoint"
)

var (
	noGuardRequiredAPI = map[string]struct{}{
		"/api/shortcut":  {},
		"/api/mail/list": {},
	}
)

func Guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noGuardRequiredAPI[r.URL.String()]; !ok {
			ak, err := key.AppKeyFromRequest(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ep, err := endpoint.FindByKey(ak)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if ep == nil {
				http.Error(w, "endpoint not exist", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
