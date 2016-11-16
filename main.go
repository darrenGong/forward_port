package main

import (
	"flag"
	"forward_port/config"
	"log"
	"uframework/log"
	"forward_port/src"
	"net"
)

const (
	LOG_SIZE = 10 * 1024 * 1024
)

var (
	configFile = flag.String("c", "F:\\go-dev\\src\\forward_port\\config\\config.json", "configuration, json format")
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

	return addr[0].String()
}

func main() {
	flag.Parse()

	Config := config.Config{}
	if err := config.ParseConfig(*configFile, &Config);
		err != nil {
		log.Fatalf("Failed to parse config: %s, err: %v\n", *configFile, err)
	}
	uflog.InitLogger(Config.LogPath, Config.LogPrefix, "", LOG_SIZE, "DEBUG")

	addr := GetAddrByInterfaceName(Config.InterfaceName)
	forwardPort.StartServer(addr, Config.Port)
}
