package main

import (
	"../../upnp"
	"bufio"
	"fmt"
	"os"
)

func main() {
	StartUP()
}

func StartUP() {
	mapping := new(upnp.Upnp)
	if ok := mapping.AddPortMapping(55789, 55789, "TCP"); ok {
		fmt.Println("端口映射成功")
	} else {
		fmt.Println("不支持upnp协议")
	}

	running := true
	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		command := string(data)
		switch command {
		case "find":
		case "quit":
			mapping.Reclaim()
			running = false
		case "see":
		case "self":
		case "cap":
		case "odp":
		case "cdp":
		case "dump":
		}
	}
}
