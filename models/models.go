package models

import "time"

type Idea struct {
	ID          string
	Title       string
	Description string
	IconPath    string
	IsImportant bool
	MetaData       []*MetaData
}

type Memory struct {
	ID          string
	Title       string
	Description string
	Date        time.Time
	IsImportant bool
	MetaData       []*MetaData
}

type InterLink struct {
	ID         string
	SourceID   string
	SourceType string // "Idea" or "Memory"
	TargetID   string
	TargetType string // "Idea" or "Memory"
	Name       string
	LinkMetaData   *MetaData  // Embedded MetaData specific to this link
}

type MetaData struct {
	ID          string
	ParentID    string
	ParentType  string // "Idea", "Memory", "InterLink", "MetaData"
	Level       int
	Title       string
	Description string
	FilePaths   []string
	PhotoPaths  []string
	SubMetaData    []*MetaData
}
