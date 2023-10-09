package main

import (
	"log"
	"os"
	"regexp"

	"github.com/gdamore/tcell/v2"
)

const textboxPos = 0
const listPos = 1

type FuzzyContentItem struct {
	Content             string
	Id                  int
	LevenshteinDistance int
	Valid               bool
}

type FuzzyFinder struct {
	Renderer     Renderer
	SearchList   []FuzzyContentItem
	FilteredList []FuzzyContentItem
	SearchPos    int
	SearchInput  []rune
}

func textSectionDimensions(w, h int) (int, int, int, int) {
	return 0, h - 2, w, h
}

func listSectionDimensions(w, h int) (int, int, int, int) {
	return 0, 0, w, h - 3
}

func newTextSection(w, h int) Section {
	sx, sy, ex, ey := textSectionDimensions(w, h)
	s := NewSection(sx, sy, ex, ey)
	s.SetCursor(sx+5, ey-1, Typing)
	return s
}

func newListSection(w, h int) Section {
	sx, sy, ex, ey := listSectionDimensions(w, h)
	s := NewSection(sx, sy, ex, ey)
	s.SetCursor(sx+2, ey-1, Selection)
	return s
}

func NewFuzzyFinder() FuzzyFinder {
	r := NewRenderer()
	x, y := r.Size()
	textbox := newTextSection(x, y)
	textbox.Content = append(textbox.Content, []rune{'>', ' '})
	list := newListSection(x, y)
	r.AddSection(textboxPos, textbox).
		AddSection(listPos, list)
	return FuzzyFinder{
		Renderer:  r,
		SearchPos: 0,
	}
}

func (ff *FuzzyFinder) SetFuzzyItems(in []string) *FuzzyFinder {
	list := make([]FuzzyContentItem, 0)
	os.Stderr.WriteString("setting items")
	for i, s := range in {
		list = append(list, FuzzyContentItem{
			Content:             s,
			LevenshteinDistance: 0,
			Valid:               true,
			Id:                  i,
		})
	}
	ff.SearchList = list
	return ff
}

func (ff FuzzyFinder) FuzzyFind() int {
	ff.recalcList()
	ff.setListContent()
	ff.Draw()
	return ff.listen()
}

func (ff FuzzyFinder) Draw() {
	ff.Renderer.Draw()
}

func (ff FuzzyFinder) recalc() {
	w, h := ff.Renderer.Size()
	ff.Renderer.Sections[listPos].ResizeSection(listSectionDimensions(w, h))
	ff.Renderer.Sections[textboxPos].ResizeSection(textSectionDimensions(w, h))
}

func (ff FuzzyFinder) listen() int {
	for {
		ff.Renderer.Screen.Show()
		switch ev := ff.Renderer.Screen.PollEvent().(type) {
		case *tcell.EventResize:
			ff.recalc()
			ff.Draw()
		case *tcell.EventKey:
			switch ev.Key() {

			//Exit keys
			case tcell.KeyEscape:
			case tcell.KeyCtrlC:
				ff.Renderer.Screen.Fini()
				os.Exit(0)

			case tcell.KeyCtrlL:
				ff.Renderer.Screen.Sync()

			//Adding a rune
			case tcell.KeyRune:
				ff.fuzzyInputHandler(ev.Rune())

			//Prompt movements
			case tcell.KeyRight:
				ff.moveBottomPromptRight()
			case tcell.KeyLeft:
				ff.moveBottomPromptLeft()
			case tcell.KeyUp:
				ff.moveTopPromptUp()
			case tcell.KeyDown:
				ff.moveTopPromptDown()

			//Removing runes
			case tcell.KeyBackspace2:
				ff.backspaceHandler()
			case tcell.KeyDelete:
				ff.deleteHandler()

			//Selecting an entry
			case tcell.KeyEnter:
				ff.Renderer.Screen.Fini()
				item, err := ff.selectItem()
				if err != nil {
					log.Fatal("item was not selected")
				}
				return item
			}
			ff.recalc()
			ff.Draw()
		}
	}
}

func (ff *FuzzyFinder) selectItem() (int, error) {
	return ff.FilteredList[ff.Renderer.Sections[listPos].EndY-ff.Renderer.Sections[listPos].Cursor.YLoc].Id, nil
}

func (ff *FuzzyFinder) filterList() {
	ff.FilteredList = []FuzzyContentItem{}
	for i := 0; i < len(ff.SearchList); i++ {
		item := &ff.SearchList[i]
		valid, _ := regexp.Match(createRegex(ff.SearchInput), []byte(string(item.Content)))
		if valid {
			ff.FilteredList = append(ff.FilteredList, *item)
		}
	}
}

func (ff *FuzzyFinder) recalcList() {
	ff.filterList()
	for i := 0; i < len(ff.FilteredList); i++ {
		item := &ff.FilteredList[i]
		item.LevenshteinDistance = levenshteinDistance([]rune(item.Content), ff.Renderer.Sections[textboxPos].Content[0])
	}
	ff.FilteredList = quickSort(ff.FilteredList, 0, len(ff.FilteredList)-1)
	ff.setListContent()
}

func (ff *FuzzyFinder) setListContent() {
	r := make([][]rune, len(ff.FilteredList))
	for i, v := range ff.FilteredList {
		line := []rune(v.Content)
		r[i] = line
	}
	ff.Renderer.Sections[listPos].Content = r
}

func (ff *FuzzyFinder) charHandler(char rune, pos int) {
	ff.SearchInput = insertRune(ff.SearchInput, pos, char)
	ff.moveBottomPromptRight()
}

func (ff *FuzzyFinder) setTextBoxContent() {
	ff.Renderer.Sections[textboxPos].Content[0] = append([]rune{'>', ' '}, ff.SearchInput...)
}

func (ff *FuzzyFinder) moveBottomPromptRight() {
	if ff.Renderer.Sections[textboxPos].Cursor.XLoc < 5+len(ff.SearchInput) {
		ff.Renderer.Sections[textboxPos].MoveCursorRight(1)
	}
}

func (ff *FuzzyFinder) moveBottomPromptLeft() {
	if ff.Renderer.Sections[textboxPos].Cursor.XLoc > 5 {
		ff.Renderer.Sections[textboxPos].MoveCursorLeft(1)
	}
}

func (ff *FuzzyFinder) moveTopPromptUp() {
	ff.Renderer.Sections[listPos].MoveCursorUp(1)
}

func (ff *FuzzyFinder) moveTopPromptDown() {
	ff.Renderer.Sections[listPos].MoveCursorDown(1)
}

func (ff *FuzzyFinder) backspaceHandler() {
	if len(ff.SearchInput) > 0 && ff.Renderer.Sections[textboxPos].Cursor.XLoc > 5 {
		ff.SearchInput = removeRune(ff.SearchInput, ff.Renderer.Sections[textboxPos].Cursor.XLoc-6)
		ff.setTextBoxContent()
		ff.moveBottomPromptLeft()
		ff.recalcList()
	}
}

func (ff *FuzzyFinder) deleteHandler() {
	stringLen := len(ff.SearchInput)
	if stringLen > 0 && ff.Renderer.Sections[textboxPos].Cursor.XLoc < 5+stringLen {
		ff.SearchInput = removeRune(ff.SearchInput, ff.Renderer.Sections[textboxPos].Cursor.XLoc-5)
		ff.setTextBoxContent()
		ff.recalcList()
	}
}

func (ff *FuzzyFinder) fuzzyInputHandler(char rune) {
	loc := ff.Renderer.Sections[textboxPos].Cursor.XLoc
	ff.charHandler(char, loc)
	ff.setTextBoxContent()
	ff.recalcList()
}

func insertRune(runes []rune, index int, in rune) []rune {
	if index > len(runes) {
		return append(runes, in)
	}
	ret := make([]rune, 0)
	ret = append(ret, runes[:index]...)
	ret = append(ret, in)
	return append(ret, runes[index:]...)
}

func removeRune(runes []rune, index int) []rune {
	ret := make([]rune, 0)
	ret = append(ret, runes[:index]...)
	return append(ret, runes[index+1:]...)
}

func createRegex(runes []rune) string {
	const interleave = "(.*)"
	inStr := string(runes)
	outStr := interleave
	for _, ch := range inStr {
		outStr = outStr + string(ch) + interleave
	}
	return outStr
}

func levenshteinDistance(a []rune, b []rune) int {

	r1, r2 := a, b
	column := make([]int, 1, 64)

	for y := 1; y <= len(r1); y++ {
		column = append(column, y)
	}

	for x := 1; x <= len(r2); x++ {
		column[0] = x

		for y, lastDiag := 1, x-1; y <= len(r1); y++ {
			oldDiag := column[y]
			cost := 0
			if r1[y-1] != r2[x-1] {
				cost = 1
			}
			column[y] = min(column[y]+1, column[y-1]+1, lastDiag+cost)
			lastDiag = oldDiag
		}
	}

	return column[len(r1)]
}

func partition(arr []FuzzyContentItem, low, high int) ([]FuzzyContentItem, int) {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[j].LevenshteinDistance < pivot.LevenshteinDistance {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return arr, i
}

func quickSort(arr []FuzzyContentItem, low, high int) []FuzzyContentItem {
	if low < high {
		var p int
		arr, p = partition(arr, low, high)
		arr = quickSort(arr, low, p-1)
		arr = quickSort(arr, p+1, high)
	}
	return arr
}
