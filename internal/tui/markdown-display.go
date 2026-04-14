package tui

import (
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type line struct {
	lineNo int
	lines  [][]string
}

type MDViewer struct {
	Renderer Renderer
	Content  string
	lines    []string
	display  []line
	width    int
}

func mdSectionDimensions(w, h int) (int, int, int, int) {
	return 0, 0, w, h
}

func newMDSection(w, h int) *Section {
	sx, sy, ex, ey := mdSectionDimensions(w, h)
	s := NewSection().
		SetStartX(sx).
		SetEndX(ex).
		SetStartY(sy).
		SetEndY(ey)
	s.SetCursor(sx+2, ey-1, Typing)
	s.Cursor.SetCursorXBoundary(sx+1, ex-1)
	s.Cursor.SetCursorYBoundary(sy+1, ey-1)
	return s
}

func NewMDViewer(content string) MDViewer {
	r := NewRenderer()
	x, y := r.Size()
	section := newMDSection(x, y)
	r.AddSection(0, *section)
	return MDViewer{
		Renderer: r,
		Content:  content,
		width:    x - 4,
	}
}

func (m *MDViewer) Display() {
	m.calcLines()
	m.calcDisplay()
	m.listen()
}

func (m *MDViewer) calcLines() {
	m.lines = strings.Split(m.Content, "\n")
}

func (m *MDViewer) calcDisplay() {
	for i, l := range m.lines {
		displayLines := line{
			lineNo: i + 1,
		}
		var temp []string
		var lenCount int

		for _, w := range strings.Split(l, " ") {
			if lenCount+len(w) >= m.width {
				displayLines.lines = append(displayLines.lines, temp)
				temp = []string{w}
				lenCount = len(w)
			} else {
				temp = append(temp, w)
				lenCount += len(w) + 1
			}
		}
		// temp = append(temp, "\n") // append new line at the end that was there originally
		if len(temp) > 0 {
			displayLines.lines = append(displayLines.lines, temp)
		}
		m.display = append(m.display, displayLines)
	}
	var content [][]rune
	for i := len(m.display) - 1; i >= 0; i-- {
		dl := m.display[i]
		for j := len(dl.lines) - 1; j >= 0; j-- {
			lineWords := dl.lines[j]
			joined := strings.Join(lineWords, " ")
			content = append(content, []rune(joined))
		}
	}
	m.Renderer.Sections[0].SetContent(content)
}

func (m *MDViewer) draw() {
	m.Renderer.Draw()
}

func (m *MDViewer) resize() {
	w, h := m.Renderer.Size()
	m.width = w - 4

	sx, sy, ex, ey := mdSectionDimensions(w, h)
	m.Renderer.Sections[0].ResizeSection(sx, sy, ex, ey)
	m.Renderer.Sections[0].Cursor.SetCursorXBoundary(sx+1, ex-1)
	m.Renderer.Sections[0].Cursor.SetCursorYBoundary(sy+1, ey-1)

	m.calcDisplay()
	m.draw()
}

func (m MDViewer) listen() {
	for {
		m.draw()
		switch ev := m.Renderer.Screen.PollEvent().(type) {
		case *tcell.EventResize:
			m.resize()
		case *tcell.EventKey:
			switch ev.Key() {

			//Exit keys
			case tcell.KeyEscape:
				m.Renderer.Screen.Fini()
				os.Exit(0)
			case tcell.KeyCtrlC:
				m.Renderer.Screen.Fini()
				os.Exit(0)

			case tcell.KeyCtrlL:
				m.Renderer.Screen.Sync()

			//Prompt movements
			case tcell.KeyRight:
				m.movePromptRight()
			case tcell.KeyLeft:
				m.movePromptLeft()
			case tcell.KeyUp:
				m.movePromptUp()
			case tcell.KeyDown:
				m.movePromptDown()
				m.draw()
			}
		}
	}
}

func (m *MDViewer) movePromptRight() {
	m.Renderer.Sections[0].MoveCursorRight(1)
}

func (m *MDViewer) movePromptLeft() {
	m.Renderer.Sections[0].MoveCursorLeft(1)
}

func (m *MDViewer) movePromptUp() {
	m.Renderer.Sections[0].MoveCursorUp(1)
}

func (m *MDViewer) movePromptDown() {
	m.Renderer.Sections[0].MoveCursorDown(1)
}
