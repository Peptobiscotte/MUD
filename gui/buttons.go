package main

import (
	"fmt"
	"net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func CreateButtonLook(ui *UI, conn net.Conn) *widget.Button {
	buttonLook := widget.NewButton("LOOK", func() {
		_, err := fmt.Fprintln(conn, "LOOK")
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonLook
}

func CreateButtonInventory(ui *UI, conn net.Conn) *widget.Button {
	buttonInventory := widget.NewButton("INVENTORY", func() {
		_, err := fmt.Fprintln(conn, "INVENTORY")
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonInventory
}

func CreateButtonUse(ui *UI, conn net.Conn) *widget.Button {
	buttonUse := widget.NewButton("USE", func() {
		_, err := fmt.Fprintln(conn, "USE "+ui.SelectedInventory)
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonUse
}

func CreateButtonDrop(ui *UI, conn net.Conn) *widget.Button {
	buttonDrop := widget.NewButton("DROP", func() {
		_, err := fmt.Fprintln(conn, "DROP "+ui.SelectedInventory)
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonDrop
}

func CreateButtonWho(ui *UI, conn net.Conn) *widget.Button {
	buttonWho := widget.NewButton("WHO", func() {
		_, err := fmt.Fprintln(conn, "WHO")
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonWho
}

func CreateButtonQuit(ui *UI, conn net.Conn) *widget.Button {
	buttonQuit := widget.NewButton("QUIT", func() {
		_, err := fmt.Fprintln(conn, "QUIT")
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonQuit
}

func CreateButtonStatus(ui *UI, conn net.Conn) *widget.Button {
	buttonStatus := widget.NewButton("STATUS", func() {
		if ui.StatusVisible {
			ui.StatusVisible = false
			ui.HpLabel.Hide()
			ui.MaxHpLabel.Hide()
			ui.StatusLabel.Hide()
			return
		}
		ui.StatusVisible = true
		_, err := fmt.Fprintln(conn, "STATUS")
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonStatus
}

func CreateButtonTake(ui *UI, conn net.Conn) *widget.Button {
	buttonTake := widget.NewButton("TAKE", func() {
		_, err := fmt.Fprintln(conn, "TAKE "+ui.SelectedItem)
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonTake
}

// func CreateButtonCheck(ui *UI, conn net.Conn) *widget.Button {
// 	buttonCheck := widget.NewButton("CHECK", func() {
// 		fmt.Println(ui.SelectedQuest)
// 	})
// 	return buttonCheck
// }

func CreateButtonQuest(ui *UI, conn net.Conn) *widget.Button {
	buttonQuest := widget.NewButton("QUEST", func() {
		_, err := fmt.Fprintln(conn, "QUEST "+ui.SelectedQuestgiver)
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonQuest
}

func CreateButtonQuests(ui *UI, conn net.Conn) *widget.Button {
	buttonQuests := widget.NewButton("QUESTS", func() {
		_, err := fmt.Fprintln(conn, "QUESTS")
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonQuests
}

func CreateButtonCompleteQuest(ui *UI, conn net.Conn) *widget.Button {
	buttonCompleteQuest := widget.NewButton("COMPLETE_QUEST", func() {
		_, err := fmt.Fprintln(conn, "COMPLETE_QUEST "+ui.SelectedQuest)
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonCompleteQuest
}

func CreateButtonAbandonQuest(ui *UI, conn net.Conn) *widget.Button {
	buttonAbandonQuest := widget.NewButton("ABANDON_QUEST", func() {
		_, err := fmt.Fprintln(conn, "ABANDON_QUEST "+ui.SelectedQuest)
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonAbandonQuest
}

func CreateButtonCreate(ui *UI, conn net.Conn) *widget.Button {
	buttonCreate := widget.NewButton("GROUP CREATE", func() {
		_, err := fmt.Fprintln(conn, "GROUP CREATE")
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonCreate
}

func CreateButtonGroupInvite(ui *UI, conn net.Conn) *widget.Button {
	buttonInvite := widget.NewButton("GROUP INVITE", func() {
		_, err := fmt.Fprintln(conn, "GROUP INVITE "+ui.SelectedPlayer)
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonInvite
}

func CreateButtonGroupLeave(ui *UI, conn net.Conn) *widget.Button {
	buttonLeave := widget.NewButton("GROUP LEAVE", func() {
		_, err := fmt.Fprintln(conn, "GROUP LEAVE ")
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonLeave
}

func CreateButtonGroupJoin(ui *UI, conn net.Conn) *widget.Button {
	buttonLeave := widget.NewButton("GROUP JOIN", func() {
		_, err := fmt.Fprintln(conn, "GROUP JOIN "+ui.SelectedInvite)
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
		_, err2 := fmt.Fprintln(conn, "GROUPMEM")
		if err2 != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
	})
	return buttonLeave
}

func CreateButtonChatGlobal(ui *UI, conn net.Conn) *widget.Button {
	buttonChatGlobal := widget.NewButton("GLOBAL", func() {
		_, err := fmt.Fprintln(conn, "CHAT GLOBAL "+ui.ChatEntry.Text)
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
		label := widget.NewLabel("GLOBAL " + ui.ChatEntry.Text)
		label.Wrapping = fyne.TextWrapWord
		ui.ChatBox.Add(label)
		ui.ChatContainer.Refresh()
		ui.ChatEntry.Text = ("")
		ui.ScrollChat.ScrollToBottom()
	})
	return buttonChatGlobal
}

func CreateButtonChatRoom(ui *UI, conn net.Conn) *widget.Button {
	buttonChatGlobal := widget.NewButton("ROOM", func() {
		_, err := fmt.Fprintln(conn, "CHAT ROOM "+ui.ChatEntry.Text)
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
		label := widget.NewLabel("ROOM " + ui.ChatEntry.Text)
		label.Wrapping = fyne.TextWrapWord
		ui.ChatBox.Add(label)
		ui.ChatContainer.Refresh()
		ui.ChatEntry.Text = ("")
		ui.ScrollChat.ScrollToBottom()
	})
	return buttonChatGlobal
}

func CreateButtonChatGroup(ui *UI, conn net.Conn) *widget.Button {
	buttonChatGlobal := widget.NewButton("GROUP", func() {
		_, err := fmt.Fprintln(conn, "CHAT GROUP "+ui.ChatEntry.Text)
		if err != nil {
			ui.Response.Text = ("Send error: " + err.Error())
		}
		label := widget.NewLabel("GROUP " + ui.ChatEntry.Text)
		label.Wrapping = fyne.TextWrapWord
		ui.ChatBox.Add(label)
		ui.ChatContainer.Refresh()
		ui.ChatEntry.SetText("")
		ui.ScrollChat.ScrollToBottom()
	})
	return buttonChatGlobal
}
