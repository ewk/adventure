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
	fmt.Println("\nSome of the things that you see include:")
	for _, item := range curRoom.Items {
		if item.Discovered == true {
			fmt.Printf("     %s\n", item.Name)
		}
	}
}

// lookAtItem prints the description of an object or feature
func lookAtItem(item string) {
	// handle special cases first
	if item == "inventory" {
		listInventory()
		return
	}

	if curRoom.Name == "Yard" && item == "eagle" {
		lookAtEagle()
		return
	}

	if val, ok := inventory[item]; ok { // check player inventory for requested item
		fmt.Println(val.Description)
	} else if val, ok := curRoom.Items[item]; ok { // check the room for requested item
		if val.Discovered == true {
			fmt.Println(val.Description)
		} else {
			fmt.Println("You cannot see that, at least not from here!")
		}

		if val.ContainsHiddenObject == true {
			if hiddenThing, ok := curRoom.Items[val.HiddenObject]; ok {
				fmt.Println(val.DiscoveryStatement)
				hiddenThing.Discovered = true
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
			fmt.Printf("You have picked up the %s.\nIt is now in your INVENTORY.\n", item)
		} else if val.TooBig == true {
			fmt.Printf("%s is too big to pick up!\nWhy don't you try to SHRINK it first?\n", item)
		} else {
			fmt.Println("You cannot pick that up!")
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
	if climbedUp {
		fmt.Printf("Make sure you CLIMB DOWN before you try to go anywhere!\n")
		return
	}
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
				if curRoom.Name == "Attic" && val.Name == "Upstairs Hallway" {
					if _, ok := inventory["paper"]; ok {
						useTheThread()
					} else {
						fmt.Printf("You might be forgetting something important, but you can always come back.\n\n")
						useTheThread()
					}
				}
				if curRoom.Name == "Upstairs Hallway" && val.Name == "Attic" {
					useTheThread()
				}
				if curRoom.Name == "Upstairs Hallway" && val.Name == "Large Bedroom" && val.Visited == false {
					bounceEnterLargeBedroom()
				}
				if curRoom.Name == "Staircase" && val.Name == "Downstairs Hallway" {
					downTheBanister()
				}
				if curRoom.Name == "Downstairs Hallway" && val.Name == "Staircase" {
					climbTheStairs()
				}
				curRoom = val // if found, the exit is the new current room

				if curRoom.Visited == false { // have we been here before?
					curRoom.Visited = true
					fmt.Println(curRoom.LongDesc)
				} else {
					fmt.Println(curRoom.Description)
				}

				fmt.Println("\nSome of the things that you see include:")

				for _, item := range curRoom.Items {
					if item.Discovered == true {
						fmt.Printf("     %s\n", item.Name)
					}
				}

				return
			}
		}
	}
	fmt.Printf("%s is not a valid exit.\n", exit)
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

// haveAllItems checks if player's inventory has all the items needed to win
// the game. Should only be called if player is in the attic.
func haveAllItems() bool {
	for _, item := range winningItems {
		_, ok := inventory[item]
		if !ok {
			return false
		}
	}

	return true
}

// help prints a subset of verbs the game understands
func help() {
	m := fmt.Sprintf(`Here are some of the commands the game understands:

	inventory :: Lists the contents of your inventory.
	mystuff :: See inventory.
	look :: Prints the long form explanation of the current room.
	look at <feature or object> :: Prints the description of an item.
	go :: "go <room>" or "go to <room>" - Proceed through
		the indicated exit to the next room.
	take :: Acquire an object, putting it into your inventory.
	grab :: See take.
	drop :: Remove an object from your inventory, dropping it in
		the current room.
	eat :: Restore your strength by eating an item.
	pull :: See take.
	whistle :: With the right item at hand, you can whistle to
		summon the family pet.
	call :: Call your parents to come and fix things for you.
	enter :: Type a secret password into a computer.
	climb :: Climb a desk. Maybe someday you can climb a mountain. Or even
		climb on the rest of the furniture ...
	use :: Make use of an item in your inventory.
	taunt :: Pick a fight!
	jump :: Get vertical!
	slide :: Travel quickly.
	yank :: See take.
	shrink :: Make a big thing a smaller thing.
	cut :: Cut an item.
	savegame :: Saves the state of the game to a file.
	loadgame :: Confirms that this really is desired, then loads
		the game state from a file.
	exit :: Save game and then exit.
	quit :: See exit.
	help :: Print this message.`)

	fmt.Println(m)
}

func shrinkObject(item string) {
	if _, ok := inventory["shrink ray"]; ok {
		if item == "shrink ray" {
			fmt.Println("You can't shrink the shrink ray")
			return
		} else if val, ok := curRoom.Items[item]; ok {
			if val.IsFeature == false {
				if val.TooBig == true {
					fmt.Println("SHRINKING!")
					val.TooBig = false
					fmt.Println("This item is now small enough to collect. You can TAKE it now.")
				} else {
					fmt.Println("I don't think that can get any smaller. Did you try to just TAKE it?")
				}
			} else {
				fmt.Println("You can't shrink this. Mom and Dad might notice!")
			}
		} else {
			fmt.Printf("There is no '%s' in this room to shrink!\n", item)
		}
	} else {
		fmt.Println("You need the shrink ray to shrink things.")
	}
}

func callTheDog(item string) {
	if _, ok := inventory[item]; ok {
		if curRoom.Name == "Staircase" {
			fmt.Println(curRoom.Items["dog"].DiscoveryStatement)
		} else {
			fmt.Printf("You hear the padding footsteps of your loyal steed.\nHe comes loping into the %s.\n", strings.ToLower(curRoom.Name))
			fmt.Println("You grab onto him and he starts running.")
			fmt.Println("When he finally slows down at the top of the stairs you jump off.\n")
			curRoom = rooms["Staircase"]
		}
	} else {
		fmt.Println("The dog can't hear you without the dog whistle")
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
	if haveAllItems() {
		fmt.Println("But you're so close you just have to get back to the attic!")

	} else {
		gameOver = true
	}

}

// eatItem searches the player's inventory for an edible item and consumes it
func eatItem(item string) {
	if val, ok := inventory[item]; ok {
		if val.IsEdible {
			fmt.Println("That was delicious! Your strength has been restored.")
		} else {
			fmt.Printf("I know you're hangry. But %s is not food!\n", item)
		}
	} else {
		fmt.Printf("%s is not in your backpack.\n", item)
	}
}

// enterThePassword types the secret password into the computer
func enterThePassword() {
	if _, ok := inventory["password"]; ok {
		if curRoom.Name == "Basement Lab" && curRoom.Items["computer"].Discovered {
			fmt.Println("TAKE the SOFTWARE you need.")
			curRoom.Items["software"].Discovered = true
		} else {
			fmt.Println("There's nothing that needs a password here.")
		}
	} else {
		fmt.Println("I don't think you know the password.")
	}
}

func climbStuff(item string) {
	if curRoom.Name == "Basement Lab" && item == "desk" {
		climbedUp = true
		fmt.Println("You climb up the desk and are face to face with the COMPUTER.")
		fmt.Println("It seems locked. Why don't you take a LOOK?")
		curRoom.Items["computer"].Discovered = true
	} else if curRoom.Name == "Large Bedroom" && item == "desk" {
		fmt.Println("You'd better not climb on your parents' desk!")
	} else if curRoom.Name == "Pantry" && item == "paper towels" {
		climbedUp = true
		fmt.Println("From up on the paper towels you can get a better look at the shelves.")
		fmt.Println("There is a box of CORN FLAKES pushed all the way back on one of the shelves.\nWeren't you looking for corn flakes?")
		curRoom.Items["corn flakes"].Discovered = true
	} else if curRoom.Name == "Dining Room" && item == "dining room table" {
		climbedUp = true
		fmt.Println("From on top of the table you can see more.")
		curRoom.Items["candelabra"].Discovered = true
		fmt.Println("There's wax everywhere but it looks like there might still be a bit of\ncandle left. LOOK at the CANDELABRA to investigate.")
	} else if item == "down" {
		climbedUp = false
		fmt.Println("You climb back down to the ground before you get dizzy!")
	} else if _, ok := curRoom.Items[item]; !ok {
			fmt.Printf("%s not found.\n", item)
	} else {
		fmt.Printf("You can't climb on that!\n")
	}
}

func climbTheStairs() {
	fmt.Println("Oof that's a lot of stairs to climb!")
	if _, ok := inventory["dog whistle"]; ok {
		fmt.Println("But you have the dog whistle!")
		callTheDog("dog whistle")
	} else {
		fmt.Println("You scream in frustration and your wailing wakes the dog up.")
		fmt.Println("He takes pity on you and picks you up by the scruff and drops you off at the top of the stairs.")
		fmt.Println("You're drenched and smell terrible now but at least you didn't have to climb those stairs")
	}
}

func cutStuff(item string) {
	if curRoom.Name == "Family Room" && item == "copper wire" || curRoom.Name == "Living Room" && item == "couch stuffing" {
		fmt.Println("snip snip")
		takeItem(item)
	} else if _, ok := curRoom.Items[item]; !ok {
		fmt.Printf("%s not found.\n", item)
	} else {
		fmt.Println("Please don't cut that.")
	}
}

func downTheBanister() {
	if _, ok := inventory["scarf"]; ok {
		fmt.Println("\nYou use the scarf to slide quickly and safely down the banister.\n")
	} else {
		fmt.Println("\nYou try to slide down the banister but your jeans don't slide down easily\nso it's more of a scooch.")
		fmt.Println("After a couple minutes of struggling you're sweaty and have worn a hole\ndown in the seat of your pants.")
		fmt.Println("You fall off the banister halfway down and tumble down the rest of the stairs.")
		fmt.Println("The dog just raises his head and looks at you while you flail helplessly.")
		fmt.Println("You land with another thud, thankfully nothing seems broken.")
		fmt.Println("You should have grabbed that silky scarf.\n")
	}
}

func lookAtEagle() {
	fmt.Printf("You look directly into the eagle's eyes.\nHe has a look on his face screaming 'YOU WANNA FIGHT, BRO?' as he flies towards you.\n\n")
	if _, ok := inventory["umbrella"]; ok {
		fmt.Println("But you have the umbrella! You can use it to hide from the eagle.\nHe'll be distracted by the bright colors.")
		fmt.Println("If you want to use the umbrella to hide from the eagle say: use umbrella")
		fmt.Println("If you want to be taken by the eagle say: taunt eagle")
	} else {
		fmt.Println("The eagle swoops down and picks you up. You can see your whole neighborhood\nfrom up here!\n\nYou manage to wriggle free and drop down the chimney. You climb down\ntowards a bit of sunlight, and exit through a small hole in the chimney\ninto the large bedroom.\n")
		curRoom = rooms["Large Bedroom"]
		lookAtRoom()
	}
}

func useTheThread() {
	if _, ok := inventory["thread"]; ok {
		if curRoom.Name == "Upstairs Hallway" {
			threadUp := fmt.Sprintf(`You throw the thread up like a lasso and it attaches to the bottom of the
ladder to the attic. You free climb up it like the Man in Black from
the Princess Bride on the Cliffs of Insanity.

You look so cool.
`)
			fmt.Println(threadUp)
		} else if curRoom.Name == "Attic" {
			threadDown := fmt.Sprintf(`You tie one end of the thread around your waist and the other around the top
rung of the attic ladder. Here goes nothing!

You leap out of the attic door and the thread acts as a bungee. It catches you
right before you smash into the floor of the upstairs hallway.

As you're hanging, catching your breath, it unravels from the ladder and you
drop with a small thud. You gather up the thread and put it in your backpack.
`)
			fmt.Println(threadDown)
		}
	} else {
		fmt.Println("You don't have the thread.")
	}
}

func bounceEnterLargeBedroom() {
	fmt.Println("The door to the large bedroom is closed and you can't reach it at this size.")
	fmt.Println("You take a running start and hurl yourself at your dad's exercise ball.")
	fmt.Println("You bounce off of it with a loud *VWOMP* and grab onto the door handle.")
	fmt.Println("You're just heavy enough to make the handle turn and the door creaks open.")
	fmt.Println("You drop to the floor and walk right in.\n")
}

func useTheUmbrella() {
	if _, ok := inventory["umbrella"]; ok && curRoom.Name == "Yard" {
		fmt.Println("You open the umbrella and are completely hidden from the eagle.")
		fmt.Println("The bright colors calm him and he no longer wants to fight.\nThe eagle flies away.")
		curRoom.Items["eagle"].Discovered = false
	} else if _, ok := inventory["umbrella"]; ok && curRoom.Name != "Yard" {
		fmt.Println("You can't open the umbrella inside!")
	} else {
		fmt.Println("You don't have an umbrella.")
	}
}

func tauntTheEagle() {
	fmt.Printf("The eagle has heard your taunts and it has made him mad!\n\n")
	fmt.Println("The eagle swoops down and picks you up. You can see your whole neighborhood\nfrom up here!\n\nYou manage to wriggle free and drop down the chimney. You climb down\ntowards a bit of sunlight, and exit through a small hole in the chimney\ninto the large bedroom.\n")
	curRoom = rooms["Large Bedroom"]
	lookAtRoom()
}

func slideDownJumpIn(userInput []string) {
	if curRoom.Name == "Small Bedroom" {
		if len(userInput) > 1 && userInput[1] == "laundry" {
			fmt.Println("HERE GOES NOTHING!")
			fmt.Println("With all of your strength you jump into the gaping opening of the laundry chute.\nDirty clothes and dust bunnies zip past as you gain speed.\nYou bang against the sides of the chute, but it's nothing too damaging.\nFrom the bottom of the chute there's another three foot drop to the laundry basket.\nThat was the farthest three feet of your life!\nThankfully the hamper is full and your landing was soft.\nYou scamper out of the basket, throwing clothes everywhere in the process.\n")
			curRoom = rooms["Basement Lab"]
			lookAtRoom()
		}

		if len(userInput) > 2 && userInput[2] == "laundry" {
			fmt.Println("HERE GOES NOTHING!")
			fmt.Println("With all of your strength you jump into the gaping opening of the laundry chute.\nDirty clothes and dust bunnies zip past as you gain speed.\nYou bang against the sides of the chute, but it's nothing too damaging.\nFrom the bottom of the chute there's another three foot drop to the laundry basket.\nThat was the farthest three feet of your life!\nThankfully the hamper is full and your landing was soft.\nYou scamper out of the basket, throwing clothes everywhere in the process.\n")
			curRoom = rooms["Basement Lab"]
			lookAtRoom()
		}
	} else {
		if userInput[0] == "slide" {
			fmt.Println("Sliiiiiide to the left *clap* Sliiiiiide to the right.")
			fmt.Println("You can't remember any more of the dance.")
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
Your parents were developing new semi-legal technology in their lab,
and you were watching them. They've told you 100 times to not watch them
while they work, but what are they going to do? You're curious.

The shrink ray! What a cool invention. Now anything can be made smaller!
They've told you not to play with the inventions 101 times, but what are they
going to do? You're curious.

So yeah, they did kick you out of the lab when they left to go run errands,
telling you 102 times to not touch anything, but you smuggled the
shrink ray out anyway.

That's the last thing you remember. You open your eyes and seem to be in a
giant cavern. Everything is so big! Wait...you're so small!

Where are you? How will you fix this? You still have your (tiny!) cell phone
in your pocket. Should you CALL your parents? No way, they'd save you but
you'd be in sooooo much trouble.

Is there anywhere you could GO TO? Is there anything you could TAKE to
help you? HELP! Why don't you try to LOOK around?`)

	fmt.Println(openingMessage)

	input := bufio.NewScanner(os.Stdin)
	fmt.Print("\n> ")

	for input.Scan() {
		// split user input at whitespace and match known commands
		action := input.Text()
		action = strings.ToLower(action)
		fmt.Println("")

		// accept just the room name as input
		r := strings.Title(action)
		if _, ok := rooms[r]; ok {
			moveToRoom(r)

			// This is repetitive, but we must also check if the player returned
			// to the attic without using the verb 'go'
			if curRoom.Name == "Attic" {
				gameOver = haveAllItems()
			}

			if gameOver {
				break
			} else {
				fmt.Print("\n> ")
				continue
			}
		}

		s := strings.Fields(action)

		if cap(s) == 0 {
			fmt.Print("\n> ")
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
		case "take", "grab", "pull", "yank":
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
			enterThePassword()
		case "climb":
			if len(s) > 1 {
				tmp := s[1:]
				item := strings.Join(tmp, " ")
				climbStuff(item)
			} else {
				fmt.Println("Climb what? The corporate ladder?")
			}
		case "use":
			if len(s) > 1 {
				if s[1] == "umbrella" {
					useTheUmbrella()
				} else if s[1] == "whistle" {
					callTheDog("dog whistle")
				} else if len(s) > 2 && s[1] == "dog" {
					callTheDog("dog whistle")
				} else {
					fmt.Println("I don't know how to USE that, can you use a more specific action?")
				}
			} else {
				fmt.Println("Use what?")
			}
		case "taunt":
			if len(s) > 1 && s[1] == "eagle" {
				tauntTheEagle()
			} else {
				fmt.Println("There's nobody here to taunt but yourself.")
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
		case "cut":
			if len(s) > 1 {
				tmp := s[1:]
				item := strings.Join(tmp, " ")
				cutStuff(item)
			} else {
				fmt.Println("Cut what?")
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

		if curRoom.Name == "Attic" && haveAllItems() {
			gameOver = true
		}

		if gameOver {
			break
		}
		fmt.Print("\n> ")
	}
	if gameOver && haveAllItems() {
		winningMessage := fmt.Sprintf(`
You dump out your backpack and everything tumbles onto the carpet.
The shampoo, the dirty socks, the candle, the couch stuffing, the copper
wire, the cornflakes, the screw, the can, the sand, and the software and
everything else you've picked up all day. Now what?

You grab the notebook and flip through it feverishly looking for
instructions. EUREKA! You've found them.

You pop the back open on the shrink ray and start dumping stuff in,
being careful to shake it frequently, just like the notebook says. This is
one time you will be following instructions!

The shrink ray starts to vibrate and buzz and you see it start to glow purple.
You plug the software into it and watching it whir to life. You grab it 
-- it's getting really hot now -- and rush over the scorched mirror.

Here goes nothing! You think you hear a car door slam in the driveway.

WHABAM! The shrink ray explodes and your ears pop. You feel dizzy, and heavy.

"What was that?"" You hear your dad yell from a distance.

"Nooooooothing" you hear yourself say. You sound louder now. Woozily you
sit up and look in the mirror. You're back to normal size!
You got away with it! Everything is going to be fine!

"WHAT IS ALL THIS MESS?!"
		
Uh-oh.
			
ROLL CREDITS`)

		fmt.Println(winningMessage)
	} else {
		losingMessage := fmt.Sprintf(`You've given up. You can't stand to be this tiny any longer! You call your
parents who race home from the store. They start lecturing you as they collect
items from around the house. They had a spare shrink ray the whole time!
They point it at you and you hear a loud hiss and buzz and your ears pop.

A purple light surrounds you as you grow back to normal size. What a relief!
Until your mother grabs you by the ear and throws you in your room.
You hear the door lock from the outside. You are grounded for eternity.

GAME OVER`)

		fmt.Println(losingMessage)
	}

}
