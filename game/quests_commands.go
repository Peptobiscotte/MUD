package game

import (
	"net"
	"strings"
)

func (g *Game) Quest(conn net.Conn, parts []string) string {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	if len(parts) < 2 {
		return "ERR 407 MISSING_MESSAGE"
	}

	player := g.Clients[conn].Player
	room := g.Rooms[player.Room]
	input := strings.Join(parts[1:], " ")

	for _, npc := range room.Npcs {
		if npc.Name == input {
			return "ERR 406 NO_QUEST_AVAILABLE"
		}
	}

	for i := range room.Questgivers {
		questgiver := &room.Questgivers[i]
		if questgiver.Name != input {
			continue
		}
		if questgiver.Quest_done {
			return "ERR 406 NO_QUEST_AVAILABLE"
		}
		quest := questgiver.Name + "'s quest"
		if _, ok := player.Quests[quest]; ok {
			return "ERR 458 QUEST_ALREADY_ACCEPTED"
		}
		player.Quests[quest] = questgiver.Description
		return "OK quest=" + quest
	}
	return "ERR 404 NPC_NOT_FOUND"
}

func (g *Game) Quests(conn net.Conn, parts []string) string {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	player := g.Clients[conn].Player
	if len(player.Quests) == 0 {
		return "ERR 405 NO_ACTIVE_QUEST"
	}

	lines := []string{}
	for quest := range player.Quests {
		lines = append(lines, quest+": "+player.Quests[quest])
	}
	return "OK quests=" + strings.Join(lines, " || ")
}

func (g *Game) CompleteQuest(conn net.Conn, parts []string) string {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	if len(parts) < 2 {
		return "ERR 407 MISSING_MESSAGE"
	}
	player := g.Clients[conn].Player
	room := g.Rooms[player.Room]
	input := strings.Join(parts[1:], " ")
	find := false

	for quest := range player.Quests {
		if input == quest {
			switch quest {
			case "fouras's quest":
				for item := range player.Inventory {
					if item == "vial of venom" {
						find = true
					}
				}
				if room.Id == 2 && find {
					delete(player.Inventory, "vial of venom")
					delete(player.Quests, quest)
					key_item := Item{
						Description: "What can i do with this...",
						Type:        "special",
					}
					player.Inventory["key fragment two"] = key_item
					for i := range room.Questgivers {
						if room.Questgivers[i].Name == "fouras" {
							room.Questgivers[i].Quest_done = true
							// room.Questgivers[i].PlayerAlreadyFinish = append(room.Questgivers[i].PlayerAlreadyFinish, player.Name)
							break
						}
					}
					return "OK complete=" + quest + "\nQuest finished, nice job!\nYou have received key fragment two"
				} else {
					return "ERR 456 CANT_FINISH_QUEST"
				}

			case "gatekeeper's quest":
				for _, room := range g.Rooms {
					if room.Id < 1 || room.Id > 4 {
						continue
					}
					for _, monster := range room.Monsters {
						if monster.Alive {
							return "ERR 456 CANT_FINISH_QUEST"
						}
					}
				}
				delete(player.Quests, quest)
				for i := range room.Questgivers {
					if room.Questgivers[i].Name == "gatekeeper" {
						room.Questgivers[i].Quest_done = true
						// room.Questgivers[i].PlayerAlreadyFinish = append(room.Questgivers[i].PlayerAlreadyFinish, player.Name)
						break
					}
				}
				return "OK complete=" + quest + "\nQuest finished, you can now enter the dungeon"

			case "Alfred the innkeeper's quest":
				for item := range player.Inventory {
					if item == "onion" {
						find = true
					}
				}
				if room.Id == 3 && find {
					delete(player.Inventory, "onion")
					delete(player.Quests, quest)
					key_item := Item{
						Description: "I won the game i guess",
						Type:        "special",
					}
					player.Inventory["trophy"] = key_item
					for i := range room.Questgivers {
						if room.Questgivers[i].Name == "Alfred the innkeeper" {
							room.Questgivers[i].Quest_done = true
							// room.Questgivers[i].PlayerAlreadyFinish = append(room.Questgivers[i].PlayerAlreadyFinish, player.Name)
							break
						}
					}
					return "OK complete=" + quest + "\nQuest finished, nice job!\nYou have received trophy"
				} else {
					return "ERR 456 CANT_FINISH_QUEST"
				}
			}
		}
	}
	return "ERR 456 CANT_FINISH_QUEST"
}

func (g *Game) AbandonQuest(conn net.Conn, parts []string) string {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	if len(parts) < 2 {
		return "ERR 407 MISSING_MESSAGE"
	}

	player := g.Clients[conn].Player
	input := strings.Join(parts[1:], " ")

	for quest := range player.Quests {
		if input == quest {
			switch quest {
			case "fouras's quest":
				delete(player.Quests, quest)
				return "OK abandon=" + quest
			case "gatekeeper's quest":
				delete(player.Quests, quest)
				return "OK abandon=" + quest
			case "Alfred the innkeeper's quest":
				delete(player.Quests, quest)
				return "OK abandon=" + quest
			}
		}
	}
	return "ERR 457 CANT_ABANDON_QUEST"
}
