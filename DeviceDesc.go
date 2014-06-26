package upnp

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
)

type DeviceDesc struct {
	upnp *Upnp
}

func (this *DeviceDesc) Send() bool {
	request := this.BuildRequest()
	response, _ := http.DefaultClient.Do(request)
	resultBody, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode == 200 {
		this.resolve(string(resultBody))
		return true
	}
	return false
}
func (this *DeviceDesc) BuildRequest() *http.Request {
	//请求头
	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("User-Agent", "preston")
	header.Set("Host", this.upnp.Geteway.Host)
	header.Set("Connection", "keep-alive")

	//请求
	request, _ := http.NewRequest("GET", "http://"+this.upnp.Geteway.Host+this.upnp.Geteway.DeviceDescUrl, nil)
	request.Header = header
	// request := http.Request{Method: "GET", Proto: "HTTP/1.1",
	// 	Host: this.upnp.Geteway.Host, Url: this.upnp.Geteway.DeviceDescUrl, Header: header}
	return request
}

func (this *DeviceDesc) resolve(resultStr string) {
	inputReader := strings.NewReader(resultStr)

	// 从文件读取，如可以如下：
	// content, err := ioutil.ReadFile("studygolang.xml")
	// decoder := xml.NewDecoder(bytes.NewBuffer(content))

	lastLabel := ""

	ISUpnpServer := false

	IScontrolURL := false
	var controlURL string //`controlURL`
	// var eventSubURL string //`eventSubURL`
	// var SCPDURL string     //`SCPDURL`

	decoder := xml.NewDecoder(inputReader)
	for t, err := decoder.Token(); err == nil && !IScontrolURL; t, err = decoder.Token() {
		switch token := t.(type) {
		// 处理元素开始（标签）
		case xml.StartElement:
			if ISUpnpServer {
				name := token.Name.Local
				lastLabel = name
			}

		// 处理元素结束（标签）
		case xml.EndElement:
			// log.Println("结束标记：", token.Name.Local)
		// 处理字符数据（这里就是元素的文本）
		case xml.CharData:
			//得到url后其他标记就不处理了
			content := string([]byte(token))

			//找到提供端口映射的服务
			if content == this.upnp.Geteway.ServiceType {
				ISUpnpServer = true
				continue
			}
			//urn:upnp-org:serviceId:WANIPConnection
			if ISUpnpServer {
				switch lastLabel {
				case "controlURL":

					controlURL = content
					IScontrolURL = true
				case "eventSubURL":
					// eventSubURL = content
				case "SCPDURL":
					// SCPDURL = content
				}
			}
		default:
			// ...
		}
	}
	this.upnp.CtrlUrl = controlURL
}
