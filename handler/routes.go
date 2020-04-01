package handler

import (
	"github.com/gorilla/mux"

	"GoMailer/middleware/auth"
	"GoMailer/middleware/cors"
)

var (
	Router    = mux.NewRouter()
	APIRouter = Router.PathPrefix("/api/").Subrouter()

	ShortcutRouter = APIRouter.Path("/shortcut").Subrouter()

	// Root properties
	// /api/*
	MailRouter = APIRouter.PathPrefix("/mail/").Subrouter()
	UserRouter = APIRouter.PathPrefix("/user/").Subrouter()

	// User properties.
	// /api/user/{UID}/*
	DialerRouter   = UserRouter.PathPrefix("/{UID}/dialer").Subrouter()
	TemplateRouter = UserRouter.PathPrefix("/{UID}/template").Subrouter()
	AppRouter      = UserRouter.PathPrefix("/{UID}/app").Subrouter()

	// User app properties.
	// /api/user/{UID}/app/{AID}/*
	EndPointRouter = AppRouter.PathPrefix("/{AID}/endpoint").Subrouter()

	// End point properties.
	// /api/user/{UID}/app/{AID}/endpoint/{EID}/*
	PreferenceRouter = EndPointRouter.PathPrefix("/{EID}/preference").Subrouter()
	ReceiverRouter   = EndPointRouter.PathPrefix("/{EID}/receiver").Subrouter()
)

func init() {
	var middleware []mux.MiddlewareFunc
	middleware = append(middleware, auth.Guard)
	middleware = append(middleware, cors.CORS(APIRouter))
	APIRouter.Use(middleware...)
}
