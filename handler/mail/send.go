package mail

import (
	"net/http"

	"GoMailer/app"
	"GoMailer/handler"
	"GoMailer/log"
)

func init() {
	router := handler.MailRouter.Path("/send").Subrouter()
	router.Methods(http.MethodPost).Handler(app.ServiceHandler(send))
}

func send(http.ResponseWriter, *http.Request) (interface{}, *app.ServerError) {
	log.Info("got...")
	return nil, nil
}
