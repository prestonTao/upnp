package main

import (
	// "bufio"
	"fmt"
	"github.com/prestonTao/upnp"
	// "os"
)

func main() {
	SearchGateway()
	ExternalIPAddr()
	AddPortMapping()
}

//搜索网关设备
func SearchGateway() {
	upnpMan := new(upnp.Upnp)
	err := upnpMan.SearchGateway()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("本机ip地址：", upnpMan.LocalHost)
		fmt.Println("upnp设备地址：", upnpMan.Geteway.Host)
	}
}

//获得公网ip地址
func ExternalIPAddr() {
	upnpMan := new(upnp.Upnp)
	err := upnpMan.ExternalIPAddr()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("外网ip地址为：", upnpMan.GetewayOutsideIP)
	}
}

//添加一个端口映射
func AddPortMapping() {
	mapping := new(upnp.Upnp)
	if err := mapping.AddPortMapping(55789, 55789, "TCP"); err == nil {
		fmt.Println("端口映射成功")
		mapping.Reclaim()
	} else {
		fmt.Println("端口映射失败")
	}

}
