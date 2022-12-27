package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

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
	eq     *Equipment
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
	name string
	desc string
	slot int
}

type Equipment struct {
	head    *Item //slot 1
	face    *Item //slot 2
	neck    *Item //slot 3
	about   *Item //slot 4
	chest   *Item //slot 5
	back    *Item //slot 6
	holdL   *Item //slot 7
	holdR   *Item //slot 8
	waist   *Item //slot 9
	legs    *Item //slot 10
	feet    *Item //slot 11
	arms    *Item //slot 12
	wristL  *Item //slot 13
	wristR  *Item //slot 14
	hands   *Item //slot 15
	fingerL *Item //slot 16
	fingerR *Item //slot 17
}

type Session struct {
	conn net.Conn
}

type World struct {
	users  []*User
	rooms  []*Room
	cmnds  []*Command
	items  []*Item
	emotes []*Emote
}

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
			desc: "Like say/yell, but heard everywhere.",
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
			cmnd: "give <item name>",
			desc: "Tries to give <item>. Has to exist in world item array.",
		},
		{
			cmnd: "i, inv, inventory",
			desc: "Displays held items.",
		},
		{
			cmnd: "listitems",
			desc: "Lists all items in world item array.",
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
	}
}

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
	}
}

func emoteHandler(input []string, usr *User, w *World) {
	hasTarget := len(input) > 1
	fail := true

	for _, e := range w.emotes {
		if strings.ToUpper(e.name) == strings.ToUpper(input[0]) {
			for _, u := range usr.room.users {
				if hasTarget {
					input[1] = strings.TrimLeft(input[1], " ")
					lenTar := len(input[1])
					tar := &User{}

					for _, tgt := range usr.room.users {
						if len(tgt.name) < lenTar {
							lenTar = len(tgt.name)
						}
						if strings.ToUpper(tgt.name[0:lenTar]) == strings.ToUpper(input[1][0:lenTar]) {
							tar = tgt
							fail = false
						}
					}
					if fail == false {
						if u != usr && u != tar {
							u.session.WriteLine(color("cyan", usr.name) + e.tPt + color("cyan", tar.name) + ".")
						}
						if u != usr && u == tar {
							u.session.WriteLine(color("cyan", usr.name) + e.tar)
						}
						if u == usr {
							usr.session.WriteLine(e.fPt + color("cyan", tar.name))
						}
					}
				}
				if !hasTarget {
					fail = false
					if u != usr {
						u.session.WriteLine(color("cyan", usr.name) + e.tP)
					} else {
						usr.session.WriteLine(e.fP)
					}
				}
			}
			if fail == true {
				usr.session.WriteLine("Emote failed. Most likely unavailable recipient.")
				return
			}
		}
	}
}

func (u *User) returnEQ() {

	equip := make([]*Item, 0)

	if eq := u.char.eq.head; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.face; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.neck; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.about; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.chest; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.back; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.holdL; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.holdR; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.waist; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.legs; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.feet; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.arms; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.wristL; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.wristR; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.hands; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.fingerL; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.fingerR; eq != nil {
		equip = append(equip, eq)
	}
	u.session.WriteLine("You are wearing:")
	if len(equip) == 0 {
		u.session.WriteLine(color("cyan", " ...Nothing!"))
	} else {
		for _, itm := range equip {
			u.session.WriteLine(color("cyan", itm.slotToString()+itm.name))
		}
	}
	u.session.WriteLine(u.getPrompt(u.room))
}

func (r *Room) east(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "east" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Println(fmt.Sprintf("Tried getting east exit in room %s, found none", r.name))
	return nil
}
func (r *Room) west(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "west" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Println(fmt.Sprintf("Tried getting west exit in room %s, found none", r.name))
	return nil
}
func (r *Room) north(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "north" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Println(fmt.Sprintf("Tried getting north exit in room %s, found none", r.name))
	return nil
}
func (r *Room) south(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "south" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Println(fmt.Sprintf("Tried getting south exit in room %s, found none", r.name))
	return nil
}
func (r *Room) up(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "up" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Println(fmt.Sprintf("Tried getting up exit in room %s, found none", r.name))
	return nil
}
func (r *Room) down(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "down" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Println(fmt.Sprintf("Tried getting down exit in room %s, found none", r.name))
	return nil
}
func (r *Room) in(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "in" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Println(fmt.Sprintf("Tried getting in exit in room %s, found none", r.name))
	return nil
}
func (r *Room) out(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "out" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Println(fmt.Sprintf("Tried getting out exit in room %s, found none", r.name))
	return nil
}
func (r *Room) through(w *World) *Room {

	for _, ex := range r.exits {
		if ex.keyword == "through" {
			return getRoomByID(ex.linkedID, w)
		}
	}
	fmt.Println(fmt.Sprintf("Tried getting through exit in room %s, found none", r.name))
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

func (r *Room) sendText(u *User) {
	u.session.WriteLine(color("blue", r.name))
	u.session.WriteLine(color("blue", "   "+r.desc))
	for _, itm := range r.items {
		u.session.WriteLine("A " + color("cyan", itm.name) + " is lying here.")
	}
	for _, user := range r.users {
		if user != u {
			u.session.WriteLine(color("cyan", user.name+" is here."))
		}
	}
	u.session.WriteLine(u.getPrompt(r))
}

func removeUserFromRoom(u *User, r *Room, w *World) {
	for _, rm := range w.rooms {
		if r == rm {
			for n, usr := range r.users {
				if usr == u {

					rm.removeUser(n)
					fmt.Println(fmt.Sprintf("%s, in room %s, removed from index #%s", u.name, r.name, fmt.Sprint(n)))
					return
				}
			}
		}
	}
	fmt.Println(fmt.Sprintf("Unable to remove %s from %s room index", u.name, r.name))
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
			moveUser(u, u.room, getRoomByID(exit.linkedID, w), dir)
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
	u.session.WriteLine(u.getPrompt(u.room))
	for _, usr := range u.room.users {
		if usr != u {
			switch dir {
			case "north", "south", "east", "west":
				usr.session.WriteLine(color("green", u.name+" slams their face into an invisible wall to the "+dir+"."))
				usr.session.WriteLine(usr.getPrompt(usr.room))
			case "up":
				usr.session.WriteLine(color("green", u.name+" climbs an invisible staircase and falls flat on their face."))
				usr.session.WriteLine(usr.getPrompt(usr.room))
			case "down":
				usr.session.WriteLine(color("green", u.name+" decends an imaginary staircase. Are we miming?"))
				usr.session.WriteLine(usr.getPrompt(usr.room))
			case "in", "out":
				usr.session.WriteLine(color("green", u.name+" makes motions as if they're trying to crawl in or out of something..."))
				usr.session.WriteLine(usr.getPrompt(usr.room))
			case "through":
				usr.session.WriteLine(color("green", u.name+" successfully penetrates the air. You clap."))
				usr.session.WriteLine(usr.getPrompt(usr.room))
			}
		}
	}
}

func moveUser(u *User, from *Room, to *Room, dir string) {
	for n, user := range from.users {
		if u == user {
			from.removeUser(n)
			fmt.Println(fmt.Sprintf("%s, in room %s, removed from index #%s", user.name, from.name, fmt.Sprint(n)))

			for _, usr := range to.users {

				usr.session.WriteLine(color("green", u.name+" arrives from the "+getOppDir(dir)+"."))
				usr.session.WriteLine(usr.getPrompt(usr.room))
			}
			to.addUser(u)
			u.room = to
			to.sendText(u)

		} else {

			user.session.WriteLine(color("green", u.name+" heads "+dir+"."))
			user.session.WriteLine(user.getPrompt(user.room))
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

func executeCmd(cmd string, usr *User, w *World) {

	args := strings.Split(cmd, " ")
	switch args[0] {
	case "say":
		msg := ""
		for i := 1; i < len(args); i++ {
			msg = msg + " " + args[i]
		}
		if len(usr.room.users) < 2 {
			usr.session.WriteLine(color("magenta", "So uh, you talking to a ghost?"))
			usr.session.WriteLine(usr.getPrompt(usr.room))
		} else {
			for _, user := range usr.room.users {
				if user != usr {

					user.session.WriteLine(color("yellow", fmt.Sprintf("%s says, \"%s.\"", usr.name, strings.TrimLeft(msg, " "))))
					user.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
			usr.session.WriteLine(color("yellow", "You say, \""+strings.TrimLeft(msg, " ")+".\""))
			usr.session.WriteLine(usr.getPrompt(usr.room))
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
						if test == false {
							rooms = append(rooms, r1)
							for _, u1 := range r1.users {
								recips = append(recips, u1)
							}
						}
					}
				}
			}
		}
		for _, recip := range recips {
			recip.session.WriteLine(color("red", fmt.Sprintf("%s yells, \"%s.\"", usr.name, msg)))
			recip.session.WriteLine(recip.getPrompt(recip.room))
		}
		usr.session.WriteLine(color("red", fmt.Sprintf("You yell, \"%s.\"", msg)))
		usr.session.WriteLine(usr.getPrompt(usr.room))
	case "shout":
		args := strings.Split(cmd, " ")
		msg := ""
		for i := 1; i < len(args); i++ {
			msg = msg + " " + args[i]
		}
		msg = strings.TrimLeft(msg, " ")
		for _, recip := range w.users {
			if recip != usr {
				recip.session.WriteLine(color("blue", fmt.Sprintf("%s shouts, \"%s.\"", usr.name, msg)))
				recip.session.WriteLine(recip.getPrompt(recip.room))
			}
		}
		usr.session.WriteLine(color("blue", fmt.Sprintf("You shout, \"%s.\"", msg)))
		usr.session.WriteLine(usr.getPrompt(usr.room))
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
						usr.session.WriteLine(usr.getPrompt(usr.room))
						return
					}
				}
				usr.session.WriteLine(color("magenta", "Not much to see."))
				usr.session.WriteLine(usr.getPrompt(usr.room))
				return
			case "":
				usr.session.WriteLine(color("magenta", "What were you trying to look at?"))
				usr.session.WriteLine(usr.getPrompt(usr.room))
				return
			default:
				usr.session.WriteLine(color("magenta", "Not much to see."))
				usr.session.WriteLine(usr.getPrompt(usr.room))
				return
			}
		}
	case "help":
		for _, cmd := range w.cmnds {
			usr.session.WriteLine(color("red", cmd.cmnd) + " - " + cmd.desc)
		}
		usr.session.WriteLine(usr.getPrompt(usr.room))
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
			usr.session.WriteLine(usr.getPrompt(usr.room))
			return
		} else {
			switch args[1] {
			case "in", "out", "through", "iI", "o", "t":
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
				return
			default:
				usr.session.WriteLine(color("magenta", fmt.Sprintf("You can't go '%s.'", args[1])))
				usr.session.WriteLine(usr.getPrompt(usr.room))
				return
			}
		}
	case "":
		usr.session.WriteLine(usr.getPrompt(usr.room))
	case "eq", "equip":
		usr.returnEQ()
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
			slot, err := strconv.Atoi(flds[2])
			if err == nil {
				fail := false
				i := &Item{
					name: name,
					desc: desc,
					slot: slot,
				}
				for _, itm := range w.items {
					if strings.ToUpper(itm.name) == strings.ToUpper(i.name) {
						fail = true
						usr.session.WriteLine("Item name exists in world.")
						return
					}
				}
				if fail == false {
					w.items = append(w.items, i)
					usr.session.WriteLine("Added: " + i.name + " - " + i.desc + " - " + fmt.Sprint(i.slot))
				}

			} else {
				fmt.Println(usr.name + " failed item creation.")
				usr.session.WriteLine("Failed.")
			}
		} else {
			usr.session.WriteLine("Not enough arguments to make an item.")
		}
		/*if i := createItem(usr); i != nil && i.name != "" && i.desc != "" {
			for _, itm := range w.items {
				if i.name == itm.name {
					fmt.Println(usr.name + " failed item creation.")
					usr.session.WriteLine("Failed.")
					return
				}
			}
			w.items = append(w.items, i)
		} else {
			fmt.Println(usr.name + " failed item creation.")
			usr.session.WriteLine("Failed.")
		}
		for _, itm := range w.items {
			usr.session.WriteLine(itm.name + " - " + itm.desc + " - " + fmt.Sprint(itm.slot))
		}*/

	case "listitems":
		for _, itm := range w.items {
			usr.session.WriteLine(itm.name + " - " + itm.desc + " - " + fmt.Sprint(itm.slot))
		}
	case "give":
		giveStr := ""
		for i := 1; i < len(args); i++ {
			giveStr = giveStr + " " + args[i]
		}
		if len(args) > 1 {
			fail := true
			for _, itm := range w.items {
				if strings.ToUpper(itm.name) == strings.ToUpper(strings.TrimLeft(giveStr, " ")) {
					usr.char.inv = append(usr.char.inv, itm)
					usr.session.WriteLine("You have received " + itm.name)
					fail = false
				}
			}
			if fail == true {
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
			for _, i := range usr.char.inv {
				usr.session.WriteLine(color("cyan", "    "+i.name))
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
				if strings.ToUpper(i.name[0:adjStr]) == strings.ToUpper(wearStr[0:adjStr]) {
					tryToWear(i, usr)
					fail = false
					return
				}
			}
			if fail == true {
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

			tryToRemove(strings.TrimLeft(remStr, " "), usr)
		} else {
			usr.session.WriteLine("What are you trying to remove?")
		}
	case "slots":
		usr.session.WriteLine(color("magenta", "head    *Item //slot 1"))
		usr.session.WriteLine(color("magenta", "face    *Item //slot 2"))
		usr.session.WriteLine(color("magenta", "neck    *Item //slot 3"))
		usr.session.WriteLine(color("magenta", "about   *Item //slot 4"))
		usr.session.WriteLine(color("magenta", "chest   *Item //slot 5"))
		usr.session.WriteLine(color("magenta", "back    *Item //slot 6"))
		usr.session.WriteLine(color("magenta", "holdL   *Item //slot 7"))
		usr.session.WriteLine(color("magenta", "holdR   *Item //slot 8"))
		usr.session.WriteLine(color("magenta", "waist   *Item //slot 9"))
		usr.session.WriteLine(color("magenta", "legs    *Item //slot 10"))
		usr.session.WriteLine(color("magenta", "feet    *Item //slot 11"))
		usr.session.WriteLine(color("magenta", "arms    *Item //slot 12"))
		usr.session.WriteLine(color("magenta", "wristL  *Item //slot 13"))
		usr.session.WriteLine(color("magenta", "wristR  *Item //slot 14"))
		usr.session.WriteLine(color("magenta", "hands   *Item //slot 15"))
		usr.session.WriteLine(color("magenta", "fingerL *Item //slot 16"))
		usr.session.WriteLine(color("magenta", "fingerR *Item //slot 17"))
		usr.session.WriteLine(usr.getPrompt(usr.room))
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
				if strings.ToUpper(i.name[0:adjLen]) == strings.ToUpper(dropStr[0:adjLen]) {
					fail = false
					usr.room.items = append(usr.room.items, i)
					removeItemFromInventory(i, usr)
					for _, u := range usr.room.users {
						if u != usr {
							u.session.WriteLine(usr.name + " drops a " + color("cyan", i.name) + " on the ground here.")
						} else {
							usr.session.WriteLine("You drop a " + color("cyan", i.name) + " on the ground here.")
						}
					}
					return
				}
			}
			if fail == true {
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
				if strings.ToUpper(itm.name[0:adjLen]) == strings.ToUpper(takeStr[0:adjLen]) {
					fail = false
					removeItemFromRoom(itm, usr.room)
					usr.char.inv = append(usr.char.inv, itm)
					for _, u := range usr.room.users {
						if u != usr {
							u.session.WriteLine(usr.name + " picks up a " + color("cyan", itm.name) + " off the ground here.")
						} else {
							usr.session.WriteLine("You pick up a " + color("cyan", itm.name) + " off the ground here.")
						}
					}
				}
				if fail == false {
					break
				}
			}
			if fail == true {
				usr.session.WriteLine(color("magenta", "You don't see that here. "+takeStr))
			}

		} else {
			usr.session.WriteLine(color("magenta", "What are you trying to take?"))
		}
	case "exa", "examine":
		if len(args) > 1 {
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
				if strings.ToUpper(u.name[0:adjLen]) == strings.ToUpper(exaStr[0:adjLen]) {
					fail = false
					itms := u.eqToItemArray()
					u.session.WriteLine(color("cyan", usr.name) + " looks you over thoroughly.")
					for _, nt := range usr.room.users {
						if nt != u && nt != usr {
							nt.session.WriteLine(color("cyan", usr.name) + " looks over " + u.name + "'s equipment.")
						}
					}
					usr.session.WriteLine(u.name + " is wearing:")
					if len(itms) == 0 {
						usr.session.WriteLine(color("cyan", " ...nothing!"))
					} else {
						for _, itm := range itms {
							usr.session.WriteLine(color("cyan", itm.slotToString()+itm.name))
						}
					}
					usr.session.WriteLine(usr.getPrompt(usr.room))
					return
				}
			}
			for _, i := range usr.room.items {
				adjLen := exaLen
				if len(i.name) < adjLen {
					adjLen = len(i.name)
				}
				if strings.ToUpper(i.name[0:adjLen]) == strings.ToUpper(exaStr[0:adjLen]) {
					fail = false
					usr.session.WriteLine("You take a closer look at a " + color("cyan", i.name) + "...")
					usr.session.WriteLine("    " + i.desc)
					usr.session.WriteLine(usr.getPrompt(usr.room))
					return
				}
			}
			for _, i := range usr.char.inv {
				adjLen := exaLen
				if len(i.name) < adjLen {
					adjLen = len(i.name)
				}
				if strings.ToUpper(i.name[0:adjLen]) == strings.ToUpper(exaStr[0:adjLen]) {
					fail = false
					usr.session.WriteLine("You take a closer look at a " + color("cyan", i.name) + "...")
					usr.session.WriteLine("    " + i.desc)
					usr.session.WriteLine(usr.getPrompt(usr.room))
					return
				}
			}
			eq := usr.eqToItemArray()
			for _, i := range eq {
				adjLen := exaLen
				if len(i.name) < adjLen {
					adjLen = len(i.name)
				}
				if strings.ToUpper(i.name[0:adjLen]) == strings.ToUpper(exaStr[0:adjLen]) {
					fail = false
					usr.session.WriteLine("You take a closer look at a " + color("cyan", i.name) + "...")
					usr.session.WriteLine("    " + i.desc)
					usr.session.WriteLine(usr.getPrompt(usr.room))
					return
				}
			}
			if fail == true {
				usr.session.WriteLine(color("magenta", "You see nothing with that name here."))
			}
			usr.session.WriteLine(usr.getPrompt(usr.room))
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
		usr.session.WriteLine(usr.getPrompt(usr.room))
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

func removeItemFromWorld(i *Item, w *World) {

	for n, item := range w.items {
		if item == i {

			w.items[n] = w.items[len(w.items)-1]
			w.items = w.items[:len(w.items)-1]
		}
	}
}

func (u *User) eqToItemArray() []*Item {
	equip := make([]*Item, 0)

	if eq := u.char.eq.head; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.face; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.neck; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.about; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.chest; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.back; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.holdL; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.holdR; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.waist; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.legs; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.feet; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.arms; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.wristL; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.wristR; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.hands; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.fingerL; eq != nil {
		equip = append(equip, eq)
	}
	if eq := u.char.eq.fingerR; eq != nil {
		equip = append(equip, eq)
	}
	return equip
}

func tryToRemove(itm string, u *User) {
	eq := u.eqToItemArray()
	tItem := &Item{}
	fail := true
	lenStr := len(itm)
	for _, i := range eq {
		adjStr := lenStr
		if len(i.name) < lenStr {
			adjStr = len(i.name)
		}
		if strings.ToUpper(i.name[0:adjStr]) == strings.ToUpper(itm[0:adjStr]) {
			tItem = i
			fail = false
		}
	}
	if fail == true {
		u.session.WriteLine("Could not remove " + itm)
		return
	}
	switch tItem.slot {
	case 1:
		if u.char.eq.head != nil {
			itm := u.char.eq.head
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.head = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your head.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their head.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 1).")
		}
	case 2:
		if u.char.eq.face != nil {
			itm := u.char.eq.face
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.face = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your face.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their face.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 2).")
		}
	case 3:
		if u.char.eq.neck != nil {
			itm := u.char.eq.neck
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.neck = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your neck.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their neck.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 3).")
		}
	case 4:
		if u.char.eq.about != nil {
			itm := u.char.eq.about
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.about = nil
			u.session.WriteLine("You unwrap a " + color("cyan", tItem.name) + " from around you.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " unwraps a " + color("cyan", tItem.name) + " from around them.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 4).")
		}
	case 5:
		if u.char.eq.chest != nil {
			itm := u.char.eq.chest
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.chest = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your chest.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their chest.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 5).")
		}
	case 6:
		if u.char.eq.back != nil {
			itm := u.char.eq.back
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.back = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your back.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their back.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 6).")
		}
	case 7:
		if u.char.eq.holdL != nil {
			itm := u.char.eq.holdL
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.holdL = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your left hand.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their left hand.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 7).")
		}
	case 8:
		if u.char.eq.holdR != nil {
			itm := u.char.eq.holdR
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.holdR = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your right hand.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their right hand.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 8).")
		}
	case 9:
		if u.char.eq.waist != nil {
			itm := u.char.eq.waist
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.waist = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your waist.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their waist.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 9).")
		}
	case 10:
		if u.char.eq.legs != nil {
			itm := u.char.eq.legs
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.legs = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your legs.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their legs.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 10).")
		}
	case 11:
		if u.char.eq.feet != nil {
			itm := u.char.eq.feet
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.feet = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your feet.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their feet.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 11).")
		}
	case 12:
		if u.char.eq.arms != nil {
			itm := u.char.eq.arms
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.arms = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your arms.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their arms.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 12).")
		}
	case 13:
		if u.char.eq.wristL != nil {
			itm := u.char.eq.wristL
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.wristL = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your left wrist.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their left wrist.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 13).")
		}
	case 14:
		if u.char.eq.wristR != nil {
			itm := u.char.eq.wristR
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.wristR = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your right wrist.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their right wrist.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 14).")
		}
	case 15:
		if u.char.eq.hands != nil {
			itm := u.char.eq.hands
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.hands = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your hands.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their hands.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 15).")
		}
	case 16:
		if u.char.eq.fingerL != nil {
			itm := u.char.eq.fingerL
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.fingerL = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your finger.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their finger.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 16).")
		}
	case 17:
		if u.char.eq.fingerR != nil {
			itm := u.char.eq.fingerR
			u.char.inv = append(u.char.inv, itm)
			u.char.eq.fingerR = nil
			u.session.WriteLine("You remove a " + color("cyan", tItem.name) + " from your finger.")
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " removes a " + color("cyan", tItem.name) + " from their finger.")
					usr.session.WriteLine(usr.getPrompt(usr.room))
				}
			}
		} else {
			u.session.WriteLine("This should not be possible - tryToRemove(case 17).")
		}
	default:
		fmt.Println("default case called on tryToRemove : " + u.name)
	}
}

func tryToWear(i *Item, u *User) {
	switch i.slot {
	case 1:
		if u.char.eq.head == nil {
			u.char.eq.head = i
			u.session.WriteLine("You place a " + color("cyan", i.name) + " on your head.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " places a " + color("cyan", i.name) + " on their head.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your head.")
		}
	case 2:
		if u.char.eq.face == nil {
			u.char.eq.face = i
			u.session.WriteLine("You place a " + color("cyan", i.name) + " on your face.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " places a " + color("cyan", i.name) + " on their face.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your face.")
		}
	case 3:
		if u.char.eq.neck == nil {
			u.char.eq.neck = i
			u.session.WriteLine("You place a " + color("cyan", i.name) + " on your neck.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " places a " + color("cyan", i.name) + " on their neck.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your neck.")
		}
	case 4:
		if u.char.eq.about == nil {
			u.char.eq.about = i
			u.session.WriteLine("You wrap a " + color("cyan", i.name) + " around yourself.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " wraps a " + color("cyan", i.name) + " around themself.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your head.")
		}
	case 5:
		if u.char.eq.chest == nil {
			u.char.eq.chest = i
			u.session.WriteLine("You place a " + color("cyan", i.name) + " on your chest.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " places a " + color("cyan", i.name) + " on their chest.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your chest.")
		}
	case 6:
		if u.char.eq.back == nil {
			u.char.eq.back = i
			u.session.WriteLine("You place a " + color("cyan", i.name) + " on your back.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " places a " + color("cyan", i.name) + " on their back.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your back.")
		}
	case 7:
		if u.char.eq.holdL == nil {
			u.char.eq.holdL = i
			u.session.WriteLine("You grab a " + color("cyan", i.name) + " in your left hand.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " grabs a " + color("cyan", i.name) + " in their left hand.")
				}
			}
		} else {
			u.session.WriteLine("You already have something in your left hand.")
		}
	case 8:
		if u.char.eq.holdR == nil {
			u.char.eq.holdR = i
			u.session.WriteLine("You grab a " + color("cyan", i.name) + " in your right hand.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " grabs a " + color("cyan", i.name) + " in their right hand.")
				}
			}
		} else {
			u.session.WriteLine("You already have something in your right hand.")
		}
	case 9:
		if u.char.eq.waist == nil {
			u.char.eq.waist = i
			u.session.WriteLine("You place a " + color("cyan", i.name) + " on your waist.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " places a " + color("cyan", i.name) + " on their waist.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your waist.")
		}
	case 10:
		if u.char.eq.legs == nil {
			u.char.eq.legs = i
			u.session.WriteLine("You place a " + color("cyan", i.name) + " on your legs.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " places a " + color("cyan", i.name) + " on their legs.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your legs.")
		}
	case 11:
		if u.char.eq.feet == nil {
			u.char.eq.feet = i
			u.session.WriteLine("You place a " + color("cyan", i.name) + " on your feet.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " places a " + color("cyan", i.name) + " on their feet.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your feet.")
		}
	case 12:
		if u.char.eq.arms == nil {
			u.char.eq.arms = i
			u.session.WriteLine("You place a " + color("cyan", i.name) + " on your arms.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " places a " + color("cyan", i.name) + " on their arms.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your arms.")
		}
	case 13:
		if u.char.eq.wristL == nil {
			u.char.eq.wristL = i
			u.session.WriteLine("You place a " + color("cyan", i.name) + " on your left wrist.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " places a " + color("cyan", i.name) + " on their left wrist.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your left wrist.")
		}
	case 14:
		if u.char.eq.wristR == nil {
			u.char.eq.wristR = i
			u.session.WriteLine("You place a " + color("cyan", i.name) + " on your right wrist.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " places a " + color("cyan", i.name) + " on their right wrist.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your right wrist.")
		}
	case 15:
		if u.char.eq.hands == nil {
			u.char.eq.hands = i
			u.session.WriteLine("You place a " + color("cyan", i.name) + " on your hands.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " places a " + color("cyan", i.name) + " on their hands.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your hands.")
		}
	case 16:
		if u.char.eq.fingerL == nil {
			u.char.eq.fingerL = i
			u.session.WriteLine("You place a " + color("cyan", i.name) + " on your finger.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " places a " + color("cyan", i.name) + " on their finger.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your left finger.")
		}
	case 17:
		if u.char.eq.fingerR == nil {
			u.char.eq.fingerR = i
			u.session.WriteLine("You place a " + color("cyan", i.name) + " on your finger.")
			removeItemFromInventory(i, u)
			for _, usr := range u.room.users {
				if usr != u {
					usr.session.WriteLine(u.name + " places a " + color("cyan", i.name) + " on their finger.")
				}
			}
		} else {
			u.session.WriteLine("You already have something on your right finger.")
		}
	default:

	}
}

func getSingleInput(u *User, question string) string {

	u.session.WriteLine(question)
	n, err := u.session.conn.Read(u.buf)
	input := ""
	if err == nil {
		input = string(u.buf[0 : n-2])
	}
	if err != nil {
		u.session.WriteLine("Error: " + fmt.Sprint(err))
	}
	u.session.WriteLine("Received Input: " + input)
	return input

}

var (
	mu       sync.Mutex
	tempSlot int
	tempName string
	tempDesc string
)

func getTempSlot() int {
	mu.Lock()
	me := tempSlot
	mu.Unlock()
	return me
}
func setTempSlot(me int) {
	mu.Lock()
	tempSlot = me
	mu.Unlock()
}
func getTempName() string {
	mu.Lock()
	me := tempName
	mu.Unlock()
	return me
}
func setTempName(me string) {
	mu.Lock()
	tempName = me
	mu.Unlock()
}
func getTempDesc() string {
	mu.Lock()
	me := tempDesc
	mu.Unlock()
	return me
}
func setTempDesc(me string) {
	mu.Lock()
	tempDesc = me
	mu.Unlock()
}

func createItem(userCreator *User) *Item {

	i := &Item{}

	intVar, err := strconv.Atoi(getSingleInput(userCreator, "EQ slot?"))
	if err == nil {

		i = &Item{
			name: getSingleInput(userCreator, "Name of Item?"),
			desc: getSingleInput(userCreator, "Desc of Item?"),
			slot: intVar,
		}
	}
	return i
}

func (u *User) initChar() *Character {
	char := &Character{
		name: u.name,
		desc: "ToDo",
		eq: &Equipment{
			head: &Item{
				name: "Red Hat",
				desc: "It's a red hat",
				slot: 1,
			},
			face: &Item{
				name: "Red Facemask",
				desc: "It's a red facemask",
				slot: 2,
			},
			neck: &Item{
				name: "Red Scarf",
				desc: "It's a red scarf",
				slot: 3,
			},
			about: &Item{
				name: "Red Cloak",
				desc: "It's a red cloak",
				slot: 4,
			},
			chest: &Item{
				name: "Red Chestplate",
				desc: "It's a red chestplate",
				slot: 5,
			},
		},
		inv: []*Item{},
	}
	return char
}

func (i *Item) slotToString() string {
	slot := ""
	switch i.slot {

	case 1:
		slot = "        Head: "
	case 2:
		slot = "        Face: "
	case 3:
		slot = "        Neck: "
	case 4:
		slot = "       About: "
	case 5:
		slot = "       Chest: "
	case 6:
		slot = "        Back: "
	case 7:
		slot = "   Held Left: "
	case 8:
		slot = "  Held Right: "
	case 9:
		slot = "       Waist: "
	case 10:
		slot = "        Legs: "
	case 11:
		slot = "        Feet: "
	case 12:
		slot = "        Arms: "
	case 13:
		slot = "  Left Wrist: "
	case 14:
		slot = " Right Wrist: "
	case 15:
		slot = "       Hands: "
	case 16:
		slot = " Left Finger: "
	case 17:
		slot = "Right Finger: "
	default:
		slot = "Invalid Slot: "
	}
	return slot
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
	w.items = []*Item{}
	w.loadEmotes()
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

func startGameLoop(clientInputChannel <-chan ClientInput) {
	//w := &World{}
	for input := range clientInputChannel {

		switch event := input.event.(type) {
		case *InputEvent:
			fmt.Println(fmt.Sprintf("%s: \"%s\"", input.user.name, event.msg))
			executeCmd(event.msg, input.user, input.world)

		case *UserJoinedEvent:
			fmt.Println("User Joined:", input.user.name)
			input.world.users = append(input.world.users, input.user)
			input.user.session.WriteLine(color("cyan", fmt.Sprintf("Welcome %s", input.user.name)))
			input.user.room.addUser(input.user)
			input.user.room.sendText(input.user)
			for _, user := range input.world.users {
				if user != input.user {
					user.session.WriteLine(color("red", fmt.Sprintf("%s has joined!", input.user.name)))
					user.session.WriteLine(user.getPrompt(user.room))
				}
			}
		case *UserLeftEvent:
			un := input.user.name
			fmt.Println("User Left:", un)
			for n, user := range input.world.users {
				if user != input.user {
					user.session.WriteLine(color("red", fmt.Sprintf("%s has left us!", un)))
					user.session.WriteLine(user.getPrompt(user.room))
				}
				if user == input.user {
					removeUserFromRoom(user, user.room, input.world)
					fmt.Println(fmt.Sprintf("%s removed from world index # %s", un, fmt.Sprint(n)))
					input.world.users[n] = input.world.users[len(input.world.users)-1]
					input.world.users = input.world.users[:len(input.world.users)-1]
				}
			}
		}
	}
}

func main() {
	ch := make(chan ClientInput)

	go startGameLoop(ch)

	err := startServer(ch)
	if err != nil {
		log.Fatal(err)
	}

}
