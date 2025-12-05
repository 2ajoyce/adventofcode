package main

import (
	"image"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/////////////////////////////////////////////
// AI Generated Visualization Code
/////////////////////////////////////////////
// Yeah, I know... but it's better than nothing!

const (
	squareSize = 10
)

var (
	paperColor = color.RGBA{R: 255, G: 255, B: 0, A: 255} // Yellow
	emptyColor = color.RGBA{R: 20, G: 20, B: 20, A: 255}  // Dark Gray
)

type vis struct {
	grid        [][]rune
	initialGrid [][]rune
	total       int

	canvas *gridCanvas
	paused bool
	loop   bool
	speed  time.Duration
}

type gridCanvas struct {
	widget.BaseWidget
	vis *vis
}

func (g *gridCanvas) CreateRenderer() fyne.WidgetRenderer {
	raster := canvas.NewRaster(g.draw)
	return &gridRenderer{
		raster: raster,
		canvas: g,
	}
}

func (g *gridCanvas) draw(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	gridRunes := g.vis.grid

	for y, row := range gridRunes {
		for x, r := range row {
			col := emptyColor
			if IsPaper(r) {
				col = paperColor
			}
			for i := 0; i < squareSize; i++ {
				for j := 0; j < squareSize; j++ {
					img.Set(x*squareSize+i, y*squareSize+j, col)
				}
			}
		}
	}
	return img
}

type gridRenderer struct {
	raster *canvas.Raster
	canvas *gridCanvas
}

func (r *gridRenderer) Layout(size fyne.Size) {
	r.raster.Resize(size)
}

func (r *gridRenderer) MinSize() fyne.Size {
	gridRunes := r.canvas.vis.grid
	return fyne.NewSize(float32(len(gridRunes[0])*squareSize), float32(len(gridRunes)*squareSize))
}

func (r *gridRenderer) Refresh() {
	canvas.Refresh(r.raster)
}

func (r *gridRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.raster}
}

func (r *gridRenderer) Destroy() {}

func RunVisualization(grid [][]rune) (int, error) {
	a := app.New()
	w := a.NewWindow("Paper Removal Simulation")

	initialGrid := make([][]rune, len(grid))
	for i, row := range grid {
		initialGrid[i] = make([]rune, len(row))
		copy(initialGrid[i], row)
	}

	// create a deep copy of the input grid to avoid mutating caller data
	gridCopy := make([][]rune, len(grid))
	for i, row := range grid {
		gridCopy[i] = make([]rune, len(row))
		copy(gridCopy[i], row)
	}

	v := &vis{
		grid:        gridCopy,
		initialGrid: initialGrid,
		paused:      true,
		loop:        false,
		speed:       200 * time.Millisecond,
	}

	v.canvas = &gridCanvas{vis: v}
	v.canvas.ExtendBaseWidget(v.canvas)

	playPauseButton := widget.NewButton("Play", nil)
	playPauseButton.OnTapped = func() {
		v.paused = !v.paused
		if v.paused {
			playPauseButton.SetText("Play")
		} else {
			playPauseButton.SetText("Pause")
		}
	}

	resetButton := widget.NewButton("Reset", func() {
		v.total = 0
		newGrid := make([][]rune, len(v.initialGrid))
		for i, row := range v.initialGrid {
			newGrid[i] = make([]rune, len(row))
			copy(newGrid[i], row)
		}
		v.grid = newGrid
		v.canvas.Refresh()
	})

	loopToggle := widget.NewCheck("Loop", func(on bool) {
		v.loop = on
	})

	speedSelect := widget.NewSelect([]string{"Slow", "Medium", "Fast"}, func(s string) {
		switch s {
		case "Slow":
			v.speed = 500 * time.Millisecond
		case "Medium":
			v.speed = 200 * time.Millisecond
		case "Fast":
			v.speed = 50 * time.Millisecond
		}
	})
	speedSelect.SetSelected("Medium")

	toolbar := container.NewHBox(playPauseButton, resetButton, loopToggle, speedSelect)
	content := container.NewBorder(toolbar, nil, nil, nil, v.canvas)

	go func() {
		ticker := time.NewTicker(v.speed)
		defer ticker.Stop()
		for {
			<-ticker.C
			if !v.paused {
				// run the simulation step on the main UI thread to avoid races
				fyne.DoAndWait(func() {
					h := CalculateHeatmap(v.grid)
					paperRemoved := false

					for y, row := range v.grid {
						for x, c := range row {
							if IsPaper(c) && h[y][x] < 4 {
								v.grid[y][x] = '.'
								v.total++
								paperRemoved = true
							}
						}
					}

					if paperRemoved {
						v.canvas.Refresh()
					}

					if !paperRemoved {
						if v.loop {
							// reset
							v.total = 0
							newGrid := make([][]rune, len(v.initialGrid))
							for i, row := range v.initialGrid {
								newGrid[i] = make([]rune, len(row))
								copy(newGrid[i], row)
							}
							v.grid = newGrid
							v.canvas.Refresh()
						} else {
							v.paused = true
							playPauseButton.SetText("Play")
						}
					}
				})
			}
			ticker.Reset(v.speed)
		}
	}()

	w.SetContent(content)
	w.ShowAndRun()

	return v.total, nil
}
