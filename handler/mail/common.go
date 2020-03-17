package mail

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/gomail.v2"

	"GoMailer/common/db"
	"GoMailer/common/utils"
	"GoMailer/log"
)

func create(mail *db.Mail) (*db.Mail, error) {
	if utils.IsBlankStr(mail.Content) {
		return nil, errors.New("mail content can not be empty")
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	affected, err := client.InsertOne(mail)
	if err != nil {
		return nil, err
	}
	if affected != 1 {
		return nil, errors.New("failed to InsertOne mail")
	}

	return mail, nil
}

func handleMail(endpointId int64, val map[string]string) (*db.Mail, error) {
	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	ed := &db.Endpoint{}
	get, err := client.Id(endpointId).Get(ed)
	if err != nil {
		return nil, err
	}
	if !get {
		return nil, errors.New("endpoint not exist")
	}

	// Prepare template.
	t := &db.Template{}
	get, err = client.Id(ed.TemplateId).Get(t)
	if err != nil {
		return nil, err
	}
	template, contentType := getDefaultTemplate(val)
	if get {
		template, contentType = t.Template, t.ContentType
	}

	// Prepare receiver.
	rs := make([]db.Receiver, 0)
	err = client.Where("endpoint_id = ?", ed.Id).Find(&rs)
	if err != nil {
		return nil, err
	}

	// Prepare dialer.
	d := &db.Dialer{}
	_, err = client.Id(ed.DialerId).Get(d)
	if err != nil {
		return nil, err
	}

	edp := &db.EndpointPreference{
		DeliverStrategy: db.DeliverStrategy_DELIVER_IMMEDIATELY.Name(),
		EnableReCaptcha: 1,
	}
	_, err = client.Where("endpoint_id = ?", ed.Id).Get(edp)
	if err != nil {
		return nil, err
	}

	mail := &db.Mail{}
	mail.EndpointId = endpointId
	mail.Content = parseContent(template, val)
	mail.State = db.MailState_STAGING.Name()
	if edp.DeliverStrategy == db.DeliverStrategy_DELIVER_IMMEDIATELY.Name() {
		// Construct message.
		msg := gomail.NewMessage()
		msg.SetHeader("From", msg.FormatAddress(d.AuthUsername, d.Name))
		rsMap := make(map[string][]string)
		for _, r := range rs {
			rsMap[r.ReceiverType] = append(rsMap[r.ReceiverType], r.Address)
		}
		for t, e := range rsMap {
			msg.SetHeader(t, e...)
		}
		msg.SetHeader("Subject", ed.Name)
		msg.SetBody(contentType, mail.Content)

		// Send mail.
		der := gomail.NewDialer(d.Host, d.Port, d.AuthUsername, d.AuthPassword)
		log.Infof("deliver mail: %+v", msg)
		err = der.DialAndSend(msg)
		if err != nil {
			mail.State = db.MailState_DELIVER_FAILED.Name()
			return nil, err
		}
		mail.DeliveryTime = time.Now()
		mail.State = db.MailState_DELIVER_SUCCESS.Name()
	}

	return mail, nil
}

func parseContent(t string, val map[string]string) string {
	for key, value := range val {
		t = strings.ReplaceAll(t, fmt.Sprintf("{{%s}}", key), fmt.Sprintf("%v", value))
	}
	return t
}

func getDefaultTemplate(val map[string]string) (string, string) {
	builder := strings.Builder{}
	for key := range val {
		builder.WriteString(fmt.Sprintf("%s:  {{%s}}\n", key, key))
	}
	return builder.String(), "text/plain"
}
