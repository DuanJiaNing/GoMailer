package app

import "GoMailer/conf"

func IsDevAppServer() bool {
	return conf.Env() == "dev"
}
