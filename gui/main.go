package main

import (
	"bufio"
	"fmt"
	"image/color"
	"net"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func connectToServer(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func readServer(
	conn net.Conn,
	ui *UI,
	myApp fyne.App,
) {
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fyne.Do(func() {
				ui.Response.SetText("Disconnected from server")
			})
			return
		}

		message = strings.TrimSpace(message)
		switch {
		case strings.Contains(message, `died`):
			handleDeath(conn, ui)

		case strings.Contains(message, `used=firescroll`):
			handleFireball(conn, ui, message)

		case strings.HasPrefix(message, `OK {"room":`):
			handleLookResponse(message, ui)

		case strings.HasPrefix(message, `OK taken=`):
			handleTakeResponse(message, ui, conn)

		case strings.HasPrefix(message, `OK combat=`):
			handleAttackResponse(message, ui, conn)

		case strings.HasPrefix(message, `OK inventory=`):
			handleInvResponse(message, ui)

		case strings.HasPrefix(message, `OK used=`):
			handleUsedResponse(message, ui, conn)

		// case strings.HasPrefix(message, `OK quests=`):
		// 	handleQuestsResponse(message, ui)

		case strings.HasPrefix(message, `OK dropped=`):
			handleDropResponse(message, ui, conn)

		case strings.HasPrefix(message, `OK status=`):
			handleStatusResponse(message, ui)

		case strings.HasPrefix(message, `OK group=`):
			handleCreateResponse(message, ui)

		case strings.HasPrefix(message, `OK connected`):
			handleConnect(message, ui)

		case strings.HasPrefix(message, `EVT GROUP JOIN`):
			handleGroupJResponse(message, ui)

		case strings.HasPrefix(message, `OK members=`):
			handleMembersResponse(message, ui)

		case strings.HasPrefix(message, `EVT GROUP LEAVE`):
			handleGroupLResponse(message, ui)

		case strings.HasPrefix(message, `EVT GROUP INVITE`):
			handleSelfInvite(message, ui)

		case strings.HasPrefix(message, `OK group left`):
			handleSelfLeaveResponse(ui)

		case strings.HasPrefix(message, `EVT ROOM PRESENCE ENTER`):
			handleRoomEnterResponse(message, ui)

		case strings.HasPrefix(message, `EVT ROOM PRESENCE LEAVE`):
			handleRoomLeaveResponse(message, ui)

		case strings.HasPrefix(message, `EVT GLOBAL CHAT`):
			handleChatGlobal(message, ui)

		case strings.HasPrefix(message, `EVT ROOM CHAT`):
			handleChatGlobal(message, ui)

		case strings.HasPrefix(message, `EVT GROUP CHAT`):
			handleChatGlobal(message, ui)

		case strings.HasPrefix(message, `EVT MONSTER DEAD`):
			handleRefresh(conn, ui)

		case strings.HasPrefix(message, `EVT PLAYER DEAD`):
			handleRefresh(conn, ui)

		case strings.HasPrefix(message, `EVT ITEM TAKEN`):
			handleRefresh(conn, ui)

		case strings.HasPrefix(message, `EVT ITEM DROPPED`):
			handleRefresh(conn, ui)

		case strings.HasPrefix(message, `OK end=`):
			handleTalkResponse(ui)

		case strings.HasPrefix(message, `OK quest=`):
			handleQuestResponse(message, ui)

		case strings.HasPrefix(message, `OK complete=`):
			handleCompletequestResponse(message, ui, conn)

		case strings.HasPrefix(message, `OK abandon=`):
			handleAbandonquestResponse(message, ui)

		case strings.HasPrefix(message, `OK bye`):
			handleQuit(myApp)

		}

		isTalkEnd := strings.HasPrefix(message, "OK end=")
		isErr := strings.HasPrefix(message, `ERR`)

		if isErr || ui.Dialogue {
			fyne.Do(func() {
				ui.Response.SetText(message)
				if ui.Dialogue && !isTalkEnd && !isErr {
					ui.AnswerEntry.Show()
				}
			})
		} else {
			fyne.Do(func() {
				ui.Response.SetText(message)
			})
		}
	}
}

func main() {
	//
	// CONNECT
	//
	conn, err := connectToServer("localhost:4242")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	myApp := app.New()
	myWindow := myApp.NewWindow("TAP")
	buttonAttack := widget.NewButton("ATTACK", nil)

	ui := &UI{
		RoomTitle:       canvas.NewText("", color.White),
		RoomDescription: canvas.NewText("", color.White),
		Response:        widget.NewLabel("Connected to server"),

		Players: []string{},

		Monsters: []string{},

		Items: []string{},

		Group: []string{},

		Invites: []string{},

		InventoryVisible: false,

		StatusVisible: false,
		HpLabel:       canvas.NewText("", color.RGBA{0, 255, 100, 255}),
		MaxHpLabel:    widget.NewLabel(""),
		StatusLabel:   widget.NewLabel(""),
		MonsterHp:     widget.NewLabel(""),
		CounterAtk:    widget.NewLabel(""),
	}
	// STYLE
	ui.RoomTitle.TextSize = 28
	ui.RoomTitle.TextStyle.Bold = true
	ui.RoomTitle.Color = color.RGBA{255, 131, 182, 255}
	ui.RoomDescription.TextSize = 24
	ui.RoomDescription.TextStyle.Italic = true
	ui.Response.Wrapping = fyne.TextWrapWord
	// CREATE LISTS
	ui.PlayersList = CreatePlayerList(ui)
	ui.GroupList = CreateGroupList(ui)
	ui.InventoryList = CreateInventoryList(ui)
	ui.MonstersList = CreateMonstersList(ui)
	ui.ItemsList = CreateItemsList(ui)
	ui.QuestsList = CreateQuestsList(ui)
	ui.NpcsList = CreateNpcsList(ui)
	ui.InviteList = CreateInviteList(ui)
	ui.QuestgiversList = CreateQuestgiversList(ui)
	// CREATE SELECTIONS
	SelectMonster(ui, buttonAttack, conn)
	SelectItem(ui)
	SelectNpc(ui)
	SelectQuest(ui)
	SelectInventory(ui)
	SelectPlayer(ui)
	SelectInvite(ui)
	SelectQuestgiver(ui)

	go readServer(
		conn,
		ui,
		myApp,
	)
	// CREATE BUTTONS
	fourButtons := container.NewVBox()

	buttonLook := CreateButtonLook(ui, conn)
	buttonInventory := CreateButtonInventory(ui, conn)
	buttonUse := CreateButtonUse(ui, conn)
	buttonDrop := CreateButtonDrop(ui, conn)
	buttonWho := CreateButtonWho(ui, conn)
	buttonTake := CreateButtonTake(ui, conn)
	// buttonCheck := CreateButtonCheck(ui, conn)
	buttonAbandonQuests := CreateButtonAbandonQuest(ui, conn)
	buttonQuest := CreateButtonQuest(ui, conn)
	buttonQuests := CreateButtonQuests(ui, conn)
	buttonCompleteQuests := CreateButtonCompleteQuest(ui, conn)
	buttonInvite := CreateButtonGroupInvite(ui, conn)
	buttonStatus := CreateButtonStatus(ui, conn)
	buttonCreate := CreateButtonCreate(ui, conn)
	buttonLeave := CreateButtonGroupLeave(ui, conn)
	buttonJoin := CreateButtonGroupJoin(ui, conn)
	buttonQuit := CreateButtonQuit(ui, conn)

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter player name")

	var buttonConnect *widget.Button

	buttonConnect = widget.NewButton("CONNECT", func() {
		name := entry.Text
		if name == "" {
			ui.Response.SetText("Please enter a player name")
			return
		}

		_, err := fmt.Fprintln(conn, "CONNECT "+name)
		if err != nil {
			ui.Response.SetText("Send error: " + err.Error())
			return
		}
		buttonConnect.Disable()
		buttonConnect.Hide()
		entry.Disable()
		entry.Hide()
	})
	ui.ConnectButton = buttonConnect

	buttonMove := widget.NewButton("MOVE", func() {
		if len(fourButtons.Objects) > 0 {
			fourButtons.Objects = nil
			fourButtons.Refresh()
			return
		}

		for _, exit := range ui.Exits {
			dir := exit
			btn := widget.NewButton(exit, func() {
				_, err := fmt.Fprintln(conn, "MOVE "+dir)
				if err != nil {
					ui.Response.SetText("Send error: " + err.Error())
				}
				fmt.Fprintln(conn, "LOOK")
				fourButtons.Objects = nil
				fourButtons.Refresh()
			})
			fourButtons.Objects = append(fourButtons.Objects, btn)
		}
		fourButtons.Refresh()
		UnselectAll(ui)
	})
	ui.MoveButton = buttonMove

	answer := widget.NewEntry()
	answer.SetPlaceHolder("Enter your response")
	answer.Hide()
	ui.AnswerEntry = answer

	buttonTalk := widget.NewButton("TALK", func() {
		if !ui.Dialogue {
			ui.Dialogue = true
			_, err := fmt.Fprintln(conn, "TALK "+ui.SelectedNpc)
			if err != nil {
				ui.Response.SetText(("Send error: " + err.Error()))
			}
			ui.SelectedNpc = ""
		} else {
			_, err := fmt.Fprintln(conn, answer.Text)
			if err != nil {
				ui.Response.SetText(("Send error: " + err.Error()))
			}
			answer.SetText("")
		}
	})
	//
	// CREATE CONTAINERS
	//
	// SPACE
	space := canvas.NewRectangle(color.Transparent)
	space.SetMinSize(fyne.NewSize(0, 15))
	//
	npcsContainer := CreateContainerNpc(ui)
	monsterContainer := CreateContainerMonster(ui)
	GroupContainer := CreateContainerGroup(ui)
	itemContainer := CreateContainerItem(ui)
	inventoryContainer := CreateContainerInventory(ui)
	playerContainer := CreateContainerPlayer(ui)
	questsContainer := CreateContainerQuest(ui)
	InviteContainer := CreateContainerInvites(ui)
	CreateContainerChat(ui, conn)
	CreateContainerResponse(ui, conn)
	questgiversContainer := CreateContainerQuestgiver(ui)
	//
	// CLEAN
	//
	StatDisplay := container.NewHBox(ui.HpLabel, ui.MaxHpLabel, ui.StatusLabel)
	MonsterStatDisplay := container.NewHBox(ui.MonsterHp, ui.CounterAtk)
	combatDisplay := container.NewBorder(StatDisplay, MonsterStatDisplay, nil, nil)
	midfirst := container.NewHBox(monsterContainer, itemContainer, playerContainer, questsContainer)
	midsecond := container.NewHBox(npcsContainer, questgiversContainer, inventoryContainer, InviteContainer)
	midcontainer := container.NewVBox(midfirst, midsecond)

	// center_res := container.NewCenter(ui.ScrollResponse)
	hContainer := container.New(layout.NewCenterLayout(), ui.RoomTitle)
	h2Container := container.New(layout.NewCenterLayout(), ui.RoomDescription)
	hZone := container.NewBorder(hContainer, h2Container, nil, nil, nil)
	headerMargin := container.NewVBox(space, hZone, space)
	buttons := container.New(layout.NewGridLayout(4), buttonLook, buttonWho, buttonConnect, buttonMove, buttonUse, buttonTalk, buttonTake, buttonAttack, buttonQuest, buttonQuests, buttonCompleteQuests, buttonAbandonQuests, buttonInventory, buttonDrop, buttonStatus, buttonCreate, buttonLeave, buttonInvite, buttonJoin, buttonQuit)
	center_zone := container.NewVBox(ui.ScrollResponse, fourButtons, entry, answer, combatDisplay, GroupContainer, ui.ChatContainer)
	content := container.NewBorder(headerMargin, buttons, nil, midcontainer, center_zone)
	myWindow.Resize(fyne.NewSize(1200, 1200))
	myWindow.SetContent(content)
	myWindow.SetFixedSize(true)

	myWindow.ShowAndRun()
}
