package main

import (
	"2ajoyce/adventofcode/2025/9/geometry"
	"2ajoyce/adventofcode/2025/9/render"
	"bufio"
	"fmt"
	"os"

	"github.com/schollz/progressbar/v3"
)

// Todo: Refactor Solve2 and supporting methods
// I almost ran out of time on this solution and I'm leaving this file a mess
// The solution 

func main() {
	// First Problem
	input := make(chan *geometry.Point)
	go ReadInput("input1.txt", input)
	result, err := Solve1(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// Second Problem
	input = make(chan *geometry.Point)
	go ReadInput("input2.txt", input)
	result, err = Solve2(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// ReadInput reads the input from the filepath and sends each line to the provided channel.
func ReadInput(filepath string, c chan *geometry.Point) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		c <- ParseInput(line)
	}
	close(c)
}

// ParseInput parses the input into the necessary data structure.
// On more complex inputs, this allows us to use lines of text as input for tests
func ParseInput(input string) *geometry.Point {
	// input is in the form "x,y"
	var x, y int
	fmt.Sscanf(input, "%d,%d", &x, &y)
	return geometry.NewPoint(x, y)
}

func Solve1(input chan *geometry.Point) (string, error) {
	total := 0
	points := []*geometry.Point{}
	for p := range input {
		points = append(points, p)
	}

	// Find the largest area between points
	// For each point, find the distance to every other point
	area := 0 // The largest area found so far
	for i := range points {
		for j := range points {
			if i == j {
				continue
			}
			a := Area(points[i], points[j])
			if a > area {
				area = a
			}
		}
	}
	total = area

	return fmt.Sprintf("%d", total), nil
}

type rectJob struct {
	i, j int
}

type rectResult struct {
	area int
	i, j int
}

func Solve2(input chan *geometry.Point) (string, error) {
	points := []*geometry.Point{}
	for p := range input {
		points = append(points, p)
	}

	// Build polygon boundary
	lines := []*geometry.Line{}
	for i := range points {
		if i == 0 {
			continue
		}
		lines = append(lines, geometry.NewLine(points[i-1], points[i]))
	}
	lines = append(lines, geometry.NewLine(points[len(points)-1], points[0]))

	// Precompute fast edge structures for PointInPolygonFast
	vEdges, hEdges := buildEdges(lines)

	// Prepare candidate pairs
	var jobsList []rectJob
	for i := range points {
		for j := range points {
			// skip same point or axis-aligned rectangles
			if i == j || points[i].X == points[j].X || points[i].Y == points[j].Y {
				continue
			}
			jobsList = append(jobsList, rectJob{i, j})
		}
	}
	jobCount := len(jobsList)
	if jobCount == 0 {
		return "0", nil
	}

	bar := progressbar.Default(int64(jobCount), "Searching for largest interior area")

	jobs := make(chan rectJob)
	results := make(chan rectResult)

	const workerCount = 8

	// Workers
	for range workerCount {
		go func() {
			for job := range jobs {
				p1 := points[job.i]
				p2 := points[job.j]

				// Check rectangle using fast geometry and polygon tests
				if !rectInsidePolygon(p1, p2, lines, vEdges, hEdges) {
					results <- rectResult{area: 0, i: job.i, j: job.j}
					continue
				}

				// Only compute area if it's actually inside
				a := Area(p1, p2)
				results <- rectResult{area: a, i: job.i, j: job.j}
			}
		}()
	}

	// Feed jobs
	go func() {
		for _, job := range jobsList {
			jobs <- job
		}
		close(jobs)
	}()

	bestArea := 0
	var bestP1, bestP2 *geometry.Point

	for range jobCount {
		res := <-results
		bar.Add(1)

		if res.area > bestArea {
			bestArea = res.area
			bestP1 = points[res.i]
			bestP2 = points[res.j]
		}
	}
	bar.Finish()

	render.DrawLinesPNG(points, lines, bestP1, bestP2)
	return fmt.Sprintf("%d", bestArea), nil
}

// vEdge and hEdge are precomputed axis-aligned edges of the polygon
// to speed up PointInPolygonFast.
type vEdge struct {
	x      int
	y1, y2 int // normalized so y1 <= y2
}

type hEdge struct {
	y      int
	x1, x2 int // normalized so x1 <= x2
}

// buildEdges constructs vertical and horizontal edge lists from the polygon lines.
func buildEdges(lines []*geometry.Line) ([]vEdge, []hEdge) {
	var vEdges []vEdge
	var hEdges []hEdge

	for _, e := range lines {
		if e.IsVerticalUnsafe() {
			y1, y2 := e.A.Y, e.B.Y
			if y1 > y2 {
				y1, y2 = y2, y1
			}
			vEdges = append(vEdges, vEdge{
				x:  e.A.X,
				y1: y1,
				y2: y2,
			})
		} else { // horizontal (lines are axis-aligned by design)
			x1, x2 := e.A.X, e.B.X
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			hEdges = append(hEdges, hEdge{
				y:  e.A.Y,
				x1: x1,
				x2: x2,
			})
		}
	}

	return vEdges, hEdges
}

// rectInsidePolygon checks if the axis-aligned rectangle with opposite
// corners p1 and p2 lies wholly inside the polygon defined by lines.
// It uses:
//   - DoesCrossUnsafe to reject rectangles whose edges cross the polygon
//   - DoesOverlapUnsafe + PointInPolygonFast to ensure overlapping edges
//     still lie inside
//   - A corner check as a cheap sanity check.
func rectInsidePolygon(p1, p2 *geometry.Point, lines []*geometry.Line, vEdges []vEdge, hEdges []hEdge) bool {
	// Edge checks for the candidate rectangle
	rect := GetLines(p1, p2)

	for _, line := range rect {
		for _, existingLine := range lines {
			// Hard crossing: rectangle boundary crosses polygon boundary
			if existingLine.DoesCrossUnsafe(line) {
				return false
			}

			// Overlap: walk along the overlapping line and ensure all
			// lattice points are inside polygon.
			if existingLine.DoesOverlapUnsafe(line) {
				if line.IsVerticalUnsafe() {
					y1, y2 := line.A.Y, line.B.Y
					if y1 > y2 {
						y1, y2 = y2, y1
					}
					for y := y1; y <= y2; y++ {
						pt := geometry.NewPoint(line.A.X, y)
						if !PointInPolygonFast(pt, vEdges, hEdges) {
							return false
						}
					}
				} else { // horizontal
					x1, x2 := line.A.X, line.B.X
					if x1 > x2 {
						x1, x2 = x2, x1
					}
					for x := x1; x <= x2; x++ {
						pt := geometry.NewPoint(x, line.A.Y)
						if !PointInPolygonFast(pt, vEdges, hEdges) {
							return false
						}
					}
				}
			}
		}
	}

	// Corner check: all four corners must be inside
	corners := []*geometry.Point{
		p1,
		p2,
		geometry.NewPoint(p1.X, p2.Y),
		geometry.NewPoint(p2.X, p1.Y),
	}
	for _, c := range corners {
		if !PointInPolygonFast(c, vEdges, hEdges) {
			return false
		}
	}

	return true
}

// PointInPolygonFast is a faster version of PointInPolygon that uses
// precomputed vertical and horizontal edges and avoids per-call validation.
// Edges are considered inside
func PointInPolygonFast(p *geometry.Point, vEdges []vEdge, hEdges []hEdge) bool {
	// On-edge = inside
	for _, e := range vEdges {
		if p.X == e.x && p.Y >= e.y1 && p.Y <= e.y2 {
			return true
		}
	}
	for _, e := range hEdges {
		if p.Y == e.y && p.X >= e.x1 && p.X <= e.x2 {
			return true
		}
	}

	// Ray cast to the right using vertical edges only
	crossings := 0
	for _, e := range vEdges {
		// Use half-open [y1, y2) to avoid double-counting vertices
		if p.Y < e.y1 || p.Y >= e.y2 {
			continue
		}
		if e.x > p.X {
			crossings++
		}
	}
	return crossings%2 == 1
}

// Area calculates the area of the rectangle defined by two points on opposite corners.
// The area is INCLUSIVE of the points.
func Area(p1 *geometry.Point, p2 *geometry.Point) int {
	width := intAbs(p1.X-p2.X) + 1
	height := intAbs(p1.Y-p2.Y) + 1
	return width * height
}

// GetLines generates the lines between two points
func GetLines(p1, p2 *geometry.Point) []*geometry.Line {
	lines := []*geometry.Line{}

	// Walk around the rectangle defined by p1 and p2 clockwise
	p3 := geometry.NewPoint(p1.X, p2.Y)
	lines = append(lines, geometry.NewLine(p1, p3))
	lines = append(lines, geometry.NewLine(p3, p2))
	p4 := geometry.NewPoint(p2.X, p1.Y)
	lines = append(lines, geometry.NewLine(p2, p4))
	lines = append(lines, geometry.NewLine(p4, p1))
	return lines
}

// DrawLines draws the provided lines to the console for visualization.
// If area is provided, locations within the area will be shaded using U+2591 (light shade)
// unless they are occupied by a point or line
func DrawLines(points []*geometry.Point, lines []*geometry.Line, area ...*geometry.Point) {
	// Find the bounds of the lines
	minX, maxX := 0, 0
	minY, maxY := 0, 0
	for _, line := range lines {
		if line.A.X < minX {
			minX = line.A.X
		}
		if line.B.X < minX {
			minX = line.B.X
		}

		if line.A.X > maxX {
			maxX = line.A.X
		}
		if line.B.X > maxX {
			maxX = line.B.X
		}
		if line.A.Y < minY {
			minY = line.A.Y
		}
		if line.B.Y < minY {
			minY = line.B.Y
		}
		if line.A.Y > maxY {
			maxY = line.A.Y
		}
		if line.B.Y > maxY {
			maxY = line.B.Y
		}
	}

	// Prepare shading bounds if an area was provided
	hasArea := false
	var shadeMinX, shadeMaxX, shadeMinY, shadeMaxY int
	if len(area) >= 2 && area[0] != nil && area[1] != nil {
		hasArea = true
		shadeMinX = min(area[0].X, area[1].X)
		shadeMaxX = max(area[0].X, area[1].X)
		shadeMinY = min(area[0].Y, area[1].Y)
		shadeMaxY = max(area[0].Y, area[1].Y)
	}

	// Draw the lines
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			pointDrawn := false
			// Draw points first (take precedence over lines)
			for _, p := range points {
				if p.X == x && p.Y == y {
					fmt.Print("#")
					pointDrawn = true
					break
				}
			}
			// If no point, draw any line that covers this coordinate
			if !pointDrawn {
				for _, line := range lines {
					if line.IsVertical() && line.A.X == x && y >= min(line.A.Y, line.B.Y) && y <= max(line.A.Y, line.B.Y) {
						fmt.Print("·")
						pointDrawn = true
						break
					}
					if line.IsHorizontal() && line.A.Y == y && x >= min(line.A.X, line.B.X) && x <= max(line.A.X, line.B.X) {
						fmt.Print("·")
						pointDrawn = true
						break
					}
				}
			}
			if !pointDrawn {
				// If an area was provided and this coordinate is inside it, shade it
				if hasArea && x >= shadeMinX && x <= shadeMaxX && y >= shadeMinY && y <= shadeMaxY {
					fmt.Print("░")
				} else {
					fmt.Print(" ")
				}
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
