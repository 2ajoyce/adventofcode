package packing

import (
	"maps"
	"sort"
)

// EvaluationResult holds the boolean answer for a single board
type EvaluationResult struct {
	Board  Board
	CanFit bool
	Reason string // optional
}

func (p *Problem) SolveWithBacktracking(board *Board) bool {
	// Initialize search state
	state := &SearchState{
		BoardWidth:  board.Width,
		BoardHeight: board.Height,
		Occupied:    make([][]bool, board.Height),
		Remaining:   make(map[int]int),
	}

	// Initialize occupied grid
	for i := range state.Occupied {
		state.Occupied[i] = make([]bool, board.Width)
	}

	// Copy piece requirements
	maps.Copy(state.Remaining, board.PieceCounts)

	// Precompute all possible placements
	placements := PrecomputePlacements(p, board)
	state.Placements = placements

	// Build index of pieceId[slice of placement indices]
	state.PiecePlacementIdx = make(map[int][]int)
	for idx, pl := range placements {
		state.PiecePlacementIdx[pl.PieceId] = append(state.PiecePlacementIdx[pl.PieceId], idx)
	}

	// Build piece order of pieceIds with positive remaining count sorted by
	// increasing number of placements (harder pieces first)
	for pieceId, count := range state.Remaining {
		if count > 0 {
			state.PieceOrder = append(state.PieceOrder, pieceId)
		}
	}
	sort.Slice(state.PieceOrder, func(i, j int) bool {
		a := state.PieceOrder[i]
		b := state.PieceOrder[j]
		return len(state.PiecePlacementIdx[a]) < len(state.PiecePlacementIdx[b])
	})

	// Kick off piece-type-based search
	return p.backtrackByPieceType(state, 0)
}

// backtrackByPieceType tries to place all pieces, one type at a time
// typeIndex is an index into state.PieceOrder
func (p *Problem) backtrackByPieceType(state *SearchState, typeIndex int) bool {
	// If we've considered all piece types, all required pieces must be placed
	if typeIndex >= len(state.PieceOrder) {
		for _, count := range state.Remaining {
			if count > 0 {
				return false
			}
		}
		return true
	}

	pieceId := state.PieceOrder[typeIndex]
	needed := state.Remaining[pieceId]

	if needed <= 0 {
		// Nothing to place for this type
		return p.backtrackByPieceType(state, typeIndex+1)
	}

	// Place needed copies of this pieceId
	return p.placeCopiesOfPiece(state, pieceId, typeIndex, 0, needed)
}

// placeCopiesOfPiece recursively chooses placements for pieceId
func (p *Problem) placeCopiesOfPiece(
	state *SearchState,
	pieceId int,
	typeIndex int,
	startIdx int,
	remaining int,
) bool {
	if remaining == 0 {
		// All copies of this piece type placed, move to next type
		return p.backtrackByPieceType(state, typeIndex+1)
	}

	indices := state.PiecePlacementIdx[pieceId]
	for i := startIdx; i < len(indices); i++ {
		plIndex := indices[i]
		placement := state.Placements[plIndex]

		if !placementHasRemaining(state, placement) {
			continue
		}
		if !p.isValidPlacement(state, placement) {
			continue
		}

		p.applyPlacement(state, placement)

		// Recurse, ensuring the next copy of this pieceId only considers
		// placements after i to avoid symmetric duplicates
		if p.placeCopiesOfPiece(state, pieceId, typeIndex, i+1, remaining-1) {
			return true
		}

		// backtrack
		p.undoPlacement(state, placement)
	}

	return false
}

// Placement represents a specific placement of a piece at a particular
// position on the board
type Placement struct {
	PieceId int
	Covers  []Point // absolute board coordinates covered
}

// SearchState represents the mutable state during backtracking
type SearchState struct {
	BoardWidth        int
	BoardHeight       int
	Occupied          [][]bool      // Occupied[y][x]
	Remaining         map[int]int   // pieceId -> remaining count
	Placements        []Placement   // all precomputed placements
	PiecePlacementIdx map[int][]int // pieceId -> indices into Placements
	PieceOrder        []int         // ordered list of pieceIds to place
}

// PrecomputePlacements enumerates all legal placements of pieces on this board
func PrecomputePlacements(problem *Problem, board *Board) []Placement {
	var placements []Placement

	// Generate placements for base pieces in all 4 rotations
	for _, piece := range problem.Pieces {
		for _, current := range piece.AllOrientations() {
			// Try all positions on the board
			for y := 0; y <= board.Height-current.Height; y++ {
				for x := 0; x <= board.Width-current.Width; x++ {
					// Calculate covered points
					var covers []Point
					for _, cell := range current.Cells {
						covers = append(covers, Point{X: x + cell.X, Y: y + cell.Y})
					}

					placements = append(placements, Placement{
						PieceId: piece.Id,
						Covers:  covers,
					})
				}
			}
		}
	}

	return placements
}

func placementHasRemaining(state *SearchState, pl Placement) bool {
	if state.Remaining[pl.PieceId] < 1 {
		return false
	}
	return true
}

// isValidPlacement checks if a placement conflicts with already occupied cells
func (p *Problem) isValidPlacement(state *SearchState, placement Placement) bool {
	for _, point := range placement.Covers {
		if point.X < 0 || point.X >= state.BoardWidth ||
			point.Y < 0 || point.Y >= state.BoardHeight {
			return false // Out of bounds
		}
		if state.Occupied[point.Y][point.X] {
			return false // Cell already occupied
		}
	}
	return true
}

// applyPlacement marks cells as occupied and decrements remaining piece count
func (p *Problem) applyPlacement(state *SearchState, placement Placement) {
	for _, point := range placement.Covers {
		state.Occupied[point.Y][point.X] = true
	}
	state.Remaining[placement.PieceId]--

}

// undoPlacement marks cells as free and increments remaining piece count
func (p *Problem) undoPlacement(state *SearchState, placement Placement) {
	for _, point := range placement.Covers {
		state.Occupied[point.Y][point.X] = false
	}
	state.Remaining[placement.PieceId]++
}
