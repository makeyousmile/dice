package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"image"
	"image/color"
	"math/rand"
	"time"
)

const (
	rollDuration = 1 * time.Second        // продолжительность анимации броска
	frameDelay   = 100 * time.Millisecond // задержка между кадрами анимации

	screenWidth  = 1024
	screenHeight = 768

	frameOX     = 0
	frameOY     = 0
	frameWidth  = 104
	frameHeight = 100
	frameCount  = 6
)

var (
	options       = []string{"1 Player", "2 Players", "3 Players", "4 Players"}
	selectedIndex = 0
	lastKeyPress  = time.Now()             // Время последнего нажатия клавиши
	delay         = 200 * time.Millisecond // Задержка между нажатиями
)

// Game struct хранит состояние игры
type Game struct {
	currentFace1  int       // текущая грань кубика (индекс в массиве изображений)
	currentFace2  int       // текущая грань кубика (индекс в массиве изображений)
	rolling       bool      // флаг, указывающий на то, что кубик в процессе броска
	startTime     time.Time // время начала анимации
	lastFrameTime time.Time // время последнего изменения кадра анимации
	width, height int
	diceImage1    *ebiten.Image
	diceImage2    *ebiten.Image
	count         int
	stage         int
}

// NewGame создаёт новую игру и загружает изображения кубика
func NewGame() *Game {
	g := &Game{}
	g.width = screenWidth
	g.height = screenHeight
	g.diceImage1, _, _ = ebitenutil.NewImageFromFile("dice.png")
	g.diceImage2, _, _ = ebitenutil.NewImageFromFile("dice.png")

	// Создание изображений граней кубика (6 граней)

	// Инициализируем случайное число для броска
	rand.Seed(time.Now().UnixNano())

	return g
}

// Update обновляет состояние игры каждый кадр
func (g *Game) Update() error {
	now := time.Now()

	// Проверяем, прошло ли достаточно времени с момента последнего нажатия
	if now.Sub(lastKeyPress) > delay {
		// Перемещение по меню вверх
		if ebiten.IsKeyPressed(ebiten.KeyUp) && selectedIndex > 0 {
			selectedIndex--
			lastKeyPress = now // Сбрасываем время последнего нажатия
		}

		// Перемещение по меню вниз
		if ebiten.IsKeyPressed(ebiten.KeyDown) && selectedIndex < len(options)-1 {
			selectedIndex++
			lastKeyPress = now // Сбрасываем время последнего нажатия
		}

		// Подтверждение выбора
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			fmt.Printf("Selected option: %s\n", options[selectedIndex])
			g.stage = 1
		}
	}

	g.count++
	// Если нажата клавиша пробел и анимация не идет, начинаем бросок кубика
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.rolling {
		g.StartRolling()
	}

	// Если идет анимация, обновляем результат броска в течение времени rollDuration
	if g.rolling {
		now := time.Now()
		if now.Sub(g.lastFrameTime) > frameDelay {
			g.currentFace1 = rand.Intn(6) // каждые frameDelay мс показываем случайное число
			g.currentFace2 = rand.Intn(6) // каждые frameDelay мс показываем случайное число
			g.lastFrameTime = now
		}
		// Если прошло больше времени, чем rollDuration, останавливаем анимацию
		if now.Sub(g.startTime) > rollDuration {
			g.rolling = false
		}
	}

	return nil
}

// Draw рисует игру на экране
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 120, 0, 50})
	// Если анимация идет, показываем текущую грань кубика
	switch g.stage {
	case 0:
		menu(screen)
	case 1:
		g.StartGame(screen)

	}

	// Выводим текст с подсказкой

}

// Layout определяет размеры окна игры
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}

// StartRolling запускает анимацию броска кубика
func (g *Game) StartRolling() {
	g.rolling = true
	g.startTime = time.Now()
	g.lastFrameTime = g.startTime
}

func (g *Game) StartGame(screen *ebiten.Image) {
	if g.rolling {
		op1 := &ebiten.DrawImageOptions{}
		op1.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
		op1.GeoM.Translate(screenWidth/2-75, screenHeight/2)
		i := (g.count / 5) % frameCount
		sx, sy := frameOX+i*frameWidth, frameOY
		screen.DrawImage(g.diceImage1.SubImage(image.Rect(sx, sy, sx+frameWidth-2, sy+frameHeight)).(*ebiten.Image), op1)

		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
		op2.GeoM.Translate(screenWidth/2+75, screenHeight/2)
		i2 := (g.count / 5) % frameCount
		sx2, sy2 := frameOX+i2*frameWidth, frameOY
		screen.DrawImage(g.diceImage1.SubImage(image.Rect(sx2, sy2, sx2+frameWidth-2, sy2+frameHeight)).(*ebiten.Image), op2)

	} else {
		// Показываем финальный результат (грань, на которой остановился кубик)
		op1 := &ebiten.DrawImageOptions{}
		op1.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
		op1.GeoM.Translate(screenWidth/2-75, screenHeight/2)

		sx, sy := frameOX+g.currentFace1*frameWidth, frameOY
		screen.DrawImage(g.diceImage1.SubImage(image.Rect(sx, sy, sx+frameWidth-2, sy+frameHeight)).(*ebiten.Image), op1)

		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
		op2.GeoM.Translate(screenWidth/2+75, screenHeight/2)

		sx2, sy2 := frameOX+g.currentFace2*frameWidth, frameOY
		screen.DrawImage(g.diceImage1.SubImage(image.Rect(sx2, sy2, sx2+frameWidth-2, sy2+frameHeight)).(*ebiten.Image), op2)
	}
	ebitenutil.DebugPrint(screen, "Press SPACE to roll the dice")
}
func menu(screen *ebiten.Image) {

	face := basicfont.Face7x13

	// Отрисовка меню
	for i, option := range options {
		x := screenWidth/2 - 50
		y := screenHeight/2 + i*30

		// Если элемент выбран, заливаем его фон другим цветом
		if i == selectedIndex {
			ebitenutil.DrawRect(screen, float64(x-10), float64(y-20), 150, 30, color.RGBA{0, 255, 0, 255}) // Зеленый фон
		}

		// Отрисовка текста опции
		text.Draw(screen, option, face, x, y, color.White)
	}
}
