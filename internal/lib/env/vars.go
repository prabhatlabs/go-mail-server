package env

import (
	"fmt"

	"github.com/spf13/viper"
)

type EnvsType struct {
	ENV          string
	DATABASE_URL string
	EMAIL_HOST   string
	EMAIL_PORT   string
	EMAIL_USER   string
	EMAIL_PASS   string
}

var Vars *EnvsType

func LoadEnv() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("prod")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	requiredKeys := []string{
		"ENV",
		"DATABASE_URL",
		"EMAIL_HOST",
		"EMAIL_PORT",
		"EMAIL_USER",
		"EMAIL_PASS",
	}

	var enverr []string

	// checking for missing envs
	for _, key := range requiredKeys {
		if !viper.IsSet(key) {
			enverr = append(enverr, fmt.Sprintf("%s is not set, ", key))
		}
	}

	if len(enverr) > 0 {
		return fmt.Errorf("missing environment variables: %v", enverr)
	}

	envs := &EnvsType{
		ENV:          viper.GetString("ENV"),
		DATABASE_URL: viper.GetString("DATABASE_URL"),
		EMAIL_HOST:   viper.GetString("EMAIL_HOST"),
		EMAIL_PORT:   viper.GetString("EMAIL_PORT"),
		EMAIL_USER:   viper.GetString("EMAIL_USER"),
		EMAIL_PASS:   viper.GetString("EMAIL_PASS"),
	}

	if envs.ENV == "" {
		envs.ENV = "dev"
	}

	Vars = Vars
	return nil
}

func IsProd() bool {
	return Vars != nil && Vars.ENV == "prod"
}
