package kontagio

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	// SpawnDelay es el tiempo en ticks entre la aparición de enemigos.
	SpawnDelay = 30
)

// Enemy representa a un enemigo en el juego.
type Enemy struct {
	x, y    float64
	health  int
	speed   float64
	pathIdx int
	// Path especifico de este enemigo
	path []float64
}

// Update actualiza la posición del enemigo en el juego.
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

// Draw dibuja al enemigo en la pantalla.
func (e *Enemy) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen,
		float32(e.x-10), float32(e.y-10),
		float32(20), float32(20),
		color.RGBA{255, 0, 0, 255}, false)
}
