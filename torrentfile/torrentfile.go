package torrentfile

import (
	"io"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
)

const hashLen = 20 // Length of SHA-1 hash

// TorrentFile holds the data from BencodeTorrent struct in a flat and organized way
type TorrentFile struct {
	Announce string
	InfoHash [hashLen]byte

	PieceLength int
	PieceHashes [][hashLen]byte

	Name   string
	Length int
}

// BencodeTorrent holds parased torrent file data
type BencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

// Open reads and parces the torrent file returning the parsed torrent info
func Open(r io.Reader) (*BencodeTorrent, error) {
	var torrent BencodeTorrent
	if err := bencode.Unmarshal(r, &torrent); err != nil {
		return nil, err
	}

	return &torrent, nil
}

// ToTorrentFile transforms bencode torrent to organized structure
func (bt BencodeTorrent) ToTorrentFile() (*TorrentFile, error) {
	infoHash, err := bt.Info.hash()
	if err != nil {
		return nil, err
	}

	piecesHash, err := bt.Info.splitPieceHashes()
	if err != nil {
		return nil, err
	}

	return &TorrentFile{
		Announce:    bt.Announce,
		InfoHash:    infoHash,
		PieceHashes: piecesHash,
		PieceLength: bt.Info.PieceLength,
		Name:        bt.Info.Name,
		Length:      bt.Info.Length,
	}, nil
}

// BuildTrackerURL creates an URL ready to request the bittorrent central
func (tf TorrentFile) BuildTrackerURL(peerID [hashLen]byte, port uint16) (string, error) {
	base, err := url.Parse(tf.Announce)
	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{string(tf.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(tf.Length)},
	}
	base.RawQuery = params.Encode()

	return base.String(), nil
}
