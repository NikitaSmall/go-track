package bitfield

// Bitfield provides a way to check if the data we get is present in our side
// lite coffee-shops stamps
type Bitfield []byte

// HasPiece tells if a bitfield has a particular index set
func (bf Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8
	byteOffset := uint(index % 8)

	return bf[byteIndex]>>(7-byteOffset)&1 != 0
}

// SetPiece sets a piece to a bitfield for a particular index
func (bf Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	byteOffset := uint(index % 8)

	bf[byteIndex] |= 1 << (7 - byteOffset)
}
