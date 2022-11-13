package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	. "sigmars-garden-solver/srcs"
)

func main() {
	if len(os.Args) == 2 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}

		byteValue, _ := ioutil.ReadAll(f)
		var tiles [][]TileType

		if err := json.Unmarshal(byteValue, &tiles); err != nil {
			panic(err)
		}

		board := NewBoard(tiles)

		actions, n, dur := board.Solve()

		fmt.Println("total checks:", n)
		fmt.Println("total duration:", dur)
		fmt.Println(SolutionToString(actions))
	} else {
		fmt.Printf("%s inputs/input1.json\n", os.Args[0])
	}
}

// var tiles = [][]TileType{
// 	{"", "", "", "", "", ""},
// 	{"", "", "", "", "", "", ""},
// 	{"", "", "", "", "", "", "", ""},
// 	{"", "", "", "", "", "", "", "", ""},
// 	{"", "", "", "", "", "", "", "", "", ""},
// 	{"", "", "", "", "", "", "", "", "", "", ""},
// 	{"", "", "", "", "", "", "", "", "", ""},
// 	{"", "", "", "", "", "", "", "", ""},
// 	{"", "", "", "", "", "", "", ""},
// 	{"", "", "", "", "", "", ""},
// 	{"", "", "", "", "", ""},
// }

// var tiles = [][]v1.TileType{
// 	{"cyan", "", "", "", "", "cyan"},
// 	{"light", "orange", "l5", "orange", "orange", "green", ""},
// 	{"blue", "", "", "", "light", "dark", "green", ""},
// 	{"", "orange", "", "", "cyan", "dark", "key", "dark", "white"},
// 	{"", "blue", "blue", "", "white", "cyan", "", "", "", "green"},
// 	{"green", "light", "l1", "key", "green", "l6", "key", "", "", "cyan", "dark"},
// 	{"", "cyan", "orange", "cyan", "blue", "green", "white", "", "blue", ""},
// 	{"", "l4", "", "", "", "green", "blue", "orange", ""},
// 	{"", "light", "", "", "orange", "blue", "l2", ""},
// 	{"", "key", "", "orange", "white", "l3", ""},
// 	{"green", "cyan", "blue", "", "", "key"},
// }
