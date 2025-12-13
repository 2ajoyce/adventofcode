package packing

import "fmt"

type Problem struct {
	Pieces []*Piece
	Boards []*Board
}

func NewProblem(pieces []*Piece, boards []*Board) *Problem {
	// Validate that every piece used in boards exists in pieces
	pieceCounts := make(map[int]bool)
	for _, p := range pieces {
		pieceCounts[p.Id] = true
	}
	for _, b := range boards {
		for pid := range b.PieceCounts {
			if !pieceCounts[pid] {
				panic("board references unknown piece Id")
			}
		}
	}

	return &Problem{
		Pieces: pieces,
		Boards: boards,
	}
}

func (p *Problem) EvaluateProblem() []EvaluationResult {
	results := make([]EvaluationResult, 0, len(p.Boards))

	for _, board := range p.Boards {
		result := p.evaluateBoard(board)
		results = append(results, result)
	}
	return results
}

func (p *Problem) evaluateBoard(b *Board) EvaluationResult {
	ok, reason := p.AreaFilter(*b)
	if !ok {
		return EvaluationResult{
			Board:  *b,
			CanFit: false,
			Reason: reason,
		}
	}

	canSolve := p.SolveWithBacktracking(b)

	return EvaluationResult{
		Board:  *b,
		CanFit: canSolve,
		Reason: "",
	}
}

// AreaFilter checks that total piece area <= board area.
func (p *Problem) AreaFilter(board Board) (bool, string) {
	boardArea := board.Width * board.Height

	piecesByID := make(map[int]*Piece, len(p.Pieces))
	for _, p := range p.Pieces {
		piecesByID[p.Id] = p
	}

	totalPieceArea := 0
	for pieceID, count := range board.PieceCounts {
		piece, ok := piecesByID[pieceID]
		if !ok {
			return false, fmt.Sprintf("area filter: unknown piece Id %d in board", pieceID)
		}
		if count < 0 {
			return false, fmt.Sprintf("area filter: negative count %d for piece %d", count, pieceID)
		}
		totalPieceArea += piece.Area * count
	}

	if totalPieceArea > boardArea {
		return false, fmt.Sprintf(
			"area filter: pieces require %d cells, but board has only %d",
			totalPieceArea, boardArea,
		)
	}

	return true, ""
}
