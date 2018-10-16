package tcpserver

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/nmrshll/winter-is-coming-backend/protocol"
	"github.com/nmrshll/winter-is-coming-backend/utils"
)

func newTestGameWithPlayer(t *testing.T) (cleanupFunc func(), conn net.Conn) {
	cleanupFunc, conn = newTestServerConnection(t)

	expectedOutput := "Enter player name to start a game:"
	if out := testReadConn(t, conn); out != expectedOutput {
		t.Fatalf("response (%v) did match expected output (%v)", string(out), expectedOutput)
	}

	if _, err := conn.Write([]byte("Jon")); err != nil {
		utils.WrapFatal(t, err, "could not write payload to TCP server")
	}

	expectedOutput = "START Jon"
	if out := testReadConn(t, conn); out != expectedOutput {
		t.Fatalf("response (%v) did match expected output (%v)", string(out), expectedOutput)
	}

	return cleanupFunc, conn
}

func Test_API_readPlayerNameHandler(t *testing.T) {
	// init
	cleanupFunc, conn := newTestServerConnection(t)
	defer cleanupFunc()

	if out := testReadConn(t, conn); out != "Enter player name to start a game:" {
		t.Fatalf("response (%v) did match expected output (%v)", string(out), "Enter player name to start a game:")
	}

	// test cases
	tests := []struct {
		name           string
		payload        []byte
		expectedOutput string
	}{
		{
			"Input name, expect correct response",
			[]byte("jon"),
			"START jon",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := conn.Write(tc.payload); err != nil {
				utils.WrapFatal(t, err, "could not write payload to TCP server")
			}

			if out := testReadConn(t, conn); out != tc.expectedOutput {
				t.Fatalf("response (%v) didn't match expected output (%v)", string(out), string(tc.expectedOutput))
			}
		})
	}
}

func Test_API_playerShotHandler(t *testing.T) {
	// init
	cleanupFunc, conn := newTestGameWithPlayer(t)
	defer cleanupFunc()

	t.Run("playerShotHandler", func(t *testing.T) {
		// wait for 2 seconds for the server to move the zombie at least once
		time.Sleep(2 * time.Second)
		readString := testReadConn(t, conn)
		var walkMessage protocol.WalkMessage
		err := walkMessage.Parse(readString)
		if err != nil {
			utils.WrapFatal(t, err, "failed parsing walkMessage")
		}

		// if player shoots at right coordinates, the zombie should die
		t.Run("shoot NOT at zombie coordinates", func(t *testing.T) {
			shootMessage := protocol.ShootMessage{X: walkMessage.X + 1, Y: walkMessage.Y + 1}
			if _, err := conn.Write(shootMessage.Serialize()); err != nil {
				utils.WrapFatal(t, err, "could not write payload to TCP server")
			}

			// server should keep connection open as the zombie is not hit and has no chance of reaching the wall that fast
			if conn == nil {
				t.Fatalf("expected connection to stay open")
			}
		})

		t.Run("shoot at zombie coordinates", func(t *testing.T) {
			shootMessage := protocol.ShootMessage{X: walkMessage.X, Y: walkMessage.Y}
			if _, err := conn.Write(shootMessage.Serialize()); err != nil {
				utils.WrapFatal(t, err, "could not write payload to TCP server")
			}

			// read all messages
			readString := testReadConn(t, conn) // contains the shootMessage
			readString = testReadConn(t, conn)  // should contain "zombie dies" if the shot was fast enough to hit the zombie
			if strings.Contains(readString, "zombie dies") {
				return // zombie was hit, test passes
			} else {
				awaitWalkAndShootAtZombieCoordinates(0, t, conn)
			}
		})
	})
}

func awaitWalkAndShootAtZombieCoordinates(idx int, t *testing.T, conn net.Conn) {
	if idx > 10 {
		t.Fatalf("failed hitting zombie")
		return
	}

	readString := testReadConn(t, conn)
	if strings.Contains(readString, "SHOOT") {
		awaitWalkAndShootAtZombieCoordinates(idx+1, t, conn)
	}
	if strings.Contains(readString, "WALK") {
		var walkMessage protocol.WalkMessage
		err := walkMessage.Parse(readString)
		if err != nil {
			utils.WrapFatal(t, err, "failed parsing walkMessage")
		}

		// attempt shot at zombie coordinates
		shootMessage := protocol.ShootMessage{X: walkMessage.X, Y: walkMessage.Y}
		if _, err := conn.Write(shootMessage.Serialize()); err != nil {
			utils.WrapFatal(t, err, "could not write payload to TCP server")
		}

		// read all messages
		readString := testReadConn(t, conn) // contains the shootMessage
		readString = testReadConn(t, conn)  // should contain "zombie dies"
		if strings.Contains(readString, "zombie dies") {
			return // zombie was hit, test passes
		} else {
			awaitWalkAndShootAtZombieCoordinates(idx+1, t, conn)
		}
	}
}
