package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"GoMailer/conf"
)

func IsDevAppServer() bool {
	return conf.Env() == "dev"
}

func JsonUnmarshalFromRequest(r *http.Request, dst interface{}) *Error {
	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return Errorf(err, "failed to parse request body")
	}

	err = json.Unmarshal(all, dst)
	if err != nil {
		return Errorf(err, "failed to parse request body")
	}

	return nil
}
