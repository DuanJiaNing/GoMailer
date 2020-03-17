package mail

import (
	"net/http"

	"GoMailer/app"
	"GoMailer/common/key"
	"GoMailer/handler"
)

func init() {
	router := handler.MailRouter.Path("/send").Subrouter()
	router.Methods(http.MethodPost).Handler(app.Handler(send))
}

func send(w http.ResponseWriter, r *http.Request) (interface{}, *app.Error) {
	ak, _ := key.AppKeyFromRequest(r)
	if err := r.ParseForm(); err != nil {
		return nil, app.Errorf(err, "failed to parse form request")
	}

	data := make(map[string]interface{})
	for k, vs := range r.Form {
		data[k] = vs[0]
	}
	mail, err := handleMail(ak.EndpointId, data)
	if mail != nil {
		_, err := create(mail)
		if err != nil {
			return nil, app.Errorf(err, "fail to store mail")
		}
	}
	if err != nil {
		return nil, app.Errorf(err, "failed to deliver mail")
	}

	return ak, nil
}
