package main

import (
	"log"

	"github.com/prabhatlabs/internals/lib"
	"github.com/prabhatlabs/internals/services/mail"
)

func main() {
	err := lib.LoadEnv()
	if err != nil {
		log.Fatalln("Env load error: ", err)
	}

	recipient := "workforprabhat1254@gmail.com"
	mail.SendMail(recipient, "Test Mail", "This is the body for the test mail, hope this works! <a href='https://x.com/prabhatlabs'>click here!</a>")
}
