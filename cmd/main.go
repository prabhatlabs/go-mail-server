package main

import (
	"github.com/prabhatlabs/internals/lib"
	"github.com/prabhatlabs/internals/services/mail"
)

func main() {
	lib.LoadEnv()

	recipient := "workforprabhat1254@gmail.com"
	mail.SendMail(recipient, "Test Mail", "This is the body for the test mail, hope this works! <a href='https://x.com/prabhatlabs'>click here!</a>")
}
