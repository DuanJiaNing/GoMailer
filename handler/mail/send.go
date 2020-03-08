package mail

import (
	"gopkg.in/gomail.v2"
	"net/http"

	"GoMailer/app"
	"GoMailer/handler"
)

func init() {
	router := handler.MailRouter.Path("/send").Subrouter()
	router.Methods(http.MethodPost).Handler(app.Handler(send))
}

func send(http.ResponseWriter, *http.Request) (interface{}, *app.Error) {
	dialer := gomail.NewDialer("smtp.qq.com", 465, "2213994603@qq.com", "athupcbmeyvvdjif")
	err := dialer.DialAndSend(getMsg())
	if err != nil {
		return nil, app.Errorf(err, "failed to send email")
	}

	return nil, nil
}

func getMsg() *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", msg.FormatAddress("2213994603@qq.com", "djnqq"))
	msg.SetHeader("To", "djn163<duan_jia_ning@163.com>")
	msg.SetHeader("Subject", "This is test mail")
	html := `
<div>
<hr>
<h1>H1 text</h1>
<h2>H2 text</h2>
<hr>
</div>
`
	msg.SetBody("text/html", html)
	return msg
}
