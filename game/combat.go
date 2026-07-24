package game

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

func (g *Game) Attack(conn net.Conn, parts []string) string {
	if len(parts) < 2 {
		return "ERR 407 MISSING_TARGET"
	}

	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]
	player_Room := g.Rooms[player.Room]
	monster_name := strings.Join(parts[1:], " ")
	times := time.Now().Format("15:04:05")

	// get index of monster
	monster_id := -1
	for i := range player_Room.Monsters {
		if player_Room.Monsters[i].Name == monster_name {
			monster_id = i
			if !player_Room.Monsters[i].Alive {
				return "ERR 420 MONSTER_DEAD"
			}
		}
	}

	for i := range player_Room.Npcs {
		if player_Room.Npcs[i].Name == monster_name {
			return "ERR 405 NPC_NOT_HOSTILE"
		}
	}
	for i := range player_Room.Questgivers {
		if player_Room.Questgivers[i].Name == monster_name {
			return "ERR 405 NPC_NOT_HOSTILE"
		}
	}

	if monster_id == -1 {
		return "ERR 404 MONSTER_NOT_FOUND"
	}

	if !player_Room.Monsters[monster_id].InCombat {
		player_Room.Monsters[monster_id].InCombat = true
		player_Room.Monsters[monster_id].Fighting = player
	}

	if player_Room.Monsters[monster_id].Fighting != player {
		return "ERR 422 MONSTER_IN_COMBAT"
	}

	// manage combat damage and status

	player_Room.Monsters[monster_id].HP = player_Room.Monsters[monster_id].HP - player.Damage
	player.Hp = player.Hp - player_Room.Monsters[monster_id].Attack
	player.Status = "combat"

	// both dead
	if player_Room.Monsters[monster_id].HP <= 0 && player.Hp <= 0 {
		player_Room.Monsters[monster_id].Alive = false
		player_Room.Monsters[monster_id].Fighting = nil
		player_Room.Monsters[monster_id].InCombat = false
		player_Room.Monsters[monster_id].Name = "dead " + player_Room.Monsters[monster_id].Name
		player.Hp = 20
		player.Status = "healthy"
		// REMOVE EQUIPPED OBJECT
		if player.Equipped {
			player.Damage -= player.Weapon.Damage
			player.Equipped = false
			player.Weapon = Item{}
		}
		// DROP ITEMS IN ROOM
		for i, desc := range player.Inventory {
			player_Room.Items[i] = desc
			delete(player.Inventory, i)
		}
		for i := range g.Clients {
			if g.Clients[i].Player == nil {
				continue
			}
			if g.Clients[i].Player.Room == player.Room && g.Clients[i].Player.Name != player_name {
				fmt.Println(times + " -> " + player_name + " EVT MONSTER DEAD")
				fmt.Println(times + " -> " + player_name + " EVT PLAYER DEAD")
				fmt.Fprintln(g.Clients[i].Conn, "EVT MONSTER DEAD "+monster_name)
				fmt.Fprintln(g.Clients[i].Conn, "EVT PLAYER DEAD "+player_name)
			}
		}
		player.Room = "1"
		// REMOVE PLAYER FROM ROOM
		for i, p := range player_Room.Players {
			if p == player {
				player_Room.Players = append(player_Room.Players[:i], player_Room.Players[i+1:]...)
				break
			}
		}
		// ADD PLAYER TO ROOM START
		g.Rooms["1"].Players = append(g.Rooms["1"].Players, player)
		return "OK combat=you both died"
	}

	// monster death
	if player_Room.Monsters[monster_id].HP <= 0 {
		old_name := player_Room.Monsters[monster_id].Name
		player_Room.Monsters[monster_id].Alive = false
		player_Room.Monsters[monster_id].Fighting = nil
		player_Room.Monsters[monster_id].InCombat = false
		player_Room.Monsters[monster_id].Name = "dead " + player_Room.Monsters[monster_id].Name
		player.Status = "healthy"
		for i := range g.Clients {
			if g.Clients[i].Player == nil {
				continue
			}
			if g.Clients[i].Player.Room == player.Room && g.Clients[i].Player.Name != player_name {
				fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT MONSTER DEAD")
				fmt.Fprintln(g.Clients[i].Conn, "EVT MONSTER DEAD "+monster_name)
			}
		}
		return "OK combat=" + old_name + " defeated"
	}

	// player death
	if player.Hp <= 0 {
		player.Hp = 20
		player.Status = "healthy"
		// REMOVE EQUIPPED OBJECT
		if player.Equipped {
			player.Damage -= player.Weapon.Damage
			player.Equipped = false
			player.Weapon = Item{}
		}
		// DROP ITEMS IN ROOM
		for i, desc := range player.Inventory {
			player_Room.Items[i] = desc
			delete(player.Inventory, i)
		}
		for i := range g.Clients {
			if g.Clients[i].Player == nil {
				continue
			}
			if g.Clients[i].Player.Room == player.Room && g.Clients[i].Player.Name != player_name {
				fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT PLAYER DEAD")
				fmt.Fprintln(g.Clients[i].Conn, "EVT PLAYER DEAD "+player_name)
			}
		}
		player.Room = "1"
		// REMOVE PLAYER FROM ROOM
		for i, p := range player_Room.Players {
			if p == player {
				player_Room.Players = append(player_Room.Players[:i], player_Room.Players[i+1:]...)
				break
			}
		}
		// ADD PLAYER TO ROOM START
		g.Rooms["1"].Players = append(g.Rooms["1"].Players, player)
		return "OK combat=you died"
	}

	// create response
	res := CombatResponse{
		Player_hp: player.Hp,
		Target_hp: player_Room.Monsters[monster_id].HP,
		Damage:    player.Damage,
		Counter:   player_Room.Monsters[monster_id].Attack,
		Status:    player.Status,
	}

	data, err := json.Marshal(res)

	if err != nil {
		return "ERR 901 SEND_FAILED"
	}

	return "OK combat=" + string(data)
}

func (g *Game) Use(conn net.Conn, parts []string) string {
	if len(parts) < 2 {
		return "ERR 407 MISSING_ITEM_NAME"
	}
	// var setup
	times := time.Now().Format("15:04:05")
	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]
	player_Room := g.Rooms[player.Room]
	item := strings.Join(parts[1:], " ")
	have_item := false
	// check if item in inventory
	for i := range player.Inventory {
		if i == item {
			have_item = true
		}
	}
	if !have_item {
		return "ERR 404 ITEM_NOT_FOUND"
	}
	// check types
	if player.Inventory[item].Type == "atk" && player.Weapon == player.Inventory[item] {
		return "ERR 428 ALREADY_EQUIPPED"
	}
	// HEALS
	if player.Inventory[item].Type == "heal" {
		player.Hp += player.Inventory[item].Heal
		if player.Hp > player.Max_hp {
			player.Hp = player.Max_hp
		}
		delete(player.Inventory, item)
		return "OK used=" + item
	}
	// ATTACK
	if player.Inventory[item].Type == "atk" {
		if !player.Equipped {
			player.Equipped = true
			player.Weapon = player.Inventory[item]
			player.Damage += player.Inventory[item].Damage
			return "OK used=" + item
		} else if player.Equipped {
			player.Damage -= player.Weapon.Damage
			player.Weapon = player.Inventory[item]
			player.Damage += player.Inventory[item].Damage
			return "OK used=" + item
		}
	}
	// SPECIAL
	if player.Inventory[item].Type == "special" {
		switch item {

		case "spider fang":
			if player_Room.Id != 3 {
				return "ERR 429 CANT_USE_NOW"
			}
			id := -1
			for i := range player_Room.Monsters {
				if player_Room.Monsters[i].Name == "steven" {
					id = i
					break
				}
			}
			if id == -1 {
				return "ERR 429 CANT_USE_NOW"
			}
			if player_Room.Monsters[id].Fighting == nil {
				return "ERR 429 CANT_USE_NOW"
			}

			if player.Status == "combat" && player_Room.Monsters[id].Fighting.Name == player_name {
				old_name := player_Room.Monsters[id].Name
				player_Room.Monsters[id].Alive = false
				player_Room.Monsters[id].Fighting = nil
				player_Room.Monsters[id].InCombat = false
				player_Room.Monsters[id].Name = "dead " + player_Room.Monsters[id].Name
				player.Status = "healthy"
				delete(player.Inventory, "spider fang")
				for i := range g.Clients {
					if g.Clients[i].Player == nil {
						continue
					}
					if g.Clients[i].Player.Room == player.Room && g.Clients[i].Player.Name != player_name {
						fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT MONSTER DEAD")
						fmt.Fprintln(g.Clients[i].Conn, "EVT MONSTER DEAD "+old_name)
					}
				}
				return "OK " + old_name + " defeated"
			}

		case "burning coal":
			if player_Room.Id != 4 {
				return "ERR 429 CANT_USE_NOW"
			}
			fragmentCount := 0
			for i := range player.Inventory { // change later to n fragments
				if i == "key fragment one" {
					fragmentCount += 1
				}
				if i == "key fragment two" {
					fragmentCount += 1
				}
				if i == "key fragment three" {
					fragmentCount += 1
				}
			}
			if fragmentCount != 3 {
				return "ERR 430 MISSING_KEY_FRAGMENT"
			}
			// change to include key nb
			key_item := Item{
				Description: "I wonder what it opens",
				Type:        "special",
			}
			delete(player.Inventory, "burning coal")
			delete(player.Inventory, "key fragment one")
			delete(player.Inventory, "key fragment two")
			delete(player.Inventory, "key fragment three")
			player.Inventory["crafted key"] = key_item
			return "OK taken=crafted key"

		case "lit torch":
			if player_Room.Id != 6 {
				return "ERR 429 CANT_USE_NOW"
			}
			g.Rooms["6"].Lit = true
			return "OK used=lit torch, room lit"

		case "magic mushroom":
			player.Max_hp += 50
			player.Hp += 50
			if player.Hp > player.Max_hp {
				player.Hp = player.Max_hp
			}
			delete(player.Inventory, "magic mushroom")
			return "OK used=magic mushroom"

		case "fireball scroll":
			if player_Room.Id != 8 {
				return "ERR 429 CANT_USE_NOW"
			}
			id := -1
			for i := range player_Room.Monsters {
				if player_Room.Monsters[i].Name == "Xorrt'ngst the ancient" {
					id = i
					break
				}
			}
			if id == -1 {
				return "ERR 429 CANT_USE_NOW"
			}

			if player.Status == "combat" && player_Room.Monsters[id].Fighting.Name == player_name {
				delete(player.Inventory, "fireball scroll")
				player_Room.Monsters[id].HP -= 500
				player.Hp -= player_Room.Monsters[id].Attack

				if player_Room.Monsters[id].HP <= 0 && player.Hp <= 0 {
					player_Room.Monsters[id].Alive = false
					player_Room.Monsters[id].Fighting = nil
					player_Room.Monsters[id].InCombat = false
					player_Room.Monsters[id].Name = "dead " + player_Room.Monsters[id].Name
					player.Hp = 20
					player.Status = "healthy"
					// REMOVE EQUIPPED OBJECT
					if player.Equipped {
						player.Damage -= player.Weapon.Damage
						player.Equipped = false
						player.Weapon = Item{}
					}
					// DROP ITEMS IN ROOM
					for i, desc := range player.Inventory {
						player_Room.Items[i] = desc
						delete(player.Inventory, i)
					}
					for i := range g.Clients {
						if g.Clients[i].Player == nil {
							continue
						}
						if g.Clients[i].Player.Room == player.Room && g.Clients[i].Player.Name != player_name {
							fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT MONSTER DEAD")
							fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT PLAYER DEAD")
							fmt.Fprintln(g.Clients[i].Conn, "EVT MONSTER DEAD "+player_Room.Monsters[id].Name)
							fmt.Fprintln(g.Clients[i].Conn, "EVT PLAYER DEAD "+player_name)
						}
					}
					player.Room = "1"
					// REMOVE PLAYER FROM ROOM
					for i, p := range player_Room.Players {
						if p == player {
							player_Room.Players = append(player_Room.Players[:i], player_Room.Players[i+1:]...)
							break
						}
					}
					// ADD PLAYER TO ROOM START
					g.Rooms["1"].Players = append(g.Rooms["1"].Players, player)
					return "OK you both died"
				}

				// monster death
				if player_Room.Monsters[id].HP <= 0 {
					old_name := player_Room.Monsters[id].Name
					player_Room.Monsters[id].Alive = false
					player_Room.Monsters[id].Fighting = nil
					player_Room.Monsters[id].InCombat = false
					player_Room.Monsters[id].Name = "dead " + player_Room.Monsters[id].Name
					player.Status = "healthy"
					for i := range g.Clients {
						if g.Clients[i].Player == nil {
							continue
						}
						if g.Clients[i].Player.Room == player.Room && g.Clients[i].Player.Name != player_name {
							fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT MONSTER DEAD")
							fmt.Fprintln(g.Clients[i].Conn, "EVT MONSTER DEAD "+player_Room.Monsters[id].Name)
						}
					}
					return "OK " + old_name + " defeated"
				}

				// player death
				if player.Hp <= 0 {
					player.Hp = 20
					player.Status = "healthy"
					// REMOVE EQUIPPED OBJECT
					if player.Equipped {
						player.Damage -= player.Weapon.Damage
						player.Equipped = false
						player.Weapon = Item{}
					}
					// DROP ITEMS IN ROOM
					for i, desc := range player.Inventory {
						player_Room.Items[i] = desc
						delete(player.Inventory, i)
					}
					for i := range g.Clients {
						if g.Clients[i].Player == nil {
							continue
						}
						if g.Clients[i].Player.Room == player.Room && g.Clients[i].Player.Name != player_name {
							fmt.Println(times + " -> " + g.Clients[i].Player.Name + " EVT PLAYER DEAD")
							fmt.Fprintln(g.Clients[i].Conn, "EVT PLAYER DEAD "+player_name)
						}
					}
					player.Room = "1"
					// REMOVE PLAYER FROM ROOM
					for i, p := range player_Room.Players {
						if p == player {
							player_Room.Players = append(player_Room.Players[:i], player_Room.Players[i+1:]...)
							break
						}
					}
					// ADD PLAYER TO ROOM START
					g.Rooms["1"].Players = append(g.Rooms["1"].Players, player)
					return "OK you died"
				}

				res := CombatResponse{
					Player_hp: player.Hp,
					Target_hp: player_Room.Monsters[id].HP,
					Damage:    500,
					Counter:   player_Room.Monsters[id].Attack,
					Status:    player.Status,
				}

				data, err := json.Marshal(res)
				if err != nil {
					return "ERR 901 SEND_FAILED"
				}

				return "OK used=firescroll " + string(data)
			}
		}
	}
	return "ERR 429 CANT_USE_NOW"
}
