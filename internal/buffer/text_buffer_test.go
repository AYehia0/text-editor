package buffer

import (
	"testing"
)

func TestBreakLineAtCursor(t *testing.T) {
	tb := NewTextBuffer(20)

	// Insert "hello world" into the first line
	for i, ch := range "hello world" {
		tb.Lines[0].Insert(i, byte(ch))
	}

	// Break the line at index 5 (after "hello")
	tb.BreakLineAtCursor(0, 5)

	// Expect two lines now
	if tb.NumLines() != 2 {
		t.Fatalf("Expected 2 lines, got %d", tb.NumLines())
	}

	// Validate content of the split lines
	line1 := flatten(tb.Lines[0])
	line2 := flatten(tb.Lines[1])
	if line1 != "hello" || line2 != " world" {
		t.Errorf("Expected split lines 'hello' and ' world', got '%s' and '%s'", line1, line2)
	}
}

func TestJoinLines(t *testing.T) {
	tb := &TextBuffer{
		Lines: []*GappedTextBuffer{
			NewGappedTextBuffer(20),
			NewGappedTextBuffer(20),
		},
		CursorRow: 1,
		CursorCol: 0,
	}

	for i, ch := range "hello" {
		tb.Lines[0].Insert(i, byte(ch))
	}
	for i, ch := range " world" {
		tb.Lines[1].Insert(i, byte(ch))
	}

	tb.JoinLines(tb.CursorRow)

	if tb.NumLines() != 1 {
		t.Fatalf("Expected 1 line after join, got %d", tb.NumLines())
	}

	result := flatten(tb.Lines[0])
	expected := "hello world"
	if result != expected {
		t.Errorf("Expected line to be '%s', got '%s'", expected, result)
	}

	// ✅ We no longer expect JoinLines to update cursor
	// So let's just assert that cursor remains untouched
	if tb.CursorRow != 1 || tb.CursorCol != 0 {
		t.Errorf("Expected cursor to remain at (1,0), got (%d,%d)", tb.CursorRow, tb.CursorCol)
	}
}

func TestJoinLinesNotAtStart(t *testing.T) {
	tb := &TextBuffer{
		Lines: []*GappedTextBuffer{
			NewGappedTextBuffer(20),
			NewGappedTextBuffer(20),
		},
		CursorRow: 1,
		CursorCol: 3, // Not at start
	}

	for i, ch := range "line1" {
		tb.Lines[0].Insert(i, byte(ch))
	}
	for i, ch := range "line2" {
		tb.Lines[1].Insert(i, byte(ch))
	}

	// Simulate conditional join (like editor does)
	if tb.CursorCol == 0 {
		tb.JoinLines(tb.CursorRow)
	}

	if tb.NumLines() != 2 {
		t.Errorf("JoinLines should not have modified lines; expected 2, got %d", tb.NumLines())
	}
}
