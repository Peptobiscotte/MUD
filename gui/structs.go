package main

import (
	"tap/game"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	RoomTitle       *canvas.Text
	RoomDescription *canvas.Text

	Players     []string
	PlayersList *widget.List

	Monsters        []string
	MonstersList    *widget.List
	SelectedMonster string

	Items        []string
	ItemsList    *widget.List
	SelectedItem string

	Inventory         []InventoryEntry
	InventoryList     *widget.List
	SelectedInventory string
	InventoryVisible  bool

	Response       *widget.Label
	ScrollResponse *container.Scroll

	Npcs        []string
	NpcsList    *widget.List
	SelectedNpc string

	Exits      []string
	MoveButton *widget.Button

	ConnectButton *widget.Button

	QuestgiversList    *widget.List
	Questgivers        []string
	SelectedQuestgiver string

	Quests        []string
	QuestsList    *widget.List
	SelectedQuest string

	Dialogue    bool
	AnswerEntry *widget.Entry
	ScrollChat  *container.Scroll

	StatusVisible bool
	HpLabel       *canvas.Text
	MaxHpLabel    *widget.Label
	StatusLabel   *widget.Label

	Group          []string
	GroupList      *widget.List
	SelectedPlayer string

	Invites        []string
	InviteList     *widget.List
	SelectedInvite string

	PlayerName    string
	Chat          []string
	ChatContainer *fyne.Container
	ChatBox       *fyne.Container
	ChatEntry     *widget.Entry

	MonsterHp  *widget.Label
	CounterAtk *widget.Label
}

type InventoryEntry struct {
	Name string
	Item game.Item
}
