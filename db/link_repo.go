package db

import (
	"idea_app/models"
)

func CreateInterLink(link *models.InterLink) error {
	_, err := DB.Exec("INSERT INTO interlinks (id, source_id, source_type, target_id, target_type, name) VALUES (?, ?, ?, ?, ?, ?)",
		link.ID, link.SourceID, link.SourceType, link.TargetID, link.TargetType, link.Name)
	return err
}

func GetInterLinksForSource(sourceID string) ([]*models.InterLink, error) {
	rows, err := DB.Query("SELECT id, source_id, source_type, target_id, target_type, name FROM interlinks WHERE source_id = ? OR target_id = ?", sourceID, sourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*models.InterLink
	for rows.Next() {
		l := &models.InterLink{}
		if err := rows.Scan(&l.ID, &l.SourceID, &l.SourceType, &l.TargetID, &l.TargetType, &l.Name); err != nil {
			return nil, err
		}
		// A Link only supports 1 singular top-level Note attached directly to it.
		notes, _ := GetNotesForParent(l.ID, "InterLink")
		if len(notes) > 0 {
			l.LinkNote = notes[0]
		}
		links = append(links, l)
	}
	return links, nil
}

func DeleteInterLink(id string) error {
	_, err := DB.Exec("DELETE FROM interlinks WHERE id = ?", id)
	return err
}
