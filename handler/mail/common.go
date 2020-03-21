package mail

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/gomail.v2"

	"GoMailer/common/db"
	"GoMailer/common/utils"
	"GoMailer/log"
)

func create(userId int64, mail *db.Mail) (*db.Mail, error) {
	if utils.IsBlankStr(mail.Content) {
		return nil, errors.New("mail content can not be empty")
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	mt := getUserMailTableName(userId)
	res, err := client.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	mtExist := false
	for _, r := range res {
		if string(r["Tables_in_gomailer"]) == mt {
			mtExist = true
		}
	}
	if !mtExist {
		_, err = client.Exec(buildSql(mt))
		if err != nil {
			return nil, err
		}
	}

	affected, err := client.Table(mt).InsertOne(mail)
	if err != nil {
		return nil, err
	}
	if affected != 1 {
		return nil, errors.New("failed to InsertOne mail")
	}

	return mail, nil
}

func buildSql(tableName string) string {
	sql := strings.Builder{}
	sql.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (", tableName))
	sql.WriteString("  `id` int unsigned NOT NULL AUTO_INCREMENT,")
	sql.WriteString("  `insert_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,")
	sql.WriteString("  `endpoint_id` int NOT NULL,")
	sql.WriteString("  `state` varchar(100) NOT NULL,")
	sql.WriteString("  `delivery_time` timestamp NULL DEFAULT NULL,")
	sql.WriteString("  `content` longtext NOT NULL,")
	sql.WriteString("  `raw` longtext NOT NULL,")
	sql.WriteString("  PRIMARY KEY (`id`)")
	sql.WriteString(") ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	return sql.String()
}

func getUserMailTableName(userId int64) string {
	return fmt.Sprintf("mail_%d", userId)
}

func handleMail(endpointId int64, raw map[string]string) (*db.Mail, error) {
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
			EnableReCaptcha: 1, // TODO
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
		mail.DeliveryTime = time.Now()
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

