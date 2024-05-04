package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hashicorp/consul/api"
)

type consul struct {
	client *api.Client
	name   string
	ip     string
	port   int
}

func (c *consul) serviceId() string {
	return fmt.Sprintf("%s-%s-%d", c.name, c.ip, c.port)
}

func (c *consul) RegisterService(chk bool) error {
	// 健康检查
	check := &api.AgentServiceCheck{
		GRPC:     fmt.Sprintf("%s:%d", c.ip, c.port), // 这里一定是外部可以访问的地址
		Timeout:  "10s",                              // 超时时间
		Interval: "10s",                              // 运行检查的频率
		// 指定时间后自动注销不健康的服务节点
		// 最小超时时间为1分钟，收获不健康服务的进程每30秒运行一次，因此触发注销的时间可能略长于配置的超时时间。
		DeregisterCriticalServiceAfter: "1m",
	}
	srv := &api.AgentServiceRegistration{
		ID:      c.serviceId(),
		Name:    c.name,
		Tags:    []string{"tiny", "hello"},
		Address: c.ip,
		Port:    c.port,
	}
	if chk {
		srv.Check = check
	}
	return c.client.Agent().ServiceRegister(srv)
}

func (c *consul) Discover(filter string) {
	// 返回的是一个 map[string]*api.AgentService
	// 其中key是服务ID，值是注册的服务信息
	serviceMap, err := c.client.Agent().ServicesWithFilter(filter)
	if err != nil {
		fmt.Printf("query service from consul failed, err:%v\n", err)
		return
	}

	for k, v := range serviceMap {
		fmt.Println(k)
		fmt.Println(v.ID, v.Service, v.Tags)
	}
}

// Deregister 注销服务
func (c *consul) Deregister() error {
	return c.client.Agent().ServiceDeregister(c.serviceId())
}

func newConsul(addr string, name string, ip string, port int) (*consul, error) {
	cfg := api.DefaultConfig()
	cfg.Address = addr
	c, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &consul{
		client: c,
		name:   name,
		ip:     ip,
		port:   port,
	}, nil

}

// GetOutboundIP 获取本机的出口IP
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

func main() {

	c, err := newConsul("192.168.56.101:8500", "greet", "127.0.0.1", 8888)
	if err != nil {
		log.Println(err)
		return
	}
	err = c.RegisterService(false)
	if err != nil {
		log.Println("hhh", err)
		return
	}

	time.Sleep(time.Second * 5)
	c.Discover("Service==`greet`")
	time.Sleep(time.Second * 5)

	fmt.Println(c.Deregister())
}
