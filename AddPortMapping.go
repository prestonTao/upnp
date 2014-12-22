package upnp

import (
	// "log"
	// "fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type AddPortMapping struct {
	upnp *Upnp
}

func (this *AddPortMapping) Send(localPort, remotePort int, protocol string) bool {
	request := this.buildRequest(localPort, remotePort, protocol)
	response, _ := http.DefaultClient.Do(request)
	resultBody, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode == 200 {
		this.resolve(string(resultBody))
		return true
	}
	return false
}
func (this *AddPortMapping) buildRequest(localPort, remotePort int, protocol string) *http.Request {
	//请求头
	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:1#AddPortMapping"`)
	header.Set("Content-Type", "text/xml")
	header.Set("Connection", "Close")
	header.Set("Content-Length", "")
	//请求体
	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:AddPortMapping`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:1"`}}

	childList1 := Node{Name: "NewExternalPort", Content: strconv.Itoa(remotePort)}
	childList2 := Node{Name: "NewInternalPort", Content: strconv.Itoa(localPort)}
	childList3 := Node{Name: "NewProtocol", Content: protocol}
	childList4 := Node{Name: "NewEnabled", Content: "1"}
	childList5 := Node{Name: "NewInternalClient", Content: this.upnp.LocalHost}
	childList6 := Node{Name: "NewLeaseDuration", Content: "0"}
	childList7 := Node{Name: "NewPortMappingDescription", Content: "mandela"}
	childList8 := Node{Name: "NewRemoteHost"}
	childTwo.AddChild(childList1)
	childTwo.AddChild(childList2)
	childTwo.AddChild(childList3)
	childTwo.AddChild(childList4)
	childTwo.AddChild(childList5)
	childTwo.AddChild(childList6)
	childTwo.AddChild(childList7)
	childTwo.AddChild(childList8)

	childOne.AddChild(childTwo)
	body.AddChild(childOne)
	bodyStr := body.BuildXML()

	//请求
	request, _ := http.NewRequest("POST", "http://"+this.upnp.Gateway.Host+this.upnp.CtrlUrl,
		strings.NewReader(bodyStr))
	request.Header = header
	request.Header.Set("Content-Length", strconv.Itoa(len([]byte(bodyStr))))
	return request
}

func (this *AddPortMapping) resolve(resultStr string) {
}
