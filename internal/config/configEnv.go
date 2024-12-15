package config

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// App的环境变量配置
type AppConfigEnv struct {
	DataBase DataBaseConfig `yaml:"mysql"`
	System   SystemConfig   `yaml:"system"`
	Redis    RedisConfig    `yaml:"redis"`
	Docker   DockerConfig   `yaml:"docker"`
	K8s      K8sConfig      `yaml:"k8s"`
}

type SystemConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
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

// 统一接口获取环境变量（API调用接口）
func NewConfigEnvProvider() AppConfigEnv {
	return appConfigEnv.GetAppConfigEnv()
}

// 指定方式读取ConfigEnv
func InitConfigEnv() {
	// 通过环境变量判断当前App启动mode
	appMode := os.Getenv("XDEMO_SYSTEM_MODE")
	appMode = strings.ToLower(appMode)
	if appMode == "" {
		appMode = "local"
	}
	switch appMode {
	case "local":
		var yamlConfig YAMLConfig
		yamlConfig.LoadInConfigEnv()
	case "container":
		var containerConfig ContainerConfig
		containerConfig.LoadInConfigEnv()
		fmt.Println("容器变量：", containerConfig.Env)
	default:
		panic(fmt.Errorf("暂不支持其他App启动模式"))
	}
	fmt.Println("App Start Mode: ", appMode)
}

var appConfigEnv ConfigEnv

// ! 环境变量配置接口
type ConfigEnv interface {
	LoadInConfigEnv()
	GetAppConfigEnv() AppConfigEnv
}

// ! YamlConfig
type YAMLConfig struct {
	Env AppConfigEnv
}

func (yc *YAMLConfig) LoadInConfigEnv() {
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	systemType := runtime.GOOS
	// 注意windows和Linux下路径的斜杠问题！
	var configPath string
	if systemType == "windows" {
		configPath = currentPath + "\\settings.yaml"
	} else if systemType == "linux" {
		configPath = currentPath + "/settings.yaml"
	}
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	// 解析yaml file到结构体中
	err = yaml.Unmarshal(yamlFile, &yc.Env)
	if err != nil {
		panic(err)
	}
	appConfigEnv = yc
}

func (yc *YAMLConfig) GetAppConfigEnv() AppConfigEnv {
	return yc.Env
}

// ! ContainerConfig
type ContainerConfig struct {
	Env AppConfigEnv
}

func (cc *ContainerConfig) LoadInConfigEnv() {
	// e.g [ALLUSERSPROFILE=C:\ProgramData APPDATA=C:\Users\Administrator\AppData\Roaming]
	envSlice := os.Environ()
	for _, v := range envSlice {
		envParts := strings.SplitN(v, "=", 2)
		if len(envParts) < 2 {
			continue
		}
		envKey := envParts[0]
		envValue := envParts[1]
		// 对应AppConfig前缀的才处理
		if strings.HasPrefix(envKey, "XDEMO_") {
			splitEnvKey := strings.SplitN(envKey, "_", 3)
			// 匹配对应的函数处理
			switch strings.ToLower(splitEnvKey[1]) {
			case "database":
				// fmt.Println("进入mysql的map key value")
				cc.setDataBaseEnvMap(splitEnvKey[2], envValue)
			case "redis":
				cc.setRedisEnvMap(splitEnvKey[2], envValue)
			case "system":
				cc.setSystemEnvMap(splitEnvKey[2], envValue)
			case "docker":
				cc.setDockerEnvMap(splitEnvKey[2], envValue)
			case "k8s":
				cc.setK8sEnvMap(splitEnvKey[2], envValue)
			}
		}
	}
}

func (cc *ContainerConfig) GetAppConfigEnv() AppConfigEnv {
	return cc.Env
}

func (cc *ContainerConfig) setDataBaseEnvMap(keySplitThird string, value string) {
	switch strings.ToLower(keySplitThird) {
	case "host":
		cc.Env.DataBase.Host = value
	case "port":
		cc.Env.DataBase.Port = value
	case "dbname":
		cc.Env.DataBase.DBName = value
	case "dbpassword":
		cc.Env.DataBase.DBPassword = value
	case "dbuser":
		cc.Env.DataBase.DBUser = value
	}
}

func (cc *ContainerConfig) setRedisEnvMap(keySplitThird string, value string) {
	switch strings.ToLower(keySplitThird) {
	case "addr":
		cc.Env.Redis.Addr = value
	case "port":
		cc.Env.Redis.Port = value
	case "db":
		convertResult, err := strconv.ParseInt(value, 0, 32)
		if err != nil {
			log.Println("Convert To Int Error ", err)
			panic(err)
		}
		cc.Env.Redis.DB = int(convertResult)
	case "password":
		cc.Env.Redis.Password = value
	case "tls":
		convertResult, err := strconv.ParseBool(value)
		if err != nil {
			log.Println("Convert To Bool Error ", err)
			panic(err)
		}
		cc.Env.Redis.TLS = convertResult
	}
}

func (cc *ContainerConfig) setSystemEnvMap(keySplitThird string, value string) {
	switch strings.ToLower(keySplitThird) {
	case "env":
		cc.Env.System.Mode = value
	case "host":
		cc.Env.System.Host = value
	case "port":
		cc.Env.System.Port = value
	}
}

func (cc *ContainerConfig) setDockerEnvMap(keySplitThird string, value string) {
	switch strings.ToLower(keySplitThird) {
	case "host":
		cc.Env.Docker.Host = value
	case "version":
		cc.Env.Docker.Version = value
	}
}

func (cc *ContainerConfig) setK8sEnvMap(keySplitThird string, value string) {
	switch strings.ToLower(keySplitThird) {
	case "mode":
		convertResult, err := strconv.ParseInt(value, 0, 32)
		if err != nil {
			log.Println("Convert To Int Error ", err)
			panic(err)
		}
		cc.Env.Redis.DB = int(convertResult)
	case "kubeconfig":
		cc.Env.K8s.KubeConfig = value
	}
}

// utils
