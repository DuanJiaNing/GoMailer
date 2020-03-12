package auth

import (
	"net/http"

	"GoMailer/common/key"
	"GoMailer/handler/userapp"
)

var (
	noGuardRequiredAPI = map[string]struct{}{
		"/api/shortcut": {},
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
			a, err := userapp.FindById(ak.AppId)
			if err != nil {
				http.Error(w, "illegal app_key, app may not exist", http.StatusUnauthorized)
				return
			}
			if a.UserId != ak.UserId {
				http.Error(w, "invalid app_key", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
