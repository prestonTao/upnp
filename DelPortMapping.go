package upnp

import (
	"io/ioutil"
	// "log"
	"net/http"
	"strconv"
	"strings"
)

type DelPortMapping struct {
	upnp *Upnp
}

func (this *DelPortMapping) Send(remotePort int, protocol string) bool {
	request := this.buildRequest(remotePort, protocol)
	response, _ := http.DefaultClient.Do(request)
	resultBody, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode == 200 {
		// log.Println(string(resultBody))
		this.resolve(string(resultBody))
		return true
	}
	return false
}
func (this *DelPortMapping) buildRequest(remotePort int, protocol string) *http.Request {
	//请求头
	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:1#DeletePortMapping"`)
	header.Set("Content-Type", "text/xml")
	header.Set("Connection", "Close")
	header.Set("Content-Length", "")

	//请求体
	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:DeletePortMapping`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:1"`}}
	childList1 := Node{Name: "NewExternalPort", Content: strconv.Itoa(remotePort)}
	childList2 := Node{Name: "NewProtocol", Content: protocol}
	childList3 := Node{Name: "NewRemoteHost"}
	childTwo.AddChild(childList1)
	childTwo.AddChild(childList2)
	childTwo.AddChild(childList3)
	childOne.AddChild(childTwo)
	body.AddChild(childOne)
	bodyStr := body.BuildXML()

	//请求
	request, _ := http.NewRequest("POST", "http://"+this.upnp.Geteway.Host+this.upnp.CtrlUrl,
		strings.NewReader(bodyStr))
	request.Header = header
	request.Header.Set("Content-Length", strconv.Itoa(len([]byte(bodyStr))))
	return request
}

func (this *DelPortMapping) resolve(resultStr string) {
}
