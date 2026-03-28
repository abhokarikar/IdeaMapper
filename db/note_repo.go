package db

import (
	"database/sql"
	"fmt"
	"idea_app/models"
)

func CreateNote(n *models.Note) error {
	if n.Level > 6 {
		return fmt.Errorf("maximum note nesting level of 6 exceeded")
	}
	_, err := DB.Exec("INSERT INTO notes (id, parent_id, parent_type, level, title, description) VALUES (?, ?, ?, ?, ?, ?)",
		n.ID, n.ParentID, n.ParentType, n.Level, n.Title, n.Description)
	return err
}

func GetNotesForParent(parentID string, parentType string) ([]*models.Note, error) {
	rows, err := DB.Query("SELECT id, parent_id, parent_type, level, title, description FROM notes WHERE parent_id = ? AND parent_type = ? ORDER BY created_at ASC", parentID, parentType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		n := &models.Note{}
		if err := rows.Scan(&n.ID, &n.ParentID, &n.ParentType, &n.Level, &n.Title, &n.Description); err != nil {
			return nil, err
		}

		// Fetch Attachments
		attRows, err := DB.Query("SELECT file_path, is_photo FROM attachments WHERE note_id = ?", n.ID)
		if err == nil {
			for attRows.Next() {
				var path string
				var isPhoto bool
				if err := attRows.Scan(&path, &isPhoto); err == nil {
					if isPhoto {
						n.PhotoPaths = append(n.PhotoPaths, path)
					} else {
						n.FilePaths = append(n.FilePaths, path)
					}
				}
			}
			attRows.Close()
		}

		// Recursively fetch subnotes if we are below max level
		if n.Level < 6 {
			subNotes, err := GetNotesForParent(n.ID, "Note")
			if err != nil {
				return nil, err
			}
			n.SubNotes = subNotes
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func GetNoteByID(id string) (*models.Note, error) {
	row := DB.QueryRow("SELECT id, parent_id, parent_type, level, title, description FROM notes WHERE id = ?", id)
	n := &models.Note{}
	if err := row.Scan(&n.ID, &n.ParentID, &n.ParentType, &n.Level, &n.Title, &n.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return n, nil
}

func DeleteNote(id string) error {
	_, err := DB.Exec("DELETE FROM notes WHERE id = ?", id)
	return err // Cascade triggers will automatically delete sub-notes and attachments
}

func AddAttachment(noteID, path string, isPhoto bool, attachID string) error {
	_, err := DB.Exec("INSERT INTO attachments (id, note_id, file_path, is_photo) VALUES (?, ?, ?, ?)",
		attachID, noteID, path, isPhoto)
	return err
}
