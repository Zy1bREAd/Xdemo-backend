package api

import (
	"context"
	"fmt"
	"log"
	config "xdemo/internal/config"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var K8sInstance *MyK8s

type MyK8s struct {
	KubeConfig *rest.Config
	ClientSet  *kubernetes.Clientset
}

func NewK8s() *MyK8s {
	if K8sInstance == nil {
		return &MyK8s{}
	}
	fmt.Println("K8sClient 已经初始化")
	return K8sInstance
}

func InitK8sClient() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Init K8s Client Failed,%s", err)
			return
		}
	}()
	K8sInstance = NewK8s()
	// 读取变量配置
	configProvider := config.NewConfigEnvProvider()
	if configProvider.K8s.KubeConfig == "" {
		configProvider.K8s.KubeConfig = clientcmd.RecommendedHomeFile
	}
	err := K8sInstance.BuildKubeConfig(configProvider.K8s.Mode, configProvider.K8s.KubeConfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(K8sInstance.KubeConfig)
	if err != nil {
		panic(err)
	}
	K8sInstance.ClientSet = clientset
	log.Println("K8s Cluster API Connect Success~")
}

// 读取kubeconfig配置。（master url = 1， kubeconfig = 2）
func (c *MyK8s) BuildKubeConfig(mode int, configPath string) error {
	switch mode {
	case 1:
		fmt.Println("使用master url进行初始化配置")
		fmt.Println("暂不支持")
	case 2:
		cfg, err := clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			log.Println("Build K8s kubeconfig Failed ", err)
			return err
		}
		c.KubeConfig = cfg
	}
	return nil
}

func (c *MyK8s) GetPodsForDefault() error {
	podClient := c.ClientSet.CoreV1().Pods("")
	podList, err := podClient.List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil
	}
	fmt.Println(podList.Items)
	for _, v := range podList.Items {
		fmt.Println(v.Name)
	}
	return nil
}
