package upnp

import (
	"log"
	"net"
	"strings"
	"time"
	// "net/http"
)

type Gateway struct {
	GatewayName   string //网关名称
	Host          string //网关ip和端口
	DeviceDescUrl string //网关设备描述路径
	Cache         string //cache
	ST            string
	USN           string
	deviceType    string //设备的urn   "urn:schemas-upnp-org:service:WANIPConnection:1"
	ControlURL    string //设备端口映射请求路径
	ServiceType   string //提供upnp服务的服务类型
}

type SearchGateway struct {
	searchMessage string
	upnp          *Upnp
}

func (this *SearchGateway) Send() bool {
	this.buildRequest()
	c := make(chan string)
	go this.send(c)
	result := <-c
	if result == "" {
		//超时了
		this.upnp.Active = false
		return false
	}
	this.resolve(result)

	this.upnp.Gateway.ServiceType = "urn:schemas-upnp-org:service:WANIPConnection:1"
	this.upnp.Active = true
	return true
}
func (this *SearchGateway) send(c chan string) {
	//发送组播消息，要带上端口，格式如："239.255.255.250:1900"
	var conn *net.UDPConn
	defer func() {
		if r := recover(); r != nil {
			//超时了
		}
	}()
	go func(conn *net.UDPConn) {
		defer func() {
			if r := recover(); r != nil {
				//没超时
			}
		}()
		//超时时间为3秒
		time.Sleep(time.Second * 3)
		c <- ""
		conn.Close()
	}(conn)
	remotAddr, err := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
	if err != nil {
		log.Println("组播地址格式不正确")
	}
	locaAddr, err := net.ResolveUDPAddr("udp", this.upnp.LocalHost+":")

	if err != nil {
		log.Println("本地ip地址格式不正确")
	}
	conn, err = net.ListenUDP("udp", locaAddr)
	defer conn.Close()
	if err != nil {
		log.Println("监听udp出错")
	}
	_, err = conn.WriteToUDP([]byte(this.searchMessage), remotAddr)
	if err != nil {
		log.Println("发送msg到组播地址出错")
	}
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Println("从组播地址接搜消息出错")
	}

	result := string(buf[:n])
	c <- result
}
func (this *SearchGateway) buildRequest() {
	this.searchMessage = "M-SEARCH * HTTP/1.1\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"ST: urn:schemas-upnp-org:service:WANIPConnection:1\r\n" +
		"MAN: \"ssdp:discover\"\r\n" + "MX: 3\r\n\r\n"
}

func (this *SearchGateway) resolve(result string) {
	this.upnp.Gateway = &Gateway{}

	lines := strings.Split(result, "\r\n")
	for _, line := range lines {
		//按照第一个冒号分为两个字符串
		nameValues := strings.SplitAfterN(line, ":", 2)
		if len(nameValues) < 2 {
			continue
		}
		switch strings.ToUpper(strings.Trim(strings.Split(nameValues[0], ":")[0], " ")) {
		case "ST":
			this.upnp.Gateway.ST = nameValues[1]
		case "CACHE-CONTROL":
			this.upnp.Gateway.Cache = nameValues[1]
		case "LOCATION":
			urls := strings.Split(strings.Split(nameValues[1], "//")[1], "/")
			this.upnp.Gateway.Host = urls[0]
			this.upnp.Gateway.DeviceDescUrl = "/" + urls[1]
		case "SERVER":
			this.upnp.Gateway.GatewayName = nameValues[1]
		default:
		}
	}
}
