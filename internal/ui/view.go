package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Cursor styles
var normalCursorStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("7")).
	Foreground(lipgloss.Color("0"))

var insertCursorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("2"))

func (m Model) View() string {
	var b strings.Builder

	height := m.height
	if height == 0 {
		height = 24
	}

	statusHeight := 2
	usableHeight := height - statusHeight

	var lines []string

	startRow := m.cursor.row - usableHeight/2
	if startRow < 0 {
		startRow = 0
	}
	endRow := startRow + usableHeight
	if endRow > len(m.tb.Lines) {
		endRow = len(m.tb.Lines)
		startRow = endRow - usableHeight
		if startRow < 0 {
			startRow = 0
		}
	}

	for row := startRow; row < endRow; row++ {
		line := m.tb.Lines[row]
		text := line.Flatten()
		runes := []rune(text)

		if row == m.cursor.row {
			col := m.cursor.col
			if col > len(runes) {
				col = len(runes)
			}
			if col >= len(runes) {
				runes = append(runes, ' ')
			}

			before := string(runes[:col])
			cursorChar := string(runes[col])
			after := string(runes[col+1:])

			if m.mode == Insert {
				text = before + insertCursorStyle.Render("|") + cursorChar + after
			} else {
				text = before + normalCursorStyle.Render(cursorChar) + after
			}
		}

		lines = append(lines, text)
	}

	// Pad lines if needed
	for len(lines) < usableHeight {
		lines = append(lines, "")
	}

	// Ensure we exactly have `usableHeight` lines before adding status line
	lines = lines[:usableHeight]

	// Add status line (the final line)
	statusLine := "-- NORMAL --"
	if m.mode == Insert {
		statusLine = "-- INSERT --"
	} else if m.mode == Command {
		statusLine = m.cmdBuffer
	}
	lines = append(lines, statusLine) // Now total lines == height

	// Final render
	for _, l := range lines {
		b.WriteString(l + "\n")
	}

	return b.String()
}
