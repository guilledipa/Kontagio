package kontagio

import (
	"image/color"
	"log"
	"math"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	InitialLives = 50
	// Offset para separar a los enemigos que spawnean asi no se amontonan.
	spawnOffset      = 20
	minNumberOfTurns = 2
	maxNumberOfTurns = 4
	SceneMenu        = "menu"
	ScenePlaying     = "playing"
	SceneGameOver    = "gameover"
	gameOverDelay    = 3 * time.Second // Delay before returning to menu after game over
)

func generateProceduralPath() []float64 {
	var currentX, currentY float64
	var path []float64
	// Start at a random Y position on the left edge
	currentY = rand.Float64()*(ScreenHeight-200) + 100 // Keep some padding from top/bottom
	path = append(path, 0, currentY)                   // Start at left edge
	// Determine the number of turns (between 2 and 4)
	numTurns := minNumberOfTurns + rand.IntN(maxNumberOfTurns) // Random number between 2 and 4
	for i := 0; i < numTurns; i++ {
		// Random segment width. Add 10 to avoid very narrow segments.
		segmentWidth := ((float64(ScreenWidth) - currentX) * rand.Float64()) + 10
		currentX += segmentWidth
		path = append(path, currentX, currentY) // Desplazamiento horizontal
		// Choose a random Y position within the screen bounds
		currentY = rand.Float64()*(ScreenHeight-100) + 50 // Keep some padding from top/bottom
		path = append(path, currentX, currentY)
	}
	// End at a random Y position on the right edge
	path = append(path, float64(ScreenWidth), currentY)
	return path
}

type Game struct {
	Towers       []*Tower
	Enemies      []*Enemy
	Projectiles  []*Projectile
	Path         []float64
	Wave         int
	Resource     int
	GameOver     bool
	Lives        int
	Scene        string    // Current game scene
	gameOverTime time.Time // Time when game over occurred
}

func (g *Game) SpawnWave() {
	for i := 0; i < 5+g.Wave*2; i++ {
		offsetX := rand.Float64() * spawnOffset * 5
		offsetY := rand.Float64() * spawnOffset * 5
		// Calcular el path unico para este enemigo relativo a la posición de
		// spawn.
		uniquePath := make([]float64, len(g.Path))
		for j := 0; j < len(g.Path); j += 2 {
			uniquePath[j] = g.Path[j] + offsetX
			uniquePath[j+1] = g.Path[j+1] + offsetY
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
	switch g.Scene {
	case SceneMenu:
		return g.updateMenu()
	case ScenePlaying:
		return g.updatePlaying()
	case SceneGameOver:
		return g.updateGameOver()
	}
	return nil
}

func (g *Game) updateMenu() error {
	// Handle menu selection
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.resetGame()
		g.Path = generateProceduralPath()
		g.Scene = ScenePlaying
	}
	return nil
}

func (g *Game) updatePlaying() error {
	if g.GameOver {
		g.Scene = SceneGameOver
		g.gameOverTime = time.Now()
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
			// Check if the new turret is too close to existing turrets
			canPlace := true
			for _, tower := range g.Towers {
				dx := float64(x) - tower.x
				dy := float64(y) - tower.y
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist < toweMinDistance {
					canPlace = false
					break
				}
			}
			if canPlace {
				g.Towers = append(g.Towers, &Tower{
					x:           float64(x),
					y:           float64(y),
					attackRange: 100,
				})
				g.Resource -= towerCost
			} else {
				log.Println("Cannot place turret: Too close to another turret")
			}
		}
	}
	return nil
}

func (g *Game) updateGameOver() error {
	// Return to menu after delay
	if time.Since(g.gameOverTime) > gameOverDelay {
		g.Scene = SceneMenu
	}
	return nil
}

func (g *Game) resetGame() {
	// Reset game state
	g.Towers = []*Tower{}
	g.Enemies = []*Enemy{}
	g.Projectiles = []*Projectile{}
	g.Wave = 0
	g.Resource = 100
	g.Lives = InitialLives
	g.GameOver = false
	g.Path = generateProceduralPath()
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.Scene {
	case SceneMenu:
		g.drawMenu(screen)
	case ScenePlaying:
		g.drawPlaying(screen)
	case SceneGameOver:
		g.drawGameOver(screen)
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	screen.Fill(color.Black)
	// Draw menu title
	ebitenutil.DebugPrintAt(screen, "Tower Defense", ScreenWidth/2-50, ScreenHeight/2-50)
	// Draw menu options
	ebitenutil.DebugPrintAt(screen, "Start", ScreenWidth/2-20, ScreenHeight/2)
	ebitenutil.DebugPrintAt(screen, "Quit", ScreenWidth/2-20, ScreenHeight/2+30)
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	// Draw towers
	for _, tower := range g.Towers {
		tower.Draw(screen)
	}
	// Draw enemies
	for _, enemy := range g.Enemies {
		enemy.Draw(screen)
	}
	// Draw projectiles
	for _, projectile := range g.Projectiles {
		projectile.Draw(screen)
	}
	// Draw UI
	ebitenutil.DebugPrint(screen, "Wave: "+strconv.Itoa(g.Wave))
	ebitenutil.DebugPrintAt(screen, "Resource: "+strconv.Itoa(g.Resource), 0, 20)
	ebitenutil.DebugPrintAt(screen, "Lives: "+strconv.Itoa(g.Lives), 0, 40)
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Game Over!", ScreenWidth/2-40, ScreenHeight/2)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
