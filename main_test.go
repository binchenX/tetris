package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBoard(t *testing.T) {
	x := 0
	y := 0
	w := 10
	h := 20

	board := newBoard(x, y, w, h)

	if board.x != x {
		t.Errorf("Expected x coordinate to be %d, but got %d", x, board.x)
	}

	if board.y != y {
		t.Errorf("Expected y coordinate to be %d, but got %d", y, board.y)
	}

	if len(board.Board) != h {
		t.Errorf("Expected board height to be %d, but got %d", h, len(board.Board))
	}

	for _, row := range board.Board {
		if len(row) != w {
			t.Errorf("Expected board width to be %d, but got %d", w, len(row))
		}
	}
}
func TestAddShapeInstance(t *testing.T) {
	b := newBoard(0, 0, 10, 20)
	// add first T to the bottom of the board
	s := ShapeInstance{
		s: T,
		d: Up,
		x: 0,
		y: 18,
	}
	b.addShapeInstance(s)
	if len(b.Shapes) != 1 {
		t.Errorf("Expected number of shapes to be 1, but got %d", len(b.Shapes))
	}

	assert.Equal(t, 1, b.Board[18][0])
	assert.Equal(t, 1, b.Board[18][1])
	assert.Equal(t, 1, b.Board[18][2])
	assert.Equal(t, 1, b.Board[19][1])

	// add first T to the bottom of the board, next to the first T
	b.addShapeInstance(ShapeInstance{
		s: T,
		d: Up,
		x: 3,
		y: 18,
	})
	if len(b.Shapes) != 2 {
		t.Errorf("Expected number of shapes to be 2, but got %d", len(b.Shapes))
	}

	assert.Equal(t, 1, b.Board[18][3])
	assert.Equal(t, 1, b.Board[18][4])
	assert.Equal(t, 1, b.Board[18][5])
	assert.Equal(t, 1, b.Board[19][4])

	// add another T
	b.addShapeInstance(ShapeInstance{
		s: T,
		d: Up,
		x: 6,
		y: 18,
	})
	if len(b.Shapes) != 3 {
		t.Errorf("Expected number of shapes to be 2, but got %d", len(b.Shapes))
	}

	assert.Equal(t, 1, b.Board[18][6])
	assert.Equal(t, 1, b.Board[18][7])
	assert.Equal(t, 1, b.Board[18][8])
	assert.Equal(t, 1, b.Board[19][7])

	// add I to the bottom of the board and to the right most make the second last row full

	// add another T
	b.addShapeInstance(ShapeInstance{
		s: I,
		d: Up,
		x: 9,
		y: 16,
	})
	if len(b.Shapes) != 4 {
		t.Errorf("Expected number of shapes to be 2, but got %d", len(b.Shapes))
	}

	// row 18 is full and thus removed and all above rows are moved down, thus board[18][9] is 1 and every else in row 18 is zero
	for i := 0; i < 9; i++ {
		assert.Equal(t, 0, b.Board[18][i])
	}
	assert.Equal(t, 1, b.Board[18][9])

	assert.Equal(t, 1, b.Board[19][1])
	assert.Equal(t, 1, b.Board[19][4])
	assert.Equal(t, 1, b.Board[19][7])
}

// add test shape I
func TestShape(t *testing.T) {
	// I shape
	assert.Equal(t, 1, I.v[0][0])
	assert.Equal(t, 1, I.v[1][0])
	assert.Equal(t, 1, I.v[2][0])
	assert.Equal(t, 1, I.v[3][0])

	// it's width is 1 and height is 4
	assert.Equal(t, 1, I.width())
	assert.Equal(t, 4, I.height())

	// T shape
	assert.Equal(t, 1, T.v[0][0])
	assert.Equal(t, 1, T.v[0][1])
	assert.Equal(t, 1, T.v[0][2])
	assert.Equal(t, 1, T.v[1][1])

	// it's width is 3 and height is 2
	assert.Equal(t, 3, T.width())
	assert.Equal(t, 2, T.height())

	// L shape
	assert.Equal(t, 1, L.v[0][0])
	assert.Equal(t, 1, L.v[1][0])
	assert.Equal(t, 1, L.v[2][0])
	assert.Equal(t, 1, L.v[2][1])

	// it's width is 2 and height is 3
	assert.Equal(t, 2, L.width())
	assert.Equal(t, 3, L.height())

	// O shape
	assert.Equal(t, 1, O.v[0][0])
	assert.Equal(t, 1, O.v[0][1])
	assert.Equal(t, 1, O.v[1][0])
	assert.Equal(t, 1, O.v[1][1])

	// it's width is 2 and height is 2
	assert.Equal(t, 2, O.width())
	assert.Equal(t, 2, O.height())

	// S shape
	assert.Equal(t, 0, S.v[0][0])
	assert.Equal(t, 1, S.v[0][1])
	assert.Equal(t, 1, S.v[0][2])
	assert.Equal(t, 1, S.v[1][0])

	// it's width is 3 and height is 2
	assert.Equal(t, 3, S.width())
	assert.Equal(t, 2, S.height())

	// Z shape
	assert.Equal(t, 1, Z.v[0][0])
	assert.Equal(t, 1, Z.v[0][1])
	assert.Equal(t, 0, Z.v[0][2])
	assert.Equal(t, 0, Z.v[1][0])

	// it's width is 3 and height is 2
	assert.Equal(t, 3, Z.width())
	assert.Equal(t, 2, Z.height())
}

func TestPlayMoveSound(t *testing.T) {
	loadSounds()
	playClearSound()

	time.Sleep(time.Second * 2)
}
