package kontagio

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	towerCost       = 50
	shootCooldown   = 15 // Cooldown time between shots (in ticks)
	toweMinDistance = 35 // Minimum distance between turrets
)

// Tower representa una torre en el juego.
type Tower struct {
	x, y        float64
	attackRange float64
	cooldown    int // Cooldown time for shooting again
}

// Update actualiza los disparos de la torre en el juego.
func (t *Tower) Update(g *Game) {
	if t.cooldown > 0 {
		t.cooldown--
	} else {
		for _, enemy := range g.Enemies {
			dx := enemy.x - t.x
			dy := enemy.y - t.y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= t.attackRange {
				// Shoot at the enemy
				g.Projectiles = append(g.Projectiles, &Projectile{
					x:      t.x,
					y:      t.y,
					target: enemy,
					speed:  4,
				})
				t.cooldown = shootCooldown
				break // Only shoot one enemy at a time
			}
		}
	}
}

// Draw dibuja la torre en la pantalla.
func (t *Tower) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen,
		float32(t.x-15), float32(t.y-15),
		float32(30), float32(30),
		color.RGBA{0, 0, 255, 255}, false)
}
