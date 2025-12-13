package packing

import (
	"fmt"
	"sort"
)

// 0,0 is top-left
type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) Point {
	return Point{X: x, Y: y}
}

type Piece struct {
	Id     int
	Cells  []Point // occupied cells in original orientation
	Width  int     // bounding box width
	Height int     // bounding box height
	Area   int     // number of occupied cells
}

func NewPiece(id int, cells []Point) *Piece {
	if len(cells) == 0 {
		return &Piece{
			Id:     id,
			Cells:  nil,
			Width:  0,
			Height: 0,
			Area:   0,
		}
	}

	minX, minY := cells[0].X, cells[0].Y
	maxX, maxY := cells[0].X, cells[0].Y
	for _, p := range cells {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	width := maxX - minX + 1
	height := maxY - minY + 1

	// normalize so top-left is (0,0)
	norm := make([]Point, len(cells))
	for i, p := range cells {
		norm[i] = Point{X: p.X - minX, Y: p.Y - minY}
	}

	return &Piece{
		Id:     id,
		Cells:  norm,
		Width:  width,
		Height: height,
		Area:   len(cells),
	}
}

// Rotate90 returns a new piece rotated 90 degrees clockwise.
func (p *Piece) Rotate90() *Piece {
	newCells := make([]Point, len(p.Cells))
	for i, cell := range p.Cells {
		// rotate within [0..Width-1]x[0..Height-1]
		newCells[i] = Point{
			X: cell.Y,
			Y: p.Width - 1 - cell.X,
		}
	}
	return NewPiece(p.Id, newCells)
}

// FlipHorizontal returns a new piece flipped horizontally (mirror on vertical axis).
func (p *Piece) FlipHorizontal() *Piece {
	newCells := make([]Point, len(p.Cells))
	for i, cell := range p.Cells {
		newCells[i] = Point{
			X: p.Width - 1 - cell.X,
			Y: cell.Y,
		}
	}
	return NewPiece(p.Id, newCells)
}

// pieceKey makes a canonical string key for a piece's shape, used to deduplicate orientations.
func pieceKey(cells []Point) string {
	if len(cells) == 0 {
		return ""
	}

	pts := make([]Point, len(cells))
	copy(pts, cells)
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].Y != pts[j].Y {
			return pts[i].Y < pts[j].Y
		}
		return pts[i].X < pts[j].X
	})

	// encode as "x0,y0;x1,y1;..."
	s := ""
	for i, p := range pts {
		if i > 0 {
			s += ";"
		}
		s += fmt.Sprintf("%d,%d", p.X, p.Y)
	}

	return s
}

// AllOrientations returns all distinct shapes obtained by rotations and horizontal flips.
func (p *Piece) AllOrientations() []*Piece {
	var versions []*Piece
	versions = append(versions, p)
	r1 := p.Rotate90()
	r2 := r1.Rotate90()
	r3 := r2.Rotate90()
	versions = append(versions, r1, r2, r3)

	// flip each rotation horizontally
	var all []*Piece
	all = append(all, versions...)
	for _, g := range versions {
		f := g.FlipHorizontal()
		all = append(all, f)
	}

	// deduplicate by shape
	seen := make(map[string]bool)
	var unique []*Piece
	for _, g := range all {
		key := pieceKey(g.Cells)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, g)
		}
	}
	return unique
}
