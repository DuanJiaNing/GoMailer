package mail

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"GoMailer/app"
	"GoMailer/common/db"
	"GoMailer/handler"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 10
)

type pageResult struct {
	PageNum  int
	PageSize int
	Total    int64
	List     []*userMail
}

type userMail struct {
	InsertTime   time.Time
	State        string
	DeliveryTime db.Time
	Content      string
	Raw          map[string]string
}

func init() {
	router := handler.MailRouter.Path("/list").Subrouter()
	router.Methods(http.MethodGet).Handler(app.Handler(list))
}

func list(w http.ResponseWriter, r *http.Request) (interface{}, *app.Error) {
	uid := r.URL.Query().Get("uid")
	userId, err := strconv.Atoi(uid)
	if err != nil {
		return nil, app.Errorf(errors.New("uid is not a number"), "uid illegal")
	}
	pageNum, pageSize := parsePageCondition(r)

	res, total, err := Find(int64(userId), pageNum, pageSize)
	if err != nil {
		return nil, app.Errorf(err, "failed to find user post")
	}

	return &pageResult{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    total,
		List:     res,
	}, nil
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
		psi = defaultPageSize
	}
	return pni, psi
}
