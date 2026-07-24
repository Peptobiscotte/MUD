package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"tap/game"
	"time"
)

func handleCommand(command string, conn net.Conn, game *game.Game) string {
	if len(command) > 1024 {
		return "ERR 499 LINE_TOO_LONG"
	}

	parts := strings.Fields(command)

	if len(parts) == 0 {
		return "ERR 407 EMPTY_COMMAND"
	}

	if parts[0] != "CONNECT" {
		if game.Clients[conn].Player == nil {
			return "ERR 480 CONNECT FIRST"
		}
	}

	if game.Clients[conn].Player != nil {
		if game.Clients[conn].Player.TalkingNpc != "" {
			return game.Talk(conn, append([]string{"TALK"}, parts...))
		}
		if game.Clients[conn].Player.Status == "combat" &&
			parts[0] != "ATTACK" &&
			parts[0] != "QUIT" &&
			parts[0] != "STATUS" &&
			parts[0] != "USE" {
			return "ERR 485 PLAYER_IN_COMBAT"
		}
	}

	switch parts[0] {
	case "LOOK":
		return game.Look(conn)

	case "WHO":
		return game.Who(conn)

	case "CONNECT":
		return game.Connect(conn, parts)

	case "MOVE":
		return game.Move(conn, parts)

	case "CHAT":
		return game.Chat(conn, parts)

	case "TALK":
		return game.Talk(conn, parts)

	case "QUIT":
		return game.Quit(conn)

	case "TAKE":
		return game.Take(conn, parts)

	case "INVENTORY":
		return game.Inventory(conn)

	case "GROUP":
		return game.Group(conn, parts)

	case "DROP":
		return game.Drop(conn, parts)

	case "STATUS":
		return game.Status(conn)

	case "WTF":
		return game.Badmanners(conn)

	case "ATTACK":
		return game.Attack(conn, parts)

	case "LOOT":
		return game.Loot(conn, parts)

	case "USE":
		return game.Use(conn, parts)

	case "QUEST":
		return game.Quest(conn, parts)

	case "QUESTS":
		return game.Quests(conn, parts)

	case "COMPLETE_QUEST":
		return game.CompleteQuest(conn, parts)

	case "ABANDON_QUEST":
		return game.AbandonQuest(conn, parts)

	case "CHEAT":
		return game.Cheat(conn)

	case "GROUPMEM":
		return game.GroupMem(conn)

	default:
		return "ERR 518 UNKNOWN_COMMAND"
	}
}

func handle(conn net.Conn, game *game.Game) {
	defer func() {
		client, exists := game.Clients[conn]

		if exists && client != nil && client.Player != nil {
			game.Quit(conn)
		}

		delete(game.Clients, conn)
		conn.Close()
		times := time.Now()
		fmt.Printf("%dh %dm %ds ", times.Hour(), times.Minute(), times.Second())
		fmt.Println("Client disconnected:", conn.RemoteAddr())
	}()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		times := time.Now().Format("15:04:05")
		command := strings.TrimSpace(message)
		clientName := "anonymous"
		if client, exists := game.Clients[conn]; exists &&
			client != nil &&
			client.Player != nil {
			clientName = client.Player.Name
		}
		fmt.Println(times + " <- " + clientName + " " + command)
		response := handleCommand(command, conn, game)
		fmt.Println(times + " -> " + clientName + " " + response)
		fmt.Fprintln(conn, response)

		if command == "QUIT" {
			return
		}
	}
}

func main() {

	g := &game.Game{
		Players: make(map[string]*game.Player),
		Clients: make(map[net.Conn]*game.Client),
		Rooms:   make(map[string]*game.Room),
		Groups:  make(map[string]*game.Group),
		GCount:  0,
	}
	// open the json and stock the data in game
	game.OpenJson(g)

	listener, err := net.Listen("tcp", ":4242") // creates listener
	if err != nil {
		fmt.Println("ERR:", err)
		return
	}

	defer listener.Close()
	fmt.Println("Server listening on port 4242")

	for { // infinite loop because we want the server to listen to new connections indefinitely
		conn, err := listener.Accept() // waits for connection to listener port
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		g.Mutex.Lock()
		if len(g.Clients) >= 5 {
			g.Mutex.Unlock()
			fmt.Fprintln(conn, "ERR server_full")
			conn.Close()
			continue
		}
		g.Mutex.Unlock()
		g.Clients[conn] = &game.Client{
			Conn: conn,
		}
		times := time.Now()
		fmt.Printf("%dh %dm %ds ", times.Hour(), times.Minute(), times.Second())
		fmt.Println("Client connected:", conn.RemoteAddr())
		fmt.Fprintln(conn, "OK hello proto=1")
		go handle(conn, g) // starts a go routine and returns to the loop
	}
}
