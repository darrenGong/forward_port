package main

import (
	"flag"
	"forward_port/config"
	"forward_port/rule"
	"forward_port/src"
	"log"
	"uframework/log"
)

const (
	LOG_SIZE = 10 * 1024 * 1024
)

var (
	configFile = flag.String("c", "config/config.json", "configuration, json format")
	ruleFile   = "etc/Rule_Egi5Th.json"
)

func main() {
	flag.Parse()

	Config := config.Config{}
	if err := config.ParseConfig(*configFile, &Config); err != nil {
		log.Fatalf("Failed to parse config: %s, err: %v\n", *configFile, err)
	}
	uflog.InitLogger(Config.LogPath, Config.LogPrefix, "", LOG_SIZE, "DEBUG")
	addr := forwardPort.GetAddrByInterfaceName(Config.InterfaceName)

	pRule := new(rule.Rule)
	if err := rule.LoadRule(ruleFile, pRule); err != nil {
		log.Fatalf("Failed to load rule json file[%s], err %v", ruleFile, err)
	}
	_, port, _ := rule.GetAddrPort(pRule.SrcAddr)
	forwardPort.StartServer(addr, port)
}
