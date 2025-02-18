package kontagio

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	SpawnDelay = 30 // Delay between enemy spawns in frames
)

// Enemy path: El camino deberia ser procedual, no estático, asi hay más
// variedad cada vez que se juega.
var mainPath = []float64{
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
	// Path especifico de este enemigo
	path []float64
}

func (e *Enemy) Update() {
	if e.pathIdx < len(e.path)-1 {
		targetX := e.path[e.pathIdx]
		targetY := e.path[e.pathIdx+1]

		dx := targetX - e.x
		dy := targetY - e.y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist > e.speed {
			e.x += dx / dist * e.speed
			e.y += dy / dist * e.speed
		} else {
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
