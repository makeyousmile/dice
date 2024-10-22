package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	options       = []string{"2 Players", "3 Players", "4 Players"}
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
	players       int
	result        []int
}
type Player struct {
	result []int
}

// NewGame создаёт новую игру и загружает изображения кубика
func NewGame() *Game {
	g := &Game{}
	g.width = screenWidth
	g.height = screenHeight
	g.diceImage1, _, _ = ebitenutil.NewImageFromFile("dice.png")
	g.diceImage2, _, _ = ebitenutil.NewImageFromFile("dice.png")
	g.result = make([]int, 6)
	g.result[0] = rand.Intn(5)
	g.result[1] = rand.Intn(5)
	g.result[2] = rand.Intn(5)
	g.result[3] = rand.Intn(5)
	g.result[4] = rand.Intn(5)
	g.result[5] = rand.Intn(5)

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
			g.players = selectedIndex
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
		b := []int{56, 99}
		score(screen, b)
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
	g.result[0] = rand.Intn(5)
	g.result[1] = rand.Intn(5)
	g.result[2] = rand.Intn(5)
	g.result[3] = rand.Intn(5)
	g.result[4] = rand.Intn(5)
	g.result[5] = rand.Intn(5)
}

func (g *Game) StartGame(screen *ebiten.Image) {
	if g.rolling {
		for i := 0; i < 5; i++ {
			op1 := &ebiten.DrawImageOptions{}
			op1.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
			op1.GeoM.Translate(300+100*float64(i), screenHeight/2)
			n := (g.count / 5) % frameCount
			sx, sy := frameOX+n*frameWidth, frameOY
			screen.DrawImage(g.diceImage1.SubImage(image.Rect(sx, sy, sx+frameWidth-2, sy+frameHeight)).(*ebiten.Image), op1)
		}

	} else {
		for i := 0; i < 5; i++ {
			// Показываем финальный результат (грань, на которой остановился кубик)
			op1 := &ebiten.DrawImageOptions{}
			op1.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
			op1.GeoM.Translate(300+100*float64(i), screenHeight/2)

			sx, sy := frameOX+g.result[i]*frameWidth, frameOY
			screen.DrawImage(g.diceImage1.SubImage(image.Rect(sx, sy, sx+frameWidth-2, sy+frameHeight)).(*ebiten.Image), op1)
		}

	}
	//ebitenutil.DebugPrint(screen, "Press SPACE to roll the dice")
}
func menu(screen *ebiten.Image) {

	face := basicfont.Face7x13
	// Отрисовка меню
	for i, option := range options {
		x := screenWidth/2 - 50
		y := screenHeight/2 + i*30

		// Если элемент выбран, заливаем его фон другим цветом
		if i == selectedIndex {

			vector.DrawFilledRect(screen, float32(x-10), float32(y-20), 150, 30, color.RGBA{150, 0, 0, 255}, true) // Зеленый фон
		}

		// Отрисовка текста опции
		text.Draw(screen, option, face, x, y, color.White)
	}
}
func score(screen *ebiten.Image, score []int) {
	face := basicfont.Face7x13
	// Отрисовка меню
	x := 10
	y := 30
	txt := ""
	txt += "Player 1:\n"
	var sum int
	for i := range score {
		sum += score[i]
		txt += fmt.Sprintf("%3d", score[i]) + "\n"
	}
	txt += "Result:" + fmt.Sprint(sum) + "\n"
	// Если элемент выбран, заливаем его фон другим цветом
	//vector.DrawFilledRect(screen, float32(x), float32(y), 150, 30, color.RGBA{0, 0, 0, 0}, true) // Зеленый фон

	// Отрисовка текста опции
	text.Draw(screen, txt, face, x, y, color.White)

}
