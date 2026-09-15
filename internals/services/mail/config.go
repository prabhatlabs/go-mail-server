package mail

import (
	"fmt"
	"net/smtp"

	"github.com/prabhatlabs/internals/lib"
)

func getAuth() smtp.Auth {
	return smtp.PlainAuth("", lib.Envs.EMAIL_USER, lib.Envs.EMAIL_PASS, lib.Envs.EMAIL_HOST)
}

func getServerAddr() string {
	return fmt.Sprintf("%s:%s", lib.Envs.EMAIL_HOST, lib.Envs.EMAIL_PORT)
}
