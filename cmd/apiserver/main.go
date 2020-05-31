package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/marinira/viber-e-shop/internal/app/apiserver"
	"log"
)

// переменная для указания путь к конфигурации в командной строки
var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config_path", "configs/apiserver.toml", "path to config file")

}
func main() {
	//парсиг
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	s := apiserver.New(config)
	err = s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
