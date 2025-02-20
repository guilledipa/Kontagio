package kontagio

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Projectile representa un proyectil disparado por una torre.
type Projectile struct {
	x, y   float64
	target *Enemy
	speed  float64
}

// Update actualiza la posición del proyectil en el juego.
func (p *Projectile) Update() {
	// Move toward the target
	dx := p.target.x - p.x
	dy := p.target.y - p.y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist > p.speed {
		p.x += dx / dist * p.speed
		p.y += dy / dist * p.speed
	} else {
		// Hit the target
		p.target.health--
	}
}

// Draw dibuja el proyectil en la pantalla.
func (p *Projectile) Draw(screen *ebiten.Image) {
	// Draw a small yellow circle for the projectile
	vector.DrawFilledCircle(screen, float32(p.x), float32(p.y),
		float32(3), color.RGBA{255, 255, 0, 255}, false)
}
