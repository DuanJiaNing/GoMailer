package shortcut

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"GoMailer/app"
	"GoMailer/common/db"
	"GoMailer/common/utils"
	"GoMailer/handler"
	"GoMailer/handler/dialer"
	"GoMailer/handler/endpoint"
	"GoMailer/handler/endpoint/preference"
	"GoMailer/handler/endpoint/receiver"
	"GoMailer/handler/template"
	"GoMailer/handler/user"
	"GoMailer/handler/userapp"
)

const errInvalidParameter = "invalid parameter"

func init() {
	handler.ShortcutRouter.Methods(http.MethodPost).Handler(app.Handler(shortcut))
}

type shortcutVO struct {
	User     *db.User
	App      *db.UserApp
	Endpoint *struct {
		Name       string
		Dialer     *db.Dialer
		Receiver   []*db.Receiver
		Template   *db.Template
		Preference *db.EndpointPreference
	}
}

// shortcut is a short way to create or update a endpoint
// 1. create user if not registered - required
// 2. create app for user if not created, else update it - required
// 3. check dialer exists or not for user, create dialer when not, update when exists
// 4. create template for user
// 5. create or update endpoint - required
// 6. create or update preference for endpoint
// 7. add or update receiver for endpoint
func shortcut(w http.ResponseWriter, r *http.Request) (interface{}, *app.Error) {
	vo := &shortcutVO{}
	aerr := app.JsonUnmarshalFromRequest(r, vo)
	if aerr != nil {
		return nil, aerr
	}

	if vo.User == nil || vo.Endpoint == nil || vo.App == nil {
		return nil, app.Errorf(errors.New("find nil value when validate parameter"), errInvalidParameter)
	}

	// 1. create user if not registered
	// vo.User is required
	user, err := handleUser(vo.User)
	if err != nil {
		return nil, err
	}

	// 2. create app for user if not created, else update it
	// vo.App is required
	userApp, err := handleUserApp(user, vo.App)
	if err != nil {
		return nil, err
	}

	// vo.Endpoint is required
	if utils.IsBlankStr(vo.Endpoint.Name) {
		return nil, app.Errorf(errors.New("endpoint name can not be empty"), errInvalidParameter)
	}
	// 3. check dialer exists or not for user, create dialer when not
	dialer, err := handleUserDialer(user, vo.Endpoint.Dialer)
	if err != nil {
		return nil, err
	}

	// 4. create template for user
	template, err := handleUserTemplate(user, vo.Endpoint.Template)
	if err != nil {
		return nil, err
	}

	// 5. create or update endpoint
	// vo.Endpoint is required
	ep, err := handleEndpoint(vo.Endpoint.Name, user, userApp, dialer, template)
	if err != nil {
		return nil, err
	}

	// 6. create or update preference for endpoint
	_, err = handleEndPointPreference(ep, vo.Endpoint.Preference)
	if err != nil {
		return nil, err
	}

	// 7. add or update receiver for endpoint
	err = handleEndPointReceiver(ep, user, userApp, vo.Endpoint.Receiver)
	if err != nil {
		return nil, err
	}

	return struct {
		UserId int64
		AppKey string
	}{UserId: user.Id, AppKey: ep.Key}, nil
}

func handleEndPointReceiver(ep *db.Endpoint, u *db.User, ua *db.UserApp, r []*db.Receiver) *app.Error {
	if len(r) == 0 {
		return nil
	}

	err := receiver.DeleteByEndpoint(ep.Id)
	if err != nil {
		return app.Errorf(err, "failed to delete all receiver for endpoint receiver update")
	}
	for _, r := range r {
		r.EndpointId = ep.Id
		r.UserId = u.Id
		r.UserAppId = ua.Id
	}
	err = receiver.PatchCreate(r)
	if err != nil {
		return app.Errorf(err, "failed to create receiver for endpoint")
	}

	return nil
}

func handleEndPointPreference(ep *db.Endpoint, p *db.EndpointPreference) (*db.EndpointPreference, *app.Error) {
	if p == nil {
		return nil, nil
	}

	p.EndpointId = ep.Id
	pre, err := preference.FindByEndpoint(ep.Id)
	if err != nil {
		return nil, app.Errorf(err, "failed to find endpoint preference")
	}
	if pre == nil {
		pre, err = preference.Create(p)
		if err != nil {
			return nil, app.Errorf(err, "failed to create endpoint preference")
		}
	} else {
		p.Id = pre.Id
		pre, err = preference.Update(p)
		if err != nil {
			return nil, app.Errorf(err, "failed to update endpoint preference")
		}
	}

	return pre, nil
}

func handleEndpoint(name string, u *db.User, ap *db.UserApp, ud *db.Dialer, ut *db.Template) (
	*db.Endpoint, *app.Error) {
	if utils.IsBlankStr(name) {
		return nil, app.Errorf(errors.New("endpoint name can not be empty"), errInvalidParameter)
	}

	ep, err := endpoint.FindByName(ap.Id, name)
	if err != nil {
		return nil, app.Errorf(err, "failed to find endpoint")
	}

	nep := &db.Endpoint{}
	nep.Name = name
	nep.UserId = u.Id
	nep.UserAppId = ap.Id
	if ud != nil {
		nep.DialerId = ud.Id
	}
	if ut != nil {
		nep.TemplateId = ut.Id
	}
	if ep == nil {
		if ud == nil || ut == nil {
			return nil, app.Errorf(err, "dialer and template is required when create endpoint")
		}
		ep, err = endpoint.Create(nep)
		if err != nil {
			return nil, app.Errorf(err, "failed to create endpoint")
		}
		key, err := endpoint.RefreshKey(ep.Id)
		if err != nil {
			return nil, app.Errorf(err, "failed to generate app key")
		}
		ep.Key = key
	} else {
		nep.Id = ep.Id
		ep, err = endpoint.Update(nep)
		if err != nil {
			return nil, app.Errorf(err, "failed to update endpoint")
		}
	}

	return ep, nil
}

func handleUserTemplate(u *db.User, t *db.Template) (*db.Template, *app.Error) {
	if t == nil {
		return nil, nil
	}

	t.UserId = u.Id
	utemplate, err := template.Create(t)
	if err != nil {
		return nil, app.Errorf(err, "failed to create template")
	}

	return utemplate, nil
}

func handleUserDialer(u *db.User, d *db.Dialer) (*db.Dialer, *app.Error) {
	if d == nil {
		return nil, nil
	}

	if utils.IsBlankStr(d.Name) {
		return nil, app.Errorf(errors.New("dialer name can not be empty"), errInvalidParameter)
	}

	d.UserId = u.Id
	udialer, err := dialer.FindByName(d.UserId, d.Name)
	if err != nil {
		return nil, app.Errorf(err, "failed to get user dialer")
	}
	if udialer == nil {
		udialer, err = dialer.Create(d)
		if err != nil {
			return nil, app.Errorf(err, "failed to create dialer")
		}
	} else {
		d.Id = udialer.Id
		udialer, err = dialer.Update(d)
		if err != nil {
			return nil, app.Errorf(err, "failed to update dialer")
		}
	}

	return udialer, nil
}

func handleUserApp(u *db.User, ua *db.UserApp) (*db.UserApp, *app.Error) {
	if ua == nil {
		return nil, app.Errorf(errors.New("app can not be empty"), errInvalidParameter)
	}
	if utils.IsBlankStr(ua.Name) {
		return nil, app.Errorf(errors.New("app name can not be empty"), errInvalidParameter)
	}
	ua.UserId = u.Id

	uapp, err := userapp.FindByName(ua.UserId, ua.Name)
	if err != nil {
		return nil, app.Errorf(err, "failed to get user app")
	}
	if uapp == nil {
		uapp, err = userapp.Create(ua)
		if err != nil {
			return nil, app.Errorf(err, "failed to create app")
		}
	} else {
		ua.Id = uapp.Id
		uapp, err = userapp.Update(ua)
		if err != nil {
			return nil, app.Errorf(err, "failed to update app")
		}
	}

	return uapp, nil
}

func handleUser(u *db.User) (*db.User, *app.Error) {
	if utils.IsBlankStr(u.Username) || utils.IsBlankStr(u.Password) {
		return nil, app.Errorf(errors.New("username or password can not be empty"), errInvalidParameter)
	}

	us, err := user.FindByName(u.Username)
	if err != nil {
		return nil, app.Errorf(err, "failed to get user")
	}
	ps := sha256.Sum256([]byte(u.Password))
	u.Password = hex.EncodeToString(ps[:])
	if us == nil {
		us, err = user.Create(u)
		if err != nil {
			return nil, app.Errorf(err, "failed to create user")
		}
	} else {
		if us.Password != u.Password {
			return nil, app.Errorf(errors.New("password incorrect"), "wrong password")
		}
	}

	return us, nil
}
