package config

type AppConfiguration struct {
	Spec specConfiguration
}

type specConfiguration struct {
	Websocket      string   `yaml:"websocket"`
	Products       []string `yaml:"products"`
	VwapWindowSize int      `yaml:"vwapWindowSize"`
}
