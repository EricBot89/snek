package main

type Snek struct {
	Head [2]int
	Tail [][2]int
	Len  int
	Dir  string
	Rune rune
	Dead bool
}

func NewSnek() Snek {
	return Snek{
		Head: [2]int{3, 3},
		Len:  1,
		Dir:  "R",
		Dead: false,
	}
}

func (s *Snek) eat_food() {
	s.Len++
}
