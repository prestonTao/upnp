package main

import (
	"bufio"
	"fmt"
	"github.com/prestonTao/upnp"
	"os"
	"strconv"
	"strings"
)

var mapping = new(upnp.Upnp)
var reader = bufio.NewReader(os.Stdin)

var localPort = 1990
var remotePort = 1990

func init() {

}

func main() {
	Start()
}

func Start() {
	if !CheckNet() {
		fmt.Println("你的路由器不支持upnp协议")
		return
	}
	fmt.Println("本机ip地址：", mapping.LocalHost)

	ExternalIPAddr()

tag:
	if !GetInput() {
		goto tag
	}
	if !AddPortMapping(localPort, remotePort) {
		goto tag
	}

	fmt.Println("--------------------------------------")
	fmt.Println("1.  stop    停止程序并回收映射的端口")
	fmt.Println("2.  add     添加一个端口映射")
	fmt.Println("3.  del     手动删除一个端口映射")
	fmt.Println("\n 注意：此程序映射的端口默认是TCP端口")
	fmt.Println("       需要映射udp端口请访问：")
	fmt.Println("       http://github.com/prestonTao/upnp")
	fmt.Println("--------------------------------------")

	running := true
	for running {
		data, _, _ := reader.ReadLine()
		commands := strings.Split(string(data), " ")
		switch commands[0] {
		case "help":

		case "stop":
			running = false
			mapping.Reclaim()
		case "add":
			goto tag
		case "del":
		tagDel:
			if !GetInput() {
				goto tagDel
			}
			DelPortMapping(localPort, remotePort)
		case "cdp":
		case "dump":
		}
	}

}

/*
	检查网络是否支持upnp协议
*/
func CheckNet() bool {
	err := mapping.SearchGateway()
	if err != nil {
		return false
	} else {
		return true
	}
}

//获得公网ip地址
func ExternalIPAddr() {
	err := mapping.ExternalIPAddr()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("外网ip地址为：", mapping.GatewayOutsideIP)
	}
}

/*
	得到用户输入的端口
*/
func GetInput() bool {
	var err error
	fmt.Println("请输入要映射的本地端口：")
	data, _, _ := reader.ReadLine()
	localPort, err = strconv.Atoi(string(data))
	if err != nil {
		fmt.Println("输入的端口号错误，请输入 0-65535 的数字")
		return false
	}
	if localPort < 0 || localPort > 65535 {
		fmt.Println("输入的端口号错误，请输入 0-65535 的数字")
		return false
	}

	fmt.Println("请输入要映射到外网的端口：")
	data, _, _ = reader.ReadLine()
	remotePort, err = strconv.Atoi(string(data))
	if err != nil {
		fmt.Println("输入的端口号错误，请输入 0-65535 的数字")
		return false
	}
	if remotePort < 0 || remotePort > 65535 {
		fmt.Println("输入的端口号错误，请输入 0-65535 的数字")
		return false
	}
	return true
}

/*
	添加一个端口映射
*/
func AddPortMapping(localPort, remotePort int) bool {
	//添加一个端口映射
	if err := mapping.AddPortMapping(localPort, remotePort, "TCP"); err == nil {
		fmt.Println("端口映射成功")
		return true
	} else {
		fmt.Println("端口映射失败")
		return false
	}
}

/*
	删除一个端口映射
*/
func DelPortMapping(localPort, remotePort int) {
	mapping.DelPortMapping(remotePort, "TCP")
}
