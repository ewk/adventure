package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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
		if val.Discovered == true {
			fmt.Println(val.Description)
		}
		if val.ContainsHiddenObject == true {
			if hiddenThing, ok := curRoom.Items[val.HiddenObject]; ok {
				fmt.Println(val.DiscoveryStatement)
				//fmt.Println(hiddenThing.Description)
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

// takeItem places an object which is small enough into the player's inventory
func takeItem(item string) {
	if val, ok := curRoom.Items[item]; ok {
		if val.Discovered == false {
			fmt.Printf("%s not found.\n", item)
			return
		}
		if val.IsFeature == false && val.TooBig == false {
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
		fmt.Printf("You dropped the %s in the %s.\n", item, curRoom.Name)
	} else {
		fmt.Printf("%s not found.\n", item)
	}
}

// listInventory lists the contents of your inventory.
func listInventory() {
	for key := range inventory {
		fmt.Println(key)
	}
}

// moveToRoom takes a requested exit and moves the player there if the exit exists
func moveToRoom(exit string) {
	b := checkExit() // verify we have items needed to leave
	if !b {
		fmt.Printf("You cannot leave because %s.\n", curRoom.ExitBlock)
		return
	}

	// If the exit requested by a user matches an entry in the list of
	// room aliases, then the room name becomes the requested exit.
	for key, val := range roomAliases {
		matched, _ := regexp.MatchString(key, exit)
		if matched == true {
			exit = val
		}
	}

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

// checkExit verifies the player has the item necessary to exit a room
func checkExit() bool {
	if len(curRoom.ExitItems) == 0 {
		return true
	}

	obj := curRoom.ExitItems[0]
	if _, ok := inventory[obj]; ok {
		return true
	}

	return false
}

// help prints a subset of verbs the game understands
func help() {
	m := fmt.Sprintf(`
	Here are some of the commands the game understands:

	inventory :: Lists the contents of your inventory.

	mystuff :: see inventory

	look :: Print the long form explanation of the current room.

	look at <feature or object> :: gives a fictionally interesting explanation of
	the feature or object. You should be able to "look at" objects in your
	inventory, as well. If you describe something in your text descriptions, you
	should be able to "look at" it to examine it.

	"go Upstairs Hallway" or "go $EXIT" or "go to $ROOM" :: proceed through
	the indicated exit to the next room.

	take :: acquire an object, putting it into your inventory.

	grab :: see take

	drop :: remove an object from your inventory, dropping it in the current room.

	eat :: restore your strength by eating an item

	pull :: see take

	whistle :: With the right item at hand, you can whistle to summon the family pet.

	call :: Call your parents to come get you.

	enter :: Type a secret password into a computer.

	climb :: Climb a mountain, climb the furniture ...

	use :: Make use of an item in your inventory.

	taunt :: Pick a fight!

	jump :: Get vertical!

	slide :: Travel quickly

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
		if val.IsFeature == false {
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
	} else {
		fmt.Println("The dog can't hear you")
	}
}

func playerJump() {
	if curRoom.Name == "Pantry" || curRoom.Name == "Upstairs Hallway" || curRoom.Name == "Basement Lab" {
		fmt.Println("You need to jump here")
	} else {
		fmt.Println("Jump all you want it's not going to do you any good")
	}
}

func callYourParents() {
	fmt.Println("Are you sure you want to do that? You'll be grounded forever")
}

// eatItem searches the player's inventory for an edible item and consumes it
func eatItem(item string) {
	if val, ok := inventory[item]; ok {
		if val.IsEdible {
			fmt.Println("That was delicious! Your strength has been restored.")
			delete(inventory, item)
		} else {
			fmt.Printf("I know you're hangry. But %s is not food!\n", item)
		}
	} else {
		fmt.Printf("%s is not in your backpack.\n", item)
	}
}

func enterThePassword(password string) {
	if _, ok := inventory[password]; ok {
		if curRoom.Name == "Basement Lab" {
			fmt.Println("TAKE the software you need")
			curRoom.Items["software"].Discovered = true

		} else {
			fmt.Println("There's nothing that needs a password here")
		}
	} else {
		fmt.Println("I don't think you know the password")
	}

}

func climbTheDesk() {
	if curRoom.Name == "Basement Lab" {
		curRoom.Items["computer"].Discovered = true
	} else if curRoom.Name == "Large Bedroom" {
		fmt.Println("You can't climb on your parent's desk!")
	} else {
		fmt.Println("There is no desk to climb here")
	}
}

func lookAtEagle() {
	if curRoom.Items["eagle"].Discovered == true {
		if _, ok := inventory["umbrella"]; ok {
			fmt.Println("If you want to use the umbrella to hide from the eagle say: use umbrella")
			fmt.Println("If you want to be taken by the eagle say: taunt eagle")
		} else {
			fmt.Println("The eagle swoops down and picks you up, you manage to wriggle free and drop down the chimney into the master bedroom")
			curRoom = rooms["Large Bedroom"]
			lookAtRoom()
		}
	} else {
		fmt.Println("Hmm...the eagle doesn't seem to be here right now")
	}

}

func useTheUmbrella() {
	if _, ok := inventory["umbrella"]; ok && curRoom.Name == "Yard" {
		fmt.Println("You open the umbrella and are completely hidden from the eagle")
		fmt.Println("Not finding lunch the eagle flies away")
		curRoom.Items["eagle"].Discovered = false
	} else if _, ok := inventory["umbrella"]; ok && curRoom.Name != "Yard" {
		fmt.Println("You can't open the umbrella inside!")
	} else {
		fmt.Println("You don't have an umbrella")
	}
}

func tauntTheEagle() {
	if curRoom.Items["eagle"].Discovered == true {
		curRoom = rooms["Large Bedroom"]
		lookAtRoom()
	} else if curRoom.Items["eagle"].Discovered == false {
		fmt.Println("The eagle has heard your taunts and it has made him mad!")
		curRoom.Items["eagle"].Discovered = true
		curRoom = rooms["Large Bedroom"]
		lookAtRoom()
	}

}

func slideDownJumpIn(userInput []string) {
	if curRoom.Name == "Large Bedroom" || curRoom.Name == "Small Bedroom" {
		if userInput[1] == "fireplace" {
			fmt.Println("GERONIMO!!!!")
			curRoom = rooms["Living Room"]
			lookAtRoom()
		} else if userInput[1] == "laundry" {
			fmt.Println("HERE GOES NOTHING")
			curRoom = rooms["Basement Lab"]
			lookAtRoom()
		}
		if len(userInput) > 2 {
			if userInput[2] == "fireplace" {
				fmt.Println("GERONIMO!!!!")
				curRoom = rooms["Living Room"]
				lookAtRoom()
			} else if userInput[2] == "laundry" {
				fmt.Println("HERE GOES NOTHING")
				curRoom = rooms["Basement Lab"]
				lookAtRoom()
			}
		}
	} else {
		if userInput[0] == "slide" {
			fmt.Println("Sliiiiiide to the left *clap* Sliiiiiide to the right.")
			fmt.Println("You can't remmeber any more of the dance.")
		}
		if userInput[0] == "jump" {
			playerJump()
		}

	}
}

// capInput is a helper function to capitalize case insensitive input
func capInput(input []string) []string {
	for i, w := range input {
		input[i] = strings.Title(strings.ToLower(w))
	}

	return input
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
Why don't you try LOOKing around.`)

	fmt.Println(openingMessage)

	input := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for input.Scan() {
		// split user input at whitespace and match known commands
		action := input.Text()
		s := strings.Fields(strings.ToLower(action))

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
					if curRoom.Name == "Yard" && item == "eagle" {
						lookAtEagle()
					} else {
						lookAtItem(item)
					}
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
					loc = capInput(loc)
					exit := strings.Join(loc, " ")
					moveToRoom(exit)
				}
			} else if len(s) > 1 { // If player says "go" ...
				loc := s[1:]
				loc = capInput(loc)
				exit := strings.Join(loc, " ")
				moveToRoom(exit)
			} else {
				fmt.Println("Go where?")
			}
		case "goto":
			fmt.Println("Go To Statement Considered Harmful!  https://xkcd.com/292")
		case "take", "grab", "pull":
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
		case "call":
			callYourParents()
		case "eat":
			if len(s) > 1 {
				tmp := s[1:]
				item := strings.Join(tmp, " ")
				eatItem(item)
			} else {
				fmt.Println("Eat what?")
			}
		case "enter":
			if len(s) > 1 && s[1] == "password" {
				enterThePassword("password")
			} else {
				fmt.Println("Nothing to enter here")
			}
		case "climb":
			if len(s) > 1 && s[1] == "desk" {
				climbTheDesk()
			} else {
				fmt.Println("There's nothing to climb")
			}
		case "use":
			if len(s) > 1 && s[1] == "umbrella" {
				useTheUmbrella()
			} else {
				fmt.Println("Use what?")
			}
		case "taunt":
			if len(s) > 1 && s[1] == "eagle" {
				tauntTheEagle()
			} else {
				fmt.Println("There's nobody here to taunt but yourself")
			}
		case "slide":
			if len(s) > 1 {
				slideDownJumpIn(s)
			} else {
				fmt.Println("Sliiiiiide to the left *clap* Sliiiiiide to the right.")
				fmt.Println("You can't remmeber and more of the dance.")
			}

		case "jump":
			if len(s) > 1 {
				slideDownJumpIn(s)
			} else {
				playerJump()
			}

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
