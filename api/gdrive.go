package api

import (
	"encoding/json"
	"fmt"
	"idea_app/db"
	"idea_app/models"
)

type BackupPayload struct {
	Ideas    []*models.Idea   `json:"ideas"`
	Memories []*models.Memory `json:"memories"`
}

// GenerateBackupBlob extracts all IsImportant=true entities and translates them to JSON
func GenerateBackupBlob() ([]byte, error) {
	ideas, err := db.GetAllIdeas()
	if err != nil {
		return nil, err
	}
	
	memories, err := db.GetAllMemories()
	if err != nil {
		return nil, err
	}
	
	var impIdeas []*models.Idea
	for _, i := range ideas {
		if i.IsImportant {
			// Fetch its notes too
			notes, _ := db.GetNotesForParent(i.ID, "Idea")
			i.Notes = notes
			impIdeas = append(impIdeas, i)
		}
	}
	
	var impMemories []*models.Memory
	for _, m := range memories {
		if m.IsImportant {
			notes, _ := db.GetNotesForParent(m.ID, "Memory")
			m.Notes = notes
			impMemories = append(impMemories, m)
		}
	}
	
	payload := BackupPayload{Ideas: impIdeas, Memories: impMemories}
	
	return json.MarshalIndent(payload, "", "  ")
}

// UploadToDrive is a placeholder for the actual OAuth2 drive.v3 API upload.
func UploadToDrive(jsonData []byte, authCode string) error {
	// 1. Exchange authCode for Token
	// 2. Initialize drive.NewService
	// 3. Create drive.File metadata
	// 4. Call driveService.Files.Create(file).Media(bytes.NewReader(jsonData)).Do()
	return fmt.Errorf("Google Drive API OAuth integration requires specific client configurations")
}
