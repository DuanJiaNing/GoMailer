package endpoint

import (
	"errors"

	"GoMailer/common/db"
	"GoMailer/common/key"
	"GoMailer/common/utils"
)

func FindByKey(key string) (*db.Endpoint, error) {
	if utils.IsBlankStr(key) {
		return nil, nil
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	u := &db.Endpoint{}
	has, err := client.Where("key = ?", key).Get(u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return u, nil
}

func RefreshKey(endpointId int64) (string, error) {
	client, err := db.NewClient()
	if err != nil {
		return "", err
	}

	ud := &db.Endpoint{Id: endpointId, Key: key.GenerateKey()}
	affected, err := client.Id(endpointId).Update(ud)
	if err != nil {
		return "", err
	}
	if affected != 1 {
		return "", errors.New("failed to Update endpoint key")
	}

	return ud.Key, nil
}
