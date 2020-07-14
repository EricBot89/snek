package main

import (
	"math/rand"
	"time"
)

type Board struct {
	Min_x  int
	Min_y  int
	Max_x  int
	Max_y  int
	Width  int
	Height int
	Food   [][2]int
}

func NewBoard() Board {
	return Board{
		Min_x:  0,
		Min_y:  0,
		Max_x:  50,
		Max_y:  50,
		Width:  50,
		Height: 50,
	}
}

func (b *Board) remove_food(f_idx int) {
	b.Food[f_idx] = b.Food[len(b.Food)-1]
	b.Food = b.Food[:len(b.Food)-1]
}

func (b *Board) add_food() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	fx := rand.Intn(b.Width)
	fy := rand.Intn(b.Height)
	b.Food = append(b.Food, [2]int{fx, fy})
}
