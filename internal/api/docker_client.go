package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"

	reqData "xdemo/internal/app/tmp"
	config "xdemo/internal/config"
)

// 常量
const (
	Default_Respository = ""
)

var DockerInstance *MyDocker
var apiClient *client.Client

type MyDocker struct {
	CTX               context.Context
	ContainerConfig   *container.Config
	HostConfig        *container.HostConfig
	NetworkConfig     *network.NetworkingConfig
	RegistryAuthToken string
}

func NewDocker() *MyDocker {
	if DockerInstance != nil {
		return DockerInstance
	}
	DockerInstance = &MyDocker{
		CTX: context.Background(),
	}
	return DockerInstance
}

func InitDocker() {
	DockerInstance = NewDocker()
	apiClient = DockerInstance.NewDockerClient()
	if apiClient == nil {
		panic("初始化Docker出现错误")
	}
	log.Println("初始化并连接Docker成功！")
}

func (d *MyDocker) NewDockerClient() *client.Client {
	// 从环境变量中读取，这很适合容器
	configProvider := config.NewConfigEnvProvider()
	c, err := client.NewClientWithOpts(client.WithHost(configProvider.Docker.Host), client.WithVersion(configProvider.Docker.Version))
	if err != nil {
		panic(err)
	}
	defer c.Close()
	return c
}

// 业务逻辑接口
func (c *MyDocker) ContainersList() ([]types.Container, error) {
	// 区分是Docker还是K8s集群中的（containerd）
	return apiClient.ContainerList(c.CTX, container.ListOptions{All: true})
}

func (c *MyDocker) ContainerInspect(cid string) (types.ContainerJSON, error) {
	// 区分是Docker还是K8s集群中的（containerd）
	return apiClient.ContainerInspect(c.CTX, cid)
}

func (c *MyDocker) ContainerInspectBaseInfo(cid string) (map[string]interface{}, error) {
	result, err := c.ContainerInspect(cid)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"c_config": result.Config,
		"c_id":     result.ID,
		"c_name":   result.Name,
		"c_state":  result.State.Status,
		"c_ip":     result.NetworkSettings.IPAddress,
	}, nil
}

// 将配置整理起来
func (c *MyDocker) ParseContainerCreateConfig(configData *reqData.ContainerCreateConfig) error {
	if c == nil {
		return fmt.Errorf("未实例化成功")
	}

	// 解析exposedPort和Volumes的数据
	for k, v := range configData.PersistenceVolumes {
		c.ContainerConfig.Volumes[k] = struct{}{}
		c.HostConfig.Binds = append(c.HostConfig.Binds, v)
	}
	bindingMap := nat.PortMap{}
	bindingSet := nat.PortSet{}
	for k, v := range configData.ExposedPorts {
		// 仅对应一个hostpath and port
		bindingSet[k] = struct{}{}
		bindingMap[k] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: v,
			},
		}

	}
	c.ContainerConfig = &container.Config{
		Env:          configData.Env,
		Hostname:     configData.HostName,
		Image:        configData.Image,
		ExposedPorts: bindingSet,
		Tty:          configData.TTY,
		AttachStdin:  configData.Interactive,
		// Volumes: ,

	}
	c.HostConfig = &container.HostConfig{
		PortBindings: bindingMap,
		// Binds: ,
	}

	return nil
}

// 容器创建基础版
// 支持常用的选项： 名字，暴露端口，使用的网络，itd常用参数，持久化等
func (c *MyDocker) ContainerCreate(cName string, i string) (string, error) {
	// 先检查image是否存在
	if i == "" {
		return "", fmt.Errorf("Image不能为空")
	}
	err := c.IsImageExist(i)
	if err != nil {
		log.Println(err)
		// re-Pull image
		result := strings.Split(i, ":")
		fmt.Println(result, i)
		if len(result) != 2 {
			log.Println("解析Image Name发生错误")
			return "", fmt.Errorf("创建容器时拉取Image解析错误")
		}
		_, err := c.ImagePull(result[0], result[1], false)
		if err != nil {
			log.Println("解析Image Name发生错误")
			return "", fmt.Errorf("创建容器时拉取Image错误")
		}
	}
	createResp, err := apiClient.ContainerCreate(c.CTX, c.ContainerConfig, c.HostConfig, c.NetworkConfig, &v1.Platform{}, cName)
	if err != nil {
		fmt.Println(err, err.Error())
		return "", err
	}
	return createResp.ID, nil
}

// 容器创建命令版，能支持更复杂的选项
func (c *MyDocker) ContainerCreateWithCmd() {

}

// 创建容器并立即启动
func (c *MyDocker) ContainerCreateAndStart(cName string, i string) (string, error) {
	cid, err := c.ContainerCreate(cName, i)
	if err != nil {
		return "", err
	}
	return cid, nil
}

func (c *MyDocker) ContainerStart(cid string) error {
	err := apiClient.ContainerStart(c.CTX, cid, container.StartOptions{})
	return err
}

func (c *MyDocker) ContainerStop(cid string, force bool) error {
	opt := container.StopOptions{}
	if force {
		*opt.Timeout = 0
		opt.Signal = "SIGKILL"
	}
	return apiClient.ContainerStop(c.CTX, cid, opt)
}

func (c *MyDocker) ContainerRestart(cid string, force bool) error {
	opt := container.StopOptions{}
	if force {
		*opt.Timeout = 0
		opt.Signal = "SIGKILL"
	}
	return apiClient.ContainerRestart(c.CTX, cid, opt)
}

func (c *MyDocker) ContainerDelete(cid string, delOpts map[string]bool) error {
	err := apiClient.ContainerRemove(c.CTX, cid, container.RemoveOptions{
		RemoveVolumes: delOpts["clean"],
		RemoveLinks:   delOpts["clean"],
		Force:         delOpts["force"],
	})
	return err
}

func (c *MyDocker) IsImageExist(i string) error {
	_, _, err := apiClient.ImageInspectWithRaw(c.CTX, i)
	if err != nil {
		// DOcker官方返回的错误分404和500。很可能没有image，当然这边500错误要单独判断
		return err
	}
	return nil
}

func (c *MyDocker) ImagePull(name string, tag string, isCover bool) (io.ReadCloser, error) {
	fmt.Println(name, tag, isCover)
	target := name + ":" + tag
	if tag == "" {
		target = name
	}
	if name == "" {
		return nil, fmt.Errorf("镜像名不能为空")
	}

	pullOpts := image.PullOptions{
		All: isCover,
	}
	return apiClient.ImagePull(c.CTX, target, pullOpts)
	// 临时buffer
	// var bf []byte
	// tempBf := make([]byte, 1024)
	// reader, err := apiClient.ImagePull(c.CTX, target, pullOpts)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }
	// defer reader.Close()
	// // 循环写入buffer中
	// for {
	// 	n, err := reader.Read(tempBf)
	// 	if err == io.EOF {
	// 		// 传输完毕
	// 		break
	// 	} else if err != nil {
	// 		log.Println("传输拉取镜像的数据时发生错误")
	// 		return nil, err
	// 	}

	// 	newBf := append(bf, tempBf[:n]...)
	// 	bf = newBf
	// }
	// return bf, nil
}

func (c *MyDocker) ContainerExecCmd(cid string, cmd []string) ([]byte, error) {
	if cid == "" {
		return nil, fmt.Errorf("容器ID或名字不能为空")
	} else if len(cmd) == 0 {
		return nil, fmt.Errorf("执行命令不能为空")
	}
	resp, err := apiClient.ContainerExecCreate(c.CTX, cid, container.ExecOptions{AttachStdin: true, AttachStdout: true, AttachStderr: true, Cmd: cmd})
	if err != nil {
		log.Println("创建EXEC Error")
		return nil, err
	}
	hjResp, err := apiClient.ContainerExecAttach(c.CTX, resp.ID, container.ExecStartOptions{})
	if err != nil {
		log.Println("进入EXEC Error")
		return nil, err
	}
	defer hjResp.Close()

	respSlice, err := io.ReadAll(hjResp.Reader)
	if err != nil {
		log.Println("输出Reader数据Error")
		return nil, err
	}
	// var bf []byte
	// hjResp.Conn.Read(bf)
	// fmt.Println("data:", bf)
	// fmt.Println(hjResp)
	// fmt.Println(hjResp.Reader.Read())
	return respSlice, nil
}

func (c *MyDocker) UpdateImageTag(source string, target string) error {
	err := apiClient.ImageTag(c.CTX, source, target)
	if err != nil {
		log.Println("Update Image Tag Error: ", err)
		return err
	}
	return nil
}

func (c *MyDocker) ImagePush(ctx context.Context, target string) (io.ReadCloser, error) {
	// 判断是否登录Docker Registry
	if c.RegistryAuthToken == "" {
		return nil, fmt.Errorf("docker Registry Not Logged")
	}
	resultReader, err := apiClient.ImagePush(c.CTX, target, image.PushOptions{
		RegistryAuth: c.RegistryAuthToken,
	})
	if err != nil {
		return nil, err
	}

	return resultReader, nil
}

func (c *MyDocker) RegistryLogin(ctx context.Context, usr string, pwd string, addr string) error {
	authConfig := registry.AuthConfig{
		Username:      usr,
		Password:      pwd,
		ServerAddress: addr,
	}
	loginResult, err := apiClient.RegistryLogin(ctx, authConfig)
	if err != nil {
		return err
	}
	log.Println("Login Status is ", loginResult.Status)
	// c.RegistryAuthToken = authBody.IdentityToken
	// 将认证信息序列化成JSON，再由此转成base64字符串
	authConfigToJSON, err := json.Marshal(authConfig)
	if err != nil {
		return fmt.Errorf("AuthConfig Marshal To JSON Failed!%s", err)
	}
	encodeData := base64.StdEncoding.EncodeToString(authConfigToJSON)
	fmt.Println(encodeData)
	c.RegistryAuthToken = encodeData
	return nil
}
