package kontagio

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var path = []float64{
	0, 240, // Start at left, middle of the screen
	320, 240, // Move to the middle
	320, 120, // Move up
	640, 120, // Move to the right
}

type Enemy struct {
	x, y    float64
	health  int
	speed   float64
	pathIdx int
}

func (e *Enemy) Update() {
	// Move along the path
	if e.pathIdx < len(path)-1 {
		targetX := path[e.pathIdx]
		targetY := path[e.pathIdx+1]
		dx := targetX - e.x
		dy := targetY - e.y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > e.speed {
			e.x += dx / dist * e.speed
			e.y += dy / dist * e.speed
		} else {
			e.x = targetX
			e.y = targetY
			e.pathIdx += 2
		}
	}
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen,
		float32(e.x-10), float32(e.y-10),
		float32(20), float32(20),
		color.RGBA{255, 0, 0, 255}, false)
}
