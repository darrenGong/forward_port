package main

import (
	"flag"
	"forward_port/config"
	"log"
	"uframework/log"
	"forward_port/src"
	"net"
	"forward_port/rule"
)

const (
	LOG_SIZE = 10 * 1024 * 1024
)

var (
	configFile = flag.String("c", ".\\config\\config.json", "configuration, json format")
	ruleFile = ".\\etc\\Rule_Egi5Th.json"
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

	pRule := new(rule.Rule)
	if err := rule.LoadRule(ruleFile, pRule); err != nil {
		log.Fatalf("Failed to load rule json file[%s]", ruleFile)
	}
	forwardPort.StartServer(addr, Config.Port)
}