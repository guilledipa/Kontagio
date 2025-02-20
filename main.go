package main

import (
	"log"

	"github.com/guilledipa/Kontagio/kontagio"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := &kontagio.Game{
		Towers:      []*kontagio.Tower{},
		Enemies:     []*kontagio.Enemy{},
		Projectiles: []*kontagio.Projectile{},
		Wave:        0,
		Resource:    100,
		Lives:       kontagio.InitialLives,
		Scene:       kontagio.SceneMenu,
	}
	ebiten.SetWindowSize(kontagio.ScreenWidth, kontagio.ScreenHeight)
	ebiten.SetWindowTitle("Tower Defense")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
