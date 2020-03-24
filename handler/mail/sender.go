package mail

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/gomail.v2"

	"GoMailer/common/db"
	"GoMailer/common/key"
	"GoMailer/log"
)

func handleMail(endpointId int64, raw map[string]string, reCaptchaKey string) (*db.Mail, error) {
	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	edp := &db.EndpointPreference{}
	get, err := client.Where("endpoint_id = ?", endpointId).Get(edp)
	if err != nil {
		return nil, err
	}
	if !get {
		edp = &db.EndpointPreference{
			DeliverStrategy: db.DeliverStrategy_DELIVER_IMMEDIATELY.Name(),
			EnableReCaptcha: 2, // Default is disable reCaptcha.
		}
	}

	if edp.EnableReCaptcha == 1 {
		ok, err := key.VerifyReCaptcha(reCaptchaKey)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New("reCaptcha verify failed")
		}
	}

	message, content, err := prepareMessage(endpointId, raw)
	if err != nil {
		return nil, err
	}

	mail := &db.Mail{}
	mail.Content = content
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	mail.Raw = string(bytes)
	mail.EndpointId = endpointId
	mail.State = db.MailState_STAGING.Name()
	if edp.DeliverStrategy == db.DeliverStrategy_DELIVER_IMMEDIATELY.Name() {
		dialer, err := getMailDialer(endpointId)
		if err != nil {
			return nil, err
		}
		log.Infof("deliver mail: %+v", message)

		err = dialer.DialAndSend(message)
		mail.DeliveryTime = db.Time(time.Now())
		if err != nil {
			mail.State = db.MailState_DELIVER_FAILED.Name()
			return nil, err
		}
		mail.State = db.MailState_DELIVER_SUCCESS.Name()
	}

	return mail, nil
}

func prepareMessage(endpointId int64, val map[string]string) (*gomail.Message, string, error) {
	client, err := db.NewClient()
	if err != nil {
		return nil, "", err
	}

	ed := &db.Endpoint{}
	get, err := client.Id(endpointId).Get(ed)
	if err != nil {
		return nil, "", err
	}
	if !get {
		return nil, "", errors.New("endpoint not exist")
	}

	d := &db.Dialer{}
	_, err = client.Id(ed.DialerId).Get(d)
	if err != nil {
		return nil, "", err
	}

	// Prepare receiver.
	rs := make([]db.Receiver, 0)
	err = client.Where("endpoint_id = ?", ed.Id).Find(&rs)
	if err != nil {
		return nil, "", err
	}

	// Prepare template.
	t := &db.Template{}
	get, err = client.Id(ed.TemplateId).Get(t)
	if err != nil {
		return nil, "", err
	}
	template, contentType := getDefaultTemplate(val)
	if get {
		template, contentType = t.Template, t.ContentType
	}

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
	content := parseContent(template, val)
	msg.SetBody(contentType, content)

	return msg, content, nil
}

func getMailDialer(endpointId int64) (*gomail.Dialer, error) {
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

	d := &db.Dialer{}
	_, err = client.Id(ed.DialerId).Get(d)
	if err != nil {
		return nil, err
	}

	return gomail.NewDialer(d.Host, d.Port, d.AuthUsername, d.AuthPassword), nil
}

func parseContent(t string, val map[string]string) string {
	for key, value := range val {
		t = strings.ReplaceAll(t, fmt.Sprintf("{{%s}}", key), fmt.Sprintf("%v", value))
	}
	t = strings.ReplaceAll(t, "}}", "\"")
	t = strings.ReplaceAll(t, "{{", "\"")
	return t
}

func getDefaultTemplate(val map[string]string) (string, string) {
	builder := strings.Builder{}
	for key := range val {
		builder.WriteString(fmt.Sprintf("%s:  {{%s}}\n", key, key))
	}
	return builder.String(), "text/plain"
}
