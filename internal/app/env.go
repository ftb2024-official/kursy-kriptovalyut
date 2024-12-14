package app

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	PgUser  string `mapstructure:"POSTGRES_USER"`
	PgPswd  string `mapstructure:"POSTGRES_PASSWORD"`
	PgDB    string `mapstructure:"POSTGRES_DB"`
	PgHost  string `mapstructure:"POSTGRES_HOST"`
	PgPort  string `mapstructure:"POSTGRES_PORT"`
	SrvAddr string `mapstructure:"SRV_ADDR"`
	SrvPort string `mapstructure:"SRV_PORT"`
	Url     string `mapstructure:"URL"`
	ApiKey  string `mapstructure:"API_KEY"`
}

func NewEnv() *Env {
	env := &Env{}
	viper.SetConfigFile("./../../.env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("failed to find .env file:", err)
	}

	err = viper.Unmarshal(env)
	if err != nil {
		log.Fatal("failed to load .env file:", err)
	}

	return env
}
