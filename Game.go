package main

import (
	"math/rand"
	"time"
)

type Game struct {
	Sneks map[string]Snek
	b     Board
}

func NewGame() Game {
	g := Game{
		b:     NewBoard(),
		Sneks: map[string]Snek{},
	}
	g.b.add_food()
	g.b.add_food()
	return g
}

func (g *Game) game_tick() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	p := rand.Intn(1000)
	if p < 5 {
		g.b.add_food()
	}
	if len(g.b.Food) > 8 {
		g.b.Food = g.b.Food[1:]
	}
	g.move_sneks()
	g.check_food()
}

func (g *Game) move_sneks() {
	for _, s := range g.Sneks {
		s.Tail = append(s.Tail, s.Head)
		if len(s.Tail) > s.Len {
			s.Tail = s.Tail[1 : s.Len+1]
		}
		switch s.Dir {
		case "U":
			s.Head[1] = (s.Head[1] - 1 + g.b.Height) % g.b.Height
		case "D":
			s.Head[1] = (s.Head[1] + 1) % g.b.Height
		case "L":
			s.Head[0] = (s.Head[0] - 1 + g.b.Width) % g.b.Width
		case "R":
			s.Head[0] = (s.Head[0] + 1) % g.b.Width
		}
	}

}

func (g *Game) check_food() {
	for _, s := range g.Sneks {
		for f_idx, cell := range g.b.Food {
			if s.Head[0] == cell[0] && s.Head[1] == cell[1] {
				s.eat_food()
				g.b.remove_food(f_idx)
			}
		}
	}
}

func (g *Game) check_loss() []string {
	var lost []string
	for name, s := range g.Sneks {
		for _, cell := range s.Tail {
			if s.Head[0] == cell[0] && s.Head[1] == cell[1] {
				lost = append(lost, name)
			}
		}
	}
	return lost
}

func (g *Game) run_snek() {
	for {
		g.game_tick()
		time.Sleep(101 * time.Millisecond)
	}
}
