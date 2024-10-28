package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"sort"
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

	fontFace, err := loadFontFace("fox.ttf", 16)
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
			g.round = g.players[g.currentPlayer].getPhase()
			continueIndex = 1
		} else {

		}

		g.StartRolling()
		g.round++
		g.rollDice()

		log.Print(g.result)
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.startTimeLoose = time.Now()
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
		g.score2(screen)
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
		g.ShowAnimateDices(screen)

	} else {
		g.ShowDices(screen)
		if g.round > 0 {

			//g.players[g.currentPlayer].score[g.round-1] = g.calculateScore(g.result)

			g.numberOfDice = g.temp
			menuContunue(screen)

		}

	}
	if g.loose {
		g.showText(screen, 330, 100, 50, "Proebal")
	}
	//ebitenutil.DebugPrint(screen, "Press SPACE to roll the dice")
	if time.Since(g.startTimeLoose) > 2*time.Second {

		g.loose = false
	}
}
func menu(screen *ebiten.Image) {

	face := basicfont.Face7x13
	// Отрисовка меню
	for i, option := range options {
		x := screenWidth/2 - 50
		y := screenHeight/2 + i*30

		// Если элемент выбран, заливаем его фон другим цветом
		if i == playerIndex {

			vector.DrawFilledRect(screen, float32(x-10), float32(y-20), 150, 30, color.RGBA{150, 0, 0, 255}, true) // Зеленый фон
		}

		// Отрисовка текста опции
		text.Draw(screen, option, face, x, y, color.White)
	}
}
func menuContunue(screen *ebiten.Image) {

	face := basicfont.Face7x13
	text.Draw(screen, "Continue?", face, screenWidth/2-55, screenHeight/2+70, color.White)
	// Отрисовка меню-
	for i, option := range options2 {
		x := screenWidth/2 - i*100
		y := screenHeight/2 + 100

		// Если элемент выбран, заливаем его фон другим цветом
		if i == continueIndex {
			vector.DrawFilledRect(screen, float32(x-10), float32(y-20), 50, 30, color.RGBA{150, 0, 0, 255}, true) // Зеленый фон
		}

		// Отрисовка текста опции
		text.Draw(screen, option, face, x, y, color.White)
	}
}

func (g *Game) score2(screen *ebiten.Image) {
	for count, player := range g.players {

		switch count {
		case 0:
			//face := basicfont.Face7x13
			// Отрисовка меню
			x := 50
			y := 50
			txt := ""
			txt += "Player 1:\n"
			var sum int
			for i := 0; i < len(player.score); i++ {
				if player.score[i] == 0 {
					continue
				}
				sum += player.score[i]
				txt += fmt.Sprintf("%3d", player.score[i]) + "\n"
			}

			txt += "Result:" + fmt.Sprint(sum) + "\n"
			// Если элемент выбран, заливаем его фон другим цветом
			if g.currentPlayer == 0 {
				vector.DrawFilledRect(screen, float32(x)-5, float32(y-17), 100, 20, color.RGBA{0, 0, 200, 0}, true) // Зеленый фон
			}

			// Отрисовка текста опции
			text.Draw(screen, txt, g.fontFace, x, y, color.White)
		case 1:
			//face := basicfont.Face7x13
			// Отрисовка меню
			x := screenWidth - 250
			y := 50
			txt := ""
			txt += "Player 2:\n"
			var sum int
			for i := 0; i < len(player.score); i++ {
				sum += player.score[i]
				txt += fmt.Sprintf("%3d", player.score[i]) + "\n"
			}

			txt += "Result:" + fmt.Sprint(sum) + "\n"
			// Если элемент выбран, заливаем его фон другим цветом
			if g.currentPlayer == 1 {
				vector.DrawFilledRect(screen, float32(x)-5, float32(y-17), 100, 20, color.RGBA{0, 0, 200, 0}, true) // Зеленый фон
			}
			// Отрисовка текста опции
			//alphaColor := color.RGBA{0, 0, 255, 0}
			text.Draw(screen, txt, g.fontFace, x, y, color.White)

		}
	}

}

func countDice(dice []int) map[int]int {
	counts := make(map[int]int)
	for _, d := range dice {
		counts[d]++
	}
	return counts
}

// Проверка на специальные комбинации (1,2,3,4,5) или (2,3,4,5,6)
func checkSpecialCombos(dice []int) int {
	sort.Ints(dice)
	if len(dice) == 5 {
		if equalSlices(dice, []int{1, 2, 3, 4, 5}) {
			return 125
		}
		if equalSlices(dice, []int{2, 3, 4, 5, 6}) {
			return 250
		}
	}
	return 0
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (g *Game) calculateScore() int {
	counts := countDice(g.result)
	score := 0
	dices := 0

	// Проверяем на специальные комбинации
	score += checkSpecialCombos(g.result)

	// Подсчет очков для троек, четверок и пятерок одинаковых значений
	for value, count := range counts {
		switch value {
		case 1:
			if count == 3 {
				score += 100
				dices += 3
			} else if count == 4 {
				score += 200
				dices += 4
			} else if count == 5 {
				score += 1000
			} else {
				score += count * 10 // за каждую 1 - 10 очков
				dices += 1 * count
			}
		case 2:
			if count == 3 {
				dices += 3
				score += 20
			} else if count == 4 {
				dices += 4
				score += 40
			} else if count == 5 {
				score += 200
			}
		case 3:
			if count == 3 {
				score += 30
				dices += 3
			} else if count == 4 {
				score += 60
				dices += 4
			} else if count == 5 {
				score += 300
			}
		case 4:
			if count == 3 {
				score += 40
				dices += 3
			} else if count == 4 {
				score += 80
				dices += 4
			} else if count == 5 {
				score += 400
			}
		case 5:
			if count == 3 {
				dices += 3
				score += 50
			} else if count == 4 {
				score += 100
				dices += 4
			} else if count == 5 {
				score += 500
			} else {
				score += count * 5
				dices += 1 * count // за каждую 5 - 5 очков
			}
		case 6:
			if count == 3 {
				dices += 3
				score += 60
			} else if count == 4 {
				dices += 4
				score += 120
			} else if count == 5 {
				score += 600
			}
		}
	}
	if score == 0 {
		g.temp = 5
		g.loose = true
		g.changePlayer()
	} else {
		g.temp = g.numberOfDice - dices
		if g.temp == 0 {
			g.temp = 5
			g.numberOfDice = 5
		}
	}

	log.Print("calculate score", score)
	return score
}
func (g *Game) rollDice() {
	//dice := make([]int, g.numberOfDice)
	//for i := range dice {
	//	dice[i] = rand.Intn(6) + 1
	//}
	var dice []int
	for i := 1; i <= g.numberOfDice; i++ {
		dice = append(dice, rand.Intn(6)+1)
	}

	g.result = dice
}

func (g *Game) changePlayer() {

	log.Print(g.round)
	g.numberOfDice = 5
	if g.currentPlayer == 0 {
		g.currentPlayer = 1
	} else {
		g.currentPlayer = 0
	}
	g.round = g.players[g.currentPlayer].getPhase()
}
func (p Player) getPhase() int {
	var count int
	for _, _ = range p.score {
		count++
	}
	return count
}

func (g *Game) rollingAnimation() {
	for i := 0; i < g.numberOfDice; i++ {

	}
}
func (g *Game) addScore() bool {
	score := g.calculateScore()
	if score != 0 {
		g.players[g.currentPlayer].score[g.round-1] = score

		return true
	} else {

		return false
		//g.players[g.currentPlayer].score[g.round-1] = g.calculateScore()
	}
}
func (g *Game) ShowDices(screen *ebiten.Image) {

	for i, _ := range g.result {
		op1 := &ebiten.DrawImageOptions{}
		op1.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
		op1.GeoM.Translate(300+100*float64(i), screenHeight/2)

		sx, sy := frameOX+(g.result[i]-1)*frameWidth, frameOY
		screen.DrawImage(g.diceImage1.SubImage(image.Rect(sx, sy, sx+frameWidth-2, sy+frameHeight)).(*ebiten.Image), op1)
	}
}
func (g *Game) ShowAnimateDices(screen *ebiten.Image) {

	for i := 0; i < g.numberOfDice; i++ {
		op1 := &ebiten.DrawImageOptions{}
		op1.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
		op1.GeoM.Translate(300+100*float64(i), screenHeight/2)
		n := (g.count / 5) % frameCount
		sx, sy := frameOX+n*frameWidth, frameOY
		screen.DrawImage(g.diceImage1.SubImage(image.Rect(sx, sy, sx+frameWidth-2, sy+frameHeight)).(*ebiten.Image), op1)
	}
}
func (g Game) showChangePlayer() {

}
func loadFontFace(path string, size float64) (font.Face, error) {
	// Чтение файла шрифта
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read font: %w", err)
	}

	// Парсинг шрифта
	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse font: %w", err)
	}

	// Загрузка шрифта с нужным размером
	const dpi = 72
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create font face: %w", err)
	}

	return face, nil
}

func (g *Game) showText(screen *ebiten.Image, x, y, size int, txt string) {
	fontFace, err := loadFontFace("fox.ttf", float64(size))
	if err != nil {
		log.Fatalf("could not load font: %v", err)
	}
	text.Draw(screen, txt, fontFace, x, y, color.White)
}
