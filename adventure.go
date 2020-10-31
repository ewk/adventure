package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
)

// Some necessary globals
const MinRooms = 15
const MinItems = 8

var rooms = make(map[string]*Room)     // map of rooms
var inventory = make(map[string]*Item) // player inventory
var curRoom *Room

// definition of a room
type Room struct {
	Name        string
	LongDesc    string
	Description string
	Items       map[string]*Item
	Visited     bool
	Exits       []string // outbound connection room names
}

// struct used for both features and objects
type Item struct {
	Name                 string
	Description          string
	Portable             bool   // False for features, true for objects
	Discovered           bool   // True for features. For some objects this starts as true, for other objects it starts as false.
	ContainsHiddenObject bool   // This is false for all objects. For some features this starts as true.
	DiscoveryStatment    string // If there's a hidden object, this described the connection. "Underneath the couch, you see a cat toy."
	HiddenObject         string // Name of hidden object if there is one
}

// loadRooms reads room definitions from local storage and creates a
// corresponding Room struct. Rooms must be defined as JSON and saved to the
// 'rooms' directory relative to the game's home directory.
func loadRooms() {
	files, e := ioutil.ReadDir("rooms")
	if e != nil {
		log.Fatal(e)
	}

	for _, f := range files {
		matched, _ := regexp.MatchString(`\.json`, f.Name())
		if matched == true {
			roomJson, e := ioutil.ReadFile("rooms/" + f.Name())

			if e != nil {
				log.Fatal(e)
			}

			var r Room
			json.Unmarshal([]byte(roomJson), &r)
			rooms[r.Name] = &r
		}
	}

	// Debug: uncomment to show imported JSON data
	/*
		for key, value := range rooms {
			fmt.Println("Key:", key, "Value:", value)
		}
	*/

	// TODO Disabled for now while development continues.
	// Panic if fewer than 15 rooms are defined.
	/*
		if len(rooms) < MinRooms {
			panic("The game must have at least 15 rooms")
		}
	*/
}

func main() {
	loadRooms()

	curRoom = rooms["Attic"]

	playGame()

}
