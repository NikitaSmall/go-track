package peers

import (
	"encoding/binary"
	"fmt"
	"net"
)

const peerSize = 6 // 4 for IP, 2 for port

// TrackerResponse holds the data from central tracker response
type TrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    []byte `bencode:"peers"`
}

// Peer holds single peer data
type Peer struct {
	IP   net.IP
	Port uint16
}

// Unmarshal parses byte slice into actual Peer slice
func Unmarshal(peerData []byte) ([]Peer, error) {
	if len(peerData)%peerSize != 0 {
		err := fmt.Errorf("Received malformed peers")
		return nil, err
	}

	numPeers := len(peerData) / peerSize

	peers := make([]Peer, numPeers)
	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peerData[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peerData[offset+4 : offset+6])
	}

	return peers, nil
}
