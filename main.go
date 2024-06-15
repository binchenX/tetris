package main

import (
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Model - Shapes

// define shape struct for teris block, and draw function
type Shape struct {
	v [][]int
}

func (s Shape) width() int {
	return len(s.v[0])
}

func (s Shape) height() int {
	return len(s.v)
}

// define direction of the shape, 0, 90, 180, 270, use enum 0, 1, 2, 3
type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type ShapeInstance struct {
	// the shape
	s Shape
	d Direction
	// location of the shapeW
	x, y int
}

// define the 6 terris shapes, S, Z, L, J, T, I

var I = Shape{
	v: [][]int{
		{1},
		{1},
		{1},
		{1},
	},
}

var T = Shape{
	v: [][]int{
		{1, 1, 1},
		{0, 1, 0},
	},
}

var L = Shape{
	v: [][]int{
		{1, 0},
		{1, 0},
		{1, 1},
	},
}

var O = Shape{
	v: [][]int{
		{1, 1},
		{1, 1},
	},
}

var S = Shape{
	v: [][]int{
		{0, 1, 1},
		{1, 1, 0},
	},
}

var Z = Shape{
	v: [][]int{
		{1, 1, 0},
		{0, 1, 1},
	},
}

// Model - Board

// define game board, it has shapes instances and current active shape
type Board struct {
	Shapes []*ShapeInstance
	Active *ShapeInstance
	// the complete board with all the shapes except the active one
	// when redraw, draw the board first, then draw the active shape
	Board [][]int
	// board top left location and width and height
	x, y int
	w, h int
	// sidebars area
	sideBar Area
	// score
	score int
	// pause status
	isPaused bool
}

// define area type with top left x,y and width and height
type Area struct {
	x, y int
	w, h int
}

const AreaWidth = 25

// create a new board with width w and height h and location x, y
func newBoard(x, y, w, h int) *Board {
	b := &Board{
		x: x,
		y: y,
		w: w,
		h: h,
		sideBar: Area{
			x: x + w + 1,
			y: y,
			w: AreaWidth,
			h: h,
		},
		isPaused: false,
	}
	b.Board = make([][]int, h)
	for i := range b.Board {
		b.Board[i] = make([]int, w)
	}
	b.Active = b.getCurrentShape()
	return b
}

// add a shape instance to the board when it colides
func (b *Board) addShapeInstance(s ShapeInstance) {
	b.Shapes = append(b.Shapes, &s)
	// update the board array
	// or of the shape with the board array
	for y, row := range s.s.v {
		for x, cell := range row {
			if cell == 1 {
				b.Board[s.y+y][s.x+x] = 1
			}
		}
	}

	// check if the row is full, if so, remove it
	for y, row := range b.Board {
		full := true
		for _, cell := range row {
			if cell == 0 {
				full = false
				break
			}
		}
		if full {
			playClearSound()
			b.score = b.score + 100
			b.Board = append(b.Board[:y], b.Board[y+1:]...)
			b.Board = append([][]int{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}, b.Board...)
		}
	}

}

// redraw the board means redraw all the shapes on the board and side bars
func (b *Board) draw(s tcell.Screen) {
	s.Clear()
	b.drawLayout(s)
	// draw the active one
	if b.Active != nil {
		drawShape(s, b.Active.s, b.Active.x, b.Active.y)
	}
	// draw the "freezed" board
	for y, row := range b.Board {
		for x, cell := range row {
			if cell == 1 {
				s.SetContent(x, y, bgBoxChar, nil, tcell.StyleDefault.Foreground(tcell.ColorSkyblue))
			}
		}
	}

	// draw the side bars
	b.drawSideBar(s)
	s.Show()
}

// draw sidebar
func (b *Board) drawSideBar(s tcell.Screen) {
	// draw a box - right side the game area
	// draw keyboard helper text
	// draw current score
	gameStatus := "playing"
	if b.isPaused {
		gameStatus = "paused"
	}

	helps := []string{
		"scores: " + strconv.Itoa(b.score),
		"status:" + gameStatus,
		"----------------",
		" >    - move left",
		" <    - move right",
		"space - rotate",
		"enter - pause/continue",
		"esc   - quit",
	}

	for i, h := range helps {
		drawText(s, b.sideBar.x, b.sideBar.y+i, b.sideBar.x+b.sideBar.w, b.sideBar.y+b.sideBar.h, tcell.StyleDefault, h)
	}
}

func (b *Board) IsLeftMost(i ShapeInstance) bool {
	return i.x == b.x
}

func (b *Board) IsRightMost(i ShapeInstance) bool {
	return i.x+i.s.width() == b.x+b.w
}

func (b *Board) IsBottomMost(i ShapeInstance) bool {
	return i.y+i.s.height() == b.y+b.h
}

// draw layout and backgroud of the board
func (b *Board) drawLayout(s tcell.Screen) {
	drawBoardBackGround(s, b.x, b.y, b.x+b.w, b.y+b.h, "Tetris Game")
}

func (b *Board) moveCurrentShapeDownSpeed() {
	if b.isPaused {
		return
	}
	for i := 0; i < 2; i++ {
		b.moveCurrentShapeDown()
	}
}

// move the current shape down
func (b *Board) moveCurrentShapeDown() {
	if b.isPaused {
		return
	}
	b.getCurrentShape()
	cx := b.Active.x
	cy := b.Active.y
	// attemp to move the instance down. If it will colide, add the shape to the board.
	if b.IsBottomMost(*b.Active) || b.colide(ShapeInstance{b.Active.s, b.Active.d, cx, cy + 1}) {
		b.addShapeInstance(*b.Active)
		b.Active = nil
	} else {
		b.Active.y++
	}
}

func (b *Board) moveCurrentShapeLeft() {
	if b.isPaused {
		return
	}
	b.getCurrentShape()
	if b.IsLeftMost(*b.Active) {
		return
	}
	cx := b.Active.x
	cy := b.Active.y
	if !b.colide(ShapeInstance{b.Active.s, b.Active.d, cx - 1, cy}) {
		b.Active.x--
	}
}

func (b *Board) moveCurrentShapeRight() {
	if b.isPaused {
		return
	}
	b.getCurrentShape()
	if b.IsRightMost(*b.Active) {
		return
	}
	cx := b.Active.x
	cy := b.Active.y
	if !b.colide(ShapeInstance{b.Active.s, b.Active.d, cx + 1, cy}) {
		b.Active.x++
	}
}

// check if the shape instance colides with other shapes in the boards
func (b *Board) colide(si ShapeInstance) bool {
	for y, row := range si.s.v {
		for x, cell := range row {
			if cell == 1 {
				if si.y+y >= b.h || si.x+x >= b.w || b.Board[si.y+y][si.x+x] == 1 {
					return true
				}
			}
		}
	}
	return false
}

// get current active shape
func (b *Board) getCurrentShape() *ShapeInstance {
	if b.Active == nil {
		b.Active = &ShapeInstance{
			s: b.generateRandomShape(),
			d: Up,
			x: b.w / 2,
			y: 0,
		}
	}
	return b.Active
}

var allShapes = []Shape{I, T, L, O, S, Z}

// generate a random shape instance, chose from I, S, Z, L, J, T
func (b *Board) generateRandomShape() Shape {
	return allShapes[rand.Intn(len(allShapes))]
}

// rorate the shape
func rototeShape(s Shape) Shape {
	// rotate the shape
	n := len(s.v)
	m := len(s.v[0])
	v := make([][]int, m)
	for i := 0; i < m; i++ {
		v[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			v[j][n-i-1] = s.v[i][j]
		}
	}
	return Shape{v: v}
}

func (b *Board) rotateCurrentShape() {
	if b.isPaused {
		return
	}
	b.getCurrentShape()
	// rotate the shape and keep x,y the same
	rs := rototeShape(b.Active.s)
	if !b.colide(ShapeInstance{rs, b.Active.d, b.Active.x, b.Active.y}) {
		b.Active.s = rs
	}
}

func (b *Board) TogglePause() {
	// toggle the pause state
	b.isPaused = !b.isPaused
}

// Main - Control
func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.Clear()
	w := flag.Int("w", 40, "width")
	h := flag.Int("h", 30, "height")
	flag.Parse()

	game := newBoard(0, 0, *w, *h)
	loadSounds()

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// put keyboard events in a channel
	eventCh := make(chan tcell.Event, 10)
	go func() {
		for {
			eventCh <- s.PollEvent()
		}
	}()

	// create a timer event
	timer := time.NewTicker(time.Millisecond * 500)
	defer timer.Stop()

	// Event loop
	for {
		// select both time events and keyboard events
		select {
		case <-timer.C:
			// Update screen
			game.moveCurrentShapeDown()
			game.draw(s)
		case ev := <-eventCh:
			// Process event
			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEnter:
					game.TogglePause()
				case tcell.KeyLeft:
					game.moveCurrentShapeLeft()
				case tcell.KeyRight:
					game.moveCurrentShapeRight()
				case tcell.KeyDown:
					game.moveCurrentShapeDownSpeed()
				case tcell.KeyRune:
					switch ev.Rune() {
					case ' ':
						game.rotateCurrentShape()
					}
				case tcell.KeyEscape, tcell.KeyCtrlC:
					return
				}
			}
		}
	}
}
