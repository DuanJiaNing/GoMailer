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
	MailRouter = APIRouter.PathPrefix("/mail/").Subrouter()
)

func init() {
	var middleware []mux.MiddlewareFunc
	middleware = append(middleware, cors.CORS(APIRouter))
	middleware = append(middleware, auth.Guard)
	APIRouter.Use(middleware...)
}
