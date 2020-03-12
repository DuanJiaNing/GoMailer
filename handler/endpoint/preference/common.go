package preference

import (
	"errors"

	"GoMailer/common/db"
	"GoMailer/common/utils"
)

func FindByEndpoint(endpointId int64) (*db.EndpointPreference, error) {
	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	u := &db.EndpointPreference{}
	has, err := client.Where("endpoint_id = ?", endpointId).Get(u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return u, nil
}

func Create(p *db.EndpointPreference) (*db.EndpointPreference, error) {
	if !utils.IsBlankStr(p.DeliverStrategy) && db.DeliverStrategy(p.DeliverStrategy) == "" {
		return nil, errors.New("deliver strategy illegal")
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	endpointExist, err := client.ID(p.EndpointId).Exist(&db.Endpoint{})
	if err != nil {
		return nil, err
	}
	if !endpointExist {
		return nil, errors.New("endpoint not exist")
	}

	affected, err := client.InsertOne(p)
	if err != nil {
		return nil, err
	}
	if affected != 1 {
		return nil, errors.New("failed to InsertOne endpoint")
	}

	return p, nil
}

func Update(p *db.EndpointPreference) (*db.EndpointPreference, error) {
	if !utils.IsBlankStr(p.DeliverStrategy) && db.DeliverStrategy(p.DeliverStrategy) == "" {
		return nil, errors.New("deliver strategy illegal")
	}

	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	_, err = client.Where("endpoint_id = ?", p.EndpointId).Update(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
