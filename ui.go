package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"io/ioutil"
	"log"
)

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
	fontFace, err := loadFontFace("Carlito-Bold.ttf", float64(size))
	if err != nil {
		log.Fatalf("could not load font: %v", err)
	}
	textColor := color.RGBA{245, 0, 0, 250}
	text.Draw(screen, txt, fontFace, x, y, textColor)
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

func (g *Game) showScore(screen *ebiten.Image) {
	for count, player := range g.players {

		switch count {
		case 0:
			//face := basicfont.Face7x13
			// Отрисовка меню
			x := 50
			y := 50
			txt := ""
			txt += "Игрок 1:\n"
			var sum int
			for i := 0; i < len(player.score); i++ {
				if player.score[i] == 0 {
					continue
				}
				sum += player.score[i]
				txt += fmt.Sprintf("%3d", player.score[i]) + "\n"
			}

			txt += "Сумма:" + fmt.Sprint(sum) + "\n"
			// Если элемент выбран, заливаем его фон другим цветом
			if g.currentPlayer == 0 {
				vector.DrawFilledRect(screen, float32(x)-5, float32(y-17), 100, 20, color.RGBA{0, 0, 200, 255}, true) // Зеленый фон
			}

			// Отрисовка текста опции
			text.Draw(screen, txt, g.fontFace, x, y, color.White)
		case 1:
			//face := basicfont.Face7x13
			// Отрисовка меню
			x := screenWidth - 100
			y := 50
			txt := ""
			txt += "Игрок 2:\n"
			var sum int
			for i := 0; i < len(player.score); i++ {
				if player.score[i] == 0 {
					continue
				}
				sum += player.score[i]
				txt += fmt.Sprintf("%3d", player.score[i]) + "\n"
			}

			txt += "Сумма:" + fmt.Sprint(sum) + "\n"
			// Если элемент выбран, заливаем его фон другим цветом
			if g.currentPlayer == 1 {
				vector.DrawFilledRect(screen, float32(x)-5, float32(y-17), 100, 20, color.RGBA{0, 0, 200, 255}, true) // Зеленый фон
			}
			// Отрисовка текста опции
			//alphaColor := color.RGBA{0, 0, 255, 0}
			text.Draw(screen, txt, g.fontFace, x, y, color.White)

		}
	}

}
