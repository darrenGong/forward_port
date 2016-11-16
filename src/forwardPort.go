package forwardPort

import (
	"net"
	"time"
	"uframework/log"
	"errors"
	"io"
)

var (
	gClientConnMap map[uint16]ForwardPort
)

type ForwardPort struct {
	SrcAddr string
	SrcPort uint16
	DstAddr string
	DstPort uint16
	SrcConn net.Conn
	DstConn net.Conn
	Timeout time.Duration

	QuitChan chan int
}

func CopyBytes(dstConn, srcConn net.Conn) error {
	lenByte, err := io.Copy(dstConn, srcConn)
	if err != nil {
		uflog.Debugf("Send error from src[%s] to dst[%s]\n",
			srcConn.LocalAddr(), dstConn.RemoteAddr())
		return err
	}
	uflog.Debugf("Send %d bytes from src[%s] to dst[%s]\n",
		lenByte, srcConn.LocalAddr(), dstConn.RemoteAddr())

	return nil
}

func (fp *ForwardPort) ForwardWork() error {
	if fp.SrcConn == nil {
		uflog.ERROR("Invalid conn[SrcConn is nil]\n")
		return errors.New("Invalid conn")
	}

	for {
		if err := CopyBytes(fp.DstConn, fp.SrcConn); err != nil {
			uflog.DEBUGF("Connection have closed src -> dst")
			fp.CloseConn()
			return err
		}

		if err := CopyBytes(fp.SrcConn, fp.DstConn); err != nil {
			uflog.DEBUGF("Connection have closed dst -> src")
			fp.CloseConn()
			return err
		}
	}
	return nil
}

func (fp *ForwardPort)CloseConn() {
	fp.DstConn.Close()
	fp.SrcConn.Close()

	fp.QuitChan <- 0
}