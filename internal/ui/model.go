package ui

import (
	"text-editor/internal/buffer"

	tea "github.com/charmbracelet/bubbletea"
)

type Mode int

const (
	Normal Mode = iota
	Insert
	Command
)

type Model struct {
	tb           *buffer.TextBuffer
	cursor       struct{ row, col int }
	mode         Mode
	cmdBuffer    string
	filePath     string
	width        int
	height       int
	scrollOffset int
}

func New(tb *buffer.TextBuffer, filePath string) Model {
	return Model{
		tb:       tb,
		filePath: filePath,
		mode:     Normal,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
