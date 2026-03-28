package models

import "time"

type Idea struct {
	ID          string
	Title       string
	Description string
	IconPath    string
	IsImportant bool
	Notes       []*Note
}

type Memory struct {
	ID          string
	Title       string
	Description string
	Date        time.Time
	IsImportant bool
	Notes       []*Note
}

type InterLink struct {
	ID         string
	SourceID   string
	SourceType string // "Idea" or "Memory"
	TargetID   string
	TargetType string // "Idea" or "Memory"
	Name       string
	LinkNote   *Note  // Embedded Note specific to this link
}

type Note struct {
	ID          string
	ParentID    string
	ParentType  string // "Idea", "Memory", "InterLink", "Note"
	Level       int
	Title       string
	Description string
	FilePaths   []string
	PhotoPaths  []string
	SubNotes    []*Note
}
