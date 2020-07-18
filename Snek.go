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
func NewSnek() *Snek {
	return &Snek{
		Head: [2]int{3, 3},
		Len:  1,
		Dir:  "R",
		Dead: false,
	}
}

func (s *Snek) eatFood() {
	s.Len++
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
}
