*This activity has been created as part of the 42 curriculum by mde-nas, mmonnier.*

# The Answer Protocol (TAP)

A small shared-world, text-based multiplayer adventure (MUD) built on top of a custom
TCP protocol, **RFC 42TAP**. The repository contains a Go server, a minimal CLI client,
and a Fyne-based GUI client, all speaking the same line-oriented protocol.

## Description

The Answer Protocol is a client/server multiplayer text adventure. Players connect over
TCP, explore a shared world of interconnected rooms, fight monsters, talk to NPCs, pick
up and manage items, form groups, chat, and complete quests. The world state (rooms,
exits, items, NPCs, monsters, quest givers) is fully data-driven and loaded from a single
JSON file (`config.json`) at server startup.

The goal of the activity is to implement a complete client/server system around a
formally specified protocol (an RFC-style document), including the parts of the spec
that are intentionally left open to the implementers — combat mechanics and quest
progression.

## Instructions

A root [Makefile](Makefile) wraps the commands above:

```bash
make install         # go mod download
make run-server       # go run ./server
make run-client        # go run ./client
make run-client-gui      # make -C gui run (Dockerized Fyne build)
make lint                 # go vet ./...
make clean                  # go clean
```

The GUI client uses the [Fyne](https://fyne.io/) toolkit, which needs a working
OpenGL/X11 environment (`libGL`, `libX11`, `libXxf86vm`, etc.) to link natively; rather
than assume those are installed, `run-client-gui` delegates to the Dockerized build in
[gui/Makefile](gui/Makefile). Run `go run ./gui` directly instead if your machine already
has the native Fyne dependencies set up.

## Resources

**References**

- Go standard library docs: [`net`](https://pkg.go.dev/net),
  [`bufio`](https://pkg.go.dev/bufio), [`encoding/json`](https://pkg.go.dev/encoding/json).
- [Fyne documentation](https://docs.fyne.io/) for the GUI client toolkit.
- Classic MUD references (LPMud/DikuMUD design notes) for general inspiration on
  room/NPC/quest data modeling.

**AI usage**

AI (Claude Code) was used during this activity for:

- **RFC-compliance audit**: reading `protocol-rfc.html` (RFC 42TAP) against the actual
  Go implementation to identify every point where server responses, error formats, or
  event messages diverge from the spec. This produced the deviation list documented in
  [Protocol Implementation](#protocol-implementation) below.

- **This README**: drafted from a full re-read of the codebase (`server/`, `game/`,
  `client/`, `gui/`, `config.json`) and the activity subject PDF, to make sure every
  claim below matches what the code actually does rather than what the spec asks for.

No AI-generated code was merged without being read and understood by the team; the
protocol deviations and known limitations listed below are the result of that review,
not assumptions.

## Architecture

The project is a single Go module (`tap`) split into four packages:

| Package  | Role                                                                 |
|----------|-----------------------------------------------------------------------|
| `server` | TCP listener, per-connection read loop, command dispatcher (`package main`) |
| `game`   | All business logic: rooms, players, combat, quests, groups, chat — no networking code |
| `client` | Minimal CLI client: forwards stdin lines verbatim to the server, prints server output |
| `gui`    | Fyne-based graphical client (buttons, room view, chat, inventory)     |

**Dispatch model**: a single `switch` statement in `handleCommand`
([server/main.go](server/main.go)) maps the first token of each line to a `game.*`
method. This was preferred over a registry/router abstraction because the command set is
small, fixed, and known at compile time — a dispatch table would add indirection without
real benefit here.

**Concurrency model**: one goroutine per TCP connection
(`go handle(conn, g)` in `server/main.go`), each reading lines with
`bufio.Reader.ReadString('\n')` in a loop. All shared game state
(`Players`, `Clients`, `Rooms`, `Groups`) lives in a single `*game.Game` struct guarded
by `g.Mutex`.

## Protocol Implementation

RFC 42TAP intentionally leaves combat and quest mechanics unspecified (§6.1); those
design choices are documented in their own sections below. Separately, an audit of the
wire-level behavior against the RFC surfaced the following deviations from the strict
spec text:

Everything else (command names, `LOOK`/`MOVE`/`TAKE`/`DROP`/`TALK`/`STATUS`/`QUEST`/
`QUESTS`/`WHO`/`GROUP *`/`CHAT *`/`QUIT` responses, `EVT ROOM PRESENCE ENTER|LEAVE`,
`EVT ROOM|GLOBAL|GROUP CHAT`, `EVT GROUP INVITE|JOIN|LEAVE`, and the standard error codes
201/301/401/402/404/405/406) follows the RFC as specified. Game-specific error codes
beyond the RFC's base table (e.g. `442 DOOR_LOCKED`, `443 CHEST_LOCKED`,
`421 CANT_LOOT_ALIVE_MONSTER`) are extensions needed by our world/combat design, applied
consistently in the `ERR <3DIGIT> <MESSAGE>` format.

## Combat System

- Players start with 100 HP (`Max_hp: 100`), tracked per-player in `game.Player`.
- Combat is **not turn-based with initiative**: each `ATTACK <target>` command resolves
  a full exchange in one step — the player deals damage equal to their current
  `Damage` stat (base 10, increased by any equipped weapon via `USE`), and the monster
  simultaneously counters for its fixed `Attack` stat. There is no separate initiative
  roll or randomness: damage values are deterministic, defined per-monster in
  `config.json`.
- No dedicated `DEFEND` or `FLEE` command exists. While `Status == "combat"`, the
  command set is restricted to `ATTACK`, `QUIT`, `STATUS`, and `USE` (see
  `handleCommand` in [server/main.go](server/main.go)) — `USE` lets a player drink a
  healing item or trigger a special item (e.g. `spider fang`, `fireball scroll`) as a
  combat action, which is our substitute for a dedicated action set.
- On death (player HP ≤ 0), the player respawns in the starting room (room `"1"`) with
  20 HP, drops their entire inventory in the room they died in, and any equipped weapon
  is unequipped. This satisfies the "respawn at a safe location with reduced health"
  requirement with a fixed respawn HP rather than a percentage-based formula.
- On monster death, the monster is renamed `"dead <name>"` and left in the room; its
  loot can be retrieved with the separate `LOOT <name>` command rather than being
  granted automatically, so that loot pickup follows the same
  "take = remove from source, no duplication" rule as room items.

**Justification**: a simultaneous-exchange model was chosen over classic turn order
because the protocol is fundamentally request/response (one `ATTACK` = one server
reply) — there is no natural place to interleave a separate "monster's turn" without
extra asynchronous events. This keeps combat resolvable in a single round-trip per
command, at the cost of the "turn-based mechanics with initiative" the subject
describes as a possible design axis.

## Quest System

- `QUEST <npc-name>` grants a quest if the named NPC is a quest giver in the player's
  current room and the player doesn't already hold that quest; quests are keyed by a
  fixed string per quest giver (e.g. `"fouras's quest"`) and stored in
  `player.Quests[quest] = description`.
- `QUESTS` lists all quests currently held by the player.
- `COMPLETE_QUEST <quest-name>` validates completion with quest-specific logic in
  [game/quests_commands.go](game/quests_commands.go):
  - `fouras's quest` — requires being in room `2` with a `vial of venom` in inventory;
    consumes the item, grants `key fragment two`.
  - `Alfred the innkeeper's quest` — requires being in room `3` with an `onion` in
    inventory; consumes the item, grants `trophy`.
  - `gatekeeper's quest` — requires every monster in rooms `1`–`4` to be dead; unlocks
    the north exit from room `5` (the dungeon gate) on completion.
- `ABANDON_QUEST <quest-name>` removes a quest from the player's list.

**Justification**: completion conditions are hard-coded per quest (inventory checks,
room checks, world-state checks) rather than driven by a generic "objective" data
structure, since the world only defines three quests and a generic quest-objective
engine would be significant complexity for no real reuse at this scale.

## World Design

The world (`config.json`) has 8 rooms forming a loop with two dead-end branches, meeting
the "at least 8 rooms, one loop, one branch" requirement:

```
        [2] Village of Devas
       / | \
   [4]  [3] (branch: forge)   (branch: tavern)
   Old   The
   Forge Wanderer's Rest
       \
       [1] Glade of the Powerfull ── loop ── [5] Dungeon Gates
                                                  |
       [8] Demon's Rest ── [7] Necropolis ── [6] Dungeon Hall
         \_______________________loop back to [1]_____________/
```

- **Rooms (8)**: Glade of the Powerfull, Village of Devas, The Wanderer's Rest, The Old
  Forge, The Dungeon Gates, The Dungeon Hall, Necropolis, The Demon's Rest. The main
  loop is `1 → 5 → 6 → 7 → 8 → 1`; rooms `2`, `3`, `4` form a branch off room `1`
  through the village.
- **NPC roles (3 distinct types, satisfies the "≥3 roles" requirement)**:
  - *Dialogue*: the great sage (riddle mini-game), big fat cat, taverier.
  - *Quest givers*: fouras, Alfred the innkeeper, gatekeeper.
  - *Enemies*: spiderling, goblin, steven, fire elemental, Jose, Skeleton King,
    Xorrt'ngst the ancient.
- **Items**: well beyond the minimum of 4 — weapons (`wooden stick`, `goblin dagger`,
  `the destroyer`), healing items (`water of the brave`, `health potion`), quest/key
  items (`vial of venom`, `onion`, `key fragment one/two/three`, `crafted key`), and
  room-locked specials (`chest`, `lit torch`).
- **Quests (3, satisfies "≥2")**: fetch-and-deliver (fouras, Alfred), and a
  clear-the-area quest (gatekeeper) that gates progress into the dungeon.
- **Exploration**: the `1-5-6-7-8-1` loop guarantees a full circuit; the door from room
  `5` to room `6` stays locked (`ERR 442 DOOR_LOCKED`) until the gatekeeper's quest is
  completed, and room `6` starts dark (`ERR 433 ROOM_TOO_DARK` on `LOOK`) until a
  `lit torch` is used.

## Server Logging

The server currently prints a handful of structured lines to stdout (`fmt.Println`) 
for connect/disconnect/listen events

## Group Contributions

- **mmonnier**: GUI client (Fyne widgets, layout, styling, chat views, group UI),
  and a significant share of the server/game logic (combat, chat, group commands).
- **mde-nas**: protocol/game logic (talk & quest flow, NPC dialogue system), world data
  design (`config.json`), and GUI integration work.

Both members contributed across the server, game logic, and GUI layers rather than
working in strictly separate silos; the split above reflects the dominant focus per
area, not exclusive ownership.

## Building and Running

The root [Makefile](Makefile) wraps the Go toolchain:

```bash
make install           # go mod download
make run-server         # go run ./server
make run-client          # go run ./client (CLI)
make run-client-gui       # make -C gui run (Dockerized Fyne build, no native deps needed)
make lint                   # go vet ./...
make clean                    # go clean
```

`go run ./gui` can be used directly instead of `make run-client-gui` on a machine with
native Fyne/OpenGL/X11 dependencies already installed. `go build ./...` compiles all
four packages.

## Testing

There is no automated test suite yet (no `*_test.go` files) — see
[Roadmap](#roadmap--known-limitations). Current manual testing procedure:

1. Start the server: `go run ./server`.
2. Connect two or more CLI/GUI clients and run `CONNECT <name>` on each.
3. **Multiplayer**: from one client, `MOVE` between rooms and confirm the others receive
   `EVT ROOM PRESENCE ENTER/LEAVE`; use `CHAT GLOBAL|ROOM|GROUP` and confirm delivery
   scoping.
4. **Combat**: `MOVE` to a room with a monster (e.g. room `1`, `spiderling`), run
   `ATTACK <name>` repeatedly, check `STATUS` between rounds, and confirm death/respawn
   behavior (inventory drop, HP reset, room reset to `1`).
5. **Quests**: `TALK`/`QUEST` a quest giver (e.g. `fouras` in room `2`), fetch the
   required item, run `COMPLETE_QUEST`, and confirm the reward is granted and `QUESTS`
   reflects the change.
