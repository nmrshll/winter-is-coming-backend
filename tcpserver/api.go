// Package tcpserver contains all request handlers, which define the actions the user can do through the communication channel
package tcpserver

import (
	"fmt"
	"strings"

	"github.com/nmrshll/winter-is-coming-backend/gamelogic"
	"github.com/nmrshll/winter-is-coming-backend/protocol"
	"github.com/nmrshll/winter-is-coming-backend/utils/errors"
)

// requestHandler is any function that handles one payload sent by the client on the tcp socket (= one read by the server)
// both input and output are in string type for this function
type requestHandler func(body string) (response string, _ error)

func readPlayerNameHandler(playerName *string) requestHandler {
	return func(body string) (response string, _ error) {
		*playerName = strings.TrimSpace(body)
		return fmt.Sprintf("START %s", *playerName), nil
	}
}

func playerShotHandler(currentGame *gamelogic.Game) requestHandler {
	return func(body string) (response string, _ error) {
		// parse and validates coordinates
		var shootMessage protocol.ShootMessage
		err := shootMessage.Parse(body)
		if err != nil {
			return "", errors.Wrap(err, "failed parsing shootMessage")
		}

		// shoot on the board
		currentGame.Shoot(shootMessage.X, shootMessage.Y)

		return fmt.Sprintf("SHOOT %v %v %v", currentGame.PlayerName, shootMessage.X, shootMessage.Y), nil
	}
}
