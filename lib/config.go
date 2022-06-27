package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"strings"
	"ui-project/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	SecretKey string `json:"secretKey"`
	Db        Db     `json:"db"`
	Log       Log    `json:"log"`
}

type Db struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Log struct {
	LogLevel int `json:"logLevel"`
	Skip     int `json:"skip"`
}

var defaultConfig Config = Config{
	Host:      "127.0.0.1",
	Port:      8888,
	SecretKey: "Chehung",
	Db: Db{
		Host:     "127.0.0.1",
		Port:     5432,
		Name:     "postgres",
		User:     "",
		Password: "",
	},
	Log: Log{
		LogLevel: 3,
		Skip:     4,
	},
}

func LoadConfig(ctx context.Context, l logger.LogUsecase) {
	var configPath string
	flag.StringVar(&configPath, "c", "", "Configuration file path.")
	flag.Parse()

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("test")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if configPath != "" {
		if content, err := ioutil.ReadFile(configPath); err != nil {
			l.LogErr(ctx, err)
			panic(err)
		} else {
			l.LogInfo(ctx, "Using config file:", configPath)
			viper.ReadConfig(bytes.NewBuffer(content))
		}
	} else {
		if err := viper.ReadInConfig(); err == nil {
			l.LogInfo(ctx, "Using config file:", viper.ConfigFileUsed())
		} else {
			l.LogInfo(ctx, "Using default config file")
			data, _ := json.Marshal(defaultConfig)
			viper.ReadConfig(bytes.NewBuffer(data))
		}
	}
}
