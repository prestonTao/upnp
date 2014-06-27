package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	youleTest()
}

func youleTest() {
	lAddr := "192.168.1.100"
	rAddr := "192.168.1.2"
	//---------------------------------------------------------
	//      搜素网关设备
	//---------------------------------------------------------

	searchDevice(lAddr+":9981", "239.255.255.250:1900")

	//---------------------------------------------------------
	//      查看设备描述
	//---------------------------------------------------------

	// readDeviceDesc(rAddr + ":1900")

	//---------------------------------------------------------
	//      查看设备状态 SOAPAction: "urn:schemas-upnp-org:service:WANIPConnection:1#GetStatusInfo"\r\n
	//---------------------------------------------------------
	// getDeviceStatusInfo(rAddr + ":1900")
	getDeviceStatusInfo(rAddr + ":56688")

	addPortMapping(rAddr + ":56688")

	time.Sleep(time.Second * 10)

	remotePort(rAddr + ":56688")
}

func simple1() {
	lAddr := "192.168.1.100"
	rAddr := "192.168.1.1"
	//---------------------------------------------------------
	//      搜素网关设备
	//---------------------------------------------------------

	searchDevice(lAddr+":9981", "239.255.255.250:1900")

	//---------------------------------------------------------
	//      查看设备描述
	//---------------------------------------------------------

	// readDeviceDesc(rAddr + ":1900")

	//---------------------------------------------------------
	//      查看设备状态 SOAPAction: "urn:schemas-upnp-org:service:WANIPConnection:1#GetStatusInfo"\r\n
	//---------------------------------------------------------
	// getDeviceStatusInfo(rAddr + ":1900")

	addPortMapping(rAddr + ":1900")

	time.Sleep(time.Second * 10)

	remotePort(rAddr + ":1900")
}

func searchDevice(localAddr, remoteAddr string) string {
	fmt.Println("搜素网关设备")
	msg := "M-SEARCH * HTTP/1.1\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"ST: urn:schemas-upnp-org:device:InternetGatewayDevice:1\r\n" +
		"MAN: \"ssdp:discover\"\r\n" +
		"MX: 3\r\n" + // seconds to delay response
		"\r\n"
	remotAddr, err := net.ResolveUDPAddr("udp", remoteAddr)
	chk(err)
	locaAddr, err := net.ResolveUDPAddr("udp", localAddr)
	chk(err)
	conn, err := net.ListenUDP("udp", locaAddr)
	chk(err)
	_, err = conn.WriteToUDP([]byte(msg), remotAddr)
	chk(err)
	buf := make([]byte, 1024)
	_, _, err = conn.ReadFromUDP(buf)
	chk(err)
	defer conn.Close()
	fmt.Println(string(buf))
	return string(buf)
}

func readDeviceDesc(rAddr string) string {
	fmt.Println("查看设备描述")
	msg := "GET /igd.xml HTTP/1.1\r\n" +
		"User-Agent: Java/1.7.0_45\r\n" +
		"Host: 192.168.1.1:1900\r\n" +
		"Accept: text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2\r\n" +
		"Connection: keep-alive\r\n\r\n"
	conn, err := net.Dial("tcp", rAddr)
	chk(err)
	_, err = conn.Write([]byte(msg))
	chk(err)
	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	chk(err)
	fmt.Println(string(buf))
	buf = make([]byte, 3048)
	_, err = conn.Read(buf)
	chk(err)
	fmt.Println(string(buf))

	// //查看设备状态
	// fmt.Println("查看设备状态")
	// statusHeader := `POST /ipc HTTP/1.1\r\n
	// Content-Type: text/xml\r\n
	// SOAPAction: "urn:schemas-upnp-org:device:InternetGatewayDevice:1#GetStatusInfo"\r\n
	// Connection: Close\r\n
	// User-Agent: Java/1.7.0_45\r\n
	// Host: 192.168.1.1:1900\r\n
	// Accept: text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2\r\n
	// Content-Length: 309\r\n
	// \r\n`
	// statusBody := `<?xml version="1.0"?><SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><SOAP-ENV:Body><m:GetStatusInfo xmlns:m="urn:schemas-upnp-org:device:InternetGatewayDevice:1"></m:GetStatusInfo></SOAP-ENV:Body></SOAP-ENV:Envelope>`
	// _, err = conn.Write([]byte(statusHeader))
	// chk(err)
	// _, err = conn.Write([]byte(statusBody))
	// chk(err)

	// buf = make([]byte, 1024)
	// _, err = conn.Read(buf)
	// chk(err)
	// fmt.Println(string(buf))
	// buf = make([]byte, 3048)
	// _, err = conn.Read(buf)
	// chk(err)
	// fmt.Println(string(buf))
	return string(buf)
}

func getDeviceStatusInfo(rAddr string) {

	fmt.Println("查看设备状态")

	readMappingBody := `<?xml version="1.0"?>
	<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
	<SOAP-ENV:Body>
	<m:GetStatusInfo xmlns:m="urn:schemas-upnp-org:service:WANIPConnection:1">
	</m:GetStatusInfo></SOAP-ENV:Body></SOAP-ENV:Envelope>`

	client := &http.Client{}
	// 第三个参数设置body部分
	reqest, _ := http.NewRequest("POST", "http://"+rAddr+"/ipc", strings.NewReader(readMappingBody))
	reqest.Proto = "HTTP/1.1"
	reqest.Host = rAddr

	reqest.Header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	reqest.Header.Set("Content-Type", "text/xml")
	reqest.Header.Set("SOAPAction", "\"urn:schemas-upnp-org:service:WANIPConnection:1#GetStatusInfo\"")

	reqest.Header.Set("Connection", "Close")
	reqest.Header.Set("Content-Length", string(len([]byte(readMappingBody))))

	response, _ := client.Do(reqest)

	body, _ := ioutil.ReadAll(response.Body)
	//bodystr := string(body)
	fmt.Println(response.StatusCode)
	if response.StatusCode == 200 {
		fmt.Println(response.Header)
		fmt.Println(string(body))
	}
}

func addPortMapping(rAddr string) {

	fmt.Println("添加一个端口映射")

	readMappingBody := `<?xml version="1.0"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
<SOAP-ENV:Body>
<m:AddPortMapping xmlns:m="urn:schemas-upnp-org:service:WANIPConnection:1">
<NewExternalPort>6991</NewExternalPort>
<NewInternalPort>6991</NewInternalPort>
<NewProtocol>TCP</NewProtocol>
<NewEnabled>1</NewEnabled>
<NewInternalClient>192.168.1.100</NewInternalClient>
<NewLeaseDuration>0</NewLeaseDuration>
<NewPortMappingDescription>test</NewPortMappingDescription>
<NewRemoteHost></NewRemoteHost>
</m:AddPortMapping>
</SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	client := &http.Client{}
	// 第三个参数设置body部分
	reqest, _ := http.NewRequest("POST", "http://"+rAddr+"/ipc", strings.NewReader(readMappingBody))
	reqest.Proto = "HTTP/1.1"
	reqest.Host = rAddr

	reqest.Header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	reqest.Header.Set("Content-Type", "text/xml")
	reqest.Header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:1#AddPortMapping"`)

	reqest.Header.Set("Connection", "Close")
	reqest.Header.Set("Content-Length", string(len([]byte(readMappingBody))))

	response, _ := client.Do(reqest)

	body, _ := ioutil.ReadAll(response.Body)
	//bodystr := string(body)
	fmt.Println(response.StatusCode)
	if response.StatusCode == 200 {
		fmt.Println(response.Header)
		fmt.Println(string(body))
	}
}

func remotePort(rAddr string) {
	fmt.Println("删除一个端口映射")

	readMappingBody := `<?xml version="1.0"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
<SOAP-ENV:Body>
<m:DeletePortMapping xmlns:m="urn:schemas-upnp-org:service:WANIPConnection:1">
<NewExternalPort>6991</NewExternalPort>
<NewProtocol>TCP</NewProtocol>
<NewRemoteHost></NewRemoteHost>
</m:DeletePortMapping>
</SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	client := &http.Client{}
	// 第三个参数设置body部分
	reqest, _ := http.NewRequest("POST", "http://"+rAddr+"/ipc", strings.NewReader(readMappingBody))
	reqest.Proto = "HTTP/1.1"
	reqest.Host = rAddr

	reqest.Header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	reqest.Header.Set("Content-Type", "text/xml")
	reqest.Header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:1#DeletePortMapping"`)

	reqest.Header.Set("Connection", "Close")
	reqest.Header.Set("Content-Length", string(len([]byte(readMappingBody))))

	response, _ := client.Do(reqest)

	body, _ := ioutil.ReadAll(response.Body)
	//bodystr := string(body)
	fmt.Println(response.StatusCode)
	if response.StatusCode == 200 {
		fmt.Println(response.Header)
		fmt.Println(string(body))
	}
}
