package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"idea_app/db"
	"idea_app/models"
)

func BuildMemoriesList(w fyne.Window) fyne.CanvasObject {
	vbox := container.NewVBox()
	scroll := container.NewVScroll(vbox)

	refreshList := func() {
		vbox.Objects = nil
		memories, err := db.GetAllMemories()
		if err != nil {
			vbox.Add(widget.NewLabel("Failed to load memories: " + err.Error()))
			return
		}
		
		var currentYear int
		for _, m := range memories {
			year := m.Date.Year()
			if year != currentYear {
				currentYear = year
				yearLabel := widget.NewLabelWithStyle(fmt.Sprintf("%d", year), fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
				vbox.Add(yearLabel)
			}
			
			dateStr := m.Date.Format("Jan 02")
			titleLabel := widget.NewLabel(fmt.Sprintf("%s - %s", dateStr, m.Title))
			
			vbox.Add(container.NewHBox(titleLabel, widget.NewButton("View", func() {
				dialog.ShowInformation("Memory Details", "Opening Memory: "+m.Title, w)
			})))
		}
		scroll.Refresh()
	}

	addBtn := widget.NewButton("Create Memory", func() {
		ShowCreateMemoryModal(w, refreshList)
	})

	topBar := container.NewHBox(
		widget.NewLabelWithStyle("My Memories", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		addBtn,
	)

	refreshList()

	return container.NewBorder(topBar, nil, nil, nil, scroll)
}

func ShowCreateMemoryModal(w fyne.Window, onSave func()) {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Memory Title")
	
	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Description")
	
	form := widget.NewForm(
		widget.NewFormItem("Title", titleEntry),
		widget.NewFormItem("Description", descEntry),
	)

	d := dialog.NewCustomConfirm("New Memory", "Save", "Cancel", form, func(save bool) {
		if save && titleEntry.Text != "" {
			mem := &models.Memory{
				ID:          models.NewID(),
				Title:       titleEntry.Text,
				Description: descEntry.Text,
				Date:        time.Now(),
				IsImportant: true, // Default to true for memories maybe?
			}
			db.CreateMemory(mem)
			onSave()
		}
	}, w)
	
	d.Show()
}
