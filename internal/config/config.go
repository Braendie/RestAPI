package config

import (
	"sync"

	"github.com/Braendie/RestAPI/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

// Config is a structure, that gets data from config.yml and uses it in this server.
type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
}

var instance *Config
var once sync.Once

// GetConfig return ready instance.
// if function starts first time, it configurates instance for using and logs it.
func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})

	return instance
}