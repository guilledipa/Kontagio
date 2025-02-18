package kontagio

import (
	"math/rand/v2"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	InitialLives = 50
	// Offset para separar a los enemigos que spawnean asi no se amontonan.
	spawnOffset = 20
)

type Game struct {
	Towers      []*Tower
	Enemies     []*Enemy
	Projectiles []*Projectile
	Wave        int
	Resource    int
	GameOver    bool
	Lives       int
}

func (g *Game) SpawnWave() {
	for i := 0; i < 5+g.Wave*2; i++ {
		offsetX := rand.Float64() * spawnOffset * 5
		offsetY := rand.Float64() * spawnOffset * 5
		// Calcular el path unico para este enemigo relativo a la posición de
		// spawn.
		uniquePath := make([]float64, len(mainPath))
		for j := 0; j < len(mainPath); j += 2 {
			uniquePath[j] = mainPath[j] + offsetX
			uniquePath[j+1] = mainPath[j+1] + offsetY
		}
		g.Enemies = append(g.Enemies, &Enemy{
			x:       uniquePath[0],
			y:       uniquePath[1],
			health:  5 + g.Wave,
			speed:   1 + float64(g.Wave)*0.2,
			pathIdx: 0,
			path:    uniquePath,
		})
	}
}

func (g *Game) RemoveDeadEnemies() []*Enemy {
	var alive []*Enemy
	for _, enemy := range g.Enemies {
		if enemy.health > 0 {
			alive = append(alive, enemy)
		} else {
			g.Resource += 10 // Recompensa, debería ser un valor variable.
		}
	}
	return alive
}

func (g *Game) RemoveEnemy(enemy *Enemy) []*Enemy {
	var remaining []*Enemy
	for _, e := range g.Enemies {
		if e != enemy {
			remaining = append(remaining, e)
		}
	}
	return remaining
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
		// enemy.x+10 is the right edge of the enemy (tengo que mover esto a
		// una variable)
		if enemy.x+10 > ScreenWidth {
			// Enemy reached the end of the path
			g.Lives--
			if g.Lives <= 0 {
				g.GameOver = true
			}
			g.Enemies = g.RemoveEnemy(enemy)
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
	ebitenutil.DebugPrintAt(screen, "Lives: "+strconv.Itoa(g.Lives), 0, 40)

	if g.GameOver {
		ebitenutil.DebugPrintAt(screen, "Game Over!", ScreenWidth/2-40, ScreenHeight/2)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
