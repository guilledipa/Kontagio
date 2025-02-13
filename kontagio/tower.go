package kontagio

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	towerCost = 50
)

type Tower struct {
	x, y        float64
	attackRange float64
}

func (t *Tower) Update(g *Game) {
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
			break // Only shoot one enemy at a time
		}
	}
}

func (t *Tower) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen,
		float32(t.x-15), float32(t.y-15),
		float32(30), float32(30),
		color.RGBA{0, 0, 255, 255}, false)
}
