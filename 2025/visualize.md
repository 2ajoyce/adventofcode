# Fyne Visualization Rules for Advent of Code

This document contains rules and patterns for creating visualizations using **Fyne v2.7.x** for Advent of Code puzzles. Following these rules will help ensure visualizations work correctly the first time.

## Version Information

- **Fyne Version**: v2.7.1
- **Go Version**: 1.25+
- **Reference**: https://pkg.go.dev/fyne.io/fyne/v2@v2.7.1

---

## 1. Project Setup

### go.mod Requirements

```go
require fyne.io/fyne/v2 v2.7.1
```

### Required Imports

```go
import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)
```

---

## 2. Application Structure

### Creating the App and Window

```go
func RunVisualization(data interface{}) error {
    a := app.New()
    w := a.NewWindow("Visualization Title")

    // Set content
    w.SetContent(yourContent)

    // Use ShowAndRun() at the end - this blocks!
    w.ShowAndRun()

    return nil
}
```

### Key Rules:

- ✅ `app.New()` creates the application
- ✅ `a.NewWindow("Title")` creates a window
- ✅ `w.ShowAndRun()` shows the window AND starts the event loop (blocking)
- ❌ Don't use `w.Show()` followed by `a.Run()` separately in simple cases

---

## 3. Custom Widget Pattern (Recommended for Grid Visualizations)

### Widget Structure

A custom widget has two parts:

1. **Widget struct** - holds state and implements `fyne.Widget`
2. **Renderer struct** - handles drawing and implements `fyne.WidgetRenderer`

### Complete Custom Widget Template

```go
// Widget struct - holds all state
type gridCanvas struct {
    widget.BaseWidget  // REQUIRED: Embed BaseWidget
    vis *visualizationState  // Reference to your state
}

// REQUIRED: Implement CreateRenderer
func (g *gridCanvas) CreateRenderer() fyne.WidgetRenderer {
    raster := canvas.NewRaster(g.draw)
    return &gridRenderer{
        raster: raster,
        canvas: g,
    }
}

// Drawing function for the raster
func (g *gridCanvas) draw(w, h int) image.Image {
    img := image.NewRGBA(image.Rect(0, 0, w, h))
    // Draw your visualization here
    return img
}

// Renderer struct - handles layout and drawing
type gridRenderer struct {
    raster *canvas.Raster
    canvas *gridCanvas
}

// REQUIRED: Layout positions and sizes the canvas objects
func (r *gridRenderer) Layout(size fyne.Size) {
    r.raster.Resize(size)
}

// REQUIRED: MinSize returns minimum size needed
func (r *gridRenderer) MinSize() fyne.Size {
    // Calculate based on your grid dimensions
    return fyne.NewSize(float32(width), float32(height))
}

// REQUIRED: Refresh redraws the widget
func (r *gridRenderer) Refresh() {
    canvas.Refresh(r.raster)  // Use canvas.Refresh(), not r.raster.Refresh()
}

// REQUIRED: Objects returns all canvas objects to draw
func (r *gridRenderer) Objects() []fyne.CanvasObject {
    return []fyne.CanvasObject{r.raster}
}

// REQUIRED: Destroy cleans up resources
func (r *gridRenderer) Destroy() {}
```

### Critical Widget Rules:

- ✅ Always embed `widget.BaseWidget` in your widget struct
- ✅ Call `widget.ExtendBaseWidget(yourWidget)` after creating the widget
- ✅ Use `canvas.Refresh(object)` to refresh canvas objects
- ❌ Don't call `Refresh()` directly on canvas primitives - use `canvas.Refresh()`
- ❌ Don't store important state in the renderer - it may be destroyed

---

## 4. Creating Your Widget Instance

```go
// Create widget instance
v := &visualizationState{
    grid: gridData,
    // ... other state
}

v.canvas = &gridCanvas{vis: v}
v.canvas.ExtendBaseWidget(v.canvas)  // CRITICAL: Must call this!
```

---

## 5. Thread Safety with fyne.DoAndWait / fyne.Do

### When Background Goroutines Modify UI

```go
go func() {
    ticker := time.NewTicker(200 * time.Millisecond)
    defer ticker.Stop()

    for {
        <-ticker.C
        if shouldUpdate {
            // REQUIRED: Use fyne.DoAndWait for UI updates from goroutines
            fyne.DoAndWait(func() {
                // Modify state here
                v.grid[y][x] = newValue
                v.canvas.Refresh()
            })
        }
    }
}()
```

### Key Rules:

- ✅ `fyne.DoAndWait(fn)` - executes on UI thread and waits for completion (Since v2.6)
- ✅ `fyne.Do(fn)` - executes on UI thread without waiting (Since v2.6)
- ❌ Never modify UI elements directly from goroutines without `fyne.Do`/`fyne.DoAndWait`

---

## 6. Canvas Primitives

### canvas.Raster - For Pixel-Perfect Grid Drawing

```go
// Method 1: Draw function returning full image
raster := canvas.NewRaster(func(w, h int) image.Image {
    img := image.NewRGBA(image.Rect(0, 0, w, h))
    // Draw pixels
    for y := 0; y < h; y++ {
        for x := 0; x < w; x++ {
            img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
        }
    }
    return img
})

// Method 2: Pixel-by-pixel generator
raster := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
    return color.RGBA{R: 255, G: 255, B: 0, A: 255}
})
```

### Other Canvas Primitives

```go
// Rectangle
rect := canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 0, A: 255})
rect.Resize(fyne.NewSize(100, 100))
rect.Move(fyne.NewPos(10, 10))

// Text
text := canvas.NewText("Hello", color.White)
text.TextSize = 14

// Image from file
img := canvas.NewImageFromFile("path/to/image.png")

// Line
line := canvas.NewLine(color.White)
line.Position1 = fyne.NewPos(0, 0)
line.Position2 = fyne.NewPos(100, 100)
```

---

## 7. Standard Widgets for Controls

### Common Control Widgets

```go
// Button
button := widget.NewButton("Click Me", func() {
    // action
})

// Check box
check := widget.NewCheck("Option", func(checked bool) {
    // handle change
})

// Select dropdown
sel := widget.NewSelect([]string{"Option1", "Option2"}, func(selected string) {
    // handle selection
})
sel.SetSelected("Option1")  // Set default

// Label
label := widget.NewLabel("Status: Ready")
label.SetText("Status: Running")  // Update text

// Slider (Since v2.0)
slider := widget.NewSlider(0, 100)
slider.OnChanged = func(value float64) {
    // handle value change
}
```

---

## 8. Container Layouts

### Common Layouts

```go
// Horizontal box - items side by side
hbox := container.NewHBox(button1, button2, button3)

// Vertical box - items stacked
vbox := container.NewVBox(label, button)

// Border layout - items at edges, one fills center
// NewBorder(top, bottom, left, right, center)
border := container.NewBorder(toolbar, nil, nil, nil, mainContent)

// Grid layout - equal sized cells
grid := container.NewGridWithColumns(3, item1, item2, item3)

// Stack - items on top of each other
stack := container.NewStack(background, foreground)

// Scroll container
scroll := container.NewScroll(largeContent)
```

### Key Layout Rules:

- ✅ `container.NewBorder(top, bottom, left, right, center)` - pass `nil` for unused edges
- ✅ Center element in Border expands to fill remaining space
- ❌ Don't use deprecated `fyne.NewContainer()` - use `container.New*` functions

---

## 9. Geometry Types

```go
// Size
size := fyne.NewSize(width, height)  // float32 values
w, h := size.Width, size.Height

// Position
pos := fyne.NewPos(x, y)  // float32 values

// Moving and resizing canvas objects
obj.Move(fyne.NewPos(10, 20))
obj.Resize(fyne.NewSize(100, 50))
```

---

## 10. Colors

```go
import "image/color"

// RGBA color (most common)
red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
transparent := color.RGBA{R: 255, G: 255, B: 255, A: 128}

// Named colors from image/color
white := color.White
black := color.Black

// Theme colors (adapt to light/dark mode)
import "fyne.io/fyne/v2/theme"
primary := theme.PrimaryColor()
```

---

## 11. Animation Pattern

### Using Ticker for Animation

```go
go func() {
    ticker := time.NewTicker(v.speed)
    defer ticker.Stop()

    for {
        <-ticker.C
        if !v.paused {
            fyne.DoAndWait(func() {
                // Update state
                // ...
                v.canvas.Refresh()
            })
        }
        ticker.Reset(v.speed)  // Allow dynamic speed changes
    }
}()
```

### Using Fyne Animation API

```go
anim := fyne.NewAnimation(time.Second, func(progress float32) {
    // progress goes from 0.0 to 1.0
    // Update visuals based on progress
    canvas.Refresh(obj)
})
anim.Start()
```

---

## 12. Complete Visualization Template

```go
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

const cellSize = 10

type visualizationState struct {
    grid        [][]rune
    initialGrid [][]rune

    canvas *gridCanvas
    paused bool
    speed  time.Duration
}

type gridCanvas struct {
    widget.BaseWidget
    vis *visualizationState
}

func (g *gridCanvas) CreateRenderer() fyne.WidgetRenderer {
    raster := canvas.NewRaster(g.draw)
    return &gridRenderer{raster: raster, canvas: g}
}

func (g *gridCanvas) draw(w, h int) image.Image {
    img := image.NewRGBA(image.Rect(0, 0, w, h))
    grid := g.vis.grid

    for y, row := range grid {
        for x, cell := range row {
            col := colorForCell(cell)
            for dy := 0; dy < cellSize; dy++ {
                for dx := 0; dx < cellSize; dx++ {
                    img.Set(x*cellSize+dx, y*cellSize+dy, col)
                }
            }
        }
    }
    return img
}

func colorForCell(cell rune) color.Color {
    switch cell {
    case '#':
        return color.RGBA{R: 255, G: 255, B: 0, A: 255}
    default:
        return color.RGBA{R: 20, G: 20, B: 20, A: 255}
    }
}

type gridRenderer struct {
    raster *canvas.Raster
    canvas *gridCanvas
}

func (r *gridRenderer) Layout(size fyne.Size) {
    r.raster.Resize(size)
}

func (r *gridRenderer) MinSize() fyne.Size {
    grid := r.canvas.vis.grid
    if len(grid) == 0 {
        return fyne.NewSize(100, 100)
    }
    return fyne.NewSize(
        float32(len(grid[0])*cellSize),
        float32(len(grid)*cellSize),
    )
}

func (r *gridRenderer) Refresh() {
    canvas.Refresh(r.raster)
}

func (r *gridRenderer) Objects() []fyne.CanvasObject {
    return []fyne.CanvasObject{r.raster}
}

func (r *gridRenderer) Destroy() {}

func RunVisualization(grid [][]rune) error {
    a := app.New()
    w := a.NewWindow("AoC Visualization")

    // Deep copy the grid
    gridCopy := make([][]rune, len(grid))
    initialCopy := make([][]rune, len(grid))
    for i, row := range grid {
        gridCopy[i] = make([]rune, len(row))
        initialCopy[i] = make([]rune, len(row))
        copy(gridCopy[i], row)
        copy(initialCopy[i], row)
    }

    v := &visualizationState{
        grid:        gridCopy,
        initialGrid: initialCopy,
        paused:      true,
        speed:       200 * time.Millisecond,
    }

    v.canvas = &gridCanvas{vis: v}
    v.canvas.ExtendBaseWidget(v.canvas)

    // Controls
    playPauseBtn := widget.NewButton("Play", nil)
    playPauseBtn.OnTapped = func() {
        v.paused = !v.paused
        if v.paused {
            playPauseBtn.SetText("Play")
        } else {
            playPauseBtn.SetText("Pause")
        }
    }

    resetBtn := widget.NewButton("Reset", func() {
        for i, row := range v.initialGrid {
            copy(v.grid[i], row)
        }
        v.canvas.Refresh()
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

    toolbar := container.NewHBox(playPauseBtn, resetBtn, speedSelect)
    content := container.NewBorder(toolbar, nil, nil, nil, v.canvas)

    // Animation loop
    go func() {
        ticker := time.NewTicker(v.speed)
        defer ticker.Stop()

        for {
            <-ticker.C
            if !v.paused {
                fyne.DoAndWait(func() {
                    // Your simulation step here
                    // Update v.grid as needed
                    v.canvas.Refresh()
                })
            }
            ticker.Reset(v.speed)
        }
    }()

    w.SetContent(content)
    w.ShowAndRun()

    return nil
}
```

---

## 13. Common Mistakes to Avoid

| ❌ Wrong                             | ✅ Correct                                      |
| ------------------------------------ | ----------------------------------------------- |
| `fyne.NewContainer()`                | `container.NewVBox()` or other `container.New*` |
| `widget.Refresh()` on canvas objects | `canvas.Refresh(obj)`                           |
| Modifying UI from goroutine directly | Wrap in `fyne.DoAndWait()`                      |
| Forgetting `ExtendBaseWidget()`      | Always call after creating custom widget        |
| Storing UI state in renderer         | Store state in widget, renderer reads from it   |
| Using `w.Show()` + `a.Run()`         | Use `w.ShowAndRun()` for simple apps            |

---

## 14. Debugging Tips

1. **Widget not appearing**: Did you call `ExtendBaseWidget()`?
2. **Crashes when updating**: Are you using `fyne.DoAndWait()` from goroutines?
3. **Widget not refreshing**: Are you calling `canvas.Refresh()` on the right object?
4. **Wrong size**: Check your `MinSize()` implementation in the renderer
5. **Layout issues**: Verify container type matches your needs (Border, VBox, HBox, etc.)

---

## 15. API Quick Reference

### fyne.WidgetRenderer Interface (Must Implement All)

```go
type WidgetRenderer interface {
    Layout(Size)              // Position and size objects
    MinSize() Size            // Return minimum dimensions
    Refresh()                 // Redraw after state change
    Objects() []CanvasObject  // Return objects to draw
    Destroy()                 // Cleanup resources
}
```

### Key Functions

| Function                     | Purpose                             |
| ---------------------------- | ----------------------------------- |
| `app.New()`                  | Create application                  |
| `a.NewWindow(title)`         | Create window                       |
| `w.SetContent(obj)`          | Set window content                  |
| `w.ShowAndRun()`             | Show window and start event loop    |
| `widget.ExtendBaseWidget(w)` | Initialize custom widget            |
| `canvas.Refresh(obj)`        | Trigger redraw of canvas object     |
| `fyne.DoAndWait(fn)`         | Execute on UI thread (blocking)     |
| `fyne.Do(fn)`                | Execute on UI thread (non-blocking) |
| `fyne.NewSize(w, h)`         | Create size (float32)               |
| `fyne.NewPos(x, y)`          | Create position (float32)           |
