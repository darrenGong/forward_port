package forwardPort

import (
	"errors"
	"io"
	"net"
	"time"
	"uframework/log"
)

var (
	gClientConnMap map[uint16]ForwardPort
)

type ForwardPort struct {
	SrcAddr  string
	SrcPort  uint16
	DstAddr  string
	DstPort  uint16
	SrcConn  net.Conn
	DstConn  net.Conn
	Timeout  time.Duration

	QuitChan chan int
}

func (fp *ForwardPort)CopyBytes(dstConn, srcConn net.Conn) error {
	for {
		lenByte, err := io.Copy(dstConn, srcConn)
		if err != nil {
			uflog.DEBUGF("Send error from src[%s] to dst[%s]\n",
				srcConn.LocalAddr(), dstConn.RemoteAddr())
			return err
		}
		uflog.DEBUGF("Send %d bytes from src[%s] to dst[%s]\n",
			lenByte, srcConn.LocalAddr(), dstConn.RemoteAddr())
		fp.CloseConn()
		return nil
	}

	return nil
}

func (fp *ForwardPort) ForwardWork() error {
	if nil == fp.SrcConn ||
		nil == fp.DstConn {
		uflog.ERROR("Invalid conn[Conn is nil]\n")
		return errors.New("Invalid conn")
	}

	go fp.CopyBytes(fp.DstConn, fp.SrcConn)
	go fp.CopyBytes(fp.SrcConn, fp.DstConn)

	return nil
}

func (fp *ForwardPort) CloseConn() {
	if fp.DstConn != nil {
		fp.DstConn.Close()
	}

	if fp.SrcConn != nil {
		fp.SrcConn.Close()
	}

	fp.QuitChan <- 0
}
