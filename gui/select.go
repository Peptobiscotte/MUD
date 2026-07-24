package main

import (
	"net"

	"fyne.io/fyne/v2/widget"
)

func SelectMonster(ui *UI, buttonAttack *widget.Button, conn net.Conn) {
	ui.MonstersList.OnSelected = func(id widget.ListItemID) {
		ui.SelectedMonster = ui.Monsters[id]
		UpdateAttackButton(buttonAttack, conn, ui)
	}
}

func SelectItem(ui *UI) {
	ui.ItemsList.OnSelected = func(id widget.ListItemID) {
		ui.SelectedItem = ui.Items[id]
	}
}

func SelectNpc(ui *UI) {
	ui.NpcsList.OnSelected = func(id widget.ListItemID) {
		ui.SelectedNpc = ui.Npcs[id]
	}
}

func SelectQuestgiver(ui *UI) {
	ui.QuestgiversList.OnSelected = func(id widget.ListItemID) {
		ui.SelectedQuestgiver = ui.Questgivers[id]
	}
}

func SelectQuest(ui *UI) {
	ui.QuestsList.OnSelected = func(id widget.ListItemID) {
		ui.SelectedQuest = ui.Quests[id]
	}
}

func SelectInventory(ui *UI) {
	ui.InventoryList.OnSelected = func(id widget.ListItemID) {
		ui.SelectedInventory = ui.Inventory[id].Name
	}
}

func SelectPlayer(ui *UI) {
	ui.PlayersList.OnSelected = func(id widget.ListItemID) {
		ui.SelectedPlayer = ui.Players[id]
	}
}

func SelectInvite(ui *UI) {
	ui.InviteList.OnSelected = func(id widget.ListItemID) {
		ui.SelectedInvite = ui.Invites[id]
	}
}
