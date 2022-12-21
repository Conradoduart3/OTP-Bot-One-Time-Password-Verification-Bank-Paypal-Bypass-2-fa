package config

import (
	"log"

	"github.com/alexflint/go-arg"
)

var Args struct {
	LOG_LEVEL          string `arg:"required,env"`
	TG_BOT_KEY         string `arg:"required,env"`
	API_KEY            string `arg:"required,env"`
	BACKEND_URL        string `arg:"required,env"`
	AUTH_URL           string
	ORIGIN_URL         string
	NEW_ORIGIN_URL     string
	NEW_USER_URL       string
	USERS_URL          string
	ADD_PERMISSION_URL string
}

func Validate() {
	if err := arg.Parse(&Args); err != nil {
		log.Fatal(err)
	}

	Args.AUTH_URL = Args.BACKEND_URL + "/api/check-admin"
	Args.ORIGIN_URL = Args.BACKEND_URL + "/api/origins"
	Args.NEW_ORIGIN_URL = Args.BACKEND_URL + "/api/origin"
	Args.NEW_USER_URL = Args.BACKEND_URL + "/api/user"
	Args.USERS_URL = Args.BACKEND_URL + "/api/users"
	Args.ADD_PERMISSION_URL = Args.BACKEND_URL + "/api/permission"
}
