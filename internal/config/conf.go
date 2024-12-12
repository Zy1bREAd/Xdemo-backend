package config

import (
	"strings"
)

func LoadInConfig(mode string) {
	mode = strings.ToLower(mode)
	switch mode {
	case "local":
		var yamlConfig YAMLConfig
		yamlConfig.LoadInConfigEnv()
	case "container":
		ReadContainerEnv()
	}
}

// func ReadContainerEnv() *YAMLConfigProvider {
// 	return &YAMLConfigProvider{
// 		DataBase: DataBaseConfig{
// 			DBUser:     os.Getenv("XDEMO_DB_USER"),
// 			DBPassword: os.Getenv("XDEMO_DB_PASSWORD"),
// 			DBName:     os.Getenv("XDEMO_DB_NAME"),
// 			Host:       os.Getenv("XDEMO_DB_HOST"),
// 			Port:       os.Getenv("XDEMO_DB_PORT"),
// 		},
// 		Redis: RedisConfig{
// 			Port: os.Getenv("XDEMO_REDIS_PORT"),
// 			Addr: os.Getenv("XDEMO_REDIS_HOST"),
// 			// DB: os.Getenv("XDEMO_REDIS_DB"),
// 			Password: os.Getenv("XDEMO_REDIS_PASSWORD"),
// 			// TLS: os.Getenv("XDEMO_REDIS_TLS"),

// 		},
// 	}
// 	// ac.DataBase.DBName = os.Getenv("XDEMO_DB_NAME")
// 	// ac.DataBase.Host = os.Getenv("XDEMO_DB_HOST")
// 	// ac.DataBase.Port = os.Getenv("XDEMO_DB_PORT")
// }
