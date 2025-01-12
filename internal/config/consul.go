package config

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	capi "github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v2"
)

// 实现ConfigProvider接口的Consul配置中心
type ConsulConfig struct {
	Env AppConfigEnv
}

// 读取Consul配置中心的KV变量作为App Env
func (c *ConsulConfig) LoadInConfigEnv() {

	InitConsul() // Consul本身的配置从Env或者配置文件中读取

	// 根据当前环境配置（dev/prod）和镜像tag获取对应的application version
	prefix := "dev" + "/" + "xdemo" + "/" + "0.2.1"
	err := c.LoadInAllKey(prefix)
	if err != nil {
		log.Fatalln(err)
		return
	}
	appConfigEnv = c
}

func (c *ConsulConfig) GetAppConfigEnv() AppConfigEnv {
	return c.Env
}

var consulClient *capi.Client

func InitConsul() {
	consulClient = NewConsulAPIClient()
	if consulClient == nil {
		log.Fatalln("Init Consul is Failed!!! Exit....")
	}
	log.Println("Consul Connect Success~")
}

func NewConsulAPIClient() *capi.Client {
	// 默认从ENV中读取Consul配置
	config := &capi.Config{
		// Address: os.Getenv("XDEMO_CONSUL_ADDR"),
		// Token:   os.Getenv("XDEMO_CONSUL_AUTH_TOKEN"),
		Address: "159.75.119.146:8500",
		Token:   "54d11d67-9a0c-96e2-eebb-32ab0dd37794",
		Scheme:  "http",
	}
	client, err := capi.NewClient(config)
	if err != nil {
		log.Println("Create Consul API Client Failed,", err)
		return nil
	}
	return client
}

// 新方法

func (c *ConsulConfig) LoadInAllKey(keyPrefix string) error {
	// 通过查找指定前缀的特定服务，获取所有符合的key value
	k := consulClient.KV()
	pairs, _, err := k.List(keyPrefix, nil)
	if err != nil {
		log.Println("遍历所有Key出错,", err)
		return err
	}
	fmt.Println(pairs)
	for _, pair := range pairs {
		// 匹配对应version的
		err := c.ParseConfig(pair)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ConsulConfig) GetKey(key string) ([]byte, error) {
	k := consulClient.KV()
	pair, _, err := k.Get(key, nil)
	if err != nil {
		log.Println("获取Key出错", err)
		return nil, err
	} else if pair == nil {
		log.Println("可能找不到对应的Key")
		return nil, fmt.Errorf("the Key Not Exist")
	}
	fmt.Println(pair.Key, pair.Value)
	return pair.Value, nil
}

func (c *ConsulConfig) ParseConfig(pair *capi.KVPair) error {
	// 以YAML配置为优先解析
	if strings.HasSuffix(pair.Key, "yaml") {
		// 解析Yaml配置
		err := yaml.Unmarshal(pair.Value, &c.Env)
		if err != nil {
			log.Println("解析YAML配置文件出错", err)
			return err
		}
	} else if strings.HasSuffix(pair.Key, "json") {
		// 解析json的配置
		err := json.Unmarshal(pair.Value, &c.Env)
		if err != nil {
			log.Println("解析JSON配置文件出错", err)
			return err
		}
	} else {
		fmt.Println(pair.Value)
	}
	return nil
}
