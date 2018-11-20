package tcpserver

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/nmrshll/winter-is-coming-backend/gamelogic"
	"github.com/nmrshll/winter-is-coming-backend/utils/cancellable"
	"github.com/nmrshll/winter-is-coming-backend/utils/errors"
)

// StartTCPServer listens to tcp connexions and handles them
// type uint16 ensures port is in correct range
//
// return values help make the server testable (they make it possible to start and stop it)
func StartTCPServer(port uint16) (stop func(), stoppedChan <-chan struct{}, _ error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, nil, err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, nil, err
	}

	fmt.Printf("Listening on tcp:%v ...\n", port)

	quitChan := make(chan struct{})
	stoppedChanBidir := make(chan struct{})
	// listen for connections and handle them
	// also listen for shutdown signal and shutdown server
	go func() {
		for {
			// handle shutdown signal
			select {
			case <-quitChan:
				close(stoppedChanBidir)
				return
			default:
			}

			// handle connections
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			go handleConn(conn)
		}
	}()

	stop = func() {
		defer listener.Close()
		close(quitChan)
	}
	return stop, stoppedChanBidir, nil
}

func writeln(conn net.Conn, format string, args ...interface{}) error {
	if format != "" {
		_, err := conn.Write([]byte(fmt.Sprintf(format, args...) + "\n"))
		if err != nil {
			return errors.Wrap(err, "failed writing out to tcp connection")
		}
	}

	return nil
}

func readConn(conn net.Conn) (out string, _ error) {
	conn.SetReadDeadline(time.Now().Add(30 * time.Hour)) // refresh timeout each time we read from the connection
	readBytes := make([]byte, 128)                       // set maximum request length to 128B to prevent flood based attacks

	readLen, err := conn.Read(readBytes)
	if err != nil {
		return "", errors.Wrap(err, "failed reading from request")
	}
	if readLen == 0 {
		return "", fmt.Errorf("connection already closed by client")
	}

	return strings.TrimSpace(string(readBytes[:readLen])), nil
}

// expectRequestType reads a string from the connection and passes it to the specified handler
func expectRequestType(conn net.Conn, handler requestHandler) error {
	readString, err := readConn(conn)
	if err != nil {
		return err
	}

	response, err := handler(readString)
	if err != nil {
		return errors.Wrap(err, "failed handling request")
	}

	if response != "" {
		return writeln(conn, response)
	}
	return nil
}

// handleConn handles one tcp connection, which means one game
func handleConn(conn net.Conn) {
	defer conn.Close() // close connection before exit

	// the order of game steps is managed here
	// the actual implementation of game logic lives in package gamelogic
	{
		// ask player for name to start a game
		writeln(conn, "Enter player name to start a game:")
		var playerName string
		err := expectRequestType(conn, readPlayerNameHandler(&playerName))
		if err != nil {
			errors.Log(err)
		}

		// create a new game. The zombie starts moving and we write to the tcp socket each time it moves.
		currentGame := gamelogic.NewGame(playerName)
		currentGame.Zombie.SubscribeToMovements(func() {
			writeln(conn, "WALK %s %d %d", currentGame.Zombie.Name(), currentGame.Zombie.X(), currentGame.Zombie.Y())
		})

		// handle user shots
		stopHandlingUserShots := cancellable.RunEvery(10*time.Millisecond, func() error {
			err := expectRequestType(conn, playerShotHandler(currentGame))
			if err != nil {
				errors.Log(err)
				return err // close connection after error
			}
			return nil
		})
		defer stopHandlingUserShots()

		// wait for outcome (zombie gets hit or reaches the wall) and end connection
		for {
			select {
			case <-currentGame.ZombieGotHitChan:
				writeln(conn, "zombie dies. You win.")
				return

			case <-currentGame.ZombieReachedWallChan:
				writeln(conn, "zombie reached the wall. Winter is here. You die.")
				return
			}
		}
	}
}
