package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Database *Database
	Server   *Server
}

type Database struct {
	Driver string
	Url    string
}

type Server struct {
	Cache    bool
	CacheDir string
	Debug    bool
	Port     int
}

var Configiure *Config = new(Config)

func init() {
	//读取配置文件
	_, err := toml.DecodeFile("config.toml", Configiure)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}
