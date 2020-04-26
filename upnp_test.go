package upnp

import (
	"fmt"
	"testing"
)

func TestUpnp_SearchGateway(t *testing.T) {
	u := new(Upnp)
	err := u.SearchGateway()
	if err != nil {
		t.Error(err)
	}
}
func TestUpnp_AddPortMapping(t *testing.T) {
	u := new(Upnp)
	defer u.Reclaim()
	if err := u.AddPortMapping(8088, 8088, TCP, "test"); err != nil {
		t.Error("err:", err)
	}
	maps := u.GetAllMapping()
	if len(maps) > 1 {
		t.Error("add wrong", maps)
	}
	if maps[0].Protocol != TCP || maps[0].LocalPort != 8088 || maps[0].RemotePort != 8088 || maps[0].Describe != "test" {
		t.Error("not correct")
	}
	if err := u.AddPortMapping(8089, 8089, TCP, "test"); err != nil {
		t.Error("err:", err)
	}
	maps = u.GetAllMapping()
	if len(maps) != 2 {
		t.Error("add wrong", maps)
	}
	if maps[1].Protocol != TCP || maps[1].LocalPort != 8089 || maps[1].RemotePort != 8089 || maps[1].Describe != "test" {
		t.Error("not correct")
	}
	// time.Sleep(100*time.Second)
}
func TestUpnp_Reclaim(t *testing.T) {
	u := new(Upnp)
	u.Reclaim()
}

func TestUpnp_GetAllMapping(t *testing.T) {
	u := new(Upnp)
	defer u.Reclaim()
	if err := u.AddPortMapping(8089, 8089, TCP, "test"); err != nil {
		t.Error("err:", err)
	}
	maps := u.GetAllMapping()
	if len(maps) != 1 {
		t.Error("get Error")
	}
	for _, rule := range maps {
		fmt.Println("Protocol:", rule.Protocol, "| Describe:", rule.Describe, "| Port:", rule.RemotePort, "->", rule.LocalPort)
	}
}
