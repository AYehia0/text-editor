package buffer

import (
	"testing"
)

// Helper
func flatten(b *GappedTextBuffer) string {
	return string(b.data[:b.gapStart]) + string(b.data[b.gapEnd:])
}

func TestInsertBasic(t *testing.T) {
	b := NewGappedTextBuffer(10)
	b.Insert(0, 'H')
	b.Insert(1, 'e')
	b.Insert(2, 'l')
	b.Insert(3, 'l')
	b.Insert(4, 'o')
	if got := flatten(b); got != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", got)
	}
}

func TestInsertMiddle(t *testing.T) {
	b := NewGappedTextBuffer(10)
	b.Insert(0, 'H')
	b.Insert(1, 'l')
	b.Insert(2, 'o')
	b.Insert(1, 'e') // Insert 'e' between H and l
	if got := flatten(b); got != "Helo" {
		t.Errorf("Expected 'Helo', got '%s'", got)
	}
}

func TestInsertGrow(t *testing.T) {
	b := NewGappedTextBuffer(2)
	for i := 0; i < 5; i++ {
		b.Insert(i, byte('a'+i))
	}
	expected := "abcde"
	if got := flatten(b); got != expected {
		t.Errorf("Expected '%s', got '%s'", expected, got)
	}
}

func TestDeleteBasic(t *testing.T) {
	b := NewGappedTextBuffer(10)
	input := "Hello"
	for i := range input {
		b.Insert(i, input[i])
	}
	b.Delete(1) // remove 'e'
	if got := flatten(b); got != "Hllo" {
		t.Errorf("Expected 'Hllo', got '%s'", got)
	}
}

func TestDeleteLastChar(t *testing.T) {
	b := NewGappedTextBuffer(10)
	b.Insert(0, 'A')
	b.Insert(1, 'B')
	b.Delete(1)
	if got := flatten(b); got != "A" {
		t.Errorf("Expected 'A', got '%s'", got)
	}
}

func TestMoveCursorLeftRight(t *testing.T) {
	b := NewGappedTextBuffer(10)
	b.Insert(0, 'a')
	b.Insert(1, 'b')
	b.Insert(2, 'c')
	b.MoveCursorTo(1)
	b.Insert(1, 'X')
	if got := flatten(b); got != "aXbc" {
		t.Errorf("Expected 'aXbc', got '%s'", got)
	}
}

func TestInsertAtStart(t *testing.T) {
	b := NewGappedTextBuffer(10)
	b.Insert(0, 'b')
	b.Insert(0, 'a')
	if got := flatten(b); got != "ab" {
		t.Errorf("Expected 'ab', got '%s'", got)
	}
}

func TestDeleteAtStart(t *testing.T) {
	b := NewGappedTextBuffer(10)
	b.Insert(0, 'a')
	b.Insert(1, 'b')
	b.Delete(0)
	if got := flatten(b); got != "b" {
		t.Errorf("Expected 'b', got '%s'", got)
	}
}

func TestDeleteAtEnd(t *testing.T) {
	b := NewGappedTextBuffer(10)
	b.Insert(0, 'a')
	b.Insert(1, 'b')
	b.Delete(1)
	if got := flatten(b); got != "a" {
		t.Errorf("Expected 'a', got '%s'", got)
	}
}

func TestMoveCursorOutOfBounds(t *testing.T) {
	b := NewGappedTextBuffer(10)
	b.Insert(0, 'a')
	b.Insert(1, 'b')
	b.MoveCursorTo(999) // should not crash
	if got := flatten(b); got != "ab" {
		t.Errorf("Expected 'ab', got '%s'", got)
	}
}

func TestMultipleEdits(t *testing.T) {
	b := NewGappedTextBuffer(10)
	word := "editor"
	for i := range word {
		b.Insert(i, word[i])
	}

	// Insert 'x' after 'e' (at pos 1)
	b.Insert(1, 'x') // → exditor

	// Delete 'd' at pos 2 (the position of 'd' after inserting 'x')
	b.Delete(2) // → exitor

	// Insert 'A' at beginning
	b.Insert(0, 'A') // → Aexitor

	expected := "Aexitor"
	if got := flatten(b); got != expected {
		t.Errorf("Expected '%s', got '%s'", expected, got)
	}
}
