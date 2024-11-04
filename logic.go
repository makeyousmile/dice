package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math/rand"
	"sort"
	"time"
)

func (g *Game) StartRolling() {
	g.rolling = true
	g.startTime = time.Now()
	g.lastFrameTime = g.startTime

}

func (g *Game) StartGame(screen *ebiten.Image) {

	if g.rolling {
		g.ShowAnimateDices(screen)

	} else {
		g.xy = 0
		g.ShowDices(screen)
		if g.round > 0 {

			//g.players[g.currentPlayer].showScore[g.round-1] = g.calculateScore(g.result)

			g.numberOfDice = g.temp
			menuContunue(screen)

		}

	}
	if g.loose {
		g.showText(screen, 410, 100, 50, "Переход хода")
	}
	//ebitenutil.DebugPrint(screen, "Press SPACE to roll the dice")
	if time.Since(g.startTimeLoose) > 2*time.Second {

		g.loose = false
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
		g.removeScore()
		g.changePlayer()

	} else {
		g.temp = g.numberOfDice - dices
		if g.temp == 0 {
			g.temp = 5
			g.numberOfDice = 5
		}
	}

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

	g.numberOfDice = 5
	if g.currentPlayer == 0 {
		g.currentPlayer = 1
	} else {
		g.currentPlayer = 0
	}
	g.round = g.players[g.currentPlayer].getPhase()
	player := g.players[g.currentPlayer]
	log.Print(player.round)

	g.turn = 0
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
		//g.players[g.currentPlayer].showScore[g.round-1] = g.calculateScore()
	}
}

func (g *Game) removeScore() {
	player := g.players[g.currentPlayer].score
	n := len(player)
	for i := 1; i < g.turn; i++ {

		log.Print(player, n-i)
		delete(player, n-i)

	}

}
