package dialer

import (
	"errors"

	"GoMailer/common/db"
	"GoMailer/common/utils"
)

func FindByName(userId int64, name string) (*db.Dialer, error) {
	if utils.IsBlankStr(name) {
		return nil, nil
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	u := &db.Dialer{}
	has, err := client.Where("name = ? AND user_id = ?", name, userId).Get(u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return u, nil
}

func Create(d *db.Dialer) (*db.Dialer, error) {
	if utils.IsBlankStr(d.Name) {
		return nil, errors.New("name can not be empty")
	}
	if utils.IsBlankStr(d.Host) {
		return nil, errors.New("host can not be empty")
	}
	if utils.IsBlankStr(d.AuthPassword) {
		return nil, errors.New("auth password can not be empty")
	}
	if utils.IsBlankStr(d.AuthUsername) {
		return nil, errors.New("auth username can not be empty")
	}
	if d.Port <= 0 {
		return nil, errors.New("port invalid")
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	userExist, err := client.ID(d.UserId).Exist(&db.User{})
	if err != nil {
		return nil, err
	}
	if !userExist {
		return nil, errors.New("user not exist")
	}

	affected, err := client.InsertOne(d)
	if err != nil {
		return nil, err
	}
	if affected != 1 {
		return nil, errors.New("failed to InsertOne dialer")
	}

	return d, nil
}

func Update(d *db.Dialer) (*db.Dialer, error) {
	if d.Port < 0 {
		return nil, errors.New("port invalid")
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	if d.UserId != 0 {
		userExist, err := client.ID(d.UserId).Exist(&db.User{})
		if err != nil {
			return nil, err
		}
		if !userExist {
			return nil, errors.New("user not exist")
		}
	}

	_, err = client.ID(d.Id).Update(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}
