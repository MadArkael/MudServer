package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	headSlot    string = "Head"
	faceSlot    string = "Face"
	neckSlot    string = "Neck"
	aboutSlot   string = "About"
	chestSlot   string = "Chest"
	backSlot    string = "Back"
	holdLSlot   string = "Left Hand"
	holdRSlot   string = "Right Hand"
	waistSlot   string = "Waist"
	legsSlot    string = "Legs"
	feetSlot    string = "Feet"
	armsSlot    string = "Arms"
	wristLSlot  string = "Left Wrist"
	wristRSlot  string = "Right Wrist"
	handsSlot   string = "Hands"
	fingerLSlot string = "Left Finger"
	fingerRSlot string = "Right Finger"
)

var clientOutputChan chan ClientOutput

type Command struct {
	cmnd string
	desc string
}

type Emote struct {
	name string
	fP   string
	fPt  string
	tar  string
	tP   string
	tPt  string
}

type Room struct {
	name  string
	desc  string
	id    int
	exits []*Exit
	users []*User
	items []*Item
}

type Exit struct {
	keyword  string
	lookMsg  string
	linkedID int
}

type InputEvent struct {
	msg string
}

type BroadcastEvent struct {
}

type UserJoinedEvent struct {
}

type UserLeftEvent struct {
	user *User
}

type ClientInput struct {
	user  *User
	event interface{}
	world *World
}

type ClientOutput struct {
	user  *User
	msg   string
	event interface{}
	world *World
}

type User struct {
	name    string
	session *Session
	room    *Room
	char    *Character
	buf     []byte
}

type Character struct {
	name   string
	user   *User
	room   *Room
	class  string
	desc   string
	status int
	str    int
	dex    int
	con    int
	intl   int
	wis    int
	cha    int
	eq     map[string]*Item
	inv    []*Item
	purse  int
	fort   int
	ref    int
	wil    int
	att    int
	dam    int
	hp     int
	mana   int
	moves  int
}

type Item struct {
	id   int
	name string
	desc string
	slot string
	loc  Location
	uID  string
}

type Container interface {
	getItems() []*Item
	getName() string
	add(item *Item)
	remove(item *Item)
	contains(item *Item) bool
}

type Location interface {
	getLocation() Location
	getName() string
}

func (u *User) getLocation() Location {
	return u
}
func (u *User) getName() string {
	return u.name
}

func (u *User) getItems() []*Item {
	return u.char.inv
}

func (u *User) add(item *Item) {
	u.char.inv = append(u.char.inv, item)
}

func (u *User) remove(item *Item) {

	for i, j := range u.char.inv {
		if j == item {
			u.char.inv = append(u.char.inv[:i], u.char.inv[i+1:]...)
		}
	}
}

func (u *User) contains(item *Item) bool {
	for _, it := range u.char.inv {
		if it == item {
			return true
		}
	}
	return false
}

func (r *Room) getLocation() Location {
	return r
}

func (r *Room) getName() string {
	return r.name
}

func (r *Room) getItems() []*Item {
	return r.items
}

func (r *Room) remove(item *Item) {

	for i, j := range r.items {
		if j == item {
			r.items = append(r.items[:i], r.items[i+1:]...)
		}
	}
}
func (r *Room) add(item *Item) {
	r.items = append(r.items, item)
}

func (r *Room) contains(item *Item) bool {
	for _, it := range r.items {
		if it == item {
			return true
		}
	}
	return false
}

type Session struct {
	conn net.Conn
}

type World struct {
	users []*User
	rooms []*Room
	cmnds []*Command

	emotes []*Emote
	eqList []string
	items  map[string]map[int]*Item
}

// todo load data from disk
func (w *World) loadEmotes() {
	w.emotes = []*Emote{
		{
			name: "nod",
			fP:   "You nod.",
			fPt:  "You nod to ",
			tar:  " nods to you.",
			tP:   " nods.",
			tPt:  " nods to ",
		},
		{
			name: "flail",
			fP:   "You flail your arms about.",
			fPt:  "You flail your arms at ",
			tar:  " flails their arms at you.",
			tP:   " flails their arms.",
			tPt:  " flails their arms at ",
		},
		{
			name: "laugh",
			fP:   "You laugh loudly.",
			fPt:  "You laugh at ",
			tar:  " laughs at you.",
			tP:   " laughs.",
			tPt:  " laughs at ",
		},
		{
			name: "smile",
			fP:   "You smile.",
			fPt:  "You smile at ",
			tar:  " smiles at you.",
			tP:   " smiles.",
			tPt:  " smiles at ",
		},
		{
			name: "bird",
			fP:   "You show everyone what you think of them. They're obviously number one!",
			fPt:  "You flip off ",
			tar:  " shows you a single digit salute.",
			tP:   " flips everyone and everything, off.",
			tPt:  " flips off ",
		},
		{
			name: "point",
			fP:   "You point at nothing in particular.",
			fPt:  "You point at ",
			tar:  " points at you.",
			tP:   " points at something you arent able to discern.",
			tPt:  " points at ",
		},
	}
}

// todo load data from disk
func (w *World) loadHelp() {
	w.cmnds = []*Command{
		{
			cmnd: "look, l, or l <dir>",
			desc: "Redisplays the room description. Can take a direction argument.",
		},
		{
			cmnd: "<exit dir>",
			desc: "Moves you in the direction specified (north, south, west, east, up, down, n,s,e,w,u,d).",
		},
		{
			cmnd: "go <exit dir>",
			desc: "Moves you in the direction specified (in, out, through, i, o, t).",
		},
		{
			cmnd: "help",
			desc: "The dialogue you're looking at right now.",
		},
		{
			cmnd: "say <text>",
			desc: "Tries to speak to other users. Does not work if they're not here.",
		},
		{
			cmnd: "yell <text>",
			desc: "Like say, except it can be heard from four rooms in any direction.",
		},
		{
			cmnd: "shout <text>",
			desc: "Like say/yell, but heard everywhere.",
		},
		{
			cmnd: "remove, rem <item name>",
			desc: "remove item worn.",
		},
		{
			cmnd: "wear <item name>",
			desc: "Tries to equip item.",
		},
		{
			cmnd: "create <item name>|<item desc>|<slot #>",
			desc: "Creates a #slot item. Must use | to separate fields.",
		},
		{
			cmnd: "slots",
			desc: "Shows what equipment belongs to what slot for create",
		},
		{
			cmnd: "new <item id or name>",
			desc: "Tries to give <item>. Has to exist in world item array.",
		},
		{
			cmnd: "i, inv, inventory",
			desc: "Displays held items.",
		},
		{
			cmnd: "listitems <arg>",
			desc: "Lists first instances w/o arg. Arg can be ID, name, or part of name. Is greedy.",
		},
		{
			cmnd: "drop <item>",
			desc: "Puts an item on the floor.",
		},
		{
			cmnd: "take <item>",
			desc: "Takes an item off the floor.",
		},
		{
			cmnd: "exa, examine <object>",
			desc: "prioritizes players, inventory, ground, then EQ.",
		},
		{
			cmnd: "emotes",
			desc: "lists available emotes",
		},
		{
			cmnd: "snatch <item id> <instance #>",
			desc: "gives you instance of item no matter where its at.",
		},
	}
}

// todo load data from disk
func (w *World) loadRooms() {
	w.rooms = []*Room{
		{
			name:  "The Entryway",
			desc:  "The entryway of the farmhouse is dark and musty, with cobwebs hanging from the ceiling and a thick layer of dust covering the floor. A creaky old staircase leads up to the second floor.",
			id:    1,
			items: []*Item{},
			exits: []*Exit{
				{
					keyword:  "east",
					lookMsg:  "The kitchen lies in that direction.",
					linkedID: 2,
				},
				{
					keyword:  "west",
					lookMsg:  "You see a garden, orchard, and meadow outside of the house.",
					linkedID: 8,
				},
			},
		},
		{
			name:  "The Kitchen",
			desc:  "The kitchen is a cluttered and cramped space, with pots and pans hanging from the ceiling and shelves lined with dusty old jars. A rickety old table sits in the center of the room, with a few broken chairs scattered around it.",
			id:    2,
			items: []*Item{},
			exits: []*Exit{
				{
					keyword:  "west",
					lookMsg:  "You see the entryway to the house in that direction.",
					linkedID: 1,
				},
				{
					keyword:  "north",
					lookMsg:  "An inviting room where one can relax lie that way.",
					linkedID: 3,
				},
				{
					keyword:  "south",
					lookMsg:  "You see a large table surrounded by chairs.",
					linkedID: 4,
				},
			},
		},
		{
			name:  "The Living Room",
			desc:  "The living room is a cozy space with a fireplace, a couple of sofas, and a coffee table. A bookcase stands in one corner, filled with dusty old volumes. The room is musty and smells of old books and wood smoke.",
			id:    3,
			items: []*Item{},
			exits: []*Exit{
				{
					keyword:  "south",
					lookMsg:  "The kitchen lies in that direction.",
					linkedID: 2,
				},
			},
		},
		{
			name:  "The Dining Room",
			desc:  "The dining room is a large, formal space with a long wooden table and matching chairs. A chandelier hangs from the ceiling, casting a dim light throughout the room. A musty old rug covers the floor, and a grandfather clock stands in the corner, ticking away the hours.",
			id:    4,
			items: []*Item{},
			exits: []*Exit{
				{
					keyword:  "north",
					lookMsg:  "The kitchen lies in that direction.",
					linkedID: 2,
				},
				{
					keyword:  "down",
					lookMsg:  "You could probably crawl under the table if you don't mind getting dirty.",
					linkedID: 5,
				},
			},
		},
		{
			name:  "Under The Table",
			desc:  "You get down on all fours, desperately looking for... looking for... you can't remember. Well, maybe if you stand up, you'll remember.",
			id:    5,
			items: []*Item{},
			exits: []*Exit{
				{
					keyword:  "up",
					lookMsg:  "The dining room from an adult perspective awaits!",
					linkedID: 4,
				},
				{
					keyword:  "down",
					lookMsg:  "You see something reflective in a large circular room, like the surface of water, below.",
					linkedID: 6,
				},
			},
		},
		{
			name:  "Before A Dimensional Portal",
			desc:  "You stand in a vast, circular chamber filled with swirling energy. The floor beneath your feet is made of smooth, polished stone, and the walls are adorned with intricate carvings and glowing symbols. In the center of the room stands a massive, shimmering portal, pulsing with otherworldly energy. The portal seems to be a gateway to another realm, filled with strange, shifting colors and patterns. As you approach, you can feel the power of the portal pulling you in, beckoning you to step through and explore the unknown dimensions that lie beyond.",
			id:    6,
			items: []*Item{},
			exits: []*Exit{
				{
					keyword:  "through",
					lookMsg:  "You see what looks to be a plaza with roads going in the cardinal directions away from it.",
					linkedID: 7,
				},
				{
					keyword:  "up",
					lookMsg:  "Back to the earthquake shelter you go!",
					linkedID: 5,
				},
			},
		},
		{
			name:  color("white", "Telnet connecting to 'isharmud.com:23' ...") + "\r\n" + color("blue", "Central Plaza"),
			desc:  "You stand in the center of a spacious plaza, its periphery adorned with potted plants and carved stone benches.  People stroll about you, clad in bright silks and chatting amongst themselves.  A bronze seal at your feet declares you to be in Mareldja, Crown on the Water.  A breeze tinged with salt and brine blows eastward, and shorebirds wheel and dive gracefully overhead. Four wide streets lead from the plaza at each of the compass points.",
			id:    7,
			items: []*Item{},
			exits: []*Exit{
				{
					keyword:  "through",
					lookMsg:  "Back through the closet into 'Spare Hroom.'",
					linkedID: 6,
				},
			},
		},
		{
			name:  "Before A Farmhouse",
			desc:  "At the end of the path, you finally reach the farmhouse. It's a quaint, two-story building with a thatched roof and a large front porch.",
			id:    8,
			items: []*Item{},
			exits: []*Exit{
				{
					keyword:  "east",
					lookMsg:  "The homes main method of entry lies in that direction.",
					linkedID: 1,
				},
				{
					keyword:  "west",
					lookMsg:  "A garden appears to be that way",
					linkedID: 9,
				},
			},
		},
		{
			name:  "The Vegetable Garden",
			desc:  "Next to the orchard is a well-tended vegetable garden, filled with rows of lettuce, tomatoes, beans, and other fresh produce. The scent of herbs and vegetables fills the air.",
			id:    9,
			items: []*Item{},
			exits: []*Exit{
				{
					keyword:  "east",
					lookMsg:  "The path comes to a halt before a dwelling.",
					linkedID: 8,
				},
				{
					keyword:  "west",
					lookMsg:  "Rows and rows of trees...",
					linkedID: 10,
				},
			},
		},
		{
			name:  "The Orchard",
			desc:  "As you continue up the path, you come upon an orchard filled with rows of fruit trees. The branches are heavy with ripe apples, pears, and cherries, and the ground is littered with fallen fruit.",
			id:    10,
			items: []*Item{},
			exits: []*Exit{
				{
					keyword:  "east",
					lookMsg:  "A garden appears to be that way",
					linkedID: 9,
				},
				{
					keyword:  "west",
					lookMsg:  "The trees end and an grassy expanse begins.",
					linkedID: 11,
				},
			},
		},
		{
			name:  "The Meadow",
			desc:  "The forest path opens up into a wide meadow, filled with tall grasses and wildflowers. The sun is warm on your skin, and the breeze carries the scent of freshly cut hay. In the distance, you can see the farmhouse nestled among the fields.",
			id:    11,
			items: []*Item{},
			exits: []*Exit{
				{
					keyword:  "east",
					lookMsg:  "Rows and rows of trees...",
					linkedID: 10,
				},
				{
					keyword:  "west",
					lookMsg:  "A path decends through a natural archway of tree branches.",
					linkedID: 12,
				},
			},
		},
		{
			name:  "The Forest Path",
			desc:  "This winding path is surrounded by tall trees, their branches forming a canopy overhead. The ground is soft and spongy beneath your feet, covered in a thick layer of fallen leaves and pine needles. The air is cool and fresh, the only sounds coming from the birds singing in the treetops and the occasional rustle of small animals in the underbrush.",
			id:    12,
			items: []*Item{},
			exits: []*Exit{
				{
					keyword:  "east",
					lookMsg:  "The trees end and an grassy expanse begins.",
					linkedID: 11,
				},
			},
		},
	}
}

func (w *World) initEQList() {

	w.eqList = append(w.eqList, headSlot)
	w.eqList = append(w.eqList, faceSlot)
	w.eqList = append(w.eqList, neckSlot)
	w.eqList = append(w.eqList, aboutSlot)
	w.eqList = append(w.eqList, chestSlot)
	w.eqList = append(w.eqList, backSlot)
	w.eqList = append(w.eqList, holdLSlot)
	w.eqList = append(w.eqList, holdRSlot)
	w.eqList = append(w.eqList, waistSlot)
	w.eqList = append(w.eqList, legsSlot)
	w.eqList = append(w.eqList, feetSlot)
	w.eqList = append(w.eqList, armsSlot)
	w.eqList = append(w.eqList, wristLSlot)
	w.eqList = append(w.eqList, wristRSlot)
	w.eqList = append(w.eqList, handsSlot)
	w.eqList = append(w.eqList, fingerLSlot)
	w.eqList = append(w.eqList, fingerRSlot)
}

func (w *World) initItems() {
	addItem(w.items, &Item{id: 1, name: "a leather cap", desc: "It's as plain as it gets, covers the melon, provides minor protection.", slot: headSlot, uID: time.Now().Format(time.RFC3339)})
}

// add items to the item map
func addItem(items map[string]map[int]*Item, item *Item) {
	name := item.name
	if _, ok := items[name]; !ok {
		items[name] = make(map[int]*Item)
	}
	instanceNo := len(items[name])
	items[name][instanceNo] = item
}

// return item instance
func getItem(items map[string]map[int]*Item, name string, instanceNo int) (*Item, error) {
	if _, ok := items[name]; !ok {
		return nil, fmt.Errorf("item with name %s not found", name)
	}
	if _, ok := items[name][instanceNo]; !ok {
		return nil, fmt.Errorf("instance %d of item with name %s not found", instanceNo, name)
	}
	return items[name][instanceNo], nil
}

func emoteHandler(input []string, usr *User, w *World) {
	hasTarget := len(input) > 1
	fail := true

	for _, e := range w.emotes {
		if strings.EqualFold(e.name, input[0]) {
			for _, u := range usr.room.users {
				if hasTarget {
					input[1] = strings.TrimLeft(input[1], " ")
					lenTar := len(input[1])
					tar := &User{}

					for _, tgt := range usr.room.users {
						if len(tgt.name) < lenTar {
							lenTar = len(tgt.name)
						}
						if strings.EqualFold(tgt.name[0:lenTar], input[1][0:lenTar]) {
							tar = tgt
							fail = false
						}
					}
					if !fail {
						if u != usr && u != tar {
							clientOutputChan <- ClientOutput{u, color("cyan", usr.name) + e.tPt + color("cyan", tar.name) + ".", &BroadcastEvent{}, w}
						}
						if u != usr && u == tar {
							clientOutputChan <- ClientOutput{u, color("cyan", usr.name) + e.tar, &BroadcastEvent{}, w}
						}
						if u == usr {
							usr.session.WriteLine(e.fPt + color("cyan", tar.name))
						}
					}
				}
				if !hasTarget {
					fail = false
					if u != usr {
						clientOutputChan <- ClientOutput{u, color("cyan", usr.name) + e.tP, &BroadcastEvent{}, w}
					} else {
						usr.session.WriteLine(e.fP)
					}
				}
			}
			if fail {
				usr.session.WriteLine("Emote failed. Most likely unavailable recipient.")
				return
			}
		}
	}
}

func (r *Room) east(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "east" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Printf("Tried getting east exit in room %s, found none\r\n", r.name)
	return nil
}
func (r *Room) west(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "west" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Printf("Tried getting west exit in room %s, found none\r\n", r.name)
	return nil
}
func (r *Room) north(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "north" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Printf("Tried getting north exit in room %s, found none\r\n", r.name)
	return nil
}
func (r *Room) south(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "south" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Printf("Tried getting south exit in room %s, found none\r\n", r.name)
	return nil
}
func (r *Room) up(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "up" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Printf("Tried getting up exit in room %s, found none\r\n", r.name)
	return nil
}
func (r *Room) down(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "down" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Printf("Tried getting down exit in room %s, found none\r\n", r.name)
	return nil
}
func (r *Room) in(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "in" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Printf("Tried getting in exit in room %s, found none\r\n", r.name)
	return nil
}
func (r *Room) out(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "out" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Printf("Tried getting out exit in room %s, found none\r\n", r.name)
	return nil
}
func (r *Room) through(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "through" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Printf("Tried getting through exit in room %s, found none\r\n", r.name)
	return nil
}

func getRoomByID(id int, w *World) *Room {
	for _, rm := range w.rooms {
		if id == rm.id {
			return rm
		}
	}
	return nil
}

// Builds room output
func (r *Room) sendText(u *User) {
	u.session.WriteLine(color("blue", r.name))
	u.session.WriteLine(color("blue", "   "+r.desc))
	itmMap := returnItemCountMap(r.items)
	for itm, cnt := range itmMap {
		if cnt > 1 {
			u.session.WriteLine(color("cyan", itm) + " (" + color("red", fmt.Sprint(cnt)) + ") is lying here.")
		} else {
			u.session.WriteLine(color("cyan", itm) + " is lying here.")
		}
	}
	for _, user := range r.users {
		if user != u {
			u.session.WriteLine(color("cyan", user.name+" is here."))
		}
	}
}

func returnItemCountMap(items []*Item) map[string]int {
	itemCounts := make(map[string]int)
	for _, item := range items {

		itemCounts[item.name]++

	}
	return itemCounts
}

func removeUserFromRoom(u *User, r *Room, w *World) {
	for _, rm := range w.rooms {
		if r == rm {
			for n, usr := range r.users {
				if usr == u {
					rm.removeUser(n)
					fmt.Printf("%s, in room %s, removed from index #%s\r\n", u.name, r.name, fmt.Sprint(n))
					return
				}
			}
		}
	}
	fmt.Printf("Unable to remove %s from %s room index\r\n", u.name, r.name)
}

func (r *Room) removeUser(i int) {
	r.users[i] = r.users[len(r.users)-1]
	r.users = r.users[:len(r.users)-1]
}

func (r *Room) addUser(u *User) {
	r.users = append(r.users, u)
}

func isMoveValid(u *User, dir string, w *World) {
	for _, exit := range u.room.exits {
		if exit.keyword == dir {
			moveUser(u, u.room, getRoomByID(exit.linkedID, w), dir, w)
			return
		}
	}
	switch dir {
	case "north", "south", "east", "west":
		u.session.WriteLine(color("magenta", "You slam your face into an invisble wall. Ouch!"))
	case "up":
		u.session.WriteLine(color("magenta", "What, is there a staircase thats invisible here?"))
	case "down":
		u.session.WriteLine(color("magenta", "In what ground orifice do you plan to stuff your body?"))
	case "in", "out":
		u.session.WriteLine(color("magenta", "You can't do that Jim."))
	case "through":
		u.session.WriteLine(color("magenta", "You successfully move through the air, or was that not your goal?"))
	}
	for _, usr := range u.room.users {
		if usr != u {
			switch dir {
			case "north", "south", "east", "west":
				clientOutputChan <- ClientOutput{usr, color("green", u.name+" slams their face into an invisible wall to the "+dir+"."), &BroadcastEvent{}, w}
			case "up":
				clientOutputChan <- ClientOutput{usr, color("green", u.name+" climbs an invisible staircase and falls flat on their face."), &BroadcastEvent{}, w}
			case "down":
				clientOutputChan <- ClientOutput{usr, color("green", u.name+" decends an imaginary staircase. Are we miming?"), &BroadcastEvent{}, w}
			case "in", "out":
				clientOutputChan <- ClientOutput{usr, color("green", u.name+" makes motions as if they're trying to crawl in or out of something..."), &BroadcastEvent{}, w}
			case "through":
				clientOutputChan <- ClientOutput{usr, color("green", u.name+" successfully penetrates the air. You clap."), &BroadcastEvent{}, w}
			}
		}
	}
}

func moveUser(u *User, from *Room, to *Room, dir string, w *World) {
	for n, user := range from.users {
		if u == user {
			from.removeUser(n)
			fmt.Printf("%s, in room %s, removed from index #%s\r\n", user.name, from.name, fmt.Sprint(n))

			for _, usr := range to.users {
				clientOutputChan <- ClientOutput{usr, color("green", u.name+" arrives from the "+getOppDir(dir)+"."), &BroadcastEvent{}, w}
			}
			to.addUser(u)
			u.room = to
			to.sendText(u)

		} else {
			clientOutputChan <- ClientOutput{user, color("green", u.name+" heads "+dir+"."), &BroadcastEvent{}, w}
		}

	}
}

func (u *User) getPrompt(r *Room) string {
	exits := ""
	for _, e := range r.exits {
		if exits == "" {
			exits = strings.ToUpper(e.keyword[0:1])
		} else {
			exits = exits + strings.ToUpper(e.keyword[0:1])
		}
	}
	return "Exits: " + exits
}

func getOppDir(dir string) string {
	opp := ""
	switch dir {
	case "north":
		opp = "south"
	case "south":
		opp = "north"
	case "east":
		opp = "west"
	case "west":
		opp = "east"
	case "up":
		opp = "area below"
	case "down":
		opp = "area above"
	case "in":
		opp = "outside"
	case "out":
		opp = "inside"
	case "through":
		opp = "other side"
	}
	return opp
}

func color(c string, text string) string {
	clr := ""

	switch c {
	case "black":
		clr = "30"
	case "red":
		clr = "31"
	case "green":
		clr = "32"
	case "yellow":
		clr = "33"
	case "blue":
		clr = "34"
	case "magenta":
		clr = "35"
	case "cyan":
		clr = "36"
	case "white":
		clr = "37"
	case "none":
		clr = "0"
	default:
		clr = "0"
		fmt.Println("Color specified for text not recognized on: " + text)
	}
	val := "\u001b[" + clr + "m" + text + "\u001b[0m"
	return val
}

func (s *Session) WriteLine(str string) error {
	_, err := s.conn.Write([]byte(str + "\r\n"))
	return err
}

func getNameFromConn(conn net.Conn) string {
	buf := make([]byte, 4096)
	conn.Write([]byte("Name?"))
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("User Disconnected During Name Creation")
		conn.Close()
		return ""
	}
	name := string(buf[0 : n-2])
	return name
}

func executeCmd(cmd string, usr *User, w *World, eventCh chan ClientOutput) {

	args := strings.Split(cmd, " ")
	switch args[0] {
	case "say":
		msg := ""
		for i := 1; i < len(args); i++ {
			msg = msg + " " + args[i]
		}
		if len(usr.room.users) < 2 {
			usr.session.WriteLine(color("magenta", "So uh, you talking to a ghost?"))
		} else {
			for _, user := range usr.room.users {
				if user != usr {
					eventCh <- ClientOutput{user, color("yellow", fmt.Sprintf("%s says, \"%s.\"", usr.name, strings.TrimLeft(msg, " "))), &BroadcastEvent{}, w}
				}
			}
			usr.session.WriteLine(color("yellow", "You say, \""+strings.TrimLeft(msg, " ")+".\""))
		}
	case "yell":
		args := strings.Split(cmd, " ")
		msg := ""
		for i := 1; i < len(args); i++ {
			msg = msg + " " + args[i]
		}
		msg = strings.TrimLeft(msg, " ")
		recips := make([]*User, 0)
		for _, recip := range usr.room.users {
			if recip != usr {
				recips = append(recips, recip)
			}
		}
		rooms := make([]*Room, 0)
		rooms = append(rooms, usr.room)

		//yell distance 1 to initiate
		for _, ext := range usr.room.exits {
			rooms = append(rooms, getRoomByID(ext.linkedID, w))
			for _, oUsr := range getRoomByID(ext.linkedID, w).users {
				if oUsr != usr {
					recips = append(recips, oUsr)
				}
			}
		}

		//i < desired yell distance
		for i := 1; i < 4; i++ {
			for _, rm := range rooms {
				for _, ex := range rm.exits {
					r1 := getRoomByID(ex.linkedID, w)
					if r1 != usr.room {
						test := false
						for _, rmm := range rooms {
							if r1 == rmm {
								test = true
							}
						}
						if !test {
							rooms = append(rooms, r1)
							recips = append(recips, r1.users...)
						}
					}
				}
			}
		}
		for _, recip := range recips {
			eventCh <- ClientOutput{recip, color("red", fmt.Sprintf("%s yells, \"%s.\"", usr.name, msg)), &BroadcastEvent{}, w}
		}
		usr.session.WriteLine(color("red", fmt.Sprintf("You yell, \"%s.\"", msg)))
	case "shout":
		args := strings.Split(cmd, " ")
		msg := ""
		for i := 1; i < len(args); i++ {
			msg = msg + " " + args[i]
		}
		msg = strings.TrimLeft(msg, " ")
		for _, recip := range w.users {
			if recip != usr {
				eventCh <- ClientOutput{recip, color("blue", fmt.Sprintf("%s shouts, \"%s.\"", usr.name, msg)), &BroadcastEvent{}, w}
			}
		}
		usr.session.WriteLine(color("blue", fmt.Sprintf("You shout, \"%s.\"", msg)))
	case "l", "look":
		if len(args) < 2 {
			usr.room.sendText(usr)
		} else {
			switch args[1] {
			case "north", "south", "east", "west", "up", "down", "in", "out", "through", "n", "s", "e", "w", "u", "d", "i", "o", "t":
				if len(strings.TrimLeft(args[1], " ")) == 1 {
					switch args[1] {
					case "n":
						args[1] = "north"
					case "s":
						args[1] = "south"
					case "e":
						args[1] = "east"
					case "w":
						args[1] = "west"
					case "u":
						args[1] = "up"
					case "d":
						args[1] = "down"
					case "i":
						args[1] = "in"
					case "o":
						args[1] = "out"
					case "t":
						args[1] = "through"
					}
				}
				for _, ext := range usr.room.exits {
					if ext.keyword == args[1] {
						usr.session.WriteLine(color("white", ext.lookMsg))
						return
					}
				}
				usr.session.WriteLine(color("magenta", "Not much to see."))
				return
			case "":
				usr.session.WriteLine(color("magenta", "What were you trying to look at?"))
				return
			default:
				usr.session.WriteLine(color("magenta", "Not much to see."))
				return
			}
		}
	case "help":
		for _, cmd := range w.cmnds {
			usr.session.WriteLine(color("red", cmd.cmnd) + " - " + cmd.desc)
		}
		return
	case "north", "south", "east", "west", "up", "down", "n", "s", "e", "w", "u", "d":
		if len(args[0]) == 1 {
			switch args[0] {
			case "n":
				args[0] = "north"
			case "s":
				args[0] = "south"
			case "e":
				args[0] = "east"
			case "w":
				args[0] = "west"
			case "u":
				args[0] = "up"
			case "d":
				args[0] = "down"

			}
		}
		isMoveValid(usr, args[0], w)
		return
	case "go":
		if len(args) < 2 {
			usr.session.WriteLine(color("magenta", "Where do you want to go?"))
		} else {
			switch args[1] {
			case "in", "out", "through", "i", "o", "t":
				if len(strings.TrimLeft(args[1], " ")) == 1 {
					switch args[1] {
					case "i":
						args[1] = "in"
					case "o":
						args[1] = "out"
					case "t":
						args[1] = "through"
					}
				}
				isMoveValid(usr, args[1], w)
			default:
				usr.session.WriteLine(color("magenta", fmt.Sprintf("You can't go '%s.'", args[1])))
			}
		}
	case "":
	case "eq", "equip":
		usr.session.WriteLine("You are wearing:")
		if len(usr.char.eq) == 0 {
			usr.session.WriteLine("    " + color("cyan", "nothing!"))
			return
		}
		for _, s := range w.eqList {
			if i := usr.char.eq[s]; i != nil {
				lenName := len(i.slot)
				adjSlot := i.slot
				// makes output :'s line up pretty
				for j := lenName; j < 12; j++ {
					adjSlot = " " + adjSlot
				}
				usr.session.WriteLine(fmt.Sprintf(color("cyan", "    %s: %s"), adjSlot, i.name))
			}
		}
	case "create":
		if len(args) >= 2 {
			str := ""
			for i := 1; i < len(args); i++ {
				str = str + " " + args[i]
			}
			flds := strings.Split(str, "|")
			if len(flds) != 3 {
				usr.session.WriteLine("Not enough arguments to make an item.")
				return
			}
			name := strings.TrimLeft(flds[0], " ")
			desc := flds[1]
			slot := flds[2]
			//slot, err := strconv.Atoi(flds[2])
			//if err == nil {
			fail := false
			id := 0
			for _, i := range w.items {
				if i[0].id > id {
					id = i[0].id
				}
			}
			i := &Item{
				id:   id + 1,
				name: name,
				desc: desc,
				slot: slot,
				uID:  time.Now().Format(time.RFC3339),
			}
			for key := range w.items {
				if strings.EqualFold(key, i.name) {
					fail = true
					usr.session.WriteLine("Item name exists in world.")

				}
			}
			if !fail {
				addItem(w.items, i)
				usr.session.WriteLine("Added: " + i.name + " - " + i.desc + " - " + fmt.Sprint(i.slot))
			}

			//} else {
			//	fmt.Println(usr.name + " failed item creation.")
			//	usr.session.WriteLine("Failed.")
			//}
		} else {
			usr.session.WriteLine("Not enough arguments to make an item.")
		}
	case "listitems":
		//has an argument to search a specific item
		if len(args) > 1 {
			//check if first arg is an int or not
			eyeD, err := strconv.Atoi(args[1])
			//not an int search by map key (item name)
			if err != nil {
				for i := 2; i < len(args); i++ {
					args[1] = args[1] + " " + args[i]
				}
				args[1] = strings.TrimLeft(args[1], " ")
				if _, ok := w.items[args[1]]; !ok {
					fail := true
					//argument supplied is not a map key value, lets try searching map key arguments
					for s, m := range w.items {
						ss := strings.Split(s, " ")
						for sss := 0; sss < len(ss); sss++ {
							if strings.EqualFold(ss[sss], args[1]) {
								fail = false
								for i := 0; i < len(m); i++ {
									if m[i].loc != nil {
										usr.session.WriteLine(fmt.Sprintf("ID: %s, Key: %s, Location: %s, Instance: %s, Address: %p", fmt.Sprint(m[i].id), s, m[i].loc.getName(), fmt.Sprint(i), m[i]))
									} else {
										usr.session.WriteLine(fmt.Sprintf("ID: %s, Key: %s, Location: nil, Instance: %s, Address: %p", fmt.Sprint(m[i].id), s, fmt.Sprint(i), m[i]))
									}
								}
							}
						}
					}
					if fail {
						usr.session.WriteLine(fmt.Sprintf("'%s' is not a valid item name or part of an item name.", args[1]))
					}
				} else {
					//argument is a map key
					for n, i := range w.items[args[1]] {
						if i.loc != nil {
							usr.session.WriteLine(fmt.Sprintf("ID: %s, Key: %s, Location: %s, Instance: %s, Address: %p", fmt.Sprint(i.id), args[1], i.loc.getName(), fmt.Sprint(n), i))
						} else {
							usr.session.WriteLine(fmt.Sprintf("ID: %s, Key: %s, Location: nil, Instance: %s, Address: %p", fmt.Sprint(i.id), args[1], fmt.Sprint(n), i))
						}
					}
				}
			} else {
				//id search
				for s, m := range w.items {
					if m[0].id == eyeD {
						for i := 0; i < len(m); i++ {
							if m[i].loc != nil {
								usr.session.WriteLine(fmt.Sprintf("ID: %s, Key: %s, Location: %s, Instance: %s, Address: %p", fmt.Sprint(m[i].id), s, m[i].loc.getName(), fmt.Sprint(i), m[i]))
							} else {
								usr.session.WriteLine(fmt.Sprintf("ID: %s, Key: %s, Location: nil, Instance: %s, Address: %p", fmt.Sprint(m[i].id), s, fmt.Sprint(i), m[i]))
							}
						}
					}
				}
			}
		} else {
			// no argument display all first instances of items
			for s, m := range w.items {
				if m[0].loc != nil {
					usr.session.WriteLine(fmt.Sprintf("ID: %s, Key: %s, 1st Loc: %s, Instances: %s, Address: %p", fmt.Sprint(m[0].id), s, m[0].loc.getName(), fmt.Sprint(len(m)), m[0]))
				} else {
					usr.session.WriteLine(fmt.Sprintf("ID: %s, Key: %s, 1st Loc: Nil, Instances: %s, Address: %p", fmt.Sprint(m[0].id), s, fmt.Sprint(len(m)), m[0]))
				}
			}
		}
	case "new":
		newStr := ""
		for i := 1; i < len(args); i++ {
			newStr = newStr + " " + args[i]
		}
		if len(args) > 1 {
			fail := true
			args[1] = strings.TrimLeft(args[1], " ")
			n, err := strconv.Atoi(args[1])
			if err != nil {
				lenStr := len(args[1])
				for s, m := range w.items {
					adjStr := lenStr
					if len(s) < adjStr {
						adjStr = len(s)
					}
					if strings.EqualFold(s[0:adjStr], newStr[0:adjStr]) {
						i := &Item{id: m[0].id, name: m[0].name, desc: m[0].desc, slot: m[0].slot, loc: usr.getLocation(), uID: fmt.Sprint(m[0].id) + "|" + time.Now().Format(time.RFC3339)}
						usr.char.inv = append(usr.char.inv, i)
						addItem(w.items, i)
						usr.session.WriteLine(fmt.Sprintf("Arg: '%s' yielded Item: '%s'", args[1], i.name))
						fail = false
						return
					}

				}
			} else {
				for s, m := range w.items {
					if m[0].id == n {
						fail = false
						item := &Item{id: m[0].id, name: m[0].name, desc: m[0].desc, slot: m[0].slot, loc: usr.getLocation(), uID: fmt.Sprint(m[0].id) + "|" + time.Now().Format(time.RFC3339)}
						usr.char.inv = append(usr.char.inv, item)
						addItem(w.items, item)
						usr.session.WriteLine(fmt.Sprintf("Arg: '%s' yielded Item: '%s' - uID: %s", fmt.Sprint(n), s, item.uID))
						return
					}
				}
			}
			if fail {
				usr.session.WriteLine("Did not find item: " + args[1])
			}
		} else {
			usr.session.WriteLine("No item specified.")
		}
	case "i", "inv", "inventory":
		usr.session.WriteLine("You are carrying...")
		if len(usr.char.inv) == 0 {
			usr.session.WriteLine(color("cyan", "    nothing!"))
		} else {
			iMap := returnItemCountMap(usr.char.inv)
			for i, cnt := range iMap {
				if cnt > 1 {
					usr.session.WriteLine(color("cyan", "    "+i) + " (" + color("red", fmt.Sprint(cnt)) + ")")
				} else {
					usr.session.WriteLine(color("cyan", "    "+i))
				}

			}
		}
	case "wear", "waer":
		wearStr := ""
		for i := 1; i < len(args); i++ {
			wearStr = wearStr + " " + args[i]
		}
		wearStr = strings.TrimLeft(wearStr, " ")
		if len(args) > 1 {
			fail := true
			lenStr := len(wearStr)
			for _, i := range usr.char.inv {
				adjStr := lenStr
				if len(i.name) < lenStr {
					adjStr = len(i.name)
				}
				if strings.EqualFold(i.name[0:adjStr], wearStr[0:adjStr]) {
					//tryToWear(i, usr, w)
					fail = false
					if usr.char.eq[i.slot] != nil {
						usr.session.WriteLine("You already have something equipped on your " + strings.ToLower(i.slot) + ".")
						return
					}
					usr.char.eq[i.slot] = i
					removeItemFromInventory(i, usr)
					usr.session.WriteLine(fmt.Sprintf("You place a %s on your %s.", color("cyan", i.name), strings.ToLower(i.slot)))
					for _, u := range usr.room.users {
						if usr != u {
							clientOutputChan <- ClientOutput{u, usr.name + " places a " + color("cyan", i.name) + " on their " + strings.ToLower(i.slot) + ".", &BroadcastEvent{}, w}
						}
					}

				}
				if !fail {
					return
				}
			}
			if fail {
				usr.session.WriteLine("You are not carrying " + args[1])
			}
		} else {
			usr.session.WriteLine("What are you trying to wear?")
		}
	case "remove", "rem":
		if len(args) > 1 {
			remStr := ""
			for i := 1; i < len(args); i++ {
				remStr = remStr + " " + args[i]
			}
			remStr = strings.TrimLeft(remStr, " ")
			lenStr := len(remStr)
			fail := true
			for _, i := range usr.char.eq {
				adjLen := lenStr
				if len(i.name) < adjLen {
					adjLen = len(i.name)
				}
				//checking for a portion of the literal name from the front back
				if strings.EqualFold(i.name[0:adjLen], remStr[0:adjLen]) {
					fail = false
					delete(usr.char.eq, i.slot)
					usr.char.inv = append(usr.char.inv, i)

					usr.session.WriteLine("You remove a " + color("cyan", i.name) + " from your " + strings.ToLower(i.slot) + ".")
					for _, u := range usr.room.users {
						if usr != u {
							clientOutputChan <- ClientOutput{u, usr.name + " removes a " + color("cyan", i.name) + " from their " + strings.ToLower(i.slot) + ".", &BroadcastEvent{}, w}
						}
					}
					return
				} else {
					//literal check failed, lets try item name arguments
					is := strings.Split(i.name, " ")
					for iss := range is {
						if strings.EqualFold(is[iss], args[1]) {
							fail = false
							delete(usr.char.eq, i.slot)
							usr.char.inv = append(usr.char.inv, i)
							usr.session.WriteLine("You remove a " + color("cyan", i.name) + " from your " + strings.ToLower(i.slot) + ".")
							for _, u := range usr.room.users {
								if usr != u {
									clientOutputChan <- ClientOutput{u, usr.name + " removes a " + color("cyan", i.name) + " from their " + strings.ToLower(i.slot) + ".", &BroadcastEvent{}, w}
								}
							}
							//dont want this to be greedy. rem leather should only remove one leather item
							return
						}
					}
				}
			}
			if fail {
				usr.session.WriteLine(fmt.Sprintf("Can't find an item like '%s' equipped.", remStr))
			}
		} else {
			usr.session.WriteLine("What are you trying to remove?")
		}
	case "slots":
		usr.session.WriteLine(color("magenta", headSlot))
		usr.session.WriteLine(color("magenta", faceSlot))
		usr.session.WriteLine(color("magenta", neckSlot))
		usr.session.WriteLine(color("magenta", aboutSlot))
		usr.session.WriteLine(color("magenta", chestSlot))
		usr.session.WriteLine(color("magenta", backSlot))
		usr.session.WriteLine(color("magenta", holdLSlot))
		usr.session.WriteLine(color("magenta", holdRSlot))
		usr.session.WriteLine(color("magenta", waistSlot))
		usr.session.WriteLine(color("magenta", legsSlot))
		usr.session.WriteLine(color("magenta", feetSlot))
		usr.session.WriteLine(color("magenta", armsSlot))
		usr.session.WriteLine(color("magenta", wristLSlot))
		usr.session.WriteLine(color("magenta", wristRSlot))
		usr.session.WriteLine(color("magenta", handsSlot))
		usr.session.WriteLine(color("magenta", fingerLSlot))
		usr.session.WriteLine(color("magenta", fingerRSlot))
	case "drop":
		if len(args) > 1 {
			dropStr := ""
			fail := true
			for i := 1; i < len(args); i++ {
				dropStr = dropStr + " " + args[i]
			}
			dropStr = strings.TrimLeft(dropStr, " ")
			lenStr := len(dropStr)
			for _, i := range usr.char.inv {
				adjLen := lenStr
				if len(i.name) < lenStr {
					adjLen = len(i.name)
				}
				if strings.EqualFold(i.name[0:adjLen], dropStr[0:adjLen]) {
					fail = false
					usr.room.items = append(usr.room.items, i)
					i.loc = usr.room.getLocation()
					removeItemFromInventory(i, usr)
					for _, u := range usr.room.users {
						if u != usr {
							eventCh <- ClientOutput{u, usr.name + " drops a " + color("cyan", i.name) + " on the ground here.", &BroadcastEvent{}, w}
						} else {
							usr.session.WriteLine("You drop a " + color("cyan", i.name) + " on the ground here.")
						}
					}
				}
				if !fail {
					break
				}
			}
			if fail {
				usr.session.WriteLine(color("magenta", "You are not carrying that. "+dropStr))
			}
		} else {
			usr.session.WriteLine(color("magenta", "What are you trying to drop?"))
		}
	case "take":
		if len(args) > 1 {
			takeStr := ""
			fail := true
			for i := 1; i < len(args); i++ {
				takeStr = takeStr + " " + args[i]
			}
			takeStr = strings.TrimLeft(takeStr, " ")
			lenStr := len(takeStr)

			for _, itm := range usr.room.items {
				adjLen := lenStr
				if len(itm.name) < lenStr {
					adjLen = len(itm.name)
				}
				if strings.EqualFold(itm.name[0:adjLen], takeStr[0:adjLen]) {
					fail = false
					removeItemFromRoom(itm, usr.room)
					usr.char.inv = append(usr.char.inv, itm)
					itm.loc = usr.getLocation()
					for _, u := range usr.room.users {
						if u != usr {
							eventCh <- ClientOutput{u, usr.name + " picks up a " + color("cyan", itm.name) + " off the ground here.", &BroadcastEvent{}, w}
						} else {
							usr.session.WriteLine("You pick up a " + color("cyan", itm.name) + " off the ground here.")
						}
					}
				}
				if !fail {
					break
				}
			}
			if fail {
				usr.session.WriteLine(color("magenta", "You don't see that here. "+takeStr))
			}

		} else {
			usr.session.WriteLine(color("magenta", "What are you trying to take?"))
		}
	case "exa", "examine":
		if len(args) > 1 {
			if args[1] == "" {
				return
			}
			exaStr := ""
			fail := true
			for i := 1; i < len(args); i++ {
				exaStr = exaStr + " " + args[i]
			}
			exaStr = strings.Trim(exaStr, " ")
			exaLen := len(exaStr)

			for _, u := range usr.room.users {
				adjLen := exaLen
				if len(u.name) < adjLen {
					adjLen = len(u.name)
				}
				if strings.EqualFold(u.name[0:adjLen], exaStr[0:adjLen]) {
					fail = false
					itms := u.char.eq
					eventCh <- ClientOutput{u, color("cyan", usr.name) + " looks you over thoroughly.", &BroadcastEvent{}, w}
					for _, nt := range usr.room.users {
						if nt != u && nt != usr {
							eventCh <- ClientOutput{nt, color("cyan", usr.name) + " looks over " + u.name + "'s equipment.", &BroadcastEvent{}, w}
						}
					}
					usr.session.WriteLine(u.name + " is wearing:")
					if len(itms) == 0 {
						usr.session.WriteLine(color("cyan", " ...nothing!"))
					}
					for _, s := range w.eqList {
						if i := itms[s]; i != nil {

							lenName := len(i.slot)
							adjSlot := i.slot
							// makes output :'s line up pretty
							for j := lenName; j < 12; j++ {
								adjSlot = " " + adjSlot
							}
							usr.session.WriteLine(fmt.Sprintf(color("cyan", "    %s: %s"), adjSlot, i.name))
						}
					}
					return
				}
			}
			for _, i := range usr.room.items {
				adjLen := exaLen
				if len(i.name) < adjLen {
					adjLen = len(i.name)
				}
				if strings.EqualFold(i.name[0:adjLen], exaStr[0:adjLen]) {
					fail = false
					usr.session.WriteLine("You take a closer look at a " + color("cyan", i.name) + "...")
					usr.session.WriteLine("    " + i.desc)
					r, room := i.loc.(*Room)
					if room {
						usr.session.WriteLine("    " + fmt.Sprint(i.id) + " - " + i.name + " - " + i.desc + " - " + i.slot + " - " + r.name)
					}
					return
				}
			}
			for _, i := range usr.char.inv {
				adjLen := exaLen
				if len(i.name) < adjLen {
					adjLen = len(i.name)
				}
				if strings.EqualFold(i.name[0:adjLen], exaStr[0:adjLen]) {
					fail = false
					usr.session.WriteLine("You take a closer look at a " + color("cyan", i.name) + "...")
					usr.session.WriteLine("    " + i.desc)
					u, user := i.loc.(*User)
					if user {
						usr.session.WriteLine("    " + fmt.Sprint(i.id) + " - " + i.name + " - " + i.desc + " - " + i.slot + " - " + u.name)
					}
					return
				}
			}
			eq := usr.char.eq
			for _, i := range eq {
				adjLen := exaLen
				if len(i.name) < adjLen {
					adjLen = len(i.name)
				}
				if strings.EqualFold(i.name[0:adjLen], exaStr[0:adjLen]) {
					fail = false
					usr.session.WriteLine("You take a closer look at a " + color("cyan", i.name) + "...")
					usr.session.WriteLine("    " + i.desc)
					return
				}
			}
			if fail {
				usr.session.WriteLine(color("magenta", "You see nothing with that name here."))
			}
		} else {
			usr.session.WriteLine(color("magenta", "What are you trying to examine?"))
		}
	case "emotes":
		output := ""
		for _, e := range w.emotes {
			if output == "" {
				output = e.name
			} else {
				output = output + ", " + e.name
			}
		}
		usr.session.WriteLine("Available Emotes: " + output)
	case "snatch":
		if len(args) < 3 {
			return
		}
		args[1] = strings.TrimLeft(args[1], " ")
		args[2] = strings.TrimLeft(args[2], " ")
		n, err := strconv.Atoi(args[1])
		n2, err2 := strconv.Atoi(args[2])
		for s, m := range w.items {
			if err == nil && err2 == nil {
				if m[0].id == n {
					if _, ok := m[n2]; !ok {
						usr.session.WriteLine(fmt.Sprintf("'%s' is not a valid instance of '%s'", fmt.Sprint(n2), s))
					} else {
						if lc, ok := m[n2].loc.(*Room); ok {
							removeItemFromRoom(m[n2], lc)
							usr.session.WriteLine(fmt.Sprintf("You snatched %s from room: %s", color("cyan", m[n2].name), color("red", lc.name)))
							for _, u := range lc.users {
								clientOutputChan <- ClientOutput{u, fmt.Sprintf("%s whisked %s away from the ground here!", color("red", usr.name), color("cyan", m[n2].name)), &BroadcastEvent{}, w}
							}
							usr.char.inv = append(usr.char.inv, m[n2])
							m[n2].loc = usr.getLocation()
							return
						}
						if lc, ok := m[n2].loc.(*User); ok {
							removeItemFromInventory(m[n2], lc)
							usr.char.inv = append(usr.char.inv, m[n2])
							usr.session.WriteLine(fmt.Sprintf("You stole %s from %s!", color("cyan", m[n2].name), color("red", lc.name)))
							clientOutputChan <- ClientOutput{lc, fmt.Sprintf("%s stole %s from your inventory!", color("red", usr.name), color("cyan", m[n2].name)), &BroadcastEvent{}, w}
							m[n2].loc = usr.getLocation()
							return
						}
					}
				}
			} else {
				usr.session.WriteLine("args 1 and 2 need to be integers")
				return
			}
		}
	case "test":
		i, err := getItem(w.items, "a leather cap", 2)
		if err != nil {
			fmt.Println(err)
		} else {
			usr.session.WriteLine(fmt.Sprintf("ID: %s, Key: %s, Loc: %s, Address: %p", fmt.Sprint(i.id), i.name, i.loc.getName(), i))
		}
	case "nod":
		emoteHandler(args, usr, w)
	case "flail":
		emoteHandler(args, usr, w)
	case "laugh":
		emoteHandler(args, usr, w)
	case "smile":
		emoteHandler(args, usr, w)
	case "bird":
		emoteHandler(args, usr, w)
	case "point":
		emoteHandler(args, usr, w)
	default:
		usr.session.WriteLine(color("magenta", fmt.Sprintf("'%s' is not recognized as a command.", args[0])))
		return
	}
}

func removeItemFromInventory(i *Item, u *User) {

	for n, item := range u.char.inv {
		if item == i {
			u.char.inv[n] = u.char.inv[len(u.char.inv)-1]
			u.char.inv = u.char.inv[:len(u.char.inv)-1]
		}
	}
}

func removeItemFromRoom(i *Item, r *Room) {

	for n, item := range r.items {
		if item == i {
			r.items[n] = r.items[len(r.items)-1]
			r.items = r.items[:len(r.items)-1]
		}
	}
}

func (u *User) initChar() *Character {
	char := &Character{
		name: u.name,
		desc: "ToDo",
		eq:   map[string]*Item{},
		inv:  []*Item{},
	}
	return char
}

func handleConnection(world *World, user *User, session *Session, conn net.Conn, inputChannel chan ClientInput) error {
	user.buf = make([]byte, 4096)
	inputChannel <- ClientInput{
		user,
		&UserJoinedEvent{},
		world,
	}

	user.char = user.initChar()
	for {
		n, err := conn.Read(user.buf)
		if err != nil {
			return err
		}
		if n == 0 {
			log.Println("Zero bytes, closing connection")
			break
		}
		input := string(user.buf[0 : n-2])
		e := ClientInput{user, &InputEvent{input}, world}
		inputChannel <- e
	}
	return nil
}

func startServer(inputChannel chan ClientInput) error {

	log.Println("Starting Server...")
	w := &World{}
	w.loadRooms()
	w.loadHelp()
	w.items = make(map[string]map[int]*Item)
	w.loadEmotes()
	w.initEQList()
	w.initItems()
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}

		go func() {
			session := &Session{conn}
			user := &User{name: getNameFromConn(conn), session: session, room: getRoomByID(1, w)}
			if err := handleConnection(w, user, session, conn, inputChannel); err != nil {
				log.Println("Error handling connection", err)
				inputChannel <- ClientInput{user, &UserLeftEvent{user}, w}
				return
			}
		}()
	}
}

func startInputLoop(clientInputChannel <-chan ClientInput) {
	for input := range clientInputChannel {

		switch event := input.event.(type) {
		case *InputEvent:
			fmt.Printf("%s: \"%s\"\r\n", input.user.name, event.msg)
			executeCmd(event.msg, input.user, input.world, clientOutputChan)

		case *UserJoinedEvent:
			fmt.Println("User Joined:", input.user.name)
			input.world.users = append(input.world.users, input.user)
			input.user.session.WriteLine(color("cyan", fmt.Sprintf("Welcome %s", input.user.name)))
			input.user.room.addUser(input.user)
			input.user.room.sendText(input.user)
			for _, user := range input.world.users {
				if user != input.user {
					clientOutputChan <- ClientOutput{user, color("red", fmt.Sprintf("%s has joined!", input.user.name)), &BroadcastEvent{}, input.world}
				}
			}
		case *UserLeftEvent:
			un := input.user.name
			fmt.Println("User Left:", un)
			for n, user := range input.world.users {
				if user != input.user {
					clientOutputChan <- ClientOutput{user, color("red", fmt.Sprintf("%s has left us!", un)), &BroadcastEvent{}, input.world}
				}
				if user == input.user {
					removeUserFromRoom(user, user.room, input.world)
					fmt.Printf("%s removed from world index # %s\r\n", un, fmt.Sprint(n))
					input.world.users[n] = input.world.users[len(input.world.users)-1]
					input.world.users = input.world.users[:len(input.world.users)-1]
				}
			}
		}
		input.user.session.WriteLine(input.user.getPrompt(input.user.room))
	}
}

func startOutputLoop(clientOutputChannel <-chan ClientOutput) {

	for output := range clientOutputChannel {
		switch output.event.(type) {
		case *BroadcastEvent:
			output.user.session.WriteLine(output.msg)
		}
		output.user.session.WriteLine(output.user.getPrompt(output.user.room))
	}
}

func main() {
	chI := make(chan ClientInput)
	clientOutputChan = make(chan ClientOutput)
	go startInputLoop(chI)
	go startOutputLoop(clientOutputChan)
	err := startServer(chI)
	if err != nil {
		log.Fatal(err)
	}
}
