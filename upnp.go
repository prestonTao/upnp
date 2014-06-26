package upnp

import (
	// "fmt"
	"log"
	"sync"
)

//对所有的端口进行管理
type MappingPortStruct struct {
	lock         *sync.Mutex
	mappingPorts map[string][][]int
}

//添加一个端口映射记录
//只对映射进行管理
func (this *MappingPortStruct) addMapping(localPort, remotePort int, protocol string) {

	this.lock.Lock()
	defer this.lock.Unlock()
	if this.mappingPorts == nil {
		one := make([]int, 0)
		one = append(one, localPort)
		two := make([]int, 0)
		two = append(two, remotePort)
		portMapping := [][]int{one, two}
		this.mappingPorts = map[string][][]int{protocol: portMapping}
		return
	}
	portMapping := this.mappingPorts[protocol]
	if portMapping == nil {
		one := make([]int, 0)
		one = append(one, localPort)
		two := make([]int, 0)
		two = append(two, remotePort)
		this.mappingPorts[protocol] = [][]int{one, two}
		return
	}
	one := portMapping[0]
	two := portMapping[1]
	one = append(one, localPort)
	two = append(two, remotePort)
	this.mappingPorts[protocol] = [][]int{one, two}
}

//删除一个映射记录
//只对映射进行管理
func (this *MappingPortStruct) delMapping(remotePort int, protocol string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.mappingPorts == nil {
		return
	}
	tmp := MappingPortStruct{lock: new(sync.Mutex)}
	mappings := this.mappingPorts[protocol]
	for i := 0; i < len(mappings[0]); i++ {
		if mappings[1][i] == remotePort {
			//要删除的映射
			break
		}
		tmp.addMapping(mappings[0][i], mappings[1][i], protocol)
	}
	this.mappingPorts = tmp.mappingPorts
}
func (this *MappingPortStruct) GetAllMapping() map[string][][]int {
	return this.mappingPorts
}

type Upnp struct {
	Active             bool              //这个upnp协议是否可用
	LocalHost          string            //本机ip地址
	GetewayInsideIP    string            //局域网网关ip
	GetewayOutsideIP   string            //网关公网ip
	OutsideMappingPort map[string]int    //映射外部端口
	InsideMappingPort  map[string]int    //映射本机端口
	Geteway            *Geteway          //网关信息
	CtrlUrl            string            //控制请求url
	MappingPort        MappingPortStruct //已经添加了的映射 {"TCP":[1990],"UDP":[1991]}
}

func (this *Upnp) init() {
	this.MappingPort = MappingPortStruct{
		lock: new(sync.Mutex),
		// mappingPorts: map[string][][]int{},
	}
	this.LocalHost = GetLocalIntenetIp()
	log.Println("获得本地ip地址:", this.LocalHost)
}

//得到局域网网关ip
func (this *Upnp) searchGateway() {
	if this.LocalHost == "" {
		this.init()
	}
	searchGateway := SearchGateway{upnp: this}
	if searchGateway.Send() {
		log.Println("成功发现网关设备")
		return
	}
	log.Println("未发现网关设备")
}

func (this *Upnp) deviceStatus() {

}

//查看设备描述，得到控制请求url
func (this *Upnp) deviceDesc() {
	if this.GetewayInsideIP == "" {
		this.searchGateway()
		if this.GetewayInsideIP == "" {
			// princ()
		}
	}
	device := DeviceDesc{upnp: this}
	device.Send()
	this.Active = true
	log.Println("获得控制请求url:", this.CtrlUrl)
}

//查看公网ip地址
func (this *Upnp) externalIPAddr() {
	if this.CtrlUrl == "" {
		this.deviceDesc()
	}
	eia := ExternalIPAddress{upnp: this}
	eia.Send()
	log.Println("获得公网ip地址为：", this.GetewayOutsideIP)
}

//添加一个端口映射
func (this *Upnp) AddPortMapping(localPort, remotePort int, protocol string) (b bool) {
	defer func(b bool) {
		err := recover()
		if err != nil {
			log.Println("upnp模块报错了", err)
			b = false
		}
	}(b)
	if this.GetewayOutsideIP == "" {
		this.externalIPAddr()
	}
	addPort := AddPortMapping{upnp: this}
	issuccess := addPort.Send(localPort, remotePort, protocol)
	if issuccess {
		this.MappingPort.addMapping(localPort, remotePort, protocol)
		log.Println("添加一个端口映射：protocol:", protocol, "local:", localPort, "remote:", remotePort)
	} else {
		this.Active = false
		log.Println("添加一个端口映射失败")
	}
	return issuccess
}

func (this *Upnp) DelPortMapping(remotePort int, protocol string) bool {
	delMapping := DelPortMapping{upnp: this}
	issuccess := delMapping.Send(remotePort, protocol)
	if issuccess {
		this.MappingPort.delMapping(remotePort, protocol)
		log.Println("删除了一个端口映射： remote:", remotePort)
	}
	return issuccess
}

//回收端口
func (this *Upnp) Reclaim() {
	mappings := this.MappingPort.GetAllMapping()
	tcpMapping, ok := mappings["TCP"]
	if ok {
		for i := 0; i < len(tcpMapping[0]); i++ {
			this.DelPortMapping(tcpMapping[1][i], "TCP")
		}
	}
	udpMapping, ok := mappings["UDP"]
	if ok {
		for i := 0; i < len(udpMapping[0]); i++ {
			this.DelPortMapping(udpMapping[0][i], "UDP")
		}
	}
}

func (this *Upnp) GetAllMapping() map[string][][]int {
	return this.MappingPort.GetAllMapping()
}
