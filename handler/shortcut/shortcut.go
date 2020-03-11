package shortcut

import (
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
	EndPoint *struct {
		Name       string
		Dialer     *db.Dialer
		Receiver   []*db.Receiver
		Template   *db.Template
		Preference *db.EndPointPreference
	}
}

// shortcut is a short way to create or update a end point
// 1. create user if not registered - required
// 2. create app for user if not created, else update it - required
// 3. check dialer exists or not for user, create dialer when not, update when exists
// 4. create template for user
// 5. create or update end point - required
// 6. create or update preference for end point
// 7. add or update receiver for end point
func shortcut(w http.ResponseWriter, r *http.Request) (interface{}, *app.Error) {
	vo := &shortcutVO{}
	aerr := app.JsonUnmarshalFromRequest(r, vo)
	if aerr != nil {
		return nil, aerr
	}

	if vo.User == nil || vo.EndPoint == nil || vo.App == nil {
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

	// vo.EndPoint is required
	if utils.IsStrBlank(vo.EndPoint.Name) {
		return nil, app.Errorf(errors.New("end point name can not be empty"), errInvalidParameter)
	}
	// 3. check dialer exists or not for user, create dialer when not
	dialer, err := handleUserDialer(user, vo.EndPoint.Dialer)
	if err != nil {
		return nil, err
	}

	// 4. create template for user
	template, err := handleUserTemplate(user, vo.EndPoint.Template)
	if err != nil {
		return nil, err
	}

	// 5. create or update end point
	// vo.EndPoint is required
	endpoint, err := handleEndPoint(vo.EndPoint.Name, user, userApp, dialer, template)
	if err != nil {
		return nil, err
	}

	// 6. create or update preference for end point
	_, err = handleEndPointPreference(endpoint, vo.EndPoint.Preference)
	if err != nil {
		return nil, err
	}

	// 7. add or update receiver for end point
	err = handleEndPointReceiver(endpoint, user, userApp, vo.EndPoint.Receiver)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func handleEndPointReceiver(ep *db.EndPoint, u *db.User, ua *db.UserApp, r []*db.Receiver) *app.Error {
	if len(r) == 0 {
		return nil
	}

	err := receiver.Delete(ep.ID)
	if err != nil {
		return app.Errorf(err, "failed to delete all receiver for end point receiver update")
	}
	for _, r := range r {
		r.EndPointID = ep.ID
		r.UserID = u.ID
		r.UserAppID = ua.ID
	}
	err = receiver.PatchCreate(r)
	if err != nil {
		return app.Errorf(err, "failed to create receiver for end point")
	}

	return nil
}

func handleEndPointPreference(ep *db.EndPoint, p *db.EndPointPreference) (*db.EndPointPreference, *app.Error) {
	if p == nil {
		return nil, nil
	}

	p.EndPointID = ep.ID
	pre, err := preference.Find(ep.ID)
	if err != nil {
		return nil, app.Errorf(err, "failed to find end point preference")
	}
	if pre == nil {
		pre, err = preference.Create(p)
		if err != nil {
			return nil, app.Errorf(err, "failed to create end point preference")
		}
	} else {
		p.ID = pre.ID
		pre, err = preference.Update(p)
		if err != nil {
			return nil, app.Errorf(err, "failed to update end point preference")
		}
	}

	return pre, nil
}

func handleEndPoint(name string, u *db.User, ap *db.UserApp, ud *db.Dialer, ut *db.Template) (
	*db.EndPoint, *app.Error) {
	if utils.IsStrBlank(name) {
		return nil, app.Errorf(errors.New("end point name can not be empty"), errInvalidParameter)
	}

	ep, err := endpoint.FindByName(name)
	if err != nil {
		return nil, app.Errorf(err, "failed to find end point")
	}

	nep := &db.EndPoint{}
	nep.Name = name
	nep.UserID = u.ID
	nep.UserAppID = ap.ID
	if ud != nil {
		nep.DialerID = ud.ID
	}
	if ut != nil {
		nep.TemplateID = ut.ID
	}
	if ep == nil {
		if ud == nil || ut == nil {
			return nil, app.Errorf(err, "dialer and template is required when create end point")
		}
		ep, err = endpoint.Create(nep)
		if err != nil {
			return nil, app.Errorf(err, "failed to create end point")
		}
	} else {
		nep.ID = ep.ID
		ep, err = endpoint.Update(nep)
		if err != nil {
			return nil, app.Errorf(err, "failed to update end point")
		}
	}

	return ep, nil
}

func handleUserTemplate(u *db.User, t *db.Template) (*db.Template, *app.Error) {
	if t == nil {
		return nil, nil
	}

	t.UserID = u.ID
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

	if utils.IsStrBlank(d.Name) {
		return nil, app.Errorf(errors.New("dialer name can not be empty"), errInvalidParameter)
	}

	d.UserID = u.ID
	udialer, err := dialer.GetByName(d.UserID, d.Name)
	if err != nil {
		return nil, app.Errorf(err, "failed to get user dialer")
	}
	if udialer == nil {
		udialer, err = dialer.Create(d)
		if err != nil {
			return nil, app.Errorf(err, "failed to create dialer")
		}
	} else {
		d.ID = udialer.ID
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
	if utils.IsStrBlank(ua.Name) {
		return nil, app.Errorf(errors.New("app name can not be empty"), errInvalidParameter)
	}
	ua.UserID = u.ID

	uapp, err := userapp.GetByName(ua.UserID, ua.Name)
	if err != nil {
		return nil, app.Errorf(err, "failed to get user app")
	}
	if uapp == nil {
		uapp, err = userapp.Create(ua)
		if err != nil {
			return nil, app.Errorf(err, "failed to create app")
		}
	} else {
		ua.ID = uapp.ID
		uapp, err = userapp.Update(ua)
		if err != nil {
			return nil, app.Errorf(err, "failed to update app")
		}
	}

	return uapp, nil
}

func handleUser(u *db.User) (*db.User, *app.Error) {
	if utils.IsStrBlank(u.Username) || utils.IsStrBlank(u.Password) {
		return nil, app.Errorf(errors.New("username or password can not be empty"), errInvalidParameter)
	}

	u, err := user.GetByName(u.Username)
	if err != nil {
		return nil, app.Errorf(err, "failed to get user")
	}
	if u == nil {
		u, err = user.Create(u)
		if err != nil {
			return nil, app.Errorf(err, "failed to create user")
		}
	}

	return u, nil
}
