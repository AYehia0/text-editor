package reader

import (
	"bufio"
	"os"
	"text-editor/internal/buffer"
)

func ReadFile(filePath string, cap int) (*buffer.TextBuffer, string, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	tb := buffer.NewTextBuffer(cap)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := buffer.NewGappedTextBuffer(128)
		text := scanner.Text()
		for i := range text {
			line.Insert(i, text[i])
		}
		tb.Lines = append(tb.Lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, "", err
	}

	// In case the file is empty, ensure there's at least one line
	if len(tb.Lines) == 0 {
		tb.Lines = []*buffer.GappedTextBuffer{buffer.NewGappedTextBuffer(128)}
	}

	return tb, filePath, nil
}

func WriteFile(filePath string, tb *buffer.TextBuffer) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range tb.Lines {
		content := line.Flatten()
		if _, err := file.WriteString(content + "\n"); err != nil {
			return err
		}
	}
	return nil
}
