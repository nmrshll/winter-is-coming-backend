package tcpserver

import (
	"net"
	"testing"

	"github.com/nmrshll/winter-is-coming-backend/utils"
)

func newTestServerConnection(t *testing.T) (cleanupFunc func(), conn net.Conn) {
	stopServerFunc, _, err := StartTCPServer(7777)
	if err != nil {
		utils.WrapFatal(t, err, "failed starting tcp server")
	}

	conn, err = net.Dial("tcp", ":7777")
	if err != nil {
		utils.WrapFatal(t, err, "could not connect to TCP server")
	}

	return func() {
		conn.Close()
		stopServerFunc()
	}, conn
}

func testReadConn(t *testing.T, conn net.Conn) (out string) {
	out, err := readConn(conn)
	if err != nil {
		utils.WrapFatal(t, err, "failed reading from connection")
	}
	return out
}

func Test_API_StartTCPServer(t *testing.T) {
	stopFunc, _, err := StartTCPServer(7777)
	if err != nil {
		utils.WrapFatal(t, err, "failed starting tcp server")
	}
	defer stopFunc()
}
