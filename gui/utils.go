package main

import (
	"fmt"
	"net"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func UpdateAttackButton(button *widget.Button, conn net.Conn, ui *UI) {
	if strings.Contains(ui.SelectedMonster, "dead") {
		button.SetText("LOOT")
		button.OnTapped = func() {
			fmt.Fprintln(conn, "LOOT "+ui.SelectedMonster)
		}
	} else {
		button.SetText("ATTACK")
		button.OnTapped = func() {
			fmt.Fprintln(conn, "ATTACK "+ui.SelectedMonster)
		}
	}
}

func UnselectAll(ui *UI) {
	fyne.Do(func() {
		ui.SelectedInventory = ""
		ui.SelectedInvite = ""
		ui.SelectedItem = ""
		ui.SelectedMonster = ""
		ui.SelectedNpc = ""
		ui.SelectedPlayer = ""
		ui.SelectedQuest = ""
		ui.SelectedQuestgiver = ""
		ui.InventoryList.UnselectAll()
		ui.InviteList.UnselectAll()
		ui.ItemsList.UnselectAll()
		ui.MonstersList.UnselectAll()
		ui.NpcsList.UnselectAll()
		ui.PlayersList.UnselectAll()
		ui.QuestsList.UnselectAll()
		ui.QuestgiversList.UnselectAll()
		ui.InventoryList.Refresh()
		ui.InviteList.Refresh()
		ui.ItemsList.Refresh()
		ui.MonstersList.Refresh()
		ui.NpcsList.Refresh()
		ui.PlayersList.Refresh()
		ui.QuestsList.Refresh()
		ui.QuestgiversList.Refresh()
	})
}
