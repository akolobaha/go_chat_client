package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Auth struct {
	Token     string `toml:"token"`
	ServerUrl string `toml:"server_url"`
}

var globalAuth *Auth

func Parse(s string) {
	c := &Auth{}
	if err := cleanenv.ReadConfig(s, c); err != nil {
		panic("Configuration file parsing error: " + err.Error())
	}

	globalAuth = c
}

func Get() *Auth {
	return globalAuth
}
