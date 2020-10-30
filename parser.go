package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// lookAtRoom repeats the long form explanation of a room.
func lookAtRoom() {
	fmt.Println(curRoom.LongDesc)
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
	} else if val, ok := curRoom.Items[item]; ok {
		fmt.Println(val.Description)
	} else {
		fmt.Printf("%s not found.\n", item)
	}
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

// moveToRoom takes a requested exit and moves the player there if the exit exists
func moveToRoom(exit string) {
	for _, e := range curRoom.Exits {
		if e == exit { // check that requested exit is valid
			for i, j := range rooms { // find the room it leads to
				if j.Name == exit {
					curRoom = rooms[i] // if found, the exit is the new current room

					if curRoom.Visited == false { // have we been here before?
						curRoom.Visited = true
						rooms[i] = curRoom
						fmt.Println(curRoom.LongDesc)
					} else {
						fmt.Println(curRoom.Description)
					}


					for _, item := range curRoom.Items {
						if item.Discovered == true {
							fmt.Println(item.Description)
						}
					}

					return
				}
			}
		}
	}
	fmt.Printf("%s is not a valid exit\n", exit)
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

	"go Upstairs Hallway" or "go $EXIT" :: proceed through the indicated exit
	to the next room.

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
				lookAtRoom()
			}
		case "go":
			if len(s) > 1 {
				loc := s[1:]
				exit := strings.Join(loc, " ")
				moveToRoom(exit)
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
			/* TODO
			   case "shrink":
			           help()
			   case "whistle":
			           help()
			   case "jump":
			           help()
			   case "attach":
			           help()
			   case "call":
			           help()
			*/
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
