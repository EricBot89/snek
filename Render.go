package main

import "github.com/nsf/termbox-go"

func draw(g *GameData) {
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlue)
	termWidth, termHeight := termbox.Size()
	offsetX := (termWidth - g.B.Width) / 3
	offsetY := (termHeight - g.B.Height) / 2
	drawBoard(g, offsetX, offsetY)
	drawGame(g, offsetX, offsetY)
	termbox.Flush()
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
			termbox.SetCell(cell[0]+x, cell[1]+y, ' ', termbox.ColorBlack, termbox.ColorRed)
		}
	}
	for _, cell := range g.B.Food {
		termbox.SetCell(cell[0]+x, cell[1]+y, ' ', termbox.ColorBlack, termbox.ColorGreen)
	}
}
