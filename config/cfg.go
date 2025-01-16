package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"kursy-kriptovalyut/pkg/logger"
)

var log = logger.NewLogger()

type Config struct {
	Cfg struct {
		Port   string `mapstructure:"srv-port"`
		PgUser string `mapstructure:"pg-user"`
		PgPswd string `mapstructure:"pg-pswd"`
		PgDB   string `mapstructure:"pg-db"`
		PgHost string `mapstructure:"pg-host"`
		PgPort int    `mapstructure:"pg-port"`
		Url    string `mapstructure:"url"`
		ApiKey string `mapstructure:"api-key"`
	} `mapstructure:"cfg"`
}

func LoadCfg() *Config {
	cfg := &Config{}

	// viper.SetConfigFile("D:\\PROGRAMMING\\GO\\go-projects\\kursy-kriptovalyut\\config\\cfg.yaml")
	viper.SetConfigFile("./cfg.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Warn("(LoadCfg) failed to find cfg.yaml file:", zap.Any("err", err.Error()))
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		log.Warn("(LoadCfg) failed to load cfg.yaml file:", zap.Any("err", err.Error()))
	}

	return cfg
}
