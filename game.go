package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"image/color"
	"log"
	"math/rand"
	"os"
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
	options       = []string{"2 игрока", "3 игрока", "4 игрока"}
	options2      = []string{"No", "Yes"}
	playerIndex   = 0
	continueIndex = 1
	lastKeyPress  = time.Now()             // Время последнего нажатия клавиши
	delay         = 200 * time.Millisecond // Задержка между нажатиями
)

// Game struct хранит состояние игры
type Game struct {
	rolling        bool      // флаг, указывающий на то, что кубик в процессе броска
	startTime      time.Time // время начала анимации
	lastFrameTime  time.Time // Время последнего изменения кадра анимации
	startTimeLoose time.Time
	width, height  int
	diceImage1     *ebiten.Image
	count          int
	stage          int
	players        map[int]Player
	currentPlayer  int
	result         []int
	round          int
	numberOfDice   int
	temp           int
	fontFace       font.Face
	loose          bool
	turn           int
}
type Player struct {
	score        map[int]int
	round        int
	numberOfDice int
	loose        bool
}

// NewGame создаёт новую игру и загружает изображения кубика
func NewGame() *Game {
	g := &Game{}
	g.width = screenWidth
	g.height = screenHeight
	g.diceImage1, _, _ = ebitenutil.NewImageFromFile("dice.png")
	g.numberOfDice = 5
	g.result = []int{1, 2, 3, 4, 5}
	g.round = 0
	g.players = make(map[int]Player)

	fontFace, err := loadFontFace("Carlito-Bold.ttf", 16)
	if err != nil {
		log.Fatalf("could not load font: %v", err)
	}
	g.fontFace = fontFace
	// Инициализируем случайное число для броска
	rand.Seed(time.Now().UnixNano())

	return g
}

// Update обновляет состояние игры каждый кадр
func (g *Game) Update() error {
	now := time.Now()
	// Проверяем, прошло ли достаточно времени с момента последнего нажатия
	switch g.round {
	case 0:
		if now.Sub(lastKeyPress) > delay {
			// Перемещение по меню вверх
			if ebiten.IsKeyPressed(ebiten.KeyUp) && playerIndex > 0 {
				playerIndex--
				lastKeyPress = now // Сбрасываем время последнего нажатия
			}

			// Перемещение по меню вниз
			if ebiten.IsKeyPressed(ebiten.KeyDown) && playerIndex < len(options)-1 {
				playerIndex++
				lastKeyPress = now // Сбрасываем время последнего нажатия
			}

			// Подтверждение выбора
			if ebiten.IsKeyPressed(ebiten.KeyEnter) {
				g.stage = 1
				for i := 0; i <= playerIndex+1; i++ {
					g.players[i] = Player{score: make(map[int]int), numberOfDice: 5}
				}
			}
		}
	default:
		if now.Sub(lastKeyPress) > delay {
			// Перемещение по меню вверх
			if ebiten.IsKeyPressed(ebiten.KeyRight) && continueIndex > 0 {
				continueIndex--
				lastKeyPress = now // Сбрасываем время последнего нажатия
			}

			// Перемещение по меню вниз
			if ebiten.IsKeyPressed(ebiten.KeyLeft) && continueIndex < len(options2)-1 {
				continueIndex++
				lastKeyPress = now // Сбрасываем время последнего нажатия
			}

			// Подтверждение выбора
			if ebiten.IsKeyPressed(ebiten.KeyEnter) {

			}
		}
	}

	g.count++

	// Если нажата клавиша пробел и анимация не идет, начинаем бросок кубика
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.rolling {

		if continueIndex == 0 {

			g.changePlayer()
			//g.round = g.players[g.currentPlayer].getPhase()
			continueIndex = 1
		} else {

		}

		g.StartRolling()
		g.round++
		g.turn++
		log.Print(g.turn)
		g.rollDice()

	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.startTimeLoose = time.Now()
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
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
			g.addScore()
		}
	} else {

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
		g.showScore(screen)
		g.StartGame(screen)

	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}
