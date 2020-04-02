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
	"GoMailer/handler/userapp"
	"GoMailer/log"
)

func init() {
	router := handler.MailRouter.Path("/send").Subrouter()
	router.Methods(http.MethodPost).Handler(app.Handler(send))
}

func send(w http.ResponseWriter, r *http.Request) (interface{}, *app.Error) {
	epk := key.EPKeyFromRequest(r)
	ep, err := endpoint.FindByKey(epk)
	if err != nil {
		return nil, app.Errorf(err, "failed to find endpoint by key")
	}

	raw, err := parseForm(r, ep)
	if err != nil {
		setupRedirectHeader(w, ep, err, r)
		return nil, app.Errorf(err, "got error when parse form")
	}
	mail, err := handleMail(ep.Id, raw, key.ReCaptchaKeyFromRequest(r))
	if mail != nil {
		_, err := Create(ep.UserId, mail)
		if err != nil {
			setupRedirectHeader(w, ep, err, r)
			return nil, app.Errorf(err, "fail to store mail")
		}
	}
	if err != nil {
		setupRedirectHeader(w, ep, err, r)
		return nil, app.Errorf(err, "failed to deliver mail")
	}

	setupRedirectHeader(w, ep, nil, r)
	return nil, nil
}

func setupRedirectHeader(w http.ResponseWriter, ep *db.Endpoint, oerr error, r *http.Request) {
	client, err := db.NewClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf("got err when set up redirect header: %v", err)
		return
	}

	epp := &db.EndpointPreference{}
	get, err := client.Where("endpoint_id = ?", ep.Id).Get(epp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf("got err when set up redirect header: %v", err)
		return
	}
	if !get {
		// No preference yet, ignore.
		return
	}
	ua, err := userapp.FindById(ep.UserAppId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf("got err when set up redirect header: %v", err)
		return
	}

	setup := func(addr string) {
		if ua.AppType == db.AppType_AMP_WEB.Name() {
			w.Header().Set("AMP-Redirect-To", addr)
			w.Header().Set("Access-Control-Expose-Headers", "AMP-Access-Control-Allow-Source-Origin, AMP-Redirect-To")
			return
		}

		w.WriteHeader(http.StatusFound)
		w.Header().Set("Location", addr)
	}
	if oerr != nil && !utils.IsBlankStr(epp.FailRedirect) {
		setup(epp.FailRedirect + "?err=" + oerr.Error())
	}
	if oerr == nil && !utils.IsBlankStr(epp.SuccessRedirect) {
		setup(epp.SuccessRedirect)
	}
}

func parseForm(r *http.Request, ep *db.Endpoint) (map[string]string, error) {
	ua, err := userapp.FindById(ep.UserAppId)
	if err != nil {
		return nil, errors.New("failed to find user app")
	}

	const defaultMaxMemory = 5 << 20         // 5 MB
	if ua.AppType == db.AppType_WEB.Name() { // application/x-www-form-urlencoded
		if err = r.ParseForm(); err != nil {
			return nil, errors.New("failed to parse form")
		}
	} else { // multipart/form-data
		if err = r.ParseMultipartForm(defaultMaxMemory); err != nil {
			return nil, errors.New("failed to parse multipart form")
		}
	}

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
