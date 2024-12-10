package config

import (
	"fmt"
	"os"
	"runtime"

	"gopkg.in/yaml.v2"
)

type YAMLConfig struct {
	DataBase DataBaseConfig `yaml:"mysql"`
	System   SystemConfig   `yaml:"system"`
	Redis    RedisConfig    `yaml:"redis"`
	Docker   DockerConfig   `yaml:"docker"`
	K8s      K8sConfig      `yaml:"k8s"`
}

type SystemConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Env  string `yaml:"env"`
}

type DataBaseConfig struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	DBName     string `yaml:"dbname"`
	DBUser     string `yaml:"dbuser"`
	DBPassword string `yaml:"dbpassword"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Port     string `yaml:"port"`
	DB       int    `yaml:"db"`
	Password string `yaml:"password"`
	TLS      bool   `yaml:"tls"`
}

type DockerConfig struct {
	Host    string `yaml:"host"`
	Version string `yaml:"version"`
}

type K8sConfig struct {
	Mode       int    `yaml:"mode"`
	KubeConfig string `yaml:"kubeconfig"`
}

func LoadConfig() *YAMLConfig {
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(currentPath)
	systemType := runtime.GOOS
	// 注意windows和Linux下路径的斜杠问题！
	var configPath string
	if systemType == "windows" {
		configPath = currentPath + "\\settings.yaml"
	} else if systemType == "linux" {
		configPath = currentPath + "/settings.yaml"
	}
	yf, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	// 解析yaml file到结构体中
	var Config YAMLConfig
	err = yaml.Unmarshal(yf, &Config)
	if err != nil {
		panic(err)
	}
	return &Config
}
