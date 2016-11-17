package forwardPort

import (
	"errors"
	"fmt"
	"forward_port/rule"
	"net"
	"sync"
	"time"
	"uframework/log"
)

var (
	gServerForwardPortMap map[uint16]*TcpServer
	gForwardPortMutex     sync.Mutex
	gStopWait             sync.WaitGroup
)

type TcpServer struct {
	Listener    net.Listener
	ForwardPort *ForwardPort
}

func Getladdr(addr string, port uint16) (string, error) {
	if "" == addr || 0 == port {
		uflog.ERRORF("Addr/Port is Empty [%s:%d]\n", addr, port)
		return "", errors.New("Param is invalid")
	}

	laddr := fmt.Sprintf("%s:%d", addr, port)
	return laddr, nil
}

func GetForwardPort(port uint16) (*ForwardPort, error) {
	srcAddr, srcPort, err := rule.ParseAddr(port, true)
	if err != nil {
		return nil, err
	}

	dstAddr, dstPort, err := rule.ParseAddr(port, false)
	if err != nil {
		return nil, err
	}

	return &ForwardPort{
		SrcAddr:  srcAddr,
		SrcPort:  srcPort,
		DstAddr:  dstAddr,
		DstPort:  dstPort,
		SrcConn:  nil,
		DstConn:  nil,
		Timeout:  2 * time.Second,
		QuitChan: make(chan int),
	}, nil
}

func StartServer(addr string, port uint16) error {
	if "" == addr || 0 == port {
		uflog.ERRORF("Addr/Port is Empty [%s:%d]\n", addr, port)
		return errors.New("Param is invalid")
	}
	uflog.INFOF("Start server, listen %s:%d\n", addr, port)

	laddr, _ := Getladdr(addr, port)
	localListener, err := net.Listen("tcp", laddr)
	if err != nil {
		uflog.ERRORF("Failed to listen [%s]\n", laddr)
		return errors.New("Failed to listen")
	}

	var chanConn chan net.Conn
	var tcpServer TcpServer
	forwardPort, _ := GetForwardPort(port)
	tcpServer.Listener = localListener
	go AcceptServer(localListener, chanConn)

	for {
		select {
		case srcConn := <-chanConn:
			uflog.INFOF("New connection [LocalAddr:%s]", srcConn.RemoteAddr().String())
			srcConn.LocalAddr().String()
			forwardPort, _ = GetForwardPort(port)
			forwardPort.SrcConn = srcConn

			dstLaddr, _ := rule.Getladdr(port, false)
			if "" == dstLaddr {
				uflog.ERRORF("Dst laddr is not exit: %d", port)
				continue
			}
			dstConn, err := net.DialTimeout("tcp", laddr, 5*time.Second)
			if err != nil {
				uflog.ERRORF("Connection failed to dst: %s", laddr)
				continue
			}
			forwardPort.DstConn = dstConn
			tcpServer.ForwardPort = forwardPort
			AddServer(port, &tcpServer)

			go forwardPort.ForwardWork()
		case <-forwardPort.QuitChan:
			forwardPort.CloseConn()
			DelServer(port, &tcpServer)
		}

	}

	return nil
}

func AcceptServer(localListener net.Listener, chanConn chan<- net.Conn) error {
	for {
		srcConn, err := localListener.Accept()
		if err != nil {
			uflog.ERRORF("Failed to accept connection, err :%v\n", err)
			return err
		}
		chanConn <- srcConn
	}
}

func AddServer(port uint16, server *TcpServer) {
	gForwardPortMutex.Lock()
	defer gForwardPortMutex.Unlock()

	gServerForwardPortMap[port] = server
	gStopWait.Add(1)
}

func DelServer(port uint16, server *TcpServer) {
	gForwardPortMutex.Lock()
	defer gForwardPortMutex.Unlock()

	delete(gServerForwardPortMap, port)
	gStopWait.Done()
}
