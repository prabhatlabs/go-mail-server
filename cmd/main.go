package main

import (
	"log"

	"github.com/prabhatlabs/internals/lib/env"
	"github.com/prabhatlabs/internals/lib/mail"
)

func main() {
	err := env.LoadEnv()
	if err != nil {
		log.Fatalln("Env load error: ", err)
	}

	recipient := "prabhat@gluckglobal.com"
	mail.SendMail(recipient, "Test Mail", "This is the body for the test mail, hope this works! <a href='https://x.com/prabhatlabs'>click here!</a>")
}
