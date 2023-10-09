package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type CursorType int64

const (
	Undefined CursorType = iota
	Typing
	Selection
)

type SectionCursor struct {
	XLoc int
	YLoc int
	Type CursorType
}

type Section struct {
	StartX      int
	StartY      int
	EndX        int
	EndY        int
	Content     [][]rune
	Cursor      SectionCursor
	TextStyle   tcell.Style
	BorderStyle tcell.Style
}

type Renderer struct {
	Sections []Section
	Screen   tcell.Screen
}

func NewCursor(x, y int, t CursorType) SectionCursor {
	return SectionCursor{
		XLoc: x,
		YLoc: y,
		Type: t,
	}
}

func NewSection(startx, starty, endx, endy int) Section {
	return Section{
		StartX: startx,
		StartY: starty,
		EndX:   endx,
		EndY:   endy,
	}
}

func (s *Section) AppendLineContent(c []rune) {
	s.Content = append(s.Content, c)
}

func (s *Section) SetContent(c [][]rune) {
	s.Content = c
}

func (s *Section) SetTextStyle(fg, bg tcell.Color) {
	s.TextStyle = tcell.StyleDefault.
		Foreground(fg).
		Background(bg)
}

func (s *Section) SetCursor(x, y int, t CursorType) {
	s.Cursor = SectionCursor{
		XLoc: x,
		YLoc: y,
		Type: t,
	}
}

func (s *Section) SetBorderStyle(fg, bg tcell.Color) {
	s.BorderStyle = tcell.StyleDefault.
		Foreground(fg).
		Background(bg)
}

func (s *Section) ResizeSection(startx, starty, endx, endy int) {
	s.StartX = startx
	s.StartY = starty
	s.EndX = endx
	s.EndY = endy
}

func (s Section) Boundaries() (bsx int, bsy int, bex int, bey int) {
	return s.StartX, s.StartY, s.EndX, s.EndY
}

func (s Section) Draw(screen tcell.Screen) {
	drawBox(screen, s.StartX, s.StartY, s.EndX, s.EndY, s.BorderStyle)
	for i, line := range s.Content {
		emitStr(screen, s, s.StartX+3, s.EndY-i-1, string(line))
	}
	if s.Cursor != (SectionCursor{}) {
		if s.Cursor.Type == Typing {
			screen.ShowCursor(s.Cursor.XLoc, s.Cursor.YLoc)
		} else if s.Cursor.Type == Selection {
			emitStr(screen, s, s.Cursor.XLoc, s.Cursor.YLoc, ">")
		}
	}
}

func (s *Section) MoveCursorUp(m int) {
	if s.Cursor != (SectionCursor{}) {
		if s.Cursor.YLoc-m > s.StartY {
			s.Cursor.YLoc -= m
		} else {
			s.Cursor.YLoc = 1
		}
	}
}

func (s *Section) MoveCursorDown(m int) {
	if s.Cursor != (SectionCursor{}) {
		if s.Cursor.YLoc+m < s.EndY {
			s.Cursor.YLoc += m
		} else {
			s.Cursor.YLoc = s.EndY - 1
		}
	}
}

func (s *Section) MoveCursorRight(m int) {
	if s.Cursor != (SectionCursor{}) {
		if s.Cursor.XLoc+m < s.EndX {
			s.Cursor.XLoc += m
		} else {
			s.Cursor.XLoc = s.EndX - 1
		}
	}
}

func (s *Section) MoveCursorLeft(m int) {
	if s.Cursor != (SectionCursor{}) {
		if s.Cursor.XLoc-m > 0 {
			s.Cursor.XLoc -= m
		} else {
			s.Cursor.XLoc = 0
		}
	}
}

func NewRenderer() Renderer {
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	return Renderer{
		Screen: s,
	}
}

func (r *Renderer) AddSection(index int, s Section) *Renderer {
	if index > len(r.Sections) {
		r.Sections = append(r.Sections, s)
	}
	temp := make([]Section, 0)
	temp = append(temp, r.Sections[:index]...)
	temp = append(temp, s)
	r.Sections = append(temp, r.Sections[index:]...)
	return r
}

func (r Renderer) Size() (w int, h int) {
	w, h = r.Screen.Size()
	return w - 1, h - 1
}

func (r Renderer) Draw() {
	for _, s := range r.Sections {
		s.Draw(r.Screen)
	}
	r.Screen.Show()
}

func emitStr(s tcell.Screen, section Section, x, y int, str string) {
	bsx, bsy, bex, bey := section.Boundaries()

	if y < bey && y > bsy {
		for _, c := range str {
			var comb []rune
			w := runewidth.RuneWidth(c)
			if w == 0 {
				comb = []rune{c}
				c = ' '
				w = 1
			}
			if x < bex && x > bsx {
				s.SetContent(x, y, c, comb, section.TextStyle)
			}
			x += w
		}
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
}
