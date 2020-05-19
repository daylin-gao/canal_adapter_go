package config

import (
"os"
"path/filepath"
"github.com/jinzhu/configor"
)

var Config = struct {
	HTTPS    bool   `default:"false" env:"HTTPS"`
	Host     string `default:"" env:"Host"`
	Certpath string `default:"" env:"Certpath"`
	Certkey  string `default:"" env:"Certkey"`
	Port     uint   `default:"80" env:"PORT"`
	DB       struct {
		Name     string `env:"DBName" default:""`
		Adapter  string `env:"DBAdapter" default:""`
		Host     string `env:"DBHost" default:""`
		Port     string `env:"DBPort" default:""`
		User     string `env:"DBUser" default:""`
		Password string `env:"DBPassword" default:""`
	}
}{}

var Root = os.Getenv("GOPATH") + "/src/canal_adapter_go"

func init() {
	if err := configor.Load(&Config, filepath.Join(Root, "config/"+"application.yml")); err != nil {
		panic(err)
	}
}
