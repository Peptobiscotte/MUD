package game

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Chat functions

func (g *Game) Chat(conn net.Conn, parts []string) string {
	if len(parts) < 2 {
		return "ERR 407 MISSING_ARGUMENT"
	}
	switch parts[1] {
	case "GLOBAL":
		return g.ChatGlobal(conn, parts)

	case "ROOM":
		return g.ChatRoom(conn, parts)

	case "GROUP":
		return g.ChatGroup(conn, parts)
	}
	return "ERR 407 WRONG ARGUMENT"
}

func (g *Game) ChatGlobal(conn net.Conn, parts []string) string {
	if len(parts) < 3 {
		return "ERR 407 MISSING_MESSAGE"
	}
	times := time.Now().Format("15:04:05")
	player_name := g.Clients[conn].Player.Name

	message := strings.Join(parts[2:], " ")

	for i := range g.Clients {
		if g.Clients[i].Player.Name != player_name {
			fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT GLOBAL CHAT " + player_name + ": " + message)
			fmt.Fprintln(g.Clients[i].Conn, "EVT GLOBAL CHAT "+player_name+": "+message)
		}
	}

	return "OK"
}

func (g *Game) ChatRoom(conn net.Conn, parts []string) string {
	if len(parts) < 3 {
		return "ERR 407 MISSING_MESSAGE"
	}
	times := time.Now().Format("15:04:05")
	player_name := g.Clients[conn].Player.Name
	player_room := g.Players[player_name].Room

	message := strings.Join(parts[2:], " ")

	for i := range g.Clients {
		if g.Clients[i].Player.Name != player_name {
			if g.Clients[i].Player.Room == player_room {
				fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT GLOBAL CHAT " + player_name + ": " + message)
				fmt.Fprintln(g.Clients[i].Conn, "EVT ROOM CHAT "+player_name+": "+message)
			}
		}
	}

	return "OK"
}

func (g *Game) ChatGroup(conn net.Conn, parts []string) string {
	if len(parts) < 3 {
		return "ERR 407 MISSING_MESSAGE"
	}

	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]
	times := time.Now().Format("15:04:05")
	if player.Group == nil {
		return "ERR 401 NOT_IN_GROUP"
	}

	message := strings.Join(parts[2:], " ")

	for i := range g.Clients {
		if g.Clients[i].Player.Name != player_name {
			if g.Clients[i].Player.Group != nil && g.Clients[i].Player.Group == player.Group {
				fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT GLOBAL CHAT " + player_name + ": " + message)
				fmt.Fprintln(g.Clients[i].Conn, "EVT GROUP CHAT "+player_name+": "+message)
			}
		}
	}

	return "OK"
}

// end Chat functions
