package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func CreatePlayerList(ui *UI) *widget.List {
	playersList := widget.NewList(
		func() int {
			return len(ui.Players)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(ui.Players[id])
		},
	)
	return playersList
}

func CreateGroupList(ui *UI) *widget.List {
	groupList := widget.NewList(
		func() int {
			return len(ui.Group)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(ui.Group[id])
		},
	)
	return groupList
}

func CreateInventoryList(ui *UI) *widget.List {
	inventoryList := widget.NewList(
		func() int {
			return len(ui.Inventory)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(ui.Inventory[id].Name)
		},
	)
	return inventoryList
}

func CreateMonstersList(ui *UI) *widget.List {
	monstersList := widget.NewList(
		func() int {
			return len(ui.Monsters)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(ui.Monsters[id])
		},
	)
	return monstersList
}

func CreateItemsList(ui *UI) *widget.List {
	itemsList := widget.NewList(
		func() int {
			return len(ui.Items)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(ui.Items[id])
		},
	)
	return itemsList
}

func CreateQuestsList(ui *UI) *widget.List {
	questsList := widget.NewList(
		func() int {
			return len(ui.Quests)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(ui.Quests[id])
		},
	)
	return questsList
}

func CreateNpcsList(ui *UI) *widget.List {
	npcsList := widget.NewList(
		func() int {
			return len(ui.Npcs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(ui.Npcs[id])
		},
	)
	return npcsList
}

func CreateInviteList(ui *UI) *widget.List {
	invite := widget.NewList(
		func() int {
			return len(ui.Invites)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(ui.Invites[id])
		},
	)
	return invite
}

func CreateQuestgiversList(ui *UI) *widget.List {
	questgiversList := widget.NewList(
		func() int {
			return len(ui.Questgivers)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(ui.Questgivers[id])
		},
	)
	return questgiversList
}
