package main

import (
	"image/color"
	"net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateContainerMonster(ui *UI) *fyne.Container {
	monsterText := canvas.NewText("Monsters", color.White)
	monsterText.TextSize = 18
	monsterText.TextStyle = fyne.TextStyle{Bold: true}

	monsterTitle := container.NewCenter(
		monsterText,
	)

	sizedList := container.NewGridWrap(
		fyne.NewSize(180, 220),
		ui.MonstersList,
	)

	monsterContainer := container.NewBorder(
		monsterTitle,
		nil,
		nil,
		nil,
		sizedList,
	)
	return monsterContainer
}

func CreateContainerGroup(ui *UI) *fyne.Container {
	GroupText := canvas.NewText("Group", color.White)
	GroupText.TextSize = 18
	GroupText.TextStyle = fyne.TextStyle{Bold: true}

	GroupTitle := container.NewCenter(
		GroupText,
	)

	sizedList := container.NewGridWrap(
		fyne.NewSize(180, 220),
		ui.GroupList,
	)

	GroupContainer := container.NewBorder(
		GroupTitle,
		nil,
		nil,
		nil,
		sizedList,
	)
	return GroupContainer
}

func CreateContainerInvites(ui *UI) *fyne.Container {
	InviteText := canvas.NewText("Group Invites", color.White)
	InviteText.TextSize = 18
	InviteText.TextStyle = fyne.TextStyle{Bold: true}

	InviteTitle := container.NewCenter(
		InviteText,
	)

	sizedList := container.NewGridWrap(
		fyne.NewSize(180, 220),
		ui.InviteList,
	)

	InviteContainer := container.NewBorder(
		InviteTitle,
		nil,
		nil,
		nil,
		sizedList,
	)
	return InviteContainer
}

func CreateContainerItem(ui *UI) *fyne.Container {
	itemText := canvas.NewText("Items", color.White)
	itemText.TextSize = 18
	itemText.TextStyle = fyne.TextStyle{Bold: true}

	itemTitle := container.NewCenter(
		itemText,
	)

	sizedList := container.NewGridWrap(
		fyne.NewSize(180, 220),
		ui.ItemsList,
	)

	itemContainer := container.NewBorder(
		itemTitle,
		nil,
		nil,
		nil,
		sizedList,
	)

	return itemContainer
}

func CreateContainerInventory(ui *UI) *fyne.Container {
	inventoryText := canvas.NewText("Inventory", color.White)
	inventoryText.TextSize = 18
	inventoryText.TextStyle = fyne.TextStyle{Bold: true}

	inventoryTitle := container.NewCenter(
		inventoryText,
	)

	sizedList := container.NewGridWrap(
		fyne.NewSize(180, 220),
		ui.InventoryList,
	)

	inventoryContainer := container.NewBorder(
		inventoryTitle,
		nil,
		nil,
		nil,
		sizedList,
	)
	return inventoryContainer
}

func CreateContainerPlayer(ui *UI) *fyne.Container {
	playerText := canvas.NewText("Players", color.White)
	playerText.TextSize = 18
	playerText.TextStyle = fyne.TextStyle{Bold: true}

	playerTitle := container.NewCenter(
		playerText,
	)

	sizedList := container.NewGridWrap(
		fyne.NewSize(180, 220),
		ui.PlayersList,
	)

	playerContainer := container.NewBorder(
		playerTitle,
		nil,
		nil,
		nil,
		sizedList,
	)

	return playerContainer
}

func CreateContainerQuest(ui *UI) *fyne.Container {
	questsText := canvas.NewText("Quests", color.White)
	questsText.TextSize = 18
	questsText.TextStyle = fyne.TextStyle{Bold: true}

	questsTitle := container.NewCenter(
		questsText,
	)

	sizedList := container.NewGridWrap(
		fyne.NewSize(180, 220),
		ui.QuestsList,
	)

	questsContainer := container.NewBorder(
		questsTitle,
		nil,
		nil,
		nil,
		sizedList,
	)

	return questsContainer
}

func CreateContainerNpc(ui *UI) *fyne.Container {
	npcsText := canvas.NewText("Npcs", color.White)
	npcsText.TextSize = 18
	npcsText.TextStyle = fyne.TextStyle{Bold: true}

	npcsTitle := container.NewCenter(
		npcsText,
	)

	sizedList := container.NewGridWrap(
		fyne.NewSize(180, 220),
		ui.NpcsList,
	)

	npcsContainer := container.NewBorder(
		npcsTitle,
		nil,
		nil,
		nil,
		sizedList,
	)
	return npcsContainer
}

func CreateContainerChat(ui *UI, conn net.Conn) {
	chatText := canvas.NewText("Chat", color.White)
	chatText.TextSize = 18
	chatText.TextStyle = fyne.TextStyle{Bold: true}

	ui.ChatEntry = widget.NewEntry()
	ui.ChatEntry.SetPlaceHolder("Start chatting")

	chatTitle := container.NewCenter(
		chatText,
	)
	ui.ChatBox = container.NewVBox()
	for i := range ui.Chat {
		label := widget.NewLabel(ui.Chat[i])
		ui.ChatBox.Add(label)
	}
	ui.ScrollChat = container.NewVScroll(ui.ChatBox)
	ui.ScrollChat.SetMinSize(fyne.NewSize(180, 220))

	globalButton := CreateButtonChatGlobal(ui, conn)
	roomButton := CreateButtonChatRoom(ui, conn)
	groupButton := CreateButtonChatGroup(ui, conn)
	buttonContainer := container.NewGridWithColumns(3, globalButton, roomButton, groupButton)
	bottomContainer := container.NewBorder(ui.ChatEntry, buttonContainer, nil, nil, nil)
	ui.ChatContainer = container.NewBorder(chatTitle, bottomContainer, nil, nil, ui.ScrollChat)
}

func CreateContainerQuestgiver(ui *UI) *fyne.Container {
	questgiversText := canvas.NewText("Questgivers", color.White)
	questgiversText.TextSize = 18
	questgiversText.TextStyle = fyne.TextStyle{Bold: true}

	questgiversTitle := container.NewCenter(
		questgiversText,
	)

	sizedList := container.NewGridWrap(
		fyne.NewSize(180, 220),
		ui.QuestgiversList,
	)

	questgiversContainer := container.NewBorder(
		questgiversTitle,
		nil,
		nil,
		nil,
		sizedList,
	)
	return questgiversContainer
}

func CreateContainerResponse(ui *UI, conn net.Conn) {
	ui.ScrollResponse = container.NewVScroll(ui.Response)
	ui.ScrollResponse.SetMinSize(fyne.NewSize(180, 80))
}
