package kontagio

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
)

type Game struct {
	Towers      []*Tower
	Enemies     []*Enemy
	Projectiles []*Projectile
	Wave        int
	Resource    int
	GameOver    bool
}

func (g *Game) SpawnWave() {
	for i := 0; i < 5+g.Wave*2; i++ {
		g.Enemies = append(g.Enemies, &Enemy{
			x:       0,
			y:       240,
			health:  5 + g.Wave,
			speed:   1 + float64(g.Wave)*0.2,
			pathIdx: 0,
		})
	}
}

func (g *Game) RemoveDeadEnemies() []*Enemy {
	var alive []*Enemy
	for _, enemy := range g.Enemies {
		if enemy.health > 0 {
			// If the enemy is still alive, add it to the "alive" slice
			alive = append(alive, enemy)
		} else {
			// If the enemy is dead, reward the player with resources
			g.Resource += 10 // Adjust the reward as needed
		}
	}
	return alive
}

func (g *Game) RemoveHitProjectiles() []*Projectile {
	var active []*Projectile
	for _, projectile := range g.Projectiles {
		if projectile.target.health > 0 {
			active = append(active, projectile)
		}
	}
	return active
}

func (g *Game) Update() error {
	if g.GameOver {
		return nil
	}

	// Spawn enemies
	if len(g.Enemies) == 0 {
		g.Wave++
		g.SpawnWave()
	}

	// Update enemies
	for _, enemy := range g.Enemies {
		enemy.Update()
		if enemy.x > ScreenWidth {
			g.GameOver = true
		}
	}

	// Update towers
	for _, tower := range g.Towers {
		tower.Update(g)
	}

	// Update projectiles
	for _, projectile := range g.Projectiles {
		projectile.Update()
	}

	// Remove dead enemies
	g.Enemies = g.RemoveDeadEnemies()

	// Remove projectiles that have hit their target
	g.Projectiles = g.RemoveHitProjectiles()

	// Place towers on mouse click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if g.Resource >= towerCost {
			g.Towers = append(g.Towers, &Tower{
				x:           float64(x),
				y:           float64(y),
				attackRange: 100,
			})
			g.Resource -= towerCost
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw towers
	for _, tower := range g.Towers {
		tower.Draw(screen)
	}

	// Draw enemies
	for _, enemy := range g.Enemies {
		enemy.Draw(screen)
	}

	for _, projectile := range g.Projectiles {
		projectile.Draw(screen)
	}

	// Draw UI
	ebitenutil.DebugPrint(screen, "Wave: "+strconv.Itoa(g.Wave))
	ebitenutil.DebugPrintAt(screen, "Resource: "+strconv.Itoa(g.Resource), 0, 20)

	if g.GameOver {
		ebitenutil.DebugPrintAt(screen, "Game Over!", ScreenWidth/2-40, ScreenHeight/2)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
