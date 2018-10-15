// Package protocol contains the definitions of the messages that are passed over the tcp socket
// this package is meant to be used both for the server and the client,
// for serializing messages on one end and parsing them on the other end
package protocol

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/davecgh/go-spew/spew"

	"github.com/nmrshll/winter-is-coming-backend/utils/errors"
)

// WalkMessage represents the payload sent by the server when a zombie walks on the board
type WalkMessage struct {
	ZombieName string
	X          int
	Y          int
}

var regexWalkMessage = regexp.MustCompile(`(?m)WALK (\S+) ([0-9]+) ([0-9]+)`)

// Parse parses a WalkMessage from the payload string sent by the server
func (msg *WalkMessage) Parse(input string) error {
	matches := regexWalkMessage.FindStringSubmatch(input)
	if len(matches) != 4 {
		return fmt.Errorf("matches expected to be length 4: %v", spew.Sdump(matches))
	}

	xInt, err := strconv.Atoi(matches[2])
	if err != nil {
		return errors.Wrap(err, "failed parsing input into int shotX")
	}
	yInt, err := strconv.Atoi(matches[3])
	if err != nil {
		return errors.Wrap(err, "failed parsing input into int shotY")
	}

	*msg = WalkMessage{
		ZombieName: matches[1],
		X:          xInt,
		Y:          yInt,
	}

	return nil
}
