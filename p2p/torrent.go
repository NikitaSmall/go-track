package p2p

import "github.com/nikitasmall/go-track/peers"

// Torrent holds everything to download
type Torrent struct {
	PieceHashes [][20]byte
	Peers       []peers.Peer
	PieceLength int
	Length      int
}

type pieceWork struct {
	Index       int
	PieceHash   [20]byte
	PieceLength int
}

type pieceResult struct {
	Index int
	Buf   []byte
}

func (t Torrent) calculateBoundariesForPiece(index int) (int, int) {
	begin := index * t.PieceLength
	end := begin + t.PieceLength

	if end > t.Length {
		end = t.Length
	}

	return begin, end
}

func (t Torrent) calculatePieceSize(index int) int {
	begin, end := t.calculateBoundariesForPiece(index)

	return end - begin
}

// Download uploads the file via the torrent and stores it in memory
func (t Torrent) Download() ([]byte, error) {
	workQueue := make(chan pieceWork, len(t.PieceHashes))
	results := make(chan pieceResult)

	for i, pieceHash := range t.PieceHashes {
		length := t.calculatePieceSize(i)
		workQueue <- pieceWork{
			Index:       i,
			PieceHash:   pieceHash,
			PieceLength: length,
		}
	}

	for _, peer := range t.Peers {
		go t.startDownloadWorker(peer, workQueue, results)
	}

	buf := make([]byte, t.Length)
	var donePieces int

	for donePieces < len(t.PieceHashes) {
		res := <-results
		begin, end := t.calculateBoundariesForPiece(res.Index)
		copy(buf[begin:end], res.Buf)
	}
	close(workQueue)

	return buf, nil
}
