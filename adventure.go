package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
)

// Some necessary globals
const MaxRooms = 15
const MaxObjects = 8

var rooms []Room                       // graph of rooms
var inventory = make(map[string]*Item) // player inventory
var curRoom Room

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
			rooms = append(rooms, r)
		}
	}

	// Panic if fewer than 15 rooms are defined.
	// TODO Disabled for now while development continues.
	/*
		if len(rooms) < 15 {
			panic("The game must have at least 15 rooms")
		}
	*/
	// Debug: uncomment to show imported JSON data
	//fmt.Printf("length=%d capacity=%d %v\n", len(rooms), cap(rooms), rooms)
	//for _, i := range rooms {
	//	for _, j := range i.Items {
	//		fmt.Println(j)
	//	}
	//}

	// Debug uncomment to print room exits
	/*
		for _, i := range rooms {
			fmt.Printf("%s:\n", i.Name)
			for _, j := range i.Exits {
				fmt.Printf("\t%s\n", j)
			}
		}
	*/
}

func main() {
	loadRooms()

	// TODO start room must be initialized with Visited = True
	curRoom = rooms[0]

	playGame()

}
