package game

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func (g *Game) Connect(conn net.Conn, parts []string) string {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	if len(parts) < 2 {
		return "ERR 202 MISSING_NAME"
	}

	times := time.Now().Format("15:04:05")
	name := parts[1]
	_, ok := g.Players[name]

	if ok {
		return "ERR 201 NAME_IN_USE"
	}

	if g.Clients[conn].Player != nil {
		return "ERR 203 ALREADY_CONNECTED"
	}

	g.Players[name] = &Player{ // creates the player
		Name:      name,
		Room:      "1",
		Inventory: make(map[string]Item),
		Hp:        100,
		Max_hp:    100,
		Status:    "healthy",
		Damage:    10,
		Quests:    make(map[string]string),
	}
	g.Clients[conn].Player = g.Players[name]
	g.Rooms["1"].Players = append(g.Rooms["1"].Players, g.Players[name])
	player_name := g.Clients[conn].Player.Name
	player_room := g.Clients[conn].Player.Room
	g.PCount = 1
	if len(g.Clients) > 0 {
		for i := range g.Clients {
			if g.Clients[i].Player == nil {
				continue
			}
			if g.Clients[i].Player.Name != player_name {
				g.PCount += 1
				if g.Clients[i].Player.Room == player_room {
					fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT ROOM PRESENCE ENTER " + player_name)
					fmt.Fprintln(g.Clients[i].Conn, "EVT ROOM PRESENCE ENTER "+player_name)
				}
			}
		}
		for i := range g.Clients {
			if g.Clients[i].Player == nil {
				continue
			}
			fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT STATS players=" + strconv.Itoa(g.PCount))
			fmt.Fprintln(g.Clients[i].Conn, "EVT STATS players="+strconv.Itoa(g.PCount))
		}
	}
	return "OK connected " + player_name
}

func (g *Game) Status(conn net.Conn) string {
	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]

	res := StatusResponse{
		Hp:     player.Hp,
		Max_hp: player.Max_hp,
		Status: player.Status,
	}

	data, err := json.Marshal(res)

	if err != nil {
		return "ERR 901 SEND_FAILED"
	}

	return "OK status=" + string(data)
}

func (g *Game) Look(conn net.Conn) string {
	player_room := g.Clients[conn].Player.Room
	room := g.Rooms[player_room]
	players := []string{} // {} at the end

	if player_room == "6" && !room.Lit {
		return "ERR 433 ROOM_TOO_DARK"
	}

	for _, p := range room.Players {
		players = append(players, p.Name)
	}
	npcs := []string{}
	monsters := []string{}
	items := []string{}
	questgivers := []string{}

	for name := range room.Items {
		items = append(items, name)
	}
	for _, npc := range room.Npcs {
		npcs = append(npcs, npc.Name)
	}
	for _, monster := range room.Monsters {
		monsters = append(monsters, monster.Name)
	}
	for _, questgiver := range room.Questgivers {
		questgivers = append(questgivers, questgiver.Name)
	}
	res := LookResponse{
		Room: RoomResponse{
			ID:          room.Id,
			Name:        room.Name,
			Description: room.Description,
			Exits:       room.Exits,
		},
		Players:     players,
		Items:       items,
		NPCs:        npcs,
		Monsters:    monsters,
		Questgivers: questgivers,
	}
	data, err := json.Marshal(res)
	if err != nil {
		return "ERR 901 SEND_FAILED"
	}
	return "OK " + string(data)
}

func (g *Game) Who(conn net.Conn) string {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	names := []string{}
	for name := range g.Players {
		names = append(names, name)
	}
	nb := strconv.Itoa(len(names))
	return "OK players=" + nb
}

func (g *Game) Move(conn net.Conn, parts []string) string {
	if len(parts) < 2 {
		return "ERR 302 MISSING_DIRECTION"
	}
	times := time.Now().Format("15:04:05")
	direction := parts[1]
	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]
	cur_room := g.Rooms[player.Room]
	old_room := player.Room

	// _, ok := player.Inventory["crafted key"]

	// if player.Room == "5" && direction == "north" && !ok {
	// 	return "ERR 442 DOOR_LOCKED"
	// }
	if player.Room == "5" && !cur_room.Questgivers[0].Quest_done && direction == "north" {
		return "ERR 442 DOOR_LOCKED"
	}

	for i := range cur_room.Exits {
		if cur_room.Exits[i] == direction {
			player.Room = cur_room.Neighbours[direction]
			// EVT presence enter and leave
			for j, client := range g.Clients {
				if client == nil || client.Player == nil {
					continue
				}

				if client.Player.Name == player_name {
					continue
				}

				if client.Player.Room == old_room {
					fmt.Println(times + " -> " + g.Clients[j].Player.Name + " EVT ROOM PRESENCE LEAVE " + player_name)
					fmt.Fprintln(
						client.Conn,
						"EVT ROOM PRESENCE LEAVE "+player_name,
					)
				}

				if client.Player.Room == player.Room {
					fmt.Println(times + " -> " + g.Clients[j].Player.Name + " EVT ROOM PRESENCE ENTER " + player_name)
					fmt.Fprintln(
						client.Conn,
						"EVT ROOM PRESENCE ENTER "+player_name,
					)
				}
			}
			// add player to room, remove from old
			g.Rooms[player.Room].Players = append(g.Rooms[player.Room].Players, g.Players[player_name])
			players := g.Rooms[old_room].Players
			for i, p := range players {
				if p == player {
					g.Rooms[old_room].Players = append(players[:i], players[i+1:]...)
					break
				}
			}
			return "OK room=" + g.Rooms[player.Room].Name
		}
	}
	return "ERR 301 NO_EXIT"
}

func (g *Game) Talk(conn net.Conn, parts []string) string {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	if len(parts) < 2 {
		return "ERR 407 MISSING_MESSAGE"
	}

	player := g.Clients[conn].Player
	room := g.Rooms[player.Room]
	input := strings.Join(parts[1:], " ")

	// not currently in a conversation must name an NPC present in the room
	if player.TalkingNpc == "" {
		for _, npc := range room.Npcs {
			if npc.Name != input {
				continue
			}
			for _, name := range npc.AlreadyTalkTo {
				if player.Name == name {
					return "OK end=Good luck with your adventure"
				}
			}
			node, ok := npc.Nodes[npc.Start]
			if !ok {
				return "ERR 404 NPC_NOT_FOUND"
			}
			player.TalkingNpc = npc.Name
			player.DialogueNode = npc.Start
			return "OK " + node.Text
		}
		return "ERR 404 NPC_NOT_FOUND"
	}

	// already talking to an npc: try to advance in the dialogue
	var npc *Npc

	for i := range room.Npcs {
		if room.Npcs[i].Name == player.TalkingNpc {
			npc = &room.Npcs[i]
			break
		}
	}

	node := npc.Nodes[player.DialogueNode]
	next, ok := node.Responses[input]
	if !ok {
		next, ok = node.Responses["wrong"]
	}
	if !ok {
		return "ERR 404 INVALID_RESPONSE"
	}

	nextNode, ok := npc.Nodes[next]
	if !ok {
		return "ERR 404 NPC_NOT_FOUND"
	}

	if nextNode.End {
		player.TalkingNpc = ""
		player.DialogueNode = ""
		if next == "correct" {
			// npc.Dialogue_done = true
			npc.AlreadyTalkTo = append(npc.AlreadyTalkTo, player.Name)
		}
	} else {
		player.DialogueNode = next
	}
	return "OK " + nextNode.Text
}

func (g *Game) Quit(conn net.Conn) string {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	client := g.Clients[conn]
	times := time.Now().Format("15:04:05")

	if client != nil && client.Player != nil {
		player_room := client.Player.Room
		player_name := client.Player.Name
		room := g.Rooms[player_room]
		for i := range g.Clients {
			if g.Clients[i].Player == nil {
				continue
			}
			if g.Clients[i].Player.Name != player_name {
				if g.Clients[i].Player.Room == player_room {
					fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT ROOM PRESENCE LEAVE " + player_name)
					fmt.Fprintln(g.Clients[i].Conn, "EVT ROOM PRESENCE LEAVE "+player_name)
				}
			}
		}
		for i, p := range room.Players {
			if p.Name == client.Player.Name {
				room.Players = append(room.Players[:i], room.Players[i+1:]...) // remove from room
				break
			}
		}
		g.GroupLeave(conn)
		delete(g.Players, client.Player.Name) // delete removes the g.players in key client.Player.Name
		g.PCount -= 1
		for i := range g.Clients {
			if g.Clients[i].Player == nil {
				continue
			}
			fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT STATS players=" + strconv.Itoa(g.PCount))
			fmt.Fprintln(g.Clients[i].Conn, "EVT STATS players="+strconv.Itoa(g.PCount))
		}
	}

	delete(g.Clients, conn)
	return "OK bye"
}

func (g *Game) Take(conn net.Conn, parts []string) string {
	times := time.Now().Format("15:04:05")
	player_room := g.Clients[conn].Player.Room
	room := g.Rooms[player_room]
	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]

	if len(parts) < 2 {
		return "ERR 407 MISSING_ITEM_NAME"
	}

	item := strings.Join(parts[1:], " ")
	_, ok := room.Items[item]
	_, ok_key := player.Inventory["crafted key"]
	if ok {
		if item == "chest" && !ok_key {
			return "ERR 443 CHEST_LOCKED"
		} else if item == "chest" && ok_key {
			g.Clients[conn].Player.Inventory["onion"] = Item{
				Description: "Stinky",
				Type:        "special",
			}
			delete(room.Items, item)
			return "OK taken=onion"
		} else {
			g.Clients[conn].Player.Inventory[item] = room.Items[item]
			delete(room.Items, item)
			for i := range g.Clients {
				if g.Clients[i].Player == nil {
					continue
				}
				if g.Clients[i].Player.Room == player_room && g.Clients[i].Player.Name != player_name {
					fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT ITEM TAKEN " + item)
					fmt.Fprintln(g.Clients[i].Conn, "EVT ITEM TAKEN "+item)
				}
			}
			return "OK taken=" + item
		}
	}
	return "ERR 404 ITEM_NOT_FOUND"
}

func (g *Game) Drop(conn net.Conn, parts []string) string {
	times := time.Now().Format("15:04:05")
	player_room := g.Clients[conn].Player.Room
	player_name := g.Clients[conn].Player.Name
	room := g.Rooms[player_room]
	if len(parts) < 2 {
		return "ERR 407 MISSING_ITEM_NAME"
	}
	item := strings.Join(parts[1:], " ")
	_, ok := g.Clients[conn].Player.Inventory[item]
	if ok {
		room.Items[item] = g.Clients[conn].Player.Inventory[item]
		delete(g.Clients[conn].Player.Inventory, item)
		for i := range g.Clients {
			if g.Clients[i].Player == nil {
				continue
			}
			if g.Clients[i].Player.Room == player_room && g.Clients[i].Player.Name != player_name {
				fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT ITEM DROPPED " + item)
				fmt.Fprintln(g.Clients[i].Conn, "EVT ITEM DROPPED "+item)
			}
		}
		return "OK dropped=" + item
	}
	return "ERR 404 ITEM_NOT_IN_INVENTORY"
}

func (g *Game) Inventory(conn net.Conn) string {
	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]
	data, err := json.Marshal(player.Inventory)
	if err != nil {
		return "ERR 901 SEND_FAILED"
	}
	return "OK inventory=" + string(data)
}

func (g *Game) Loot(conn net.Conn, parts []string) string {
	if len(parts) < 2 {
		return "ERR 407 MISSING_TARGET_NAME"
	}
	target := strings.Join(parts[1:], " ")
	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]
	player_room := g.Rooms[player.Room]

	// get the id of monster
	id := -1

	for i := range player_room.Monsters {
		if target == player_room.Monsters[i].Name {
			id = i
		}
	}
	// check if possible to loot
	if id == -1 {
		return "ERR 404 MONSTER_NOT_FOUND"
	}

	if len(player_room.Monsters[id].Items) == 0 {
		return "ERR 404 NOTHING_TO_LOOT"
	}

	if player_room.Monsters[id].Alive {
		return "ERR 421 CANT_LOOT_ALIVE_MONSTER"
	}

	// loot
	for name, desc := range player_room.Monsters[id].Items {
		player.Inventory[name] = desc
		delete(player_room.Monsters[id].Items, name)
		return "OK taken=" + name
	}
	return ""
}

func (g *Game) Cheat(conn net.Conn) string {
	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]
	key_item := Item{
		Description: "I wonder what it opens",
		Type:        "special",
	}
	attack_item := Item{
		Description: "CHEATING WEAPON",
		Type:        "atk",
		Damage:      1000,
	}
	player.Inventory["crafted key"] = key_item
	player.Inventory["cheating_gun"] = attack_item
	cur_room := g.Rooms[player.Room]
	if player.Room == "5" {
		cur_room.Questgivers[0].Quest_done = true
	}
	return "OK CHEAT MODE ACTIVATED"
}

func (g *Game) Badmanners(conn net.Conn) string {
	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]
	if g.Rooms[player.Room].Questgivers[0].Quest_done {
		return "OK"
	}
	return "ERR"
}
