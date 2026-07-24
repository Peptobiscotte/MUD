# To-Do
- [OK] command cant be other than connect if not connected
## GROUP INVITE
- [OK]cant invite itself
- [OK]cant invite player in a group
## GROUP JOIN
- [OK]send event to all group members when player joins
## GROUP LEAVE
- [OK]remove group member when leaving game
## CHAT ALL
- [OK]implement global chat
- [OK]implement room chat
- [OK]implement group chat
## STATUS
- [OK]implement status command and modification to player struct
## LOOK
- [OK]look -> update players when delete
## COMBAT SYSTEM
- [OK]create attackable npcs in json
- [OK]can only attack monsters
- [OK]combat status limits commands
- []if player dies, reset mob hp? -> add mob max hp on struct?
- []IF BOTH DED CHECK PLAYER LIFE AND DED TOO
- [WIP]design combat system
## LOOT SYSTEM
- [OK]can loot multiple items from dead mob
- [?]remove dead mob from room when empty ???
## QUEST SYSTEM
- [OK]QUEST command to accept a quest from npc
- [OK]QUESTS command to list all accepted quests -> modify player struct to accept Quests map[string]string
- [OK]auto-validate versus manual completion by coming back to the quest-giver
- []quest chains -> quest 2 from npc becomes available only if quest 1 is done
- [OK]create QUESTS
- [OK]complete + abandon quest
- [OK]cannot take a finish quest
- [OK]second quest implemented gardener of gate
- []condtion to enter the dungeon - quest of the gardener finished
## TALK
- [OK]no hardcode on TALK
- [OK]works for first room but not for others, change code should run with any room
- [OK]fix problem when wrong answer of the riddle, continue to try (or not)
- [OK]create test with other rooms
- []handle when already talk
## JSON
- [WIP]think about the next step (monsters, how many quest, purpose of the game etc)


## GUI
- []can click and click again to display/hide
- []find a way to update look and quests

## MOVE
- []wait for look?

## CONNECT
- [OK]handle when no name input

## TALK
- []global state for dialogue

upodate fireball
