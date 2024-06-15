package main

import "github.com/gdamore/tcell/v2"

// View - Draw
const (
	FullBlock         = '█'
	LightShade        = '░'
	MediumShade       = '▒'
	DarkShade         = '▓'
	UpperHalfBlock    = '▀'
	LowerHalfBlock    = '▄'
	LeftHalfBlock     = '▌'
	RightHalfBlock    = '▐'
	Square            = '□'
	BlackSquare       = '■'
	WhiteSmallSquare  = '▫'
	BlackSmallSquare  = '▪'
	MediumBlackSquare = '⬛'
	MediumWhiteSquare = '⬜'
)

// Board Background block
const bgBoxChar = FullBlock

// Board background color
var bgBoxStyle = tcell.StyleDefault.Foreground(tcell.ColorPurple)

// The block for Tetris game shape
const ShapeBlock = BlackSquare

var ShapeStyle = tcell.StyleDefault.Foreground(tcell.ColorSkyblue).Background(tcell.ColorWhite)

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func drawBoardBackGround(s tcell.Screen, x1, y1, x2, y2 int, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// var st tcell.Style
	// ps := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)
	// gs := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGreen)

	// Fill background
	for row := y1; row < y2; row++ {
		for col := x1; col < x2; col++ {
			// alternate the background color for debug
			// if (row+y1)%2 == 0 {
			// 	st = ps
			// }
			// if (col+x1)%2 == 0 {
			// 	st = gs
			// }
			// num := row % 10
			// ascii := rune(num + 48)
			// s.SetContent(col, row, ' ', nil, style)
			s.SetContent(col, row, bgBoxChar, nil, bgBoxStyle)
		}
	}

	// Draw borders
	// for col := x1; col <= x2; col++ {
	// 	s.SetContent(col, y1, tcell.RuneHLine, nil, style)
	// 	s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	// }
	// for row := y1 + 1; row < y2; row++ {
	// 	s.SetContent(x1, row, tcell.RuneVLine, nil, style)
	// 	s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	// }

	// Only draw corners if necessary
	// if y1 != y2 && x1 != x2 {
	// 	s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
	// 	s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
	// 	s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
	// 	s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	// }

	//drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
}

func drawShape(s tcell.Screen, shape Shape, sx, sy int) {
	for y, row := range shape.v {
		for x, cell := range row {
			if cell == 1 {
				s.SetContent(sx+x, sy+y, ShapeBlock, nil, ShapeStyle)
			}
		}
	}
}
