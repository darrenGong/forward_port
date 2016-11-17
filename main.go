package main

import (
	"flag"
	"forward_port/config"
	"log"
	"uframework/log"
	"forward_port/src"
	"net"
	"forward_port/rule"
	"strings"
)

const (
	LOG_SIZE = 10 * 1024 * 1024
)

var (
	configFile = flag.String("c", "F:\\go-dev\\src\\forward_port\\config\\config.json", "configuration, json format")
	ruleFile = "F:\\go-dev\\src\\forward_port\\etc\\Rule_Egi5Th.json"
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

func main() {
	flag.Parse()

	Config := config.Config{}
	if err := config.ParseConfig(*configFile, &Config);
		err != nil {
		log.Fatalf("Failed to parse config: %s, err: %v\n", *configFile, err)
	}
	uflog.InitLogger(Config.LogPath, Config.LogPrefix, "", LOG_SIZE, "DEBUG")
	//addr := GetAddrByInterfaceName(Config.InterfaceName)
	addr := "127.0.0.1"

	pRule := new(rule.Rule)
	if err := rule.LoadRule(ruleFile, pRule); err != nil {
		log.Fatalf("Failed to load rule json file[%s], err %v", ruleFile, err)
	}
	_, port, _ := rule.GetAddrPort(pRule.SrcAddr)
	forwardPort.StartServer(addr, port)
}
