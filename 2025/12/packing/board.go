package packing

type Board struct {
	Width       int
	Height      int
	PieceCounts map[int]int // piece.Id[count]
}

func NewBoard(width, height int) *Board {
	return &Board{
		Width:       width,
		Height:      height,
		PieceCounts: make(map[int]int),
	}
}

func (b *Board) AddPiece(piece *Piece, count int) {
	b.PieceCounts[piece.Id] = count
}
