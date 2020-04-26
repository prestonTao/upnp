package upnp

import (
	// "fmt"
	"errors"
	"log"
	"sync"
)

/*
 * 得到网关
 */

//对所有的端口进行管理
type MappingPortStruct struct {
	lock         *sync.Mutex
	mappingPorts []*MappingPort
}

func (this *MappingPortStruct) findPort(remotePort uint16, protocol Protocol) int {
	for index, port := range this.mappingPorts {
		if port.Protocol == protocol && port.RemotePort == remotePort {
			return index
		}
	}
	return -1
}

//添加一个端口映射记录
//只对映射进行管理
func (this *MappingPortStruct) addMapping(localPort, remotePort uint16, protocol Protocol, describe string) {

	this.lock.Lock()
	defer this.lock.Unlock()
	if this.mappingPorts == nil {
		this.mappingPorts = make([]*MappingPort, 1)
		this.mappingPorts[0] = NewMappingPort(protocol, describe, localPort, remotePort)
	} else {
		this.mappingPorts = append(this.mappingPorts, NewMappingPort(protocol, describe, localPort, remotePort))
	}
	return
}

//删除一个映射记录
//只对映射进行管理
func (this *MappingPortStruct) delMapping(remotePort uint16, protocol Protocol) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.mappingPorts == nil {
		return
	}
	if index := this.findPort(remotePort, protocol); index > 0 {
		this.mappingPorts = append(this.mappingPorts[:index], this.mappingPorts[index+1:]...)
	}
}
func (this *MappingPortStruct) GetAllMapping() []*MappingPort {
	return this.mappingPorts
}

type Upnp struct {
	Active             bool              //这个upnp协议是否可用
	LocalHost          string            //本机ip地址
	GatewayInsideIP    string            //局域网网关ip
	GatewayOutsideIP   string            //网关公网ip
	OutsideMappingPort map[string]int    //映射外部端口
	InsideMappingPort  map[string]int    //映射本机端口
	Gateway            *Gateway          //网关信息
	CtrlUrl            string            //控制请求url
	MappingPort        MappingPortStruct //已经添加了的映射 {"TCP":[1990],"UDP":[1991]}
}

//得到本地联网的ip地址
//得到局域网网关ip
func (this *Upnp) SearchGateway() (err error) {
	defer func(err error) {
		if errTemp := recover(); errTemp != nil {
			log.Println("upnp模块报错了", errTemp)
			err = errTemp.(error)
		}
	}(err)

	if this.LocalHost == "" {
		this.MappingPort = MappingPortStruct{
			lock: new(sync.Mutex),
		}
		this.LocalHost = GetLocalIntenetIp()
	}
	searchGateway := SearchGateway{upnp: this}
	if searchGateway.Send() {
		return nil
	}
	return errors.New("未发现网关设备")
}

func (this *Upnp) deviceStatus() {

}

//查看设备描述，得到控制请求url
func (this *Upnp) deviceDesc() (err error) {
	if this.GatewayInsideIP == "" {
		if err := this.SearchGateway(); err != nil {
			return err
		}
	}
	device := DeviceDesc{upnp: this}
	device.Send()
	this.Active = true
	// log.Println("获得控制请求url:", this.CtrlUrl)
	return
}

//查看公网ip地址
func (this *Upnp) ExternalIPAddr() (err error) {
	if this.CtrlUrl == "" {
		if err := this.deviceDesc(); err != nil {
			return err
		}
	}
	eia := ExternalIPAddress{upnp: this}
	eia.Send()
	return nil
	// log.Println("获得公网ip地址为：", this.GatewayOutsideIP)
}

//添加一个端口映射
func (this *Upnp) AddPortMapping(localPort, remotePort uint16, protocol Protocol, describe string) (err error) {
	defer func(err error) {
		if errTemp := recover(); errTemp != nil {
			log.Println("upnp模块报错了", errTemp)
			err = errTemp.(error)
		}
	}(err)
	if this.GatewayOutsideIP == "" {
		if err := this.ExternalIPAddr(); err != nil {
			return err
		}
	}
	addPort := AddPortMapping{upnp: this}
	if issuccess := addPort.Send(localPort, remotePort, protocol, describe); issuccess {
		this.MappingPort.addMapping(localPort, remotePort, protocol, describe)
		// log.Println("添加一个端口映射：protocol:", protocol, "local:", localPort, "remote:", remotePort)
		return nil
	} else {
		this.Active = false
		// log.Println("添加一个端口映射失败")
		return errors.New("添加一个端口映射失败")
	}
}

func (this *Upnp) DelPortMapping(remotePort uint16, protocol Protocol) bool {
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
	if mappings == nil {
		return
	}
	for _, port := range mappings {
		this.DelPortMapping(port.RemotePort, port.Protocol)
	}
}

func (this *Upnp) GetAllMapping() []*MappingPort {
	return this.MappingPort.GetAllMapping()
}
