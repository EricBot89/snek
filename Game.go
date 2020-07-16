package main

import (
	"math/rand"
	"sync"
	"time"
)

type Game struct {
	Sneks map[string]Snek
	B     Board

	m sync.RWMutex
}

type GameData struct {
	Sneks []Snek
	B     Board
}

func NewGameData(g *Game) GameData {
	var sneks []Snek
	for _, Snek := range g.Sneks {
		sneks = append(sneks, Snek)
	}
	return GameData{
		B:     g.B,
		Sneks: sneks,
	}
}

func NewGame() *Game {
	g := Game{
		B:     NewBoard(),
		Sneks: map[string]Snek{},
	}
	g.B.add_food()
	g.B.add_food()
	return &g
}

func (g *Game) game_tick() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	p := rand.Intn(1000)
	g.m.Lock()
	if p < 5 {
		g.B.add_food()
	}
	if len(g.B.Food) > 8 {
		g.B.Food = g.B.Food[1:]
	}
	g.move_sneks()
	g.check_food()
	g.m.Unlock()
}

func (g *Game) move_sneks() {
	for name, _ := range g.Sneks {
		s := g.Sneks[name]
		s.Tail = append(s.Tail, s.Head)
		if len(s.Tail) > s.Len {
			s.Tail = s.Tail[1 : s.Len+1]
		}
		switch s.Dir {
		case "U":
			s.Head[1] = (s.Head[1] - 1 + g.B.Height) % g.B.Height
		case "D":
			s.Head[1] = (s.Head[1] + 1) % g.B.Height
		case "L":
			s.Head[0] = (s.Head[0] - 1 + g.B.Width) % g.B.Width
		case "R":
			s.Head[0] = (s.Head[0] + 1) % g.B.Width
		}
		g.Sneks[name] = s
	}

}

func (g *Game) check_food() {
	for name, s := range g.Sneks {
		for f_idx, cell := range g.B.Food {
			if s.Head[0] == cell[0] && s.Head[1] == cell[1] {
				s = g.Sneks[name]
				s.eat_food()
				g.B.remove_food(f_idx)
				g.Sneks[name] = s
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
		time.Sleep(100 * time.Millisecond)
	}
}
