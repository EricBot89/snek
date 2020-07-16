package main

import "github.com/nsf/termbox-go"

func draw(g *GameData) {
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	for _, s := range g.Sneks {
		termbox.SetCell(s.Head[0], s.Head[1], s.Rune, termbox.ColorBlack, termbox.ColorWhite)
		for _, cell := range s.Tail {
			termbox.SetCell(cell[0], cell[1], ' ', termbox.ColorBlack, termbox.ColorRed)
		}
	}
	for _, cell := range g.B.Food {
		termbox.SetCell(cell[0], cell[1], ' ', termbox.ColorBlack, termbox.ColorGreen)
	}
	termbox.Flush()
}
