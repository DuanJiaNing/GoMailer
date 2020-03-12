package user

import (
	"errors"

	"GoMailer/common/db"
	"GoMailer/common/utils"
)

func FindByName(username string) (*db.User, error) {
	if utils.IsBlankStr(username) {
		return nil, nil
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	u := &db.User{}
	has, err := client.Where("username = ?", username).Get(u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return u, nil
}

func Create(user *db.User) (*db.User, error) {
	if utils.IsBlankStr(user.Username) {
		return nil, errors.New("username can not be empty")
	}
	if utils.IsBlankStr(user.Password) || len(user.Password) < 6 {
		return nil, errors.New("password invalid, min length is 6")
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	one, err := client.InsertOne(user)
	if err != nil {
		return nil, err
	}
	if one != 1 {
		return nil, errors.New("failed to InsertOne user")
	}

	return user, nil
}
