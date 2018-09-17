package tophost

import (
	"fmt"
	"strings"
)

// Input defines a set of input variables to be considered
type Input struct {
	IPSrc     string
	IPDst     string
	Peer      string
	Direction string
	Packets   uint64
	Bytes     uint64
}

func (input *Input) getIdentifier() string {
	return fmt.Sprintf("%s~~%s~~%s", input.IPSrc, input.IPDst, input.Peer)
}

func splitIdentifier(identifier string) Input {
	parts := strings.Split(identifier, "~~")
	if len(parts) != 3 {
		return Input{}
	}
	return Input{
		IPSrc: parts[0],
		IPDst: parts[1],
		Peer:  parts[2],
	}
}
