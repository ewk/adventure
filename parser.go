package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// lookAtRoom repeats the long form explanation of a room.
func lookAtRoom(cur Room) {
	fmt.Println(cur.LongDesc)
}

// lookAtItem prints the description of an object or feature
func lookAtItem(item string) {
	// This is a fairly generic function. It needs to not only describe
	// an inventory object, but also any feature that has a text description

	// TODO updates for new Item approach
	// The item's description is printed.
	// If ContainsHiddenObject == true:
	//    print DiscoveryStatement
	//    set ContainsHiddenObject = false
	//    get the HiddenObject name
	//    look up that name in the room's Items map
	//    print that object's Description
	//    set that object's Discovered bool to true

	if val, ok := inventory[item]; ok {
		fmt.Println(val.Description)
	} else {
		fmt.Printf("%s not in inventory\n", item)
	}
	// TODO same check for rooms
}

// takeObject puts an item into the player's inventory
func takeItem(itemName string) {
	// TODO add object to player inventory and remove it from room items hash
	// inventory[itemName] = &Item{Name: itemName}
	// can't take item if portable == false
}

// dropObject removes an item from the player's inventory
func dropObject(itemName string) {
	delete(inventory, itemName)
	// TODO save object in room where it was dropped and print output
}

// listInventory lists the contents of your inventory.
func listInventory() {
	fmt.Printf("length=%d %v\n", len(inventory), inventory)

	for key, value := range inventory {
		fmt.Println("Key:", key, "Value:", value)
	}
}

// help prints a subset of verbs the game understands
func help() {
	m := fmt.Sprintf(`
	Here are some of the commands the game understands:

	inventory :: Lists the contents of your inventory.

	look :: Print the long form explanation of the current room.

	look at <feature or object> :: gives a fictionally interesting explanation of
	the feature or object. You should be able to "look at" objects in your
	inventory, as well. If you describe something in your text descriptions, you
	should be able to "look at" it to examine it.

	"go north" OR "north" OR "go dank-smelling staircase" OR "dank-smelling
	staircase" :: proceed through the indicated exit to the next room (note that ALL
	FOUR of these forms of movement are required, and thus require you to describe
	the exits appropriately). You might also decide to implement other room-travel
	verbs such as "jump north" as appropriate.

	take :: acquire an object, putting it into your inventory.

	drop:: remove an object from your inventory, dropping it in the current room.

	savegame :: saves the state of the game to a file.

	loadgame :: confirms with the user that this really is desired, then loads
	the game state from the file.

	help :: Print this message
	`)

	fmt.Println(m)
}

func playGame() {
	// TODO remove dummy data
	curRoom := rooms[0]
	inventory["spoon"] = &Item{Name: "spoon", Description: "A utensil"}
	inventory["candle"] = &Item{Name: "candle", Description: "To light the way"}

	input := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for input.Scan() {
		// split user input at whitespace and match known commands
		action := input.Text()
		s := strings.Fields(action)

		if cap(s) == 0 {
			continue
		}

		switch s[0] {
		case "look":
			if len(s) > 1 && s[1] == "at" {
				if len(s) < 3 {
					fmt.Println("What would you like to look at?")
					break
				} else {
					item := s[2]
					lookAtItem(item)
				}
			} else {
				lookAtRoom(curRoom)
			}
		case "north":
			fmt.Println("You said \"north\".")
			// TODO use $ROOM.exit
		case "south":
			fmt.Println("You said \"south\".")
			// TODO use $ROOM.exit
		case "east":
			fmt.Println("You said \"east\".")
			// TODO use $ROOM.exit
		case "west":
			fmt.Println("You said \"west\".")
			// TODO use $ROOM.exit
		case "go":
			if len(s) > 1 {
				d := s[1]
				fmt.Println("You said \"go\"", d)
				// TODO use $ROOM.exit
			} else {
				fmt.Println("Go where?")
			}
		case "take":
			if len(s) > 1 {
				item := s[1]
				fmt.Println("You said \"take\"", item)
				//TODO takeObject(item)
			} else {
				fmt.Println("Take what?")
			}
		case "drop":
			if len(s) > 1 {
				item := s[1]
				dropObject(item)
			} else {
				fmt.Println("Drop what?")
			}
		case "inventory":
			listInventory()
		case "savegame":
			fmt.Println("You said \"savegame\".")
			// TODO implement save state
		case "loadgame":
			fmt.Println("You said \"loadgame\".")
			// TODO implement load state
		case "help":
			help()
		default:
			fmt.Println("Not a valid command:", action)
		}
		fmt.Print("> ")
	}
}
