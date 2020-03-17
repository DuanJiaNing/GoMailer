package mail

import (
	"errors"
	"net/http"
	"strings"

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

	data := make(map[string]string)
	allBlank := true
	for k, vs := range r.Form {
		if k == "app_key" {
			continue
		}

		str := strings.TrimSpace(vs[0])
		if len(str) > 0 {
			allBlank = false
		}
		data[k] = str
	}
	if allBlank {
		return nil, app.Errorf(errors.New("invalid parameter"), "not allow to send empty content")
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
