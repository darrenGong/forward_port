package forwardPort

import (
	"log"
	"net"
	"strings"
)

func GetAddrByInterfaceName(interfaceName string) string {
	Interface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		log.Fatalf("Failed to parse interface name[%s]\n", interfaceName)
	}

	addr, err := Interface.Addrs()
	if err != nil {
		log.Fatalf("Failed to parse addr, err: %v\n", err)
	}

	return strings.Split(addr[0].String(), "/")[0]
}
