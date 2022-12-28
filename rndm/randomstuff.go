package rndm

type Room struct {
	Name        string
	Description string
	Mobiles     []Mobile
	Items       []Item
	Exits       map[string]*Room
	Objects     []Object
}

kitchen := Room{
	Name:        "Kitchen",
	Description: "The kitchen is a cluttered and cramped space, with pots and pans hanging from the ceiling and shelves lined with dusty old jars. A rickety old table sits in the center of the room, with a few broken chairs scattered around it.",
	Mobiles:     []Mobile{ghostlyChef},
	Items:       []Item{},
	Exits: map[string]*Room{
		"South": &entryway,
		"East":  &diningRoom,
	},
	Objects: []Object{},
}

// Define a player struct with a field for the current room
type Player struct {
	CurrentRoom *Room
}

// Define a function to move the player to a new room
func (p *Player) Move(direction string) {
	// Check if the current room has an exit in the specified direction
	nextRoom, ok := p.CurrentRoom.Exits[direction]
	if !ok {
		fmt.Println("You can't go that way.")
		return
	}

	// Update the player's current room
	p.CurrentRoom = nextRoom

	// Print the description of the new room
	fmt.Println(p.CurrentRoom.Description)
}

// Example usage:
player := Player{CurrentRoom: &kitchen}
player.Move("East") // Move the player to the dining room

type Mobile struct {
	Name        string
	Description string
}

// Create a new instance of the Mobile struct
ghostlyChef := Mobile{
	Name:        "Ghostly Chef",
	Description: "A ghostly chef, who haunts the kitchen and is always muttering to himself as he stirs pots of thin air.",
}

kitchen.Mobiles = append(kitchen.Mobiles, ghostlyChef)

fmt.Println(ghostlyChef.Name) // prints "Ghostly Chef"
fmt.Println(ghostlyChef.Description) // prints "A ghostly chef, who haunts the kitchen and is always muttering to himself as he stirs pots of thin air."


type Player struct {
	Name        string
	CurrentRoom *Room
	Inventory   []Item
}

func (p *Player) Move(direction string) {
	// Check if the current room has an exit in the specified direction
	nextRoom, ok := p.CurrentRoom.Exits[direction]
	if !ok {
		fmt.Println("You can't go that way.")
		return
	}

	// Update the player's current room
	p.CurrentRoom = nextRoom

	// Print the description of the new room
	fmt.Println(p.CurrentRoom.Description)
}

func (p *Player) PickUpItem(itemName string) {
	// Check if the item is present in the current room
	var itemToPickUp *Item
	for _, item := range p.CurrentRoom.Items {
		if item.Name == itemName {
			itemToPickUp = &item
			break
		}
	}
	if itemToPickUp == nil {
		fmt.Println("There is no such item here.")
		return
	}

	// Remove the item from the room and add it to the player's inventory
	p.CurrentRoom.Items = removeFromSlice(p.CurrentRoom.Items, itemToPickUp)
	p.Inventory = append(p.Inventory, *itemToPickUp)
	fmt.Println("You pick up the " + itemName + ".")
}

func (p *Player) UseItem(itemName string) {
	// Check if the item is present in the player's inventory
	var itemToUse *Item
	for _, item := range p.Inventory {
		if item.Name == itemName {
			itemToUse = &item
			break
		}
	}
	if itemToUse == nil {
		fmt.Println("You don't have that item.")
		return
	}

	// Use the item
	fmt.Println("You use the " + itemName + ".")
}

func (p *Player) TalkToMobile(mobileName string) {
	// Check if the mobile is present in the current room
	var mobileToTalkTo *Mobile
	for _, mobile := range p.CurrentRoom.Mobiles {
		if mobile.Name == mobileName {
			mobileToTalkTo = &mobile
			break
		}
	}
	if mobileToTalkTo == nil {
		fmt.Println("There is no such mobile here.")
		return
	}

	// Talk to the mobile
	fmt.Println("You talk to the " + mobileName + ".")
}

func (p *Player) InteractWithObject(objectName string) {
	// Check if the object is present in the current room
	var objectToInteractWith *Object
	for _, object := range p.CurrentRoom.Objects {
		if object.Name == objectName {
			objectToInteractWith = &object
			break
		}
	}
	if objectToInteractWith == nil {
		fmt.Println("There is no such object here.")
		return
	}
	fmt.Println("You interact with the " + objectName + ".")
}

// Helper function to remove an item from a slice
func removeFromSlice(slice []Item, item *Item) []Item {
	for i, v := range slice {
		if &v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func handleCommand(p *Player, command string) {
	// Parse the command and execute the appropriate action
	switch command {
	case "north":
		p.Move("North")
	case "south":
		p.Move("South")
	case "east":
		p.Move("East")
	case "west":
		p.Move("West")
	case "pick up":
		p.PickUpItem("itemName")
	case "use":
		p.UseItem("itemName")
	case "talk to":
		p.TalkToMobile("mobileName")
	case "interact with":
		p.InteractWithObject("objectName")
	default:
		fmt.Println("Invalid command.")
	}
}

//command line stuff import "os"
func main() {
  
    // The first argument
    // is always program name
    myProgramName := os.Args[0]
  
    // this will take 4
    // command line arguments
    cmdArgs := os.Args[4]
  
    // getting the arguments
    // with normal indexing
    gettingArgs := os.Args[2]
  
    toGetAllArgs := os.Args[1:]
  
    // it will display
    // the program name
    fmt.Println(myProgramName)
      
    fmt.Println(cmdArgs)
      
    fmt.Println(gettingArgs)
      
    fmt.Println(toGetAllArgs)
}
// iterate over struct types and values
import (
	"fmt"
	"reflect"
)

type Person struct {
	Name   string
	Age    int
	Gender string
	Single bool
}

func main() {
	ubay := Person{
		Name:   "John",
		Gender: "Female",
		Age:    17,
		Single: false,
	}
	values := reflect.ValueOf(ubay)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		fmt.Println(types.Field(i).Index[0], types.Field(i).Name, values.Field(i))
	}
}

//Does not like values as item pointers
/*func (u *User) listEQ() {
	e := &Equipment{
		head:    u.char.eq.head,
		face:    u.char.eq.face,
		neck:    u.char.eq.neck,
		about:   u.char.eq.about,
		chest:   u.char.eq.chest,
		back:    u.char.eq.back,
		holdL:   u.char.eq.holdL,
		holdR:   u.char.eq.holdR,
		waist:   u.char.eq.waist,
		legs:    u.char.eq.legs,
		feet:    u.char.eq.feet,
		arms:    u.char.eq.arms,
		wristL:  u.char.eq.wristL,
		wristR:  u.char.eq.wristR,
		hands:   u.char.eq.hands,
		fingerL: u.char.eq.fingerL,
		fingerR: u.char.eq.fingerR,
	}
	values := reflect.ValueOf(e)
	types := values.Type()

	for i := 0; i < values.NumField(); i++ {

		fmt.Println(types.Field(i).Index[0], types.Field(i).Name, values.Field(i))
	}
}*/

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
//blocks input from other users. Need to learn how to suspend main user thread to put them on a temporary one to use functions like this
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

// find successfully iterates through the equipment struct and can return an item, however item pointers in equipment struct need to be removed so that its fields are addressable for
// the set function of field to work which causes problems as everything is designed up to this point around pointers

func findAndRemoveItem(equipment Equipment, itemName string, u *User) *Item {
	v := reflect.ValueOf(equipment)
	lenStr := len(itemName)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr {
			item := field.Interface().(*Item)
			if item != nil {
				if len(item.Name) < lenStr {
					lenStr = len(item.Name)
				}
				if item != nil && strings.EqualFold(item.Name[0:lenStr], itemName[0:lenStr]) {
					fmt.Printf("reflect found %s from equipment struct.\r\n", item.Name)
					field.Set(reflect.ValueOf(nil))
					u.char.inv = append(u.char.inv, item)
					return item
				}
			}
		}
	}
	fmt.Printf("reflect returned nil. " + itemName + "\r\n")
	return nil
}


func EquipItem(equipment *Equipment, item *Item) error {
	v := reflect.ValueOf(equipment).Elem()
	field := v.FieldByName(item.Slot)
	if !field.IsValid() || field.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid slot: %s", item.Slot)
	}
	field.Set(reflect.ValueOf(item))
	return nil
}

func RemoveItem(equipment *Equipment, item *Item) error {
	v := reflect.ValueOf(equipment).Elem()
	field := v.FieldByName(item.Slot)
	if !field.IsValid() || field.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid slot: %s", item.Slot)
	}
	currentItem := field.Interface().(*Item)
	if currentItem != item {
		return fmt.Errorf("item not equipped in slot: %s", item.Slot)
	}
	field.Set(reflect.ValueOf(nil))
	item.Slot = ""
	return nil
}

type ItemLocation struct {
    ID       string
    Location *Location
}

type Location struct {
    Name string
}

type ItemHashTable struct {
    table []*ItemLocation
    size  int
}

// Hash function to map an item ID to an index in the underlying array
func (h *ItemHashTable) hash(key string) int {
    hash := 0
    for i := 0; i < len(key); i++ {
        hash += int(key[i])
    }
    return hash % h.size
}

// Method to insert an item location into the hash table
func (h *ItemHashTable) Insert(key string, value *Location) {
    index := h.hash(key)
    h.table[index] = &ItemLocation{ID: key, Location: value}
}

// Method to look up the location of an item by its ID
func (h *ItemHashTable) Lookup(key string) (*Location, bool) {
    index := h.hash(key)
    itemLocation := h.table[index]
    if itemLocation != nil && itemLocation.ID == key {
        return itemLocation.Location, true
    }
    return nil, false
}

// Create a new hash table with a size of 10
itemLocations := &ItemHashTable{make([]*ItemLocation, 10), 10}

// Insert some item locations into the hash table
itemLocations.Insert("item1", &Location{Name: "Room 1"})
itemLocations.Insert("item2", &Location{Name: "Room 2"})
itemLocations.Insert("item3", &Location{Name: "Room 2"})
itemLocations.Insert("item4", &Location{Name: "Room 3"})

// Look up the location of an item by its ID
itemID := "item2"
location, ok := itemLocations.Lookup(itemID)
if ok {
    fmt.Println(location.Name)  // Output: "Room 2"
} else {
    fmt.Println("Item not found")
}

func GetItemsInLocation(location *Location, itemLocations *ItemHashTable) []*ItemLocation {
    items := []*ItemLocation{}
    for i := 0; i < itemLocations.size; i++ {
        itemLocation := itemLocations.table[i]
        if itemLocation != nil && itemLocation.Location.Name == location.Name {
            items = append(items, itemLocation)
        }
    }
    return items
}

room2 := &Location{Name: "Room 2"}
itemsInRoom2 := GetItemsInLocation(room2, itemLocations)

for _, item := range itemsInRoom2 {
    fmt.Println(item.ID)
}

type LocationIndex struct {
    index map[string][]string
}

// Method to add an item to the index
func (l *LocationIndex) AddItem(location string, itemID string) {
    l.index[location] = append(l.index[location], itemID)
}

// Method to retrieve all of the items in a specific location
func (l *LocationIndex) GetItemsInLocation(location string) []string {
    return l.index[location]
}