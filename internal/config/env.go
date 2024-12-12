package config

import (
	"fmt"
	"os"
	"strings"
)

//! 单例模式（适合现在全服务都塞一块的情况下）
// type ContainerAppConfig struct{
// 	RedisPort string
// 	RedisAddr string
// 	RedisDB
// }

// ! 依赖注入方式(仅适合微服务)
// e.g. StructObject.Redis().GetListenAddr().Env()  ==> 获取到Redis的变量
type ConfigProvider interface {
	// 对应哪个App的提供商(后期再拆开成一个单独的接口？)
	MySQL() ConfigProvider
	Redis() ConfigProvider

	// 获取值
	Env() string
}

type MySQLConfigProvider interface {
	ConfigProvider
	// 获取对应Key
	System() MySQLConfigProvider
	ListenHost() MySQLConfigProvider
	ListenPort() MySQLConfigProvider
	DBName() MySQLConfigProvider
}

type SystemConfigProvider interface {
	ConfigProvider
}

type ContainerConfigProvider struct {
	EnvForApp     string
	App           string
	DefaultEnvMap map[string]string
}

// 从容器Env中读取变量存入AppConfigEnv结构体中
func ReadContainerEnv() {

}

func InitAppConfig(mode string, app string) ConfigProvider {
	mode = strings.ToLower(mode)
	switch mode {
	case "yaml":
		fmt.Println("还没实现")
	case "container":
		return &ContainerConfigProvider{
			DefaultEnvMap: map[string]string{
				"SYSTEM_HOST": "127.0.0.1",
				"SYSTEM_PORT": "8081",
			},
		}
	}
	return nil
}

// 主要判断用哪种方式获取环境变量
// Q1：为什么这里要返回指向这个对象（实际会实现接口）的指针呢？
// A1：因为返回指针是为了能够修改到这个对象内部的状态（成员的值）。如果返回的是ContainerConfigProvider对象本身（值传递），那么在方法内部对ccp对象的修改不会影响到外部调用者所拥有的对象，因为 Go 语言中值传递会复制对象。
func (ccp *ContainerConfigProvider) ListenHost() MySQLConfigProvider {
	EnvKeyName := fmt.Sprintf("XDEMO_%s_LISTEN_HOST", ccp.App)
	v, exits := os.LookupEnv(EnvKeyName)
	if exits {
		ccp.EnvForApp = v
	} else {
		ccp.EnvForApp = ""
	}
	return ccp
}

func (ccp *ContainerConfigProvider) ListenPort() MySQLConfigProvider {
	EnvKeyName := fmt.Sprintf("XDEMO_%s_LISTEN_PORT", ccp.App)
	v, exits := os.LookupEnv(EnvKeyName)
	if exits {
		ccp.EnvForApp = v
	} else {
		ccp.EnvForApp = ""
	}
	return ccp
}

func (ccp *ContainerConfigProvider) DBName() MySQLConfigProvider {
	EnvKeyName := fmt.Sprintf("XDEMO_%s_DBNAME", ccp.App)
	v, exits := os.LookupEnv(EnvKeyName)
	if exits {
		ccp.EnvForApp = v
	} else {
		ccp.EnvForApp = ""
	}
	return ccp
}

func (ccp *ContainerConfigProvider) MySQL() ConfigProvider {
	ccp.App = "MYSQL"
	return ccp
}

func (ccp *ContainerConfigProvider) Redis() ConfigProvider {
	ccp.App = "REDIS"
	return ccp
}

func (ccp *ContainerConfigProvider) System() MySQLConfigProvider {
	ccp.App = "SYSTEM"
	return ccp
}

func (ccp *ContainerConfigProvider) Env() string {
	return ccp.EnvForApp
}
