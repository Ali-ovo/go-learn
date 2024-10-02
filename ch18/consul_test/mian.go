package main

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

func Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.189.128:8500"

	client, err := api.NewClient(cfg)

	if err != nil {
		panic(err)
	}

	// 检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           "http://192.168.189.128:8021/health",
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	// client.Agent().ServiceDeregister(id)

	if err != nil {
		panic(err)
	}

	return nil
}

func AllServices() error {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.189.128:8500"

	client, err := api.NewClient(cfg)

	if err != nil {
		panic(err)
	}

	data, err := client.Agent().Services()
	if err != nil {
		panic(err)
	}

	for key, values := range data {
		fmt.Println(key, values)
	}

	return nil
}

func FilterService() {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.189.128:8500"

	client, err := api.NewClient(cfg)

	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(`Service == "go-test2"`)

	if err != nil {
		panic(err)
	}

	for key, _ := range data {
		fmt.Println(key)
	}
}

func main() {

	// _ = Register("192.168.189.128", 8021, "user-web", []string{"ali"}, "user-web")

	// AllServices()

	FilterService()
}
