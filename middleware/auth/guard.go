package auth

import (
	"net/http"

	"GoMailer/common/key"
	"GoMailer/common/utils"
	"GoMailer/handler/endpoint"
)

var (
	noGuardRequiredAPI = map[string]struct{}{
		"/api/shortcut": {},
	}
)

func Guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noGuardRequiredAPI[r.URL.Path]; !ok {
			epk := key.EPKeyFromRequest(r)
			if utils.IsBlankStr(epk) {
				http.Error(w, "epKey is required", http.StatusUnauthorized)
				return
			}
			ep, err := endpoint.FindByKey(epk)
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
