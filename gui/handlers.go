package main

import (
	"encoding/json"
	"fmt"
	"net"

	// "slices"
	"strconv"
	"strings"
	"tap/game"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	// "golang.org/x/text/message"
)

func handleTakeResponse(message string, ui *UI, conn net.Conn) {
	itemName := strings.TrimPrefix(message, "OK taken=")
	itemName, _, _ = strings.Cut(itemName, "@")

	fyne.Do(func() {
		for i, item := range ui.Items {
			if item == itemName {
				ui.Items = append(
					ui.Items[:i],
					ui.Items[i+1:]...,
				)
				break
			}
		}
		fmt.Fprintln(conn, "INVENTORY")
		ui.SelectedItem = ""
		ui.Response.Text = message
		UnselectAll(ui)
	})
}

func handleConnect(message string, ui *UI) {
	name := strings.TrimPrefix(message, "OK connected ")
	fyne.Do(func() {
		ui.PlayerName = name
	})
}

func handleUsedResponse(message string, ui *UI, conn net.Conn) {
	fyne.Do(func() {
		fmt.Fprintln(conn, "INVENTORY")
		fmt.Fprintln(conn, "STATUS")
		ui.Response.Text = message
		UnselectAll(ui)
	})
}

func handleDropResponse(message string, ui *UI, conn net.Conn) {
	fyne.Do(func() {
		fmt.Fprintln(conn, "INVENTORY")
		fmt.Fprintln(conn, "LOOK")
		ui.Response.Text = message
		UnselectAll(ui)
	})
}

func handleStatusResponse(message string, ui *UI) {
	jsonData := strings.TrimPrefix(message, "OK status=")
	jsonData, _, _ = strings.Cut(jsonData, "@")

	var response game.StatusResponse

	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		fmt.Println("STATUS JSON error:", err)
		return
	}

	fyne.Do(func() {
		ui.HpLabel.Text = ("Hp = " + strconv.Itoa(response.Hp))
		ui.MaxHpLabel.SetText("MaxHP = " + strconv.Itoa(response.Max_hp))
		ui.StatusLabel.SetText("Status = " + response.Status)

		ui.HpLabel.Show()
		ui.MaxHpLabel.Show()
		ui.StatusLabel.Show()

		ui.HpLabel.Refresh()
		ui.MaxHpLabel.Refresh()
		ui.StatusLabel.Refresh()

		ui.Response.Text = message
		UnselectAll(ui)
	})
}

func handleAttackResponse(message string, ui *UI, conn net.Conn) {
	if strings.Contains(message, "target_hp") {
		jsonData := strings.TrimPrefix(message, "OK combat=")
		jsonData, _, _ = strings.Cut(jsonData, "@")
		var response game.CombatResponse

		err := json.Unmarshal([]byte(jsonData), &response)
		if err != nil {
			fmt.Println("STATUS JSON error:", err)
			return
		}
		fyne.Do(func() {
			if strings.Contains(message, "defeated") {
				fmt.Fprintln(conn, "LOOK")
				ui.SelectedMonster = ""
				ui.MonstersList.UnselectAll()
			}

			ui.MonsterHp.SetText("Monster Hp = " + strconv.Itoa(response.Target_hp))
			ui.CounterAtk.SetText("Monster dmg = " + strconv.Itoa(response.Counter))
			ui.MonsterHp.Show()
			ui.CounterAtk.Show()
			fmt.Fprintln(conn, "STATUS")
			ui.MonstersList.Refresh()
			ui.Response.Text = message
			UnselectAll(ui)
		})
	} else {
		fyne.Do(func() {
			if strings.Contains(message, "defeated") {
				fmt.Fprintln(conn, "LOOK")
				ui.SelectedMonster = ""
				ui.MonstersList.UnselectAll()
			}
			fmt.Fprintln(conn, "STATUS")
			ui.Response.Text = message
			UnselectAll(ui)
		})
	}
}

func handleCreateResponse(message string, ui *UI) {
	Inv := strings.TrimPrefix(message, "OK group=")
	fyne.Do(func() {
		for i, inv := range ui.Invites {
			if inv == Inv {
				ui.Invites = append(ui.Invites[:i], ui.Invites[i+1:]...)
				break
			}
		}
		ui.Group = append(ui.Group, ui.PlayerName)
		ui.Response.Text = message
		UnselectAll(ui)
	})
}

func handleInvResponse(message string, ui *UI) {
	jsonData := strings.TrimPrefix(message, "OK inventory=")
	jsonData, _, _ = strings.Cut(jsonData, "@")

	var response map[string]game.Item

	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		fmt.Println("INVENTORY JSON error:", err)
		return
	}

	inventory := make([]InventoryEntry, 0, len(response))
	for name, item := range response {
		inventory = append(inventory, InventoryEntry{
			Name: name,
			Item: item,
		})
	}

	fyne.Do(func() {
		ui.Inventory = inventory
		ui.SelectedInventory = ""
		ui.Response.Text = message
		UnselectAll(ui)
	})
}

func handleLookResponse(
	message string,
	ui *UI,
) {
	jsonData := strings.TrimPrefix(message, "OK ")
	jsonData, _, _ = strings.Cut(jsonData, "@")

	var response game.LookResponse

	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		fmt.Println("LOOK JSON error:", err)
		return
	}

	fyne.Do(func() {
		ui.RoomTitle.Text = response.Room.Name

		ui.RoomDescription.Text = response.Room.Description

		ui.Players = response.Players

		ui.Items = response.Items

		ui.Monsters = response.Monsters

		ui.Exits = response.Room.Exits
		ui.MoveButton.Enable()
		ui.ConnectButton.Enable()

		ui.Questgivers = response.Questgivers

		ui.Npcs = response.NPCs

		ui.Dialogue = false
		ui.AnswerEntry.Hide()

		ui.MonsterHp.Hide()
		ui.CounterAtk.Hide()
		UnselectAll(ui)
	})
}

// func handleQuestsResponse(
// 	message string,
// 	ui *UI,
// ) {
// 	text := strings.TrimPrefix(message, "OK quests=")

// 	lines := strings.Split(text, "\n")
// 	quests := make([]string, len(lines))
// 	for i, line := range lines {
// 		quests[i] = strings.SplitN(line, ": ", 2)[0]
// 	}

// 	fyne.Do(func() {
// 		ui.Quests = quests
// 		UnselectAll(ui)
// 	})
// }

func handleGroupJResponse(
	message string,
	ui *UI,
) {
	text := strings.TrimPrefix(message, "EVT GROUP JOIN")

	fyne.Do(func() {
		ui.Group = append(ui.Group, text)
		ui.GroupList.Refresh()
	})
}

func handleGroupLResponse(
	message string,
	ui *UI,
) {
	text := strings.TrimPrefix(message, "EVT GROUP LEAVE")

	fyne.Do(func() {
		for i, member := range ui.Group {
			if member == text {
				ui.Group = append(ui.Group[:i], ui.Group[i+1:]...)
				break
			}
		}
		ui.GroupList.Refresh()
	})
}

func handleRoomEnterResponse(
	message string,
	ui *UI,
) {
	text := strings.TrimPrefix(message, "EVT ROOM PRESENCE ENTER ")

	fyne.Do(func() {
		ui.Players = append(ui.Players, text)
		ui.PlayersList.Refresh()
	})
}

func handleRoomLeaveResponse(
	message string,
	ui *UI,
) {
	text := strings.TrimPrefix(message, "EVT ROOM PRESENCE LEAVE ")

	fyne.Do(func() {
		for i, member := range ui.Players {
			if member == text {
				ui.Players = append(ui.Players[:i], ui.Players[i+1:]...)
				break
			}
		}
		ui.PlayersList.Refresh()
	})
}

func handleSelfLeaveResponse(
	ui *UI,
) {
	fyne.Do(func() {
		ui.Group = nil
		ui.GroupList.Refresh()
	})
}

func handleSelfInvite(
	message string,
	ui *UI,
) {
	text := strings.TrimPrefix(message, "EVT GROUP INVITE ")
	fyne.Do(func() {
		ui.Invites = append(ui.Invites, text)
		ui.InviteList.Refresh()
	})
}

func handleMembersResponse(
	message string,
	ui *UI,
) {
	text := strings.TrimPrefix(message, "OK members= ")
	parts := strings.Fields(text)
	fyne.Do(func() {
		ui.Group = nil
		ui.Group = append(ui.Group, parts...)
		ui.GroupList.Refresh()
	})
}

func handleChatGlobal(
	message string,
	ui *UI,
) {
	text := strings.TrimPrefix(message, "EVT ")

	fyne.Do(func() {
		ui.Chat = append(ui.Chat, text)
		label := widget.NewLabel(text)
		label.Wrapping = fyne.TextWrapWord
		ui.ChatBox.Add(label)
		ui.ChatContainer.Refresh()
		ui.ScrollChat.ScrollToBottom()
	})
}

func handleTalkResponse(
	ui *UI,
) {
	fyne.Do(func() {
		ui.Dialogue = false
		ui.AnswerEntry.Hide()
		ui.SelectedNpc = ""
		UnselectAll(ui)
	})
}

func handleQuestResponse(
	message string,
	ui *UI,
) {
	quest := strings.TrimPrefix(message, "OK quest=")
	quest, _, _ = strings.Cut(quest, "@")
	quest, _ = strings.CutSuffix(quest, " ")

	fyne.Do(func() {
		ui.SelectedQuestgiver = ""
		ui.Quests = append(ui.Quests, quest)
		UnselectAll(ui)
	})
}

func handleCompletequestResponse(
	message string,
	ui *UI,
	conn net.Conn,
) {
	text := strings.TrimPrefix(message, "OK complete=")
	quest := strings.SplitN(text, "\n", 2)[0]

	fyne.Do(func() {
		for i, q := range ui.Quests {
			if q == quest {
				ui.Quests = append(ui.Quests[:i], ui.Quests[i+1:]...)
				break
			}
		}
		fmt.Fprintln(conn, "INVENTORY")
		UnselectAll(ui)
	})
}

func handleAbandonquestResponse(
	message string,
	ui *UI,
) {
	text := strings.TrimPrefix(message, "OK abandon=")

	fyne.Do(func() {
		for i, q := range ui.Quests {
			if q == text {
				ui.Quests = append(ui.Quests[:i], ui.Quests[i+1:]...)
				break
			}
		}
		UnselectAll(ui)
	})
}

func handleDeath(
	conn net.Conn,
	ui *UI,
) {
	fmt.Fprintln(conn, "LOOK")
	fmt.Fprintln(conn, "STATUS")
	fyne.Do(func() {
		ui.HpLabel.Refresh()
		ui.MaxHpLabel.Refresh()
		ui.StatusLabel.Refresh()
	})
}

func handleFireball(
	conn net.Conn,
	ui *UI,
	message string,
) {
	if strings.Contains(message, "target_hp") {
		jsonData := strings.TrimPrefix(message, "OK used=firescroll ")
		jsonData, _, _ = strings.Cut(jsonData, "@")
		var response game.CombatResponse

		err := json.Unmarshal([]byte(jsonData), &response)
		if err != nil {
			fmt.Println("STATUS JSON error:", err)
			return
		}
		fyne.Do(func() {
			if strings.Contains(message, "defeated") {
				fmt.Fprintln(conn, "LOOK")
				ui.SelectedMonster = ""
				ui.MonstersList.UnselectAll()
			}

			ui.MonsterHp.SetText("Monster Hp = " + strconv.Itoa(response.Target_hp))
			ui.CounterAtk.SetText("Monster dmg = " + strconv.Itoa(response.Counter))
			ui.MonsterHp.Show()
			ui.CounterAtk.Show()
			fmt.Fprintln(conn, "STATUS")
			ui.InventoryList.Refresh()
			ui.MonstersList.Refresh()
			ui.Response.Text = message
			UnselectAll(ui)
		})
	}
	fmt.Fprintln(conn, "STATUS")
	fyne.Do(func() {
		ui.HpLabel.Refresh()
		ui.MaxHpLabel.Refresh()
		ui.StatusLabel.Refresh()
		ui.MonsterHp.Refresh()
		ui.CounterAtk.Refresh()
	})
}

func handleQuit(MyApp fyne.App) {
	MyApp.Quit()
}

func handleRefresh(conn net.Conn, ui *UI) {
	fmt.Fprintln(conn, "LOOK")
	UnselectAll(ui)
}
