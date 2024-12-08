package internal

import (
	"strings"
)

type Gridmask interface {
	Locations() []Coord
	Symbol() rune
	SetLocations(locations []Coord)
}

type gridmask struct {
	locations []Coord
	symbol    rune
}

func (g *gridmask) Locations() []Coord {
	return g.locations
}
func (g *gridmask) Symbol() rune {
	return g.symbol
}
func (g *gridmask) SetLocations(locations []Coord) {
	g.locations = locations
}

func NewGridmask(locations []Coord, symbol rune) Gridmask {
	return &gridmask{locations: locations, symbol: symbol}
}

type Gridmap interface {
	Width() int
	Height() int
	String(masks ...Gridmask) string
	Guards() []Guard
	GuardLocations() []Coord
	SetGuards(guards []Guard) Gridmap
	Obstructions() []Object
	SetObstructions(obstructions []Object) Gridmap
	ObstructionLocations() []Coord
	ValidateCoord(coord Coord) bool
	Clone() Gridmap
}

type gridMap struct {
	width        int
	height       int
	guards       []Guard
	obstructions []Object
}

func (g *gridMap) Width() int {
	return g.width
}
func (g *gridMap) Height() int {
	return g.height
}

func (g *gridMap) String(masks ...Gridmask) string {
	// Create a 2D slice filled with empty space characters
	grid := make([][]string, g.height)
	for i := range grid {
		grid[i] = make([]string, g.width)
		for j := range grid[i] {
			grid[i][j] = "."
		}
	}

	// Apply guards to the grid
	for _, guard := range g.guards {
		if g.ValidateCoord(guard.Location()) {
			x, y := guard.Location().X(), guard.Location().Y()
			grid[y][x] = guard.String()
		}
	}

	// Apply obstructions to the grid
	for _, obstruction := range g.obstructions {
		x, y := obstruction.Location().X(), obstruction.Location().Y()
		if x >= 0 && x < g.width && y >= 0 && y < g.height {
			grid[y][x] = string(obstruction.String())
		}
	}
	for _, mask := range masks {
		for _, coord := range mask.Locations() {
			x, y := coord.X(), coord.Y()
			if x >= 0 && x < g.width && y >= 0 && y < g.height {
				grid[y][x] = string(mask.Symbol())
			}
		}
	}

	// Construct the final string representation from the grid
	var sb strings.Builder
	for _, row := range grid {
		sb.WriteString(strings.Join(row, ""))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (g *gridMap) Guards() []Guard {
	return g.guards
}
func (g *gridMap) GuardLocations() []Coord {
	var locations []Coord
	for _, guard := range g.guards {
		locations = append(locations, guard.Location())
	}
	return locations
}
func (g *gridMap) SetGuards(guards []Guard) Gridmap {
	g.guards = guards
	return g
}
func (g *gridMap) Obstructions() []Object {
	return g.obstructions
}
func (g *gridMap) SetObstructions(obstructions []Object) Gridmap {
	g.obstructions = obstructions
	return g
}
func (g *gridMap) ObstructionLocations() []Coord {
	var locations []Coord
	for _, obstruction := range g.obstructions {
		locations = append(locations, obstruction.Location())
	}
	return locations
}

// ValidateCoord checks if a given coordinate is within the grid bounds.
func (g *gridMap) ValidateCoord(coord Coord) bool {
	if coord.X() < 0 || coord.X() >= g.width {
		return false
	}
	if coord.Y() < 0 || coord.Y() >= g.height {
		return false
	}
	return true
}

func (g *gridMap) Clone() Gridmap {
	newGrid := NewGridmap(g.width, g.height)
	newGuards := make([]Guard, len(g.guards))
	copy(newGuards, g.guards)
	newGrid.SetGuards(newGuards)
	newObs := make([]Object, len(g.obstructions))
	copy(newObs, g.obstructions)
	newGrid.SetObstructions(newObs)
	return newGrid
}

func NewGridmap(width int, height int) Gridmap {
	m := gridMap{
		width:        width,
		height:       height,
		guards:       []Guard{},
		obstructions: []Object{},
	}
	return &m
}
