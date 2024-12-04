package day_04

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func Main() {
	input := readInput()

	fmt.Printf("\n====== DAY 04 ======\n")
	fmt.Printf("%d = Number of 'XMAS' in input\n", 0)

	var ridl wordMap
	ridl.readIntoDataFromText(input)
	/*	for i := range ridl.data {
			fmt.Printf("%s\n", string(ridl.data[i]))
		}
	*/
	ridl.findAllWords([]byte("XMAS"))
}

//---------structs declaration---------

type wordMap struct {
	data  [][]byte
	words []word
}

type vector struct {
	xDisplacement int
	yDisplacement int
}

type word struct {
	anchor    coords
	direction vector
	length    int
}

type coords struct {
	// coords should be positive all the time - not implemented
	x int
	y int
}

//---------methods decalration---------

func (w *wordMap) readIntoDataFromText(wordMapAsText []byte) {
	// function to parse input Text into two-dimensional array - better datastructures
	lines := bytes.Split(wordMapAsText, []byte("\n"))
	w.data = lines
}

func (w *wordMap) findAllWords(searchWord []byte) {
	// remove invalid words from the wordList

	// make List of all Start chars
	var wordStartsCoords []coords
	for yCord, Line := range w.data {
		for xCord, char := range Line {
			if char == searchWord[0] {
				wordStartsCoords = append(wordStartsCoords, coords{xCord, yCord})
			}
		}
	}
	// create word for each starting coordinate
}

//---------functions declaration---------

func readInput() []byte {
	data, err := os.ReadFile("./day_04/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	return data
}
