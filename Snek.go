package main

type Snek struct {
	Head [2]int
	Tail [][2]int
	Len  int
	Dir  string
	Rune rune
	Dead bool
}

func (s *Snek) eat_food() {
	s.Len++
}
