package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Создаем новый объект игры
	game := NewGame()
	// Запускаем игру
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Dice Roll with 2D Animation")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
