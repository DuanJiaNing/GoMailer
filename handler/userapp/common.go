package userapp

import (
	"errors"

	"GoMailer/common/db"
	"GoMailer/common/utils"
)

func FindById(appId int64) (*db.UserApp, error) {
	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	app := &db.UserApp{}
	has, err := client.Id(appId).Get(app)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("app not exist")
	}

	return app, nil
}

func FindByName(userId int64, name string) (*db.UserApp, error) {
	if utils.IsBlankStr(name) {
		return nil, nil
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	u := &db.UserApp{}
	has, err := client.Where("name = ? AND user_id = ?", name, userId).Get(u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return u, nil
}

func Create(ua *db.UserApp) (*db.UserApp, error) {
	if utils.IsBlankStr(ua.Name) {
		return nil, errors.New("name can not be empty")
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	userExist, err := client.ID(ua.UserId).Exist(&db.User{})
	if err != nil {
		return nil, err
	}
	if !userExist {
		return nil, errors.New("user not exist")
	}

	affected, err := client.InsertOne(ua)
	if err != nil {
		return nil, err
	}
	if affected != 1 {
		return nil, errors.New("failed to InsertOne user app")
	}

	return ua, nil
}

func Update(ua *db.UserApp) (*db.UserApp, error) {
	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	if ua.UserId != 0 {
		userExist, err := client.ID(ua.UserId).Exist(&db.User{})
		if err != nil {
			return nil, err
		}
		if !userExist {
			return nil, errors.New("user not exist")
		}
	}

	_, err = client.ID(ua.Id).Update(ua)
	if err != nil {
		return nil, err
	}
	return ua, nil
}
