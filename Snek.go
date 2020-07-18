package main

//Snek its a snek
type Snek struct {
	Head [2]int
	Tail [][2]int
	Len  int
	Dir  string
	Rune rune
	Dead bool
}

//NewSnek makes a new snek
func NewSnek() Snek {
	return Snek{
		Head: [2]int{3, 3},
		Len:  1,
		Dir:  "R",
		Dead: false,
	}
}

func (s *Snek) eatFood() {
	s.Len++
}
