package auth

import (
	"net/http"

	"GoMailer/common/key"
	"GoMailer/common/utils"
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
		if _, ok := noGuardRequiredAPI[r.URL.Path]; !ok {
			ak := key.AppKeyFromRequest(r)
			if utils.IsBlankStr(ak) {
				http.Error(w, "app_key is required", http.StatusUnauthorized)
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
