package template

import (
	"errors"

	"GoMailer/common/db"
	"GoMailer/common/utils"
)

func Create(t *db.Template) (*db.Template, error) {
	if utils.IsBlankStr(t.Template) {
		return nil, errors.New("template can not be empty")
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	userExist, err := client.ID(t.UserId).Exist(&db.User{})
	if err != nil {
		return nil, err
	}
	if !userExist {
		return nil, errors.New("user not exist")
	}

	affected, err := client.InsertOne(t)
	if err != nil {
		return nil, err
	}
	if affected != 1 {
		return nil, errors.New("failed to InsertOne user app")
	}

	return t, nil
}
