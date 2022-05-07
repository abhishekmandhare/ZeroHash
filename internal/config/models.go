package config

type AppConfiguration struct {
	Spec specConfiguration
}

type specConfiguration struct {
	Websocket string `yaml:"websocket"`
}