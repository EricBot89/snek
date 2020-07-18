package main

import (
	"math/rand"
	"time"
)

//Board the game board, holds food
type Board struct {
	MinX   int
	MinY   int
	MaxX   int
	MaxY   int
	Width  int
	Height int
	Food   [][2]int
}

//NewBoard makes a new board
func NewBoard() Board {
	return Board{
		MinX:   0,
		MinY:   0,
		MaxX:   50,
		MaxY:   25,
		Width:  50,
		Height: 25,
	}
}

func (b *Board) removeFood(fIdx int) {
	b.Food[fIdx] = b.Food[len(b.Food)-1]
	b.Food = b.Food[:len(b.Food)-1]
}

func (b *Board) addFood() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	fx := rand.Intn(b.Width)
	fy := rand.Intn(b.Height)
	b.Food = append(b.Food, [2]int{fx, fy})
}
