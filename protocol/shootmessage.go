package protocol

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nmrshll/winter-is-coming-backend/utils/errors"
)

// ShootMessage represents the payload sent by the server when a player shoots at coordinates on the board
type ShootMessage struct {
	X int
	Y int
}

// Serialize converts a ShootMessage into its byte representation to write to the tcp connection
func (msg *ShootMessage) Serialize() []byte {
	return []byte(fmt.Sprintf("%v %v", msg.X, msg.Y))
}

// Parse parses a ShootMessage from the payload string sent by the client
func (msg *ShootMessage) Parse(input string) error {
	trimmedSplit := strings.Split(strings.TrimSpace(input), " ")
	if len(trimmedSplit) != 2 {
		return fmt.Errorf("failed parsing input into 2 ints")
	}

	shotX, err := strconv.Atoi(trimmedSplit[0])
	if err != nil {
		return errors.Wrap(err, "failed parsing input into int shotX")
	}
	shotY, err := strconv.Atoi(trimmedSplit[1])
	if err != nil {
		return errors.Wrap(err, "failed parsing input into int shotY")
	}

	*msg = ShootMessage{
		X: shotX,
		Y: shotY,
	}
	return nil
}
