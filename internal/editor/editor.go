package editor

import (
	"io/ioutil"
	
)

// Editor represents the main editor functionality
type Editor struct {
	currentFile string
	content     string
}

// New creates a new editor instance
func New() *Editor {
	return &Editor{}
}

// LoadFile loads content from a file
func (e *Editor) LoadFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	e.currentFile = filename
	e.content = string(content)
	return nil
}

// SaveFile saves content to a file
func (e *Editor) SaveFile(filename string) error {
	err := ioutil.WriteFile(filename, []byte(e.content), 0644)
	if err != nil {
		return err
	}
	e.currentFile = filename
	return nil
}

// SetContent sets the editor content
func (e *Editor) SetContent(content string) {
	e.content = content
}

// GetContent returns the editor content
func (e *Editor) GetContent() string {
	return e.content
}
