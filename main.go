package main

import (
	"log"

	"github.com/guilledipa/Kontagio/kontagio"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := &kontagio.Game{
		Resource: kontagio.InitialResources,
		Lives:    kontagio.InitialLives,
		Scene:    kontagio.SceneMenu,
	}
	ebiten.SetWindowSize(kontagio.ScreenWidth, kontagio.ScreenHeight)
	ebiten.SetWindowTitle("Kontagio: Tower Defense")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
