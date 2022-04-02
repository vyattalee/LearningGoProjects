package conf

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
)

// 常量
const (
	ConfigFile = "client.yaml" // 配置文件
)

// 变量
var (
	GvaConfig Config // 全局配置
)

type Config struct {
	server Server `yaml:	"server"	mapstructure: "server"`
	tls    bool   `yaml:	"tls"	mapstructure: "tls"`
}

//type Config struct {
//	server string `yaml:	"server"	mapstructure: "server"`
//	tls    bool   `yaml:	"tls"	mapstructure: "tls"`
//}

type Server struct {
	address string `yaml:	"address"	mapstructure: "address"`
	port    uint16 `yaml:	"port"	mapstructure: "port"`
}

func (config *Config) GetServerAddress() string {

	return (config.server.address) + ":" + strconv.Itoa(int(config.server.port)) //string(config.server.port)
	//return config.server
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

// LoadConfig reads configuration from yaml file or environment variables.
func LoadConfig(path string) (config Config, err error) {

	if path == "" {
		//获取项目的执行路径
		path, err = os.Getwd()
		if err != nil {
			panic(err)
		}
		path += "/conf"
	}

	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	//viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	//cfg  := Config{}

	// enable live reloading of config
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config", e.Name, "changed ...")
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Println("error: Failed to reload config: ", err)
			return
		}
	})

	viper.WatchConfig()

	//err = viper.Unmarshal(&config)

	//log.Println("config:", config)
	return config, nil
}

// LoadConfig reads configuration from file or environment variables.
func LoadEnvConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("client.env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

// ViperInit 初始化viper配置解析包，函数可接受命令行参数
func InitConfig() {
	var configFile string
	// 读取配置文件优先级: 命令行 > 默认值
	flag.StringVar(&configFile, "c", ConfigFile, "配置")
	if len(configFile) == 0 {
		// 读取默认配置文件
		panic("配置文件不存在！")
	}
	// 读取配置文件
	v := viper.New()
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("配置解析失败:%s\n", err))
	}
	// 动态监测配置文件
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件发生改变")
		if err := v.Unmarshal(&GvaConfig); err != nil {
			panic(fmt.Errorf("配置重载失败:%s\n", err))
		}
	})
	if err := v.Unmarshal(&GvaConfig); err != nil {
		panic(fmt.Errorf("配置重载失败:%s\n", err))
	}
	// 设置配置文件
	//GvaConfig.ConfigFile = configFile
}
