package config

import (
	"github.com/namsral/flag"
)

type Config struct {
	TemplateDirectory string
	ServerPort        int
}

func Init() *Config {

	cfg := &Config{}

	flag.StringVar(&cfg.TemplateDirectory, "template", "./route/", "Server port")
	flag.IntVar(&cfg.ServerPort, "port", 8080, "Server port")
	flag.Parse()

	return cfg

}
