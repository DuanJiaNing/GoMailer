package mail

import (
	"errors"
	"net/http"
	"strings"

	"GoMailer/app"
	"GoMailer/common/db"
	"GoMailer/common/key"
	"GoMailer/common/utils"
	"GoMailer/handler"
	"GoMailer/handler/endpoint"
	"GoMailer/log"
)

func init() {
	router := handler.MailRouter.Path("/send").Subrouter()
	router.Methods(http.MethodPost).Handler(app.Handler(send))
}

func send(w http.ResponseWriter, r *http.Request) (interface{}, *app.Error) {
	if err := r.ParseForm(); err != nil {
		return nil, app.Errorf(err, "failed to parse form")
	}
	epk := key.EPKeyFromRequest(r)
	ep, err := endpoint.FindByKey(epk)
	if err != nil {
		return nil, app.Errorf(err, "failed to find endpoint by key")
	}

	raw, err := parseForm(r)
	if err != nil {
		sendRedirect(w, ep.Id, err)
		return nil, app.Errorf(err, "got error when parse form")
	}
	mail, err := handleMail(ep.Id, raw, key.ReCaptchaKeyFromRequest(r))
	if mail != nil {
		_, err := Create(ep.UserId, mail)
		if err != nil {
			sendRedirect(w, ep.Id, err)
			return nil, app.Errorf(err, "fail to store mail")
		}
	}
	if err != nil {
		sendRedirect(w, ep.Id, err)
		return nil, app.Errorf(err, "failed to deliver mail")
	}

	sendRedirect(w, ep.Id, nil)
	return nil, nil
}

func sendRedirect(w http.ResponseWriter, endpointId int64, err error) {
	client, ierr := db.NewClient()
	if ierr != nil {
		http.Error(w, ierr.Error(), http.StatusInternalServerError)
		log.Errorf("got err when set up redirect header: %v", ierr)
		return
	}

	ep := &db.EndpointPreference{}
	get, ierr := client.Where("endpoint_id = ?", endpointId).Get(ep)
	if ierr != nil {
		http.Error(w, ierr.Error(), http.StatusInternalServerError)
		log.Errorf("got err when set up redirect header: %v", ierr)
		return
	}
	if !get {
		// No preference yet, ignore.
		return
	}

	if err != nil && !utils.IsBlankStr(ep.FailRedirect) {
		w.WriteHeader(http.StatusFound)
		w.Header().Set("Location", ep.FailRedirect+"?err="+err.Error())
	}
	if err == nil && !utils.IsBlankStr(ep.SuccessRedirect) {
		w.WriteHeader(http.StatusFound)
		w.Header().Set("Location", ep.SuccessRedirect)
	}
}

func parseForm(r *http.Request) (map[string]string, error) {
	data := make(map[string]string)
	allBlank := true
	for k, vs := range r.Form {
		if k == key.EPKeyName || k == key.ReCaptchaTokenKeyName {
			continue
		}

		str := strings.TrimSpace(vs[0])
		if len(str) > 0 {
			allBlank = false
		}
		data[k] = str
	}
	if allBlank {
		return nil, errors.New("not allow to send empty content")
	}

	return data, nil
}
