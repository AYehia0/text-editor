package ui

import (
	"slices"
	"text-editor/internal/reader"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) clampCursor() {
	if m.cursor.row < 0 {
		m.cursor.row = 0
	}
	if m.cursor.row >= len(m.tb.Lines) {
		m.cursor.row = len(m.tb.Lines) - 1
	}
	lineLen := m.tb.Lines[m.cursor.row].Len()

	if m.mode == Normal {
		// In normal mode, don't allow cursor after the last character
		if lineLen > 0 && m.cursor.col >= lineLen {
			m.cursor.col = lineLen - 1
		}
	} else {
		if m.cursor.col > lineLen {
			m.cursor.col = lineLen
		}
	}

	if m.cursor.col < 0 {
		m.cursor.col = 0
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch m.mode {
		case Normal:
			switch key {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "i":
				m.mode = Insert

			case ":":
				m.mode = Command
				m.cmdBuffer = ":"

			case "h":
				m.cursor.col--
				m.clampCursor()

			case "l":
				m.cursor.col++
				m.clampCursor()

			case "j":
				m.cursor.row++
				m.clampCursor()

			case "k":
				m.cursor.row--
				m.clampCursor()

			case "x":
				line := m.tb.Lines[m.cursor.row]
				if line.Len() > 0 && m.cursor.col < line.Len() {
					line.Delete(m.cursor.col)
				} else if line.Len() == 0 && m.tb.NumLines() > 1 {
					// Delete current empty line
					m.tb.Lines = slices.Delete(m.tb.Lines, m.cursor.row, m.cursor.row+1)
					if m.cursor.row > 0 {
						m.cursor.row--
					}
					m.cursor.col = 0
				}
				m.clampCursor()
			}

		case Insert:
			switch key {
			case "esc":
				m.mode = Normal
				m.clampCursor()

			case "enter":
				m.tb.BreakLineAtCursor(m.cursor.row, m.cursor.col)
				m.cursor.row++
				m.cursor.col = 0
				m.clampCursor()

			case "backspace":
				if m.cursor.row == 0 && m.cursor.col == 0 {
					// Do nothing
					break
				}
				if m.cursor.col > 0 {
					m.cursor.col--
					m.tb.Lines[m.cursor.row].Delete(m.cursor.col)
				} else if m.cursor.row > 0 {
					// Join with previous line
					prevRow := m.cursor.row - 1
					prevLen := m.tb.Lines[prevRow].Len()
					m.tb.JoinLines(m.cursor.row)
					m.cursor.row = prevRow
					m.cursor.col = prevLen
				}
				m.clampCursor()

			default:
				if len(key) == 1 {
					line := m.tb.Lines[m.cursor.row]
					if m.cursor.col > line.Len() {
						m.cursor.col = line.Len()
					}
					line.Insert(m.cursor.col, key[0])
					m.cursor.col++
				}
				m.clampCursor()
			}

		case Command:
			switch key {
			case "enter":
				if m.cmdBuffer == ":w" {
					err := reader.WriteFile(m.filePath, m.tb)
					if err != nil {
						m.cmdBuffer = ":w [error]"
					} else {
						m.cmdBuffer = ":w [ok]"
					}
				}
				m.mode = Normal
			case "backspace":
				if len(m.cmdBuffer) > 1 { // keep ':' at start
					m.cmdBuffer = m.cmdBuffer[:len(m.cmdBuffer)-1]
				}
			case "esc":
				m.mode = Normal
			default:
				if len(key) == 1 {
					m.cmdBuffer += key
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}
