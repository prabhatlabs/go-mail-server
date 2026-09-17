package mail

import (
	"fmt"
	"net/smtp"

	"github.com/prabhatlabs/go-mail-server/internal/lib/env"
)

func getAuth() smtp.Auth {
	return smtp.PlainAuth("", env.Vars.EMAIL_USER, env.Vars.EMAIL_PASS, env.Vars.EMAIL_HOST)
}

func getServerAddr() string {
	return fmt.Sprintf("%s:%s", env.Vars.EMAIL_HOST, env.Vars.EMAIL_PORT)
}
