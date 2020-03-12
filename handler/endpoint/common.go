package endpoint

import (
	"errors"

	"GoMailer/common/db"
	"GoMailer/common/utils"
)

func FindByName(appId int64, name string) (*db.Endpoint, error) {
	if utils.IsStrBlank(name) {
		return nil, nil
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	u := &db.Endpoint{}
	has, err := client.Where("name = ? AND user_app_id = ?", name, appId).Get(u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return u, nil
}

func Create(ep *db.Endpoint) (*db.Endpoint, error) {
	if utils.IsStrBlank(ep.Name) {
		return nil, errors.New("name can not be empty")
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	userExist, err := client.ID(ep.UserId).Exist(&db.User{})
	if err != nil {
		return nil, err
	}
	if !userExist {
		return nil, errors.New("user not exist")
	}
	userAppExist, err := client.ID(ep.UserAppId).Exist(&db.UserApp{})
	if err != nil {
		return nil, err
	}
	if !userAppExist {
		return nil, errors.New("user app not exist")
	}
	dialerExist, err := client.ID(ep.DialerId).Exist(&db.Dialer{})
	if err != nil {
		return nil, err
	}
	if !dialerExist {
		return nil, errors.New("dialer not exist")
	}
	// Nullable.
	if ep.TemplateId != 0 {
		templateExist, err := client.ID(ep.TemplateId).Exist(&db.Template{})
		if err != nil {
			return nil, err
		}
		if !templateExist {
			return nil, errors.New("template not exist")
		}
	}

	affected, err := client.InsertOne(ep)
	if err != nil {
		return nil, err
	}
	if affected != 1 {
		return nil, errors.New("failed to InsertOne endpoint")
	}

	return ep, nil
}

func Update(ep *db.Endpoint) (*db.Endpoint, error) {
	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	if ep.UserId != 0 {
		userExist, err := client.ID(ep.UserId).Exist(&db.User{})
		if err != nil {
			return nil, err
		}
		if !userExist {
			return nil, errors.New("user not exist")
		}
	}
	if ep.UserAppId != 0 {
		userAppExist, err := client.ID(ep.UserAppId).Exist(&db.UserApp{})
		if err != nil {
			return nil, err
		}
		if !userAppExist {
			return nil, errors.New("user app not exist")
		}
	}
	if ep.DialerId != 0 {
		dialerExist, err := client.ID(ep.DialerId).Exist(&db.Dialer{})
		if err != nil {
			return nil, err
		}
		if !dialerExist {
			return nil, errors.New("dialer not exist")
		}
	}
	if ep.TemplateId != 0 {
		templateExist, err := client.ID(ep.TemplateId).Exist(&db.Template{})
		if err != nil {
			return nil, err
		}
		if !templateExist {
			return nil, errors.New("template not exist")
		}
	}

	_, err = client.ID(ep.Id).Update(ep)
	if err != nil {
		return nil, err
	}
	return ep, nil
}
