package upnp

import (
	// "log"
	"errors"
	"net"
	"strings"
)

//获取本机能联网的ip地址
func GetLocalIntenetIp() string {
	/*
	  获得所有本机地址
	  判断能联网的ip地址
	*/

	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		panic(errors.New("不能连接网络"))
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}
