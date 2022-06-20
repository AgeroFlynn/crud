// Package config provides config struct for this application and
package config

import (
	"errors"
	"fmt"
	"github.com/AgeroFlynn/crud/foundation/yaml"
	"github.com/ardanlabs/conf/v3"
	"os"
	"time"
)

type Config struct {
	conf.Version
	Web struct {
		ReadTimeout     time.Duration `conf:"default:5s" yaml:"readTimeout"`
		WriteTimeout    time.Duration `conf:"default:10s" yaml:"writeTimeout"`
		IdleTimeout     time.Duration `conf:"default:120s" yaml:"idleTimeout"`
		ShutdownTimeout time.Duration `conf:"default:20s" yaml:"shutdownTimeout"`
		APIHost         string        `conf:"default:0.0.0.0:3000" yaml:"APIHost"`
		DebugHost       string        `conf:"default:0.0.0.0:4000" yaml:"debugHost"`
	}
	Auth struct {
		KeysFolder string `conf:"default:resources/keys/" yaml:"keysFolder"`
		ActiveKID  string `conf:"default:54bb2165-71e1-41a6-af3e-7da4a0e1e2c1" yaml:"activeKID"`
	}
	DB struct {
		User        string `conf:"default:postgres"`
		Password    string `conf:"default:postgres,mask"`
		Host        string `conf:"default:localhost"`
		Name        string `conf:"default:postgres"`
		MaxIdleCons int    `conf:"default:0" yaml:"maxIdleCons"`
		MaxOpenCons int    `conf:"default:0" yaml:"maxOpenCons"`
		DisableTLS  bool   `conf:"default:true"`
	}
}

func NewConfigFromFile() error {
	data, err := os.ReadFile("resources/config/configuration.yaml")
	if err != nil {
		return fmt.Errorf("error with reading configuration file: %w", err)
	}

	const prefix = "APP"
	var cfg = Config{
		Version: conf.Version{
			Build: "build",
			Desc:  "copyright information here",
		},
	}

	help, err := conf.Parse(prefix, &cfg, yaml.WithData(data))
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	fmt.Println(conf.String(&cfg))
	return nil
}
