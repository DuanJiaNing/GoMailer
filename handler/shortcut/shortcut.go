package shortcut

import (
	"net/http"

	"GoMailer/app"
	"GoMailer/handler"
)

func init() {
	router := handler.APIRouter.Path("/shortcut").Subrouter()
	router.Methods(http.MethodPost).Handler(app.Handler(shortcut))
}

func shortcut(http.ResponseWriter, *http.Request) (interface{}, *app.Error) {
	return nil, nil
}
