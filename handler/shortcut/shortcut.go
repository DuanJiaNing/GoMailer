package shortcut

import (
	"errors"
	"net/http"

	"GoMailer/app"
	"GoMailer/common/db"
	"GoMailer/common/utils"
	"GoMailer/handler"
	happ "GoMailer/handler/app"
	"GoMailer/handler/dialer"
	"GoMailer/handler/user"
)

func init() {
	router := handler.APIRouter.Path("/shortcut").Subrouter()
	router.Methods(http.MethodPost).Handler(app.Handler(shortcut))
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

func shortcut(w http.ResponseWriter, r *http.Request) (interface{}, *app.Error) {
	vo := &shortcutVO{}
	aerr := app.JsonUnmarshalFromRequest(r, vo)
	if aerr != nil {
		return nil, aerr
	}

	if aerr = validateParameter(vo); aerr != nil {
		return nil, aerr
	}

	// 1. create user if not registered
	u, err := user.GetByName(vo.User.Username)
	if err != nil {
		return nil, app.Errorf(err, "failed to get user")
	}
	if u == nil {
		u, err = user.Create(vo.User)
		if err != nil {
			return nil, app.Errorf(err, "failed to create user")
		}
	}

	// 2. create app for user if not created, else update it
	uapp, err := happ.GetByName(vo.App.AppName)
	if err != nil {
		return nil, app.Errorf(err, "failed to get user app")
	}
	if uapp == nil {
		uapp, err = happ.Create(vo.App)
		if err != nil {
			return nil, app.Errorf(err, "failed to create app")
		}
	}

	// 3. check dialer exists or not for user, create dialer when not
	udialer, err := dialer.GetByName(u.ID, vo.EndPoint.Dialer.Name)
	if err != nil {
		return nil, app.Errorf(err, "failed to get user dialer")
	}
	if udialer == nil {
		udialer, err = dialer.Create(vo.EndPoint.Dialer)
		if err != nil {
			return nil, app.Errorf(err, "failed to create dialer")
		}
	}

	// 4. create template for user
	// 5. create or update end point
	// 6. create or update preference for end point
	// 7. add or update receiver for end point

	return nil, nil
}

func validateParameter(vo *shortcutVO) *app.Error {
	const errInvalidParameter = "invalid parameter"
	if vo.User == nil || vo.EndPoint == nil || vo.App == nil {
		return app.Errorf(errors.New("find nil value when validate parameter"), errInvalidParameter)
	}
	if utils.IsStrBlank(vo.User.Username) || utils.IsStrBlank(vo.User.Password) {
		return app.Errorf(errors.New("username or password can not be empty"), errInvalidParameter)
	}
	if utils.IsStrBlank(vo.App.AppName) {
		return app.Errorf(errors.New("app name can not be empty"), errInvalidParameter)
	}
	if utils.IsStrBlank(vo.EndPoint.Name) {
		return app.Errorf(errors.New("end point name can not be empty"), errInvalidParameter)
	}
	if vo.EndPoint.Dialer != nil && utils.IsStrBlank(vo.EndPoint.Dialer.Name) {
		return app.Errorf(errors.New("dialer name can not be empty"), errInvalidParameter)
	}

	return nil
}
