package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"go4.org/strutil"
)

const (
	headSlot    string = "Head"
	faceSlot    string = "Face"
	neckSlot    string = "Neck"
	aboutSlot   string = "About"
	chestSlot   string = "Chest"
	backSlot    string = "Back"
	holdBSlot   string = "Both Hands"
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

	serverName         string = "Ark's Chatrooms"
	serverPort         int    = 8080
	serverYellDistance int    = 4
)

var InputChannel chan ClientInput
var OutputChan chan ClientOutput
var w *World

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
	gold   int
	fort   int
	ref    int
	wil    int
	att    int
	dam    int
	hp     int
	mana   int
	moves  int
	exp    int
}

type Effects struct {
	str   int
	dex   int
	con   int
	intl  int
	wis   int
	cha   int
	fort  int
	ref   int
	wil   int
	att   int
	dam   int
	hp    int
	mana  int
	moves int
	exp   int
}

type Item struct {
	id   int
	name string
	desc string
	slot string
	loc  Location
	uID  string
	ac   int
	dmg  string
	dmgi int
	eff  *Effects
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
		{
			name: "tip",
			fP:   "You tip your hat. Been watching westerns?",
			fPt:  "You tip your hat to ",
			tar:  " tips their hat to you.",
			tP:   " tips their hat.",
			tPt:  " tips their hat to ",
		},
		{
			name: "grin",
			fP:   "You grin.",
			fPt:  "You grin at ",
			tar:  " grins at you.",
			tP:   " grins.",
			tPt:  " grins at ",
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
			cmnd: "create",
			desc: "Starts an item creation prompt.",
		},
		{
			cmnd: "slots",
			desc: "Shows what equipment belongs to what slot for create",
		},
		{
			cmnd: "new <item id or name>",
			desc: "Tries to give you <item>. Has to exist in world item array.",
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
		{
			cmnd: "give <item> <person>",
			desc: "Tries to give item to person.",
		},
		{
			cmnd: "who",
			desc: "Lists all users online.",
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
	w.eqList = append(w.eqList, holdBSlot)
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
	addItem(w.items, &Item{id: 1, name: "a leather cap", desc: "It's as plain as it gets, covers the melon, provides minor protection.", slot: headSlot, uID: time.Now().Format(time.RFC3339), ac: 2})
	addItem(w.items, &Item{id: 2, name: "a spiked chain flail", desc: "You could do some serious damage with this thing.", slot: holdRSlot, uID: time.Now().Format(time.RFC3339), dmg: "6d3", dmgi: 2})
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
					var tar *User

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
							OutputChan <- ClientOutput{u, color("cyan", usr.name) + e.tPt + color("cyan", tar.name) + ".", &BroadcastEvent{}, w}
						}
						if u != usr && u == tar {
							OutputChan <- ClientOutput{u, color("cyan", usr.name) + e.tar, &BroadcastEvent{}, w}
						}
						if u == usr {
							usr.session.WriteLine(e.fPt + color("cyan", tar.name))
						}
					}
				}
				if !hasTarget {
					fail = false
					if u != usr {
						OutputChan <- ClientOutput{u, color("cyan", usr.name) + e.tP, &BroadcastEvent{}, w}
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
				OutputChan <- ClientOutput{usr, color("green", u.name+" slams their face into an invisible wall to the "+dir+"."), &BroadcastEvent{}, w}
			case "up":
				OutputChan <- ClientOutput{usr, color("green", u.name+" climbs an invisible staircase and falls flat on their face."), &BroadcastEvent{}, w}
			case "down":
				OutputChan <- ClientOutput{usr, color("green", u.name+" decends an imaginary staircase. Are we miming?"), &BroadcastEvent{}, w}
			case "in", "out":
				OutputChan <- ClientOutput{usr, color("green", u.name+" makes motions as if they're trying to crawl in or out of something..."), &BroadcastEvent{}, w}
			case "through":
				OutputChan <- ClientOutput{usr, color("green", u.name+" successfully penetrates the air. You clap."), &BroadcastEvent{}, w}
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
				OutputChan <- ClientOutput{usr, color("green", u.name+" arrives from the "+getOppDir(dir)+"."), &BroadcastEvent{}, w}
			}
			to.addUser(u)
			u.room = to
			u.session.WriteLine("You go " + dir + ".")
			to.sendText(u)

		} else {
			OutputChan <- ClientOutput{user, color("green", u.name+" heads "+dir+"."), &BroadcastEvent{}, w}
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

func getNameFromConn(conn net.Conn) (string, error) {
	buf := make([]byte, 4096)
	name := ""
	conn.Write([]byte(fmt.Sprintf("Welcome to %s\r\n", serverName)))
	for len(name) < 3 || len(name) > 15 {
		conn.Write([]byte("What are you called?"))
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("User Disconnected During Name Creation")
			conn.Close()
			return "", err
		}
		name = string(buf[0 : n-2])
		if len(name) < 3 || len(name) > 15 {
			conn.Write([]byte("Names need to be 3 - 15 characters\r\n"))
		}
	}
	return strings.ToUpper(name[:1]) + name[1:], nil
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
					eventCh <- ClientOutput{user, fmt.Sprintf("%s says, \"%s"+color("yellow", ".")+"\"", color("cyan", usr.name), color("yellow", strings.TrimLeft(msg, " "))), &BroadcastEvent{}, w}
				}
			}
			usr.session.WriteLine("You say, \"" + color("yellow", strings.TrimLeft(msg, " ")+".") + "\"")
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
		for i := 1; i < serverYellDistance; i++ {
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
			eventCh <- ClientOutput{recip, fmt.Sprintf("%s yells, \"%s.\"", color("cyan", usr.name), color("red", msg)), &BroadcastEvent{}, w}
		}
		usr.session.WriteLine(fmt.Sprintf("You yell, \"%s.\"", color("red", msg)))
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
	case "l", "look", "exa", "examine":
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
				//examining a char probably
				for _, u := range usr.room.users {
					if strutil.ContainsFold(u.name, args[1]) && u != usr {
						exaCharacter(usr, u)
						return
					}
					if strutil.ContainsFold(u.name, args[1]) && u == usr {
						usr.session.WriteLine(color("magenta", "I recommend just typing 'eq' or looking in a mirror."))
						return
					}
				}
				for _, i := range usr.room.items {
					if strutil.ContainsFold(i.name, args[1]) {
						exaItem(usr, i, "room")
						return
					}
				}
				for _, i := range usr.char.inv {
					if strutil.ContainsFold(i.name, args[1]) {
						exaItem(usr, i, "inv")
						return
					}
				}
				for _, i := range usr.char.eq {
					if strutil.ContainsFold(i.name, args[1]) {
						exaItem(usr, i, "eq")
						return
					}
				}
				usr.session.WriteLine(color("magenta", "You see nothing with that name here."))
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
						i := &Item{}
						i.cloneItem(m[0])
						i.loc = usr.getLocation()
						usr.char.inv = append(usr.char.inv, i)
						addItem(w.items, i)
						usr.session.WriteLine(fmt.Sprintf("Arg: '%s' yielded Item: '%s' - uID: %s", fmt.Sprint(n), s, i.uID))
						fail = false
						return
					}
				}
			} else {
				for s, m := range w.items {
					if m[0].id == n {
						fail = false
						item := &Item{}
						item.cloneItem(m[0])
						item.loc = usr.getLocation()
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
		if len(args) > 1 {
			for _, i := range usr.char.inv {
				if strutil.ContainsFold(i.name, args[1]) {
					if usr.char.eq[i.slot] != nil {
						usr.session.WriteLine("You already have something equipped on your " + strings.ToLower(i.slot) + ".")
						return
					}
					if strutil.ContainsFold(i.slot, "hand") {
						if i.slot == holdBSlot {
							if usr.char.eq[holdLSlot] != nil || usr.char.eq[holdRSlot] != nil {
								usr.session.WriteLine("You already have something equipped in your hands.")
								return
							}
						}
						if usr.char.eq[holdBSlot] != nil {
							usr.session.WriteLine("You already have something equipped in your hands.")
							return
						}
					}
					usr.char.eq[i.slot] = i
					usr.char.inv = removeItemFromSlice(i, usr.char.inv)
					if strutil.ContainsFold(i.slot, "hand") {
						if i.slot == holdBSlot {
							usr.session.WriteLine(fmt.Sprintf("You grab hold of %s in %s.", color("cyan", i.name), strings.ToLower(i.slot)))

						} else {
							usr.session.WriteLine(fmt.Sprintf("You grab hold of %s in your %s.", color("cyan", i.name), strings.ToLower(i.slot)))
						}
					} else {
						usr.session.WriteLine(fmt.Sprintf("You place %s on your %s.", color("cyan", i.name), strings.ToLower(i.slot)))
					}
					for _, u := range usr.room.users {
						if usr != u {
							if strutil.ContainsFold(i.slot, "hand") {
								if i.slot == holdBSlot {
									OutputChan <- ClientOutput{u, usr.name + " holds " + color("cyan", i.name) + " in " + strings.ToLower(i.slot) + ".", &BroadcastEvent{}, w}
								} else {
									OutputChan <- ClientOutput{u, usr.name + " holds " + color("cyan", i.name) + " in their " + strings.ToLower(i.slot) + ".", &BroadcastEvent{}, w}
								}
							} else {
								OutputChan <- ClientOutput{u, usr.name + " places " + color("cyan", i.name) + " on their " + strings.ToLower(i.slot) + ".", &BroadcastEvent{}, w}
							}
						}
					}
					return
				}
			}
			usr.session.WriteLine("You are not carrying " + args[1])

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
							OutputChan <- ClientOutput{u, usr.name + " removes a " + color("cyan", i.name) + " from their " + strings.ToLower(i.slot) + ".", &BroadcastEvent{}, w}
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
									OutputChan <- ClientOutput{u, usr.name + " removes a " + color("cyan", i.name) + " from their " + strings.ToLower(i.slot) + ".", &BroadcastEvent{}, w}
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
		for _, s := range w.eqList {
			usr.session.WriteLine(color("magenta", s))
		}
	case "drop":
		if len(args) > 1 {
			if args[1] != "" {
				dropStr := args[1]
				dropStr = strings.TrimSpace(dropStr)
				for _, i := range usr.char.inv {
					if strutil.ContainsFold(i.name, dropStr) {
						dropItem(usr, i)
						return
					}
				}
				usr.session.WriteLine(color("magenta", "You are not carrying that. "+dropStr))
				return
			}
			usr.session.WriteLine(color("magenta", "You must supply an item to drop."))
		} else {
			usr.session.WriteLine(color("magenta", "What are you trying to drop?"))
		}
	case "take", "taek":
		if len(args) > 1 {
			if args[1] != "" {
				takeStr := args[1]
				takeStr = strings.TrimSpace(takeStr)
				for _, itm := range usr.room.items {
					if strutil.ContainsFold(itm.name, takeStr) {
						usr.room.items = takeItem(usr, itm, usr.room.items)
						return
					}
				}
				usr.session.WriteLine(color("magenta", "You don't see that here. "+takeStr))
				return
			}
			usr.session.WriteLine(color("magenta", "You must supply an item to drop."))
		} else {
			usr.session.WriteLine(color("magenta", "What are you trying to take?"))
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
							lc.items = removeItemFromSlice(m[n2], lc.items)
							usr.char.inv = append(usr.char.inv, m[n2])
							m[n2].loc = usr.getLocation()
							usr.session.WriteLine(fmt.Sprintf("You snatched %s from room: %s", color("cyan", m[n2].name), color("red", lc.name)))
							for _, u := range lc.users {
								OutputChan <- ClientOutput{u, fmt.Sprintf("%s whisked %s away from the ground here!", color("red", usr.name), color("cyan", m[n2].name)), &BroadcastEvent{}, w}
							}
							return
						}
						if lc, ok := m[n2].loc.(*User); ok {
							found := false
							for _, i := range lc.char.inv {
								if m[n2] == i {
									lc.char.inv = removeItemFromSlice(m[n2], lc.char.inv)
									found = true
								}
							}
							if found {
								usr.char.inv = append(usr.char.inv, m[n2])
								m[n2].loc = usr.getLocation()
								usr.session.WriteLine(fmt.Sprintf("You stole %s from %s!", color("cyan", m[n2].name), color("red", lc.name)))
								OutputChan <- ClientOutput{lc, fmt.Sprintf("%s stole %s from your inventory!", color("red", usr.name), color("cyan", m[n2].name)), &BroadcastEvent{}, w}
								return
							} else {
								delete(lc.char.eq, m[n2].slot)
								usr.char.inv = append(usr.char.inv, m[n2])
								m[n2].loc = usr.getLocation()
								usr.session.WriteLine(fmt.Sprintf("You stole %s from %s!", color("cyan", m[n2].name), color("red", lc.name)))
								OutputChan <- ClientOutput{lc, fmt.Sprintf("%s stole %s from your inventory!", color("red", usr.name), color("cyan", m[n2].name)), &BroadcastEvent{}, w}
								return
							}
						}
					}
				}
			} else {
				usr.session.WriteLine("args 1 and 2 need to be integers")
				return
			}
		}
		/*
			i, err := getItem(w.items, "a leather cap", 2)
			if err != nil {
				fmt.Println(err)
			} else {
				usr.session.WriteLine(fmt.Sprintf("ID: %s, Key: %s, Loc: %s, Address: %p", fmt.Sprint(i.id), i.name, i.loc.getName(), i))
			}
		*/
	case "test":

	case "give":
		if len(args) == 3 {
			give(usr, args[2], args[1])
			return
		}
		usr.session.WriteLine("Give requires 3 arguments: Give <object> person.")
	case "who":
		usr.session.WriteLine(fmt.Sprintf(color("blue", "%d")+" users are online.", len(w.users)))
		for _, u := range w.users {
			usr.session.WriteLine("    " + color("blue", u.name))
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
	case "tip":
		emoteHandler(args, usr, w)
	case "grin":
		emoteHandler(args, usr, w)
	default:
		usr.session.WriteLine(color("magenta", fmt.Sprintf("'%s' is not recognized as a command.", args[0])))
		return
	}
}

func (i *Item) cloneItem(itemToClone *Item) {
	i.id = itemToClone.id
	i.name = itemToClone.name
	i.desc = itemToClone.desc
	i.slot = itemToClone.slot
	i.uID = fmt.Sprint(itemToClone.id) + "|" + time.Now().Format(time.RFC3339)
	i.ac = itemToClone.ac
	i.dmg = itemToClone.dmg
	i.dmgi = itemToClone.dmgi
	i.eff = itemToClone.eff
}

func (i *Item) isWeapon() bool {
	if i.dmg != "" && i.dmg != "0" || i.dmgi != 0 {
		return true
	} else {
		return false
	}
}

func (i *Item) isArmor() bool {
	if i.ac != 0 {
		return true
	} else {
		return false
	}
}

func (i *Item) rollDamage() int {
	both := strings.Split(i.dmg, "d")
	qtyDice, err := strconv.Atoi(both[0])
	diceSides, err2 := strconv.Atoi(both[1])
	roll := i.dmgi
	if err == nil && err2 == nil {
		for i := 0; i < qtyDice; i++ {
			rand.Seed(time.Now().UnixMicro())
			rand := rand.Intn(diceSides) + 1
			roll += rand
		}
	}
	return roll
}

// removes itemToTake from sliceOfItems, handles output to users, returns the updated sliceOfItems
func takeItem(userTaker *User, itemToTake *Item, sliceOfItems []*Item) []*Item {
	sliceOfItems = removeItemFromSlice(itemToTake, sliceOfItems)
	userTaker.char.inv = append(userTaker.char.inv, itemToTake)
	itemToTake.loc = userTaker.getLocation()
	for _, u := range userTaker.room.users {
		if u != userTaker {
			OutputChan <- ClientOutput{u, fmt.Sprintf("%s picks up a %s off the ground here.", userTaker.name, color("cyan", itemToTake.name)), &BroadcastEvent{}, w}
		} else {
			userTaker.session.WriteLine(fmt.Sprintf("You pick up a %s off the ground here.", color("cyan", itemToTake.name)))
		}
	}
	return sliceOfItems
}

// removes itemToDrop from userDropper and places it in the userDropper's room
func dropItem(userDropper *User, itemToDrop *Item) {
	userDropper.room.items = append(userDropper.room.items, itemToDrop)
	itemToDrop.loc = userDropper.room.getLocation()
	//removeItemFromInventory(i, usr)
	userDropper.char.inv = removeItemFromSlice(itemToDrop, userDropper.char.inv)
	for _, u := range userDropper.room.users {
		if u != userDropper {
			OutputChan <- ClientOutput{u, userDropper.name + " drops a " + color("cyan", itemToDrop.name) + " on the ground here.", &BroadcastEvent{}, w}
		} else {
			userDropper.session.WriteLine("You drop a " + color("cyan", itemToDrop.name) + " on the ground here.")
		}
	}
}

// user examiner examines examinee
func exaCharacter(examiner *User, examinee *User) {
	itms := examinee.char.eq
	OutputChan <- ClientOutput{examinee, color("cyan", examiner.name) + " looks you over thoroughly.", &BroadcastEvent{}, w}
	for _, nt := range examiner.room.users {
		if nt != examinee && nt != examiner {
			OutputChan <- ClientOutput{nt, color("cyan", examiner.name) + " looks over " + examinee.name + "'s equipment.", &BroadcastEvent{}, w}
		}
	}
	examiner.session.WriteLine(examinee.name + " is wearing:")
	if len(itms) != 0 {
		for _, s := range w.eqList {
			if i := itms[s]; i != nil {

				lenName := len(i.slot)
				adjSlot := i.slot
				// makes output :'s line up pretty
				for j := lenName; j < 12; j++ {
					adjSlot = " " + adjSlot
				}
				examiner.session.WriteLine(fmt.Sprintf(color("cyan", "    %s: %s"), adjSlot, i.name))
			}
		}
		return
	}
	examiner.session.WriteLine(color("cyan", " ...nothing!"))
}

// examiner looks at itemExamined, itemlocation changes the first string written out to the client
func exaItem(examiner *User, itemExamined *Item, itemlocation string) {
	exaloc := ""
	switch itemlocation {
	case "room":
		exaloc = "You take a closer look at " + color("cyan", itemExamined.name) + " in the room."
	case "inv":
		exaloc = "You take a closer look at " + color("cyan", itemExamined.name) + " in your inventory."
	case "eq":
		exaloc = "You take a closer look at " + color("cyan", itemExamined.name) + " you have equipped."
	}
	examiner.session.WriteLine(exaloc)
	examiner.session.WriteLine("    " + itemExamined.desc)
	if itemExamined.isWeapon() {
		examiner.session.WriteLine(fmt.Sprintf("    %s is a weapon with a damage-roll of %s+%d held in the %s", itemExamined.name, itemExamined.dmg, itemExamined.dmgi, strings.ToLower(itemExamined.slot)))
	}
	if itemExamined.isArmor() {
		if strutil.ContainsFold(itemExamined.slot, "hand") {
			examiner.session.WriteLine(fmt.Sprintf("    %s is a piece of armor with an AC rating of %d, held in the %s.", itemExamined.name, itemExamined.ac, strings.ToLower(itemExamined.slot)))
			return
		}
		examiner.session.WriteLine(fmt.Sprintf("    %s is a piece of armor with an AC rating of %d, worn on the %s.", itemExamined.name, itemExamined.ac, strings.ToLower(itemExamined.slot)))
	}
}

// tries to give itemGiven to userTo from userFrom. tries to match str arguments to user and item
func give(userFrom *User, userTo string, itemGiven string) {

	var item *Item
	var target *User
	for _, j := range userFrom.char.inv {
		if strutil.ContainsFold(j.name, itemGiven) {
			item = j
		}
	}
	for _, u := range userFrom.room.users {
		if strutil.ContainsFold(u.name, userTo) && u != userFrom {
			target = u
		}
	}
	if item != nil {
		if target != nil {
			userFrom.char.inv = removeItemFromSlice(item, userFrom.char.inv)
			target.char.inv = append(target.char.inv, item)
			item.loc = target.getLocation()
			OutputChan <- ClientOutput{target, fmt.Sprintf("%s gives you %s.", color("cyan", userFrom.name), color("cyan", item.name)), &BroadcastEvent{}, w}
			userFrom.session.WriteLine(fmt.Sprintf("You give %s to %s.", color("cyan", item.name), color("cyan", target.name)))
			for _, u := range userFrom.room.users {
				if u != target && u != userFrom {
					OutputChan <- ClientOutput{u, fmt.Sprintf("%s gives %s to %s.", color("cyan", userFrom.name), color("cyan", item.name), color("cyan", target.name)), &BroadcastEvent{}, w}
				}
			}
			return
		}
		userFrom.session.WriteLine("You don't see that person here.")
		return
	}
	userFrom.session.WriteLine("You don't have that item in your inventory.")
}

func removeItemFromSlice(itemToRemove *Item, sliceOfItems []*Item) []*Item {
	for n, i := range sliceOfItems {
		if i == itemToRemove {
			return append(sliceOfItems[:n], sliceOfItems[n+1:]...)
		}
	}
	fmt.Println(fmt.Printf("func removeItemFromSlice did not remove %s.", itemToRemove.name))
	return sliceOfItems
}

// creation promps for an item. UTO == user thread only
func createItemUTO(u *User) *Item {
	itm := &Item{}
	questions := []string{"Name of Item?(string)", "Item description?(string)", "Equipment slot of item?(string)", "Armor class value of item?(int)", "Damage of item (2d4)?", "Flat Damage(int)?"}
	var answer []string
	for i := 0; i < len(questions); i++ {

		u.session.WriteLine(questions[i])
		n, err := u.session.conn.Read(u.buf)
		if err != nil {
			u.session.WriteLine("Item creation failed")
			fmt.Println(err)
		}
		answer = append(answer, string(u.buf[0:n-2]))
	}
	itm.id = len(w.items) + 1
	itm.name = answer[0]
	itm.desc = answer[1]
	itm.slot = answer[2]
	found := false
	b := itm.slot
	for _, s := range w.eqList {
		if strutil.ContainsFold(s, itm.slot) {
			itm.slot = s
			u.session.WriteLine(fmt.Sprintf("Item slot corrected to %s.", s))
			found = true
			if found {
				break
			}
		}
	}
	if !found {
		u.session.WriteLine(fmt.Sprintf("%s is not a valid item slot, aborting.", b))
		return nil
	}
	tempac, errA := strconv.Atoi(answer[3])
	if errA != nil {
		u.session.WriteLine("Item creation failed on AC.")
		return nil
	}
	itm.ac = tempac
	itm.dmg = answer[4]
	tempdmgi, errD := strconv.Atoi(answer[5])
	if errD != nil {
		u.session.WriteLine("Item creation failed on flat damage.")
		return nil
	}
	itm.dmgi = tempdmgi
	itm.uID = fmt.Sprint(itm.id) + "|" + time.Now().Format(time.RFC3339)
	u.session.WriteLine(fmt.Sprintf("Success. Type 'new %d' to get a copy of newly created item.", itm.id))
	u.session.WriteLine(u.getPrompt(u.room))
	return itm
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
		switch input {
		case "create":
			itm := createItemUTO(user)
			if itm != nil {
				addItem(w.items, itm)
			}
		default:
			inputChannel <- ClientInput{user, &InputEvent{input}, world}
		}
	}
	return nil
}

func startServer(inputChannel chan ClientInput) error {

	log.Println("Starting Server...")
	w = &World{}
	w.loadRooms()
	w.loadHelp()
	w.items = make(map[string]map[int]*Item)
	w.loadEmotes()
	w.initEQList()
	w.initItems()
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		fmt.Printf("Incoming connection from %s\r\n", conn.RemoteAddr())
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}

		go func() {
			session := &Session{conn: conn}
			name, err := getNameFromConn(conn)
			if err != nil {
				log.Println("Error handling connection", err)
				return
			}
			user := &User{name: name, session: session, room: getRoomByID(1, w)}
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
			executeCmd(event.msg, input.user, input.world, OutputChan)

		case *UserJoinedEvent:
			fmt.Println("User Joined:", input.user.name)
			input.world.users = append(input.world.users, input.user)
			input.user.session.WriteLine(fmt.Sprintf("Welcome %s. Type help for a list of commands.", color("cyan", input.user.name)))
			input.user.room.addUser(input.user)
			input.user.room.sendText(input.user)
			for _, user := range input.world.users {
				if user != input.user {
					OutputChan <- ClientOutput{user, color("red", fmt.Sprintf("%s has joined!", input.user.name)), &BroadcastEvent{}, input.world}
				}
			}
		case *UserLeftEvent:
			un := input.user.name
			fmt.Println("User Left:", un)
			for n, user := range input.world.users {
				if user != input.user {
					OutputChan <- ClientOutput{user, color("red", fmt.Sprintf("%s has left us!", un)), &BroadcastEvent{}, input.world}
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
	InputChannel = make(chan ClientInput)
	OutputChan = make(chan ClientOutput)
	go startInputLoop(InputChannel)
	go startOutputLoop(OutputChan)
	err := startServer(InputChannel)
	if err != nil {
		log.Fatal(err)
	}
}
