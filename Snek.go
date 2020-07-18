package main

import (
	"github.com/nsf/termbox-go"
)

//SnekColors possible color values for sneks
var (
	SnekColors = [4]termbox.Attribute{termbox.ColorRed, termbox.ColorMagenta, termbox.ColorCyan, termbox.ColorYellow}
)

//Snek its a snek
type Snek struct {
	Head   [2]int
	Tail   [][2]int
	Len    int
	Dir    string
	Rune   rune
	Dead   bool
	Color  int
	Name   string
	Points int
}

//NewSnek makes a new snek
func NewSnek(name string) *Snek {
	return &Snek{
		Head:   [2]int{3, 3},
		Len:    1,
		Dir:    "R",
		Dead:   false,
		Color:  0,
		Name:   name,
		Points: 0,
	}
}

func (s *Snek) eatFood() {
	s.Len++
	s.Points += 1000
}

func (s *Snek) move(h int, w int) {
	s.Tail = append(s.Tail, s.Head)
	if len(s.Tail) > s.Len {
		s.Tail = s.Tail[1 : s.Len+1]
	}
	switch s.Dir {
	case "U":
		s.Head[1] = (s.Head[1] - 1 + h) % h
	case "D":
		s.Head[1] = (s.Head[1] + 1) % h
	case "L":
		s.Head[0] = (s.Head[0] - 1 + w) % w
	case "R":
		s.Head[0] = (s.Head[0] + 1) % w
	}
	s.Points++
}

func (s *Snek) chSnek(keyPress termbox.Event) {
	if keyPress.Key == termbox.KeyArrowUp && s.Dir != "D" {
		s.Dir = "U"
	}
	if keyPress.Key == termbox.KeyArrowDown && s.Dir != "U" {
		s.Dir = "D"
	}
	if keyPress.Key == termbox.KeyArrowLeft && s.Dir != "R" {
		s.Dir = "L"
	}
	if keyPress.Key == termbox.KeyArrowRight && s.Dir != "L" {
		s.Dir = "R"
	}
	if keyPress.Key == termbox.KeySpace {
		s.Color = (s.Color + 1) % len(SnekColors)
	}
}
