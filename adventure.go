package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

// Some necessary globals
const MaxConnections = 6
const MaxRooms = 15
const MaxObjects = 8

var rooms []Room                         // graph of rooms
var inventory = make(map[string]*Object) // player inventory

// definition of a room
type Room struct {
	Name      string
	RoomType  string
	Feature1  string
	Feature2  string
	LongDesc  string
	ShortDesc string
	Visited   bool
	Count     int                   // number of connections
	Out       [MaxConnections]*Room // outbound connections
	objects   map[string]*Object    // objects in this room
}

// game play object
type Object struct {
	Name      string
	ShortDesc string
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
	// Debug: uncomment to show that JSON data is now a struct in rooms array
	//fmt.Printf("length=%d capacity=%d %v\n", len(rooms), cap(rooms), rooms)
}

func main() {
	loadRooms()

	playGame()

}
