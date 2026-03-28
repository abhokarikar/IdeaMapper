package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func BuildMainTabs(a fyne.App, w fyne.Window) fyne.CanvasObject {
	ideasScreen := BuildIdeasList(w)
	memoriesScreen := BuildMemoriesList(w)
	settingsScreen := BuildSettings(w)

	tabs := container.NewAppTabs(
		container.NewTabItem("Ideas", ideasScreen),
		container.NewTabItem("Memories", memoriesScreen),
		container.NewTabItem("Settings", settingsScreen),
	)
	tabs.SetTabLocation(container.TabLocationBottom)

	return tabs
}

func BuildIdeasList(w fyne.Window) fyne.CanvasObject {
	return container.NewCenter(widget.NewLabel("Ideas View Setup (In Progress)"))
}

func BuildMemoriesList(w fyne.Window) fyne.CanvasObject {
	return container.NewCenter(widget.NewLabel("Memories View Setup (In Progress)"))
}

func BuildSettings(w fyne.Window) fyne.CanvasObject {
	return container.NewCenter(widget.NewLabel("Settings (Google Drive Backup)"))
}
