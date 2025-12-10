package render

import (
	"2ajoyce/adventofcode/2025/9/geometry"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

// ai generated code - not reviewed for correctness
// This replaced the terminal output which could not render the large size images needed
// for the problem input

// DrawLinesPNG renders the given points + lines (+ optional area)
// to "output.png", scaling the world to a reasonable image size so
// huge coordinate ranges don't produce enormous images.
func DrawLinesPNG(points []*geometry.Point, lines []*geometry.Line, area ...*geometry.Point) error {
	if len(points) == 0 && len(lines) == 0 {
		return nil
	}

	// --------------------------------------------------
	// 1. Compute world bounds
	// --------------------------------------------------

	const intMax = int(^uint(0) >> 1)
	const intMin = -intMax - 1

	minX, maxX := intMax, intMin
	minY, maxY := intMax, intMin

	for _, l := range lines {
		if l == nil {
			continue
		}
		minX = min(minX, l.A.X, l.B.X)
		maxX = max(maxX, l.A.X, l.B.X)
		minY = min(minY, l.A.Y, l.B.Y)
		maxY = max(maxY, l.A.Y, l.B.Y)
	}
	for _, p := range points {
		if p == nil {
			continue
		}
		minX = min(minX, p.X)
		maxX = max(maxX, p.X)
		minY = min(minY, p.Y)
		maxY = max(maxY, p.Y)
	}

	// nothing at all
	if minX > maxX || minY > maxY {
		return nil
	}

	worldW := maxX - minX + 1
	worldH := maxY - minY + 1

	if worldW <= 0 {
		worldW = 1
	}
	if worldH <= 0 {
		worldH = 1
	}

	// --------------------------------------------------
	// 2. Choose image size + scale
	// --------------------------------------------------

	const maxDim = 1500 // max width or height in pixels
	const margin = 20   // border around drawing

	// scale factor from world units → pixels
	scale := float64(maxDim) / float64(max(worldW, worldH))

	imgW := int(math.Ceil(float64(worldW)*scale)) + margin*2
	imgH := int(math.Ceil(float64(worldH)*scale)) + margin*2

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// world → image mapping
	tx := func(x int) int {
		return int(math.Round((float64(x-minX) * scale))) + margin
	}
	ty := func(y int) int {
		return int(math.Round((float64(y-minY) * scale))) + margin
	}

	// pick a pixel radius for points that doesn't vanish when scale is small
	pointRadius := int(math.Max(2, scale*0.7))

	// --------------------------------------------------
	// 3. Optional shaded area (under everything else)
	// --------------------------------------------------

	if len(area) >= 2 && area[0] != nil && area[1] != nil {
		ax0 := tx(area[0].X)
		ay0 := ty(area[0].Y)
		ax1 := tx(area[1].X)
		ay1 := ty(area[1].Y)

		x0 := min(ax0, ax1)
		y0 := min(ay0, ay1)
		x1 := max(ax0, ax1)
		y1 := max(ay0, ay1)

		shade := color.RGBA{200, 200, 200, 255}

		for y := y0; y <= y1; y++ {
			for x := x0; x <= x1; x++ {
				img.Set(x, y, shade)
			}
		}
	}

	// --------------------------------------------------
	// 4. Draw lines
	// --------------------------------------------------

	lineColor := color.RGBA{0, 0, 0, 255}

	for _, ln := range lines {
		if ln == nil {
			continue
		}
		x0 := tx(ln.A.X)
		y0 := ty(ln.A.Y)
		x1 := tx(ln.B.X)
		y1 := ty(ln.B.Y)

		drawLine(img, x0, y0, x1, y1, lineColor)
	}

	// --------------------------------------------------
	// 5. Draw points on top
	// --------------------------------------------------

	pointColor := color.RGBA{255, 0, 0, 255}

	for _, p := range points {
		if p == nil {
			continue
		}
		drawCircle(img, tx(p.X), ty(p.Y), pointRadius, pointColor)
	}

	// --------------------------------------------------
	// 6. Save
	// --------------------------------------------------

	out, err := os.Create("output.png")
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, img)
}

// ----------------- helpers -----------------

func min(v int, rest ...int) int {
	m := v
	for _, x := range rest {
		if x < m {
			m = x
		}
	}
	return m
}

func max(v int, rest ...int) int {
	m := v
	for _, x := range rest {
		if x > m {
			m = x
		}
	}
	return m
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// Integer Bresenham line
func drawLine(img *image.RGBA, x0, y0, x1, y1 int, col color.Color) {
	dx := abs(x1 - x0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	dy := -abs(y1 - y0)
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx + dy

	for {
		img.Set(x0, y0, col)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func drawCircle(img *image.RGBA, cx, cy, r int, col color.Color) {
	if r < 1 {
		r = 1
	}
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				img.Set(cx+x, cy+y, col)
			}
		}
	}
}
