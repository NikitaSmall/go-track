package handshake

import (
	"fmt"
	"io"
)

// Handshake details
type Handshake struct {
	PSTR     string
	InfoHash [20]byte
	PeerID   [20]byte
}

// Read parses a handshake from a stream
func Read(r io.Reader) (*Handshake, error) {
	lengthBuf := make([]byte, 1)
	if _, err := io.ReadFull(r, lengthBuf); err != nil {
		return nil, err
	}

	pstrlen := int(lengthBuf[0])
	if pstrlen == 0 {
		return nil, fmt.Errorf("pstrlen cannot be 0")
	}

	handshakeBuf := make([]byte, 48+pstrlen)
	if _, err := io.ReadFull(r, handshakeBuf); err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte

	copy(infoHash[:], handshakeBuf[pstrlen+8:pstrlen+8+20])
	copy(peerID[:], handshakeBuf[pstrlen+8+20:])

	return &Handshake{
		PSTR:     string(handshakeBuf[0:pstrlen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}, nil
}

// Serialize handshake data into byte slice
func (h Handshake) Serialize() []byte {
	pstrlen := len(h.PSTR)
	bufLen := 49 + pstrlen

	buf := make([]byte, bufLen)
	buf[0] = byte(pstrlen)

	copy(buf[1:], h.PSTR)

	// Leave 8 reserved bytes
	copy(buf[1+pstrlen+8:], h.InfoHash[:])
	copy(buf[1+pstrlen+8+20:], h.PeerID[:])

	return buf
}
