package conf

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	server Server `yaml: "server"	mapstructure: "server"`
	tls    bool   `yaml:	"tls"	mapstructure: "tls"`
}

type Server struct {
	address string `yaml:	"address"	mapstructure: "address"`
	port    string `yaml:	"port"	mapstructure: "port"`
}

func (config *Config) GetServerAddress() string {

	return config.server.address + ":" + config.server.port
}

func (config *Config) GetTLS() bool {
	return config.tls
}

// LoadConfigPointer reads configuration from yaml file or environment variables.
func LoadConfigPointer(path string) (config *Config, err error) {

	if path == "" {
		//获取项目的执行路径
		path, err = os.Getwd()
		if err != nil {
			panic(err)
		}
		path += "/conf"
	}

	var cfg Config

	log.Println("config path:", path)
	viper.AddConfigPath(path)
	viper.SetConfigName("client")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)
	log.Println("config:", cfg)
	return &cfg, err
}

// LoadConfigPointer reads configuration from yaml file or environment variables.
func LoadConfig(path string) (config Config, err error) {

	if path == "" {
		//获取项目的执行路径
		path, err = os.Getwd()
		if err != nil {
			panic(err)
		}
		path += "/conf"
	}

	viper.AddConfigPath(path)
	viper.SetConfigName("client")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	//log.Println("config:", config)
	return
}
