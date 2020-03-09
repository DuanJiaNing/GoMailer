package endpoint

import (
	"GoMailer/common/db"
)

func FindByName(name string) (*db.EndPoint, error) {
	return nil, nil
}

func Create(ep *db.EndPoint) (*db.EndPoint, error) {
	return nil, nil
}

func Update(ep *db.EndPoint) (*db.EndPoint, error) {
	return nil, nil
}

func FindPreference(endpointID int64) (*db.EndPointPreference, error) {
	return nil, nil
}

func CreatePreference(p *db.EndPointPreference) (*db.EndPointPreference, error) {
	return nil, nil
}

func UpdatePreference(p *db.EndPointPreference) (*db.EndPointPreference, error) {
	return nil, nil
}

func DeleteReceiver(endpointID int64) error {
	return nil
}

func PatchCreateReceiver(receiver []*db.Receiver) error {
	return nil
}
