package main

import (
	"fmt"
	"net"
)

type networkElement struct {
	hostname    string
	deviceRole  string
	location    string
	connDetails Connection
}

type Connection struct {
	mgmtIP string
	user   string
	pwd    string
	osType string
}

func main() {
	xe := Connection{net.ParseIP("192.168.124.11").String(), "admin", "admin", "ios-xe"}
	ele := networkElement{"XE-1.cisco.com", "WAN-Edge", "San Jose", xe}

	fmt.Println("Hostname: ", ele.hostname)
	fmt.Println("Management IP: ", ele.connDetails.mgmtIP)
	fmt.Println("Device Type: ", ele.connDetails.osType)
}
