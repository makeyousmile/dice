package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Создаем новый объект игры
	game := NewGame()
	// Запускаем игру
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowDecorated(false) // Убираем рамки и панель заголовка
	ebiten.SetWindowTitle("Dice Roll")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
