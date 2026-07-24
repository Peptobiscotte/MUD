package game

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// group functions

func (g *Game) Group(conn net.Conn, parts []string) string {
	if len(parts) < 2 {
		return "ERR 407 MISSING_ARGUMENT"
	}
	switch parts[1] {
	case "CREATE":
		return g.GroupCreate(conn, parts)

	case "INVITE":
		return g.GroupInvite(conn, parts)

	case "JOIN":
		return g.GroupJoin(conn, parts)

	case "LEAVE":
		return g.GroupLeave(conn)
	}

	return ""
}

func (g *Game) GroupCreate(conn net.Conn, parts []string) string {
	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]
	if player.Group != nil {
		return "ERR 402 ALREADY_IN_GROUP"
	}
	g.GCount += 1
	group_name := player_name + "party" + strconv.Itoa(g.GCount)
	g.Groups[group_name] = &Group{
		ID:      group_name,
		Leader:  g.Players[player_name],
		Members: make(map[string]*Player),
	}
	g.Groups[group_name].Members[player_name] = player
	player.Group = g.Groups[group_name]
	return "OK group=" + group_name
}

func (g *Game) GroupInvite(conn net.Conn, parts []string) string {
	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]
	times := time.Now().Format("15:04:05")

	if len(parts) < 3 {
		return "ERR 407 MISSING_PLAYER_NAME"
	}

	if player.Group == nil {
		return "ERR 401 NOT_IN_GROUP"
	}

	if parts[2] == player_name {
		return "ERR 408 AUTO_INVITE"
	}

	_, ok := g.Players[parts[2]]
	if !ok {
		return "ERR 409 PLAYER_NOT_FOUND"
	}
	var invite_conn net.Conn
	var invite_name string

	for i := range g.Clients {
		if g.Clients[i].Player.Name == parts[2] {
			invite_conn = g.Clients[i].Conn
			invite_name = g.Clients[i].Player.Name
		}
	}

	if g.Players[invite_name].Group != nil {
		return "ERR 402 ALREADY_IN_GROUP"
	}
	groupName := player.Group.ID
	fmt.Println(times + " -> " + g.Clients[invite_conn].Player.Name + " EVT GROUP INVITE " + groupName)
	fmt.Fprintln(invite_conn, "EVT GROUP INVITE "+groupName)
	g.Players[invite_name].Invites = append(g.Players[invite_name].Invites, groupName)
	return "OK"
}

func (g *Game) GroupJoin(conn net.Conn, parts []string) string {
	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]
	times := time.Now().Format("15:04:05")

	if len(parts) < 3 {
		return "ERR 407 MISSING_GROUP_NAME"
	}

	if player.Invites == nil {
		return "ERR 410 NO_INVITES"
	}

	var is_valid bool

	for i := range player.Invites {
		if player.Invites[i] == parts[2] {
			is_valid = true
			player.Invites = append(player.Invites[:i], player.Invites[i+1:]...)
			break
		}
	}
	if is_valid {
		var group_names []string
		var group_conn []net.Conn
		group := g.Groups[parts[2]]
		for i := range group.Members {
			group_names = append(group_names, group.Members[i].Name)
		}
		for j := range group_names {
			for k := range g.Clients {
				if g.Clients[k].Player.Name == group_names[j] {
					group_conn = append(group_conn, g.Clients[k].Conn)
				}
			}
		}
		for l := range group_conn {
			fmt.Println(times + " -> " + " EVT GROUP JOIN " + player_name)
			fmt.Fprintln(group_conn[l], "EVT GROUP JOIN "+player_name)
		}
		group.Members[player.Name] = player
		player.Group = group
		return "OK group=" + parts[2]
	}
	return "ERR 411 INVALID_INVITE"
}

func (g *Game) GroupLeave(conn net.Conn) string {
	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]
	times := time.Now().Format("15:04:05")

	if player.Group == nil {
		return "ERR 401 NOT_IN_GROUP"
	}
	group := player.Group
	group_name := group.ID
	delete(player.Group.Members, player_name)
	if player.Group.Leader.Name == player_name && len(player.Group.Members) > 0 {
		for _, member := range player.Group.Members {
			player.Group.Leader = member
			break
		}
	}
	if len(group.Members) > 0 {
		var members_name []string
		var members_conn []net.Conn

		for i := range group.Members {
			members_name = append(members_name, group.Members[i].Name)
		}
		for j := range members_name {
			for k := range g.Clients {
				if g.Clients[k].Player.Name == members_name[j] {
					members_conn = append(members_conn, g.Clients[k].Conn)
				}
			}
		}
		for l := range members_conn {
			fmt.Println(times + " -> " + " EVT GROUP LEAVE " + player_name)
			fmt.Fprintln(members_conn[l], "EVT GROUP LEAVE "+player_name)
		}
	}
	if len(player.Group.Members) == 0 {
		delete(g.Groups, group_name)
	}
	player.Group = nil
	return "OK group left"
}

func (g *Game) GroupMem(conn net.Conn) string {
	player_name := g.Clients[conn].Player.Name
	player := g.Players[player_name]
	group := player.Group
	res := "OK members="

	for i := range group.Members {
		res += " " + group.Members[i].Name
	}
	return res
}

// end group functions
