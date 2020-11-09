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
	if val, ok := inventory[item]; ok {
		fmt.Println(val.Description)
	} else if val, ok := curRoom.Items[item]; ok {
		fmt.Println(val.Description)
		if val.ContainsHiddenObject == true {
			if hiddenThing, ok := curRoom.Items[val.HiddenObject]; ok {
				fmt.Println(val.DiscoveryStatement)
				fmt.Println(hiddenThing.Description)
				hiddenThing.Discovered = true
			} else {
				fmt.Println("Oops, we forgot to hide something there.")
			}
			val.ContainsHiddenObject = false
		}
	} else {
		fmt.Printf("%s not found.\n", item)
	}
}

// takeItem place a portable item, which is small enough, into the player's inventory
func takeItem(item string) {
	if val, ok := curRoom.Items[item]; ok {
		if val.Discovered == false {
			fmt.Printf("%s not found.\n", item)
			return
		}
		if val.Portable == true && val.TooBig == false {
			inventory[item] = val
			delete(curRoom.Items, item) // remove item from room after picking it up
			fmt.Printf("You have picked up the %s.\nIt is now in your inventory\n", item)
		} else if val.TooBig == true {
			fmt.Printf("%s is too big to pick up!\nWhy don't you try shrinking it first?\n", item)
		}
	} else {
		fmt.Printf("%s not found.\n", item)
	}
}

// dropItem drops an item in the current room and removes the item from the player's inventory
func dropObject(item string) {
	if val, ok := inventory[item]; ok {
		curRoom.Items[item] = val
		delete(inventory, item)
	} else {
		fmt.Printf("%s not found.\n", item)
	}
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
			if val, ok := rooms[exit]; ok {
				curRoom = val // if found, the exit is the new current room

				if curRoom.Visited == false { // have we been here before?
					curRoom.Visited = true
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
	fmt.Printf("%s is not a valid exit\n", exit)
}

// help prints a subset of verbs the game understands
func help() {
	m := fmt.Sprintf(`
	Here are some of the commands the game understands:

	inventory :: Lists the contents of your inventory.

	mystuff :: see take

	look :: Print the long form explanation of the current room.

	look at <feature or object> :: gives a fictionally interesting explanation of
	the feature or object. You should be able to "look at" objects in your
	inventory, as well. If you describe something in your text descriptions, you
	should be able to "look at" it to examine it.

	"go Upstairs Hallway" or "go $EXIT" or "go to $ROOM" :: proceed through
	the indicated exit to the next room.

	take :: acquire an object, putting it into your inventory.

	grab: see take

	drop:: remove an object from your inventory, dropping it in the current room.

	savegame :: saves the state of the game to a file.

	loadgame :: confirms with the user that this really is desired, then loads
	the game state from the file.

	exit :: save game and then exit.

	quit :: see exit

	help :: Print this message
	`)

	fmt.Println(m)
}

func shrinkObject(item string) {
	if val, ok := curRoom.Items[item]; ok {
		if val.Portable == true {
			if val.TooBig == true {
				fmt.Println("SHRINKING!")
				val.TooBig = false
				fmt.Println("This item is now small enough to collect, pick it up to add it to inventory")
			} else {
				fmt.Println("I don't think that can get any smaller, did you try just picking it up?")
			}
		} else {
			fmt.Println("You can't shrink this, Mom and Dad might notice.")
		}
	}
}

func callTheDog(item string) {
	if val, ok := inventory[item]; ok {
		fmt.Printf("Yeah you have the %s, but you don't know how to use it yet\n", item)
		if val.Name == "dog whistle" {
			fmt.Println("Whistle for the dog")
		}
		if val.Name == "squeeky toy" {
			fmt.Println("Squeek that ball")
		}
	} else {
		fmt.Println("The dog can't hear you")
	}
}

func playerJump(currentRoom string) {
	if currentRoom == "Pantry" || currentRoom == "Upstairs Hallway" || currentRoom == "Basement Lab" {
		fmt.Println("You need to jump here")
	} else {
		fmt.Println("Jump all you want it's not going to do you any good")
	}
}

func callYourParents() {
	fmt.Println("Are you sure you want to do that? You'll be grounded forever")
}

func playGame() {
	openingMessage := fmt.Sprintf(`
It was a bright and sunny afternoon. Everything was going fine.
Your parents were developing new semi-legal technology in their lab, and you were watching them.
They've told you 100 times to not watch them while they work, but what are they going to do?
You're curious. The shrink ray! What a cool invention.
You can take anything and make it...like...smaller.
They've told you not to PLAY with the inventions 101 times, but what are they going to do?
You're curious.
So yeah, they did kick you out of the lab when they left to go run errands, telling you 102 times to not play with the inventions,
but you smuggled that shrink ray out anyway. That's the last thing you remember...
Where are you?
Why don't you try LOOKing around.
	`)

	fmt.Println(openingMessage)

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
					tmp := s[2:]
					item := strings.Join(tmp, " ")
					lookAtItem(item)
				}
			} else {
				lookAtRoom()
			}
		case "go":
			// if the word after "go" is "to" ...
			if len(s) > 1 && s[1] == "to" {
				// ... but no destination is provided
				if len(s) < 3 {
					fmt.Println("Go where?")
					break
				} else { // I want to go there!
					loc := s[2:]
					exit := strings.Join(loc, " ")
					moveToRoom(exit)
				}
			} else if len(s) > 1 { // If player says "go" ...
				loc := s[1:]
				exit := strings.Join(loc, " ")
				moveToRoom(exit)
			} else {
				fmt.Println("Go where?")
			}
		case "goto":
			fmt.Println("Go To Statement Considered Harmful!  https://xkcd.com/292")
		case "take", "grab":
			if len(s) > 1 {
				tmp := s[1:]
				item := strings.Join(tmp, " ")
				takeItem(item)
			} else {
				fmt.Println("Take what?")
			}
		case "drop":
			if len(s) > 1 {
				tmp := s[1:]
				item := strings.Join(tmp, " ")
				dropObject(item)
			} else {
				fmt.Println("Drop what?")
			}
		case "inventory", "mystuff":
			listInventory()

		case "shrink":
			if len(s) > 1 {
				tmp := s[1:]
				item := strings.Join(tmp, " ")
				shrinkObject(item)
			} else {
				fmt.Println("Shrink what?")
			}
		case "whistle":
			callTheDog("dog whistle")
		case "squeek":
			callTheDog("squeeky toy")
		case "jump":
			playerJump(curRoom.Name)
		case "attach":
			help()
		case "call":
			callYourParents()

		case "savegame":
			saveGame()
		case "exit", "quit":
			saveGame()
			return
		case "loadgame":
			if len(s) > 1 {
				f := s[1:]
				g := strings.Join(f, " ")
				loadGame(g)
			} else {
				fmt.Println("Please specify a saved game to load.")
			}
		case "help":
			help()
		default:
			fmt.Println("Not a valid command:", action)
		}
		fmt.Print("> ")
	}
}
