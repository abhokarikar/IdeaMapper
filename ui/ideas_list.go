package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"idea_app/db"
	"idea_app/models"
)

func BuildIdeasList(w fyne.Window) fyne.CanvasObject {
	vbox := container.NewVBox()
	scroll := container.NewVScroll(vbox)

	refreshList := func() {
		vbox.Objects = nil
		ideas, err := db.GetAllIdeas()
		if err != nil {
			vbox.Add(widget.NewLabel("Failed to load ideas: " + err.Error()))
			return
		}
		
		for _, idea := range ideas {
			card := buildIdeaCard(w, idea)
			vbox.Add(card)
		}
		scroll.Refresh()
	}

	addBtn := widget.NewButton("Create Idea", func() {
		ShowCreateIdeaModal(w, refreshList)
	})

	topBar := container.NewHBox(
		widget.NewLabelWithStyle("My Ideas", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		addBtn,
	)

	refreshList()

	return container.NewBorder(topBar, nil, nil, nil, scroll)
}

func buildIdeaCard(w fyne.Window, idea *models.Idea) fyne.CanvasObject {
	title := widget.NewLabelWithStyle(idea.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	desc := widget.NewLabel(idea.Description)
	
	card := widget.NewCard("", "", container.NewVBox(title, desc))
	
	btn := widget.NewButton("Open", func() {
		// In a real app, this would replace the window content with the IdeaDetail view
		dialog.ShowInformation("Idea Details", "Opening Idea: "+idea.Title, w)
	})
	
	return container.NewVBox(card, btn, widget.NewSeparator())
}

func ShowCreateIdeaModal(w fyne.Window, onSave func()) {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Idea Title")
	
	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Description")
	
	form := widget.NewForm(
		widget.NewFormItem("Title", titleEntry),
		widget.NewFormItem("Description", descEntry),
	)

	d := dialog.NewCustomConfirm("New Idea", "Save", "Cancel", form, func(save bool) {
		if save && titleEntry.Text != "" {
			idea := &models.Idea{
				ID:          models.NewID(),
				Title:       titleEntry.Text,
				Description: descEntry.Text,
				IsImportant: false,
			}
			db.CreateIdea(idea)
			onSave()
		}
	}, w)
	
	d.Show()
}
