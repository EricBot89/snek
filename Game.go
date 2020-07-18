package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

//Game stuct for running snek
type Game struct {
	Sneks map[string]Snek
	B     Board

	m sync.RWMutex
}

//GameData struct for sending game data over the wire
type GameData struct {
	Sneks []Snek
	B     Board
}

//NewGameData returns a struct for the wire
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

//NewGame Game Constructor
func NewGame() *Game {

	g := Game{
		B:     NewBoard(),
		Sneks: map[string]Snek{},
	}
	g.B.addFood()
	g.B.addFood()
	return &g
}

func (g *Game) gameTick() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	p := rand.Intn(100)
	g.m.Lock()
	if p < 5 {
		g.B.addFood()
	}
	if len(g.B.Food) > 8 {
		g.B.Food = g.B.Food[1:]
	}
	g.moveSneks()
	g.checkFood()
	g.checkLoss()
	g.m.Unlock()
}

func (g *Game) moveSneks() {
	for name := range g.Sneks {
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

func (g *Game) checkFood() {
	for name, s := range g.Sneks {
		for fIdx, cell := range g.B.Food {
			if s.Head[0] == cell[0] && s.Head[1] == cell[1] {
				s = g.Sneks[name]
				s.eatFood()
				g.B.removeFood(fIdx)
				g.Sneks[name] = s
			}
		}
	}
}

func (g *Game) checkLoss() {
	var tailCells [][2]int
	for name := range g.Sneks {
		s := g.Sneks[name]
		for _, cell := range s.Tail {
			tailCells = append(tailCells, cell)
		}
	}
	for name := range g.Sneks {
		s := g.Sneks[name]
		for _, cell := range tailCells {
			if s.Head[0] == cell[0] && s.Head[1] == cell[1] {
				s.Dead = true
				log.Println(name + " Died!")
				break
			}
		}
		g.Sneks[name] = s
	}
}

func (g *Game) runSnek() {
	for {
		g.gameTick()
		time.Sleep(100 * time.Millisecond)
	}
}
