package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"

	bencode "github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Name string `bencode:"name"`

	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`

	Length int `bencode:"length"`
}

func (i bencodeInfo) hash() ([hashLen]byte, error) {
	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, i); err != nil {
		return [hashLen]byte{}, err
	}

	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfo) splitPieceHashes() ([][hashLen]byte, error) {
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][hashLen]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}
