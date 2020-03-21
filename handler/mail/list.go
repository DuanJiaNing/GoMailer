package mail

import (
	"GoMailer/app"
	"GoMailer/handler"
	"net/http"
	"strconv"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 10
)

func init() {
	router := handler.MailRouter.Path("/list").Subrouter()
	router.Methods(http.MethodGet).Handler(app.Handler(list))
}

func list(w http.ResponseWriter, r *http.Request) (interface{}, *app.Error) {
	// TODO
	/*	uid := r.URL.Query().Get("uid")
		userId, err := strconv.Atoi(uid)
		if err != nil {
			return nil, app.Errorf(errors.New("uid is not a number"), "uid illegal")
		}
		pageNum, pageSize := parsePageCondition(r)

		client, err := db.NewClient()
		if err != nil {
			return nil, app.Errorf(err,"got error when list post")
		}
		client.Where("endpoint_id = ?", ep.Id).FindAndCount()
	*/
	return nil, nil
}

func parsePageCondition(r *http.Request) (int, int) {
	pn := r.URL.Query().Get("pn")
	ps := r.URL.Query().Get("ps")

	pni, err := strconv.Atoi(pn)
	if err != nil {
		pni = defaultPageNum
	}
	psi, err := strconv.Atoi(ps)
	if err != nil {
		pni = defaultPageSize
	}
	return pni, psi
}
