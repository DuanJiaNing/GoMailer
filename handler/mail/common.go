package mail

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"GoMailer/common/db"
	"GoMailer/common/utils"
)

func Find(userId int64, pageNum int, pageSize int) ([]*userMail, int64, error) {
	mt, err := handleUserMailTable(userId)
	if err != nil {
		return nil, 0, err
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, 0, err
	}
	total, err := client.Table(mt).Count(&db.Mail{})
	if err != nil {
		return nil, 0, err
	}

	var ms []*db.Mail
	err = client.Table(mt).Limit(pageSize, (pageNum-1)*pageSize).Desc("insert_time").Find(&ms)
	if err != nil {
		return nil, 0, err
	}
	ums := make([]*userMail, 0, len(ms))
	for _, m := range ms {
		raw := make(map[string]string)
		err := json.Unmarshal([]byte(m.Raw), &raw)
		if err != nil {
			return nil, 0, err
		}
		ums = append(ums, &userMail{
			InsertTime:   m.InsertTime,
			State:        m.State,
			DeliveryTime: m.DeliveryTime,
			Content:      m.Content,
			Raw:          raw,
		})
	}

	return ums, total, nil
}

func Create(userId int64, mail *db.Mail) (*db.Mail, error) {
	if utils.IsBlankStr(mail.Content) {
		return nil, errors.New("mail content can not be empty")
	}

	mt, err := handleUserMailTable(userId)
	if err != nil {
		return nil, err
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
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

func handleUserMailTable(userId int64) (string, error) {
	client, err := db.NewClient()
	if err != nil {
		return "", err
	}

	mt := getUserMailTableName(userId)
	res, err := client.Query("SHOW TABLES")
	if err != nil {
		return "", err
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
			return "", err
		}
	}

	return mt, nil
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
