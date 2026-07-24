package game

import (
	"encoding/json"
	"fmt"
	"os"
)

func OpenJson(game *Game) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}

	var rooms []Room
	erro := json.Unmarshal(data, &rooms)
	if erro != nil {
		fmt.Println("error unmarshalling", erro)
		return
	}
	for i := range rooms {
		room := rooms[i]
		id := fmt.Sprint(room.Id)
		game.Rooms[id] = &room
	}
}
