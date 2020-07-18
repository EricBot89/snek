package main

import (
	"sort"
	"strconv"

	"github.com/nsf/termbox-go"
)

func draw(g *GameData) {
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	termWidth, termHeight := termbox.Size()
	offsetX := (termWidth - g.B.Width) / 3
	offsetY := (termHeight - g.B.Height) / 2
	drawPlayerList(g, offsetX, offsetY)
	drawBoard(g, offsetX, offsetY)
	drawGame(g, offsetX, offsetY)
	termbox.Flush()
}

func drawPlayerList(g *GameData, x int, y int) {
	printWord("PLAYERS", x+g.B.Width+4, y, termbox.ColorWhite)
	printWord("│", x+g.B.Width+16, y, termbox.ColorWhite)
	printWord("POINTS", x+g.B.Width+17, y, termbox.ColorWhite)
	sort.SliceStable(g.Sneks, func(i, j int) bool { return g.Sneks[i].Points > g.Sneks[j].Points })
	for idx, s := range g.Sneks {
		printWord(s.Name, x+g.B.Width+4, y+idx+1, SnekColors[s.Color])
		printWord("│", x+g.B.Width+16, y+idx+1, termbox.ColorWhite)
		printWord(strconv.Itoa(s.Points), x+g.B.Width+17, y+idx+1, SnekColors[s.Color])
	}
}

func printWord(word string, x int, y int, c termbox.Attribute) {
	for idx, char := range word {
		termbox.SetCell(x+idx, y, char, c, termbox.ColorBlack)
	}
}

func drawBoard(g *GameData, x int, y int) {
	for i := 0; i < g.B.Width+2; i++ {
		for j := 0; j < g.B.Height+2; j++ {
			var char rune
			char = 0x2593
			switch {
			case i == 0 && j == 0:
				char = 0x2591
			case i == 0 && j == g.B.Height+1:
				char = 0x2591
			case i == g.B.Width+1 && j == 0:
				char = 0x2591
			case i == g.B.Width+1 && j == g.B.Height+1:
				char = 0x2591
			case i == 0 || i == g.B.Width+1:
				char = 0x2591
			case j == 0 || j == g.B.Height+1:
				char = 0x2591
			}
			termbox.SetCell(i+x-1, j+y-1, char, termbox.ColorBlack, termbox.ColorBlue|termbox.ColorRed)
		}
	}
}

func drawGame(g *GameData, x int, y int) {
	for _, s := range g.Sneks {
		termbox.SetCell(s.Head[0]+x, s.Head[1]+y, s.Rune, termbox.ColorBlack, termbox.ColorWhite)
		for _, cell := range s.Tail {
			termbox.SetCell(cell[0]+x, cell[1]+y, ' ', termbox.ColorBlack, SnekColors[s.Color])
		}
	}
	for _, cell := range g.B.Food {
		termbox.SetCell(cell[0]+x, cell[1]+y, ' ', termbox.ColorBlack, termbox.ColorGreen)
	}
}
