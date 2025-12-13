package main

import (
	"2ajoyce/adventofcode/2025/12/packing"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// First Problem
	problem := ReadInput("input1.txt")
	result, err := Solve1(problem)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// ReadInput reads the input from the filepath and returns a Problem.
func ReadInput(filepath string) *packing.Problem {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)

	pieceId := 0               // Set at start of piece
	pieceRowIdx := 0           // reset at start of piece
	cells := []packing.Point{} // reset at start of piece

	pieces := []*packing.Piece{} // not reset
	var boards []*packing.Board  // not reset

	for scanner.Scan() {
		line := scanner.Text()
		// End of piece definition
		if len(line) == 0 {
			piece := packing.NewPiece(pieceId, cells)
			pieces = append(pieces, piece)
			continue
		}
		// Piece header
		if line[0] >= '0' && line[0] <= '9' && line[len(line)-1] == ':' {
			pieceId = int(line[0] - '0')
			pieceRowIdx = 0
			cells = []packing.Point{}
			continue
		}
		// Piece cell definition
		if line[0] == '#' || line[0] == '.' {
			for colIdx, ch := range line {
				if ch == '#' {
					cells = append(cells, packing.NewPoint(colIdx, pieceRowIdx))
				}
			}
			pieceRowIdx++
			continue
		}
		// Board definition
		if strings.Contains(line, "x") {
			board := ParseBoard(line, pieces)
			boards = append(boards, board)
			continue
		}

	}
	return packing.NewProblem(pieces, boards)
}

func ParseBoard(line string, pieces []*packing.Piece) *packing.Board {
	// Example: "12x5: 1 0 1 0 2 2"
	split := strings.Split(line, ":")

	dimStr := strings.TrimSpace(split[0])
	var width, height int
	fmt.Sscanf(dimStr, "%dx%d", &width, &height)
	board := packing.NewBoard(width, height)

	// Split piece counts into individual strings
	pieceCountsStr := strings.Fields(strings.TrimSpace(split[1]))
	for pid, countStr := range pieceCountsStr {
		var count int
		fmt.Sscanf(countStr, "%d", &count)
		board.AddPiece(pieces[pid], count)
	}

	return board
}

func Solve1(problem *packing.Problem) (string, error) {
	total := 0
	result := problem.EvaluateProblem()
	for _, res := range result {
		if res.CanFit {
			total++
		}
	}
	fmt.Printf("Tried: %d, Passed: %d\n", len(result), total)
	return fmt.Sprintf("%d", total), nil
}
