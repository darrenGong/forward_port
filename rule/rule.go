package rule

import (
	"io/ioutil"
	"uframework/log"
	"errors"
	"encoding/json"
	"strings"
	"strconv"
)

type Rule struct {
	SrcAddr string // 127.0.0.1:80
	DstAddr string // 10.10.1.53:8088
}

var (
	gRuleMap = make(map[uint16]*Rule, 100)
)

func GetAddrPort(laddr string) (string, uint16, error) {
	strArray := strings.Split(laddr, ":")
	if len(strArray[0]) != 2 {
		uflog.ERRORF("Falied to split src addr[%s]\n", laddr)
		return "", 0, errors.New("Falied to split src addr")
	}
	addr := strArray[0]
	port, err := strconv.ParseInt(strArray[1], 10, 16)
	if err != nil {
		uflog.ERRORF("Failed to parse int[%s]\n", laddr)
		return "", 0, errors.New("Failed to parse int")
	}

	return addr, uint16(port), nil
}

func ParseAddr(port uint16, isSrc bool) (string, uint16, error) {
	laddr, _ := Getladdr(port, isSrc)
	if "" == laddr {
		uflog.ERRORF("Port is not exist: %d", port)
		return "", 0, nil
	}

	addr, port, err := GetAddrPort(laddr)
	if err != nil {
		uflog.ERRORF("Failed to parse laddr[%s], err: %v", laddr, err)
		return "", 0, err
	}

	return addr, port, nil
}

func Getladdr(port uint16, isSrc bool) (string, error) {
	laddr := gRuleMap[port].SrcAddr
	if !isSrc {
		laddr = gRuleMap[port].DstAddr
	}

	return laddr, nil
}

func LoadRule(ruleFile string, rule *Rule) error {
	bytes, err := ioutil.ReadFile(ruleFile)
	if err != nil {
		uflog.ERRORF("Failed to read rule file[%s], err: %v", ruleFile, err)
		return errors.New("Failed to read rule file")
	}

	if err := json.Unmarshal(bytes, rule); err != nil {
		uflog.ERRORF("Failed to unmarshal json file[%s], err: %v\n", ruleFile, err)
		return errors.New("Failed to unmarshal json file")
	}
	addr, port, err := GetAddrPort(rule.SrcAddr)
	if err != nil {
		uflog.ERRORF("Failed to get src port[%s], err: %v", addr, err)
		return err
	}
	gRuleMap[port] = rule

	return nil
}