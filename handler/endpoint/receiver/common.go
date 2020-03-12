package receiver

import (
	"errors"

	"GoMailer/common/db"
	"GoMailer/common/utils"
)

func DeleteByEndpoint(endpointId int64) error {
	client, err := db.NewClient()
	if err != nil {
		return err
	}

	_, err = client.Where("endpoint_id = ?", endpointId).Delete(&db.Receiver{})
	if err != nil {
		return err
	}
	return nil
}

func PatchCreate(receiver []*db.Receiver) error {
	client, err := db.NewClient()
	if err != nil {
		return err
	}

	for _, r := range receiver {
		userExist, err := client.ID(r.UserId).Exist(&db.User{})
		if err != nil {
			return err
		}
		if !userExist {
			return errors.New("user not exist")
		}
		userAppExist, err := client.ID(r.UserAppId).Exist(&db.UserApp{})
		if err != nil {
			return err
		}
		if !userAppExist {
			return errors.New("user app not exist")
		}
		endpointExist, err := client.ID(r.EndpointId).Exist(&db.Endpoint{})
		if err != nil {
			return err
		}
		if !endpointExist {
			return errors.New("endpoint not exist")
		}
		if utils.IsBlankStr(r.Address) {
			return errors.New("address can not be empty")
		}
		if db.ReceiverType(r.ReceiverType) == "" {
			return errors.New("receiver type illegal")
		}
	}

	_, err = client.Insert(receiver)
	if err != nil {
		return err
	}
	return nil
}
