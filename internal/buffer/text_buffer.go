package buffer

import "slices"

type TextBuffer struct {
	Lines     []*GappedTextBuffer
	CursorRow int
	CursorCol int
}

func NewTextBuffer(cap int) *TextBuffer {
	return &TextBuffer{
		Lines:     []*GappedTextBuffer{NewGappedTextBuffer(cap)},
		CursorRow: 0,
		CursorCol: 0,
	}
}

// Get the number of lines in the buffer
func (b *TextBuffer) NumLines() int {
	return len(b.Lines)
}

// Break the current line at the cursor position
func (tb *TextBuffer) BreakLineAtCursor(row, col int) {
	line := tb.Lines[row]
	line.MoveCursorTo(col)

	left := NewGappedTextBuffer(cap(line.data))
	right := NewGappedTextBuffer(cap(line.data))

	copy(left.data, line.data[:line.gapStart])
	left.gapStart = line.gapStart
	left.gapEnd = cap(left.data)

	rightText := line.data[line.gapEnd:]
	copy(right.data[cap(right.data)-len(rightText):], rightText)
	right.gapStart = 0
	right.gapEnd = cap(right.data) - len(rightText)

	tb.Lines[row] = left
	tb.Lines = append(tb.Lines[:row+1], append([]*GappedTextBuffer{right}, tb.Lines[row+1:]...)...)
}

// Join lines at the cursor position
func (tb *TextBuffer) JoinLines(cursorRow int) {
	if cursorRow <= 0 {
		return // can't join at top line
	}

	curr := tb.Lines[cursorRow]
	prev := tb.Lines[cursorRow-1]

	// Move cursor in prev to end for appending
	prev.MoveCursorTo(prev.Len())

	left := curr.data[:curr.gapStart]
	right := curr.data[curr.gapEnd:]

	for _, b := range left {
		prev.Insert(prev.Len(), b)
	}
	for _, b := range right {
		prev.Insert(prev.Len(), b)
	}

	// Remove current line
	tb.Lines = slices.Delete(tb.Lines, cursorRow, cursorRow+1)
}
