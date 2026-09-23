package db

import (
	"idea_app/models"
	"strings"
	"testing"
)

func TestMainDB(t *testing.T) {
	// 1. Setup in-memory DB
	err := Initialize(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize DB: %v", err)
	}
	defer Close()

	// 2. Smoke Test: Create Idea
	idea := &models.Idea{
		ID:          models.NewID(),
		Title:       "Test Idea",
		Description: "Smoke Test",
	}
	if err := CreateIdea(idea); err != nil {
		t.Errorf("Failed to create idea: %v", err)
	}

	// 3. SQL Injection Test
	// Name contains malicious SQL. If prepared statements work, it will literally save this string.
	maliciousTitle := "Hack'; DROP TABLE ideas;--"
	maliciousIdea := &models.Idea{
		ID:    models.NewID(),
		Title: maliciousTitle,
	}
	if err := CreateIdea(maliciousIdea); err != nil {
		t.Errorf("SQL Injection test failed incorrectly: %v", err)
	}
	
	// Verify table still exists and data is safe
	ideas, err := GetAllIdeas()
	if err != nil || len(ideas) != 2 {
		t.Errorf("Ideas table was dropped or compromised by SQL Injection")
	}

	// 4. Nested MetaData Logic & Enforcement (Corner Cases)
	n1 := &models.Note{ID: models.NewID(), ParentID: idea.ID, ParentType: "Idea", Level: 1, Title: "L1"}
	n2 := &models.Note{ID: models.NewID(), ParentID: n1.ID, ParentType: "Note", Level: 2, Title: "L2"}
	n3 := &models.Note{ID: models.NewID(), ParentID: n2.ID, ParentType: "Note", Level: 3, Title: "L3"}
	n4 := &models.Note{ID: models.NewID(), ParentID: n3.ID, ParentType: "Note", Level: 4, Title: "L4"}
	n5 := &models.Note{ID: models.NewID(), ParentID: n4.ID, ParentType: "Note", Level: 5, Title: "L5"}
	n6 := &models.Note{ID: models.NewID(), ParentID: n5.ID, ParentType: "Note", Level: 6, Title: "L6"}
	
	for _, n := range []*models.Note{n1, n2, n3, n4, n5, n6} {
		if err := CreateNote(n); err != nil {
			t.Errorf("Failed to insert valid Note level %d: %v", n.Level, err)
		}
	}

	// Negative Test: Exceeding max levels
	n7 := &models.Note{ID: models.NewID(), ParentID: n6.ID, ParentType: "Note", Level: 7, Title: "L7"}
	if err := CreateNote(n7); err == nil {
		t.Errorf("Allowed insertion of Level 7 Note!")
	} else if !strings.Contains(err.Error(), "maximum note nesting level") {
		t.Errorf("Unexpected error when testing Level 7 boundary: %v", err)
	}

	// Fetch MetaData hierarchy (Smoke test recursion)
	loadedMetaData, err := GetMetaDataForParent(idea.ID, "Idea")
	if err != nil || len(loadedMetaData) == 0 {
		t.Errorf("Failed to retrieve nested metaData: %v", err)
	}
	
	// 5. Negative Test: NPEs and Nil data handling
	err = DeleteIdea("") // Deleting non-existent/empty ID
	if err != nil {
		t.Errorf("Empty deletion should safely perform NOOP, got error: %v", err)
	}
	
	// 6. Cascade Delete Test
	AddAttachment(n6.ID, "/test/path/img.png", true, models.NewID())
	
	// Trigger cascade
	if err := DeleteIdea(idea.ID); err != nil {
		t.Errorf("Failed to delete idea: %v", err)
	}

	// Verify cascade happened
	n6Deleted, err := GetNoteByID(n6.ID)
	if err != nil {
		t.Errorf("DB error querying deleted note: %v", err)
	}
	if n6Deleted != nil {
		t.Errorf("Cascade delete failed, subnote L6 still exists")
	}

	// Verify attachments were deleted natively by Foreign Keys
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM attachments").Scan(&count)
	if count != 0 {
		t.Errorf("Cascade delete failed to drop attachments, count is %d", count)
	}
}
