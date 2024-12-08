package day_04

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

/*
Solution one:	scan from each of the 8 directions
Solution two:	branch out of each
*/

func Main() {
	input := readInput()
	var tsk1 wordMap
	tsk1.readIntoDataFromText(input)
	tsk1.findAllWords([]byte("XMAS"))

	var tsk2 wordMap
	tsk2.readIntoDataFromText(input)
	tsk2.findAllWords([]byte("MAS"))
	crosses := findAllXmadeOfMAS(tsk2)
	tsk2.confirmCrossCenters(crosses)

	fmt.Printf("\n====== DAY 04 ======\n")
	fmt.Printf("%d = Number of 'XMAS' in input\n", len(tsk1.words))
	fmt.Printf("%d = Number of X-MAS in output\n", len(crosses))

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

// words can only exist as part of a wordMap (is not jet enforeced unfortunatyly)
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
	for _, coord := range wordStartsCoords {
		directions := []vector{
			{-1, -1},
			{0, -1},
			{1, -1},
			{-1, 0},
			{1, 0},
			{-1, 1},
			{0, 1},
			{1, 1},
		}
		for _, direction := range directions {
			_ = w.addNewWord(coord, direction, len(searchWord)) // does not matter, if a potential word reaches eof
		}
	}

	// remove invalid words from the wordList (all invalid, because out of bound already removed)
	var buff []word
	for _, word := range w.words {
		s, _ := w.wordToString(word)
		pattern := string(searchWord)
		re := regexp.MustCompile(pattern)
		if re.MatchString(s) {
			// only collect valid words
			buff = append(buff, word)
		}
	}
	w.words = buff

	/*for _, word := range w.words {
		s, _ := w.wordToString(word)
		fmt.Printf("%s @[%3d|%3d] \n", s, word.anchor.x, word.anchor.y)
	}*/

}

func (w *wordMap) addNewWord(startCord coords, direction vector, length int) error {
	newWord := word{startCord, direction, length}
	_, eof := w.wordToString(newWord)
	if eof != nil {
		return eof
	}
	w.words = append(w.words, newWord)
	return nil
}

func (w *wordMap) wordToString(wrd word) (string, error) {
	// serves as readWord() function
	buff := bytes.Buffer{}
	for i := 0; i < wrd.length; i++ {
		charCordX := wrd.anchor.x + i*wrd.direction.xDisplacement
		charCordY := wrd.anchor.y + i*wrd.direction.yDisplacement
		if charCordY >= len(w.data) || charCordY < 0 || charCordX >= len(w.data[charCordY]) || charCordX < 0 {
			return buff.String(), io.EOF // Out of Bound is similar to EOF :)
		}
		char := w.data[charCordY][charCordX]
		buff.WriteByte(char)
	}
	return buff.String(), nil
}

func (w *wordMap) wordIsPartOfACross(wrd word, center coords) bool {
	/* a Cross:
	.S.M.S
	..A.A.
	.S.M.S
	..A.A.
	.S.M.S
	*/
	//TODO: if wrd.length is uneven return false
	//TODO: get word center
	if !w.wordExists(wrd) {
		return false
	}
	return w.crossExists(center)

}

func (w *wordMap) crossExists(center coords) bool {
	// true if two of four partners exist
	// hardcoded for MAS ; add word in args to make it universal
	var centerShiftDirs = [4]vector{
		{-1, -1},
		{-1, 1},
		{1, -1},
		{1, 1},
	}
	exitstingWords := 0
	for _, dir := range centerShiftDirs {
		potCrossMemb := word{
			anchor: coords{
				x: center.x + dir.xDisplacement,
				y: center.y + dir.yDisplacement,
			},
			direction: vector{
				xDisplacement: dir.xDisplacement * (-1),
				yDisplacement: dir.yDisplacement * (-1),
			},
			length: 3,
		}
		if w.wordExists(potCrossMemb) {
			exitstingWords += 1
		}
	}
	return exitstingWords == 2
}

func (w *wordMap) wordExists(searchedWord word) bool {
	// see, if the word was stored in the Wordlist
	for _, storedWord := range w.words {
		if storedWord.equalsByValue(searchedWord) {
			return true
		}
	}
	return false
}

// just needed for printf debuging
func (w *wordMap) confirmCrossCenters(centers []coords) {
	for _, coord := range centers {
		//fmt.Printf("%d exists %t\n", coord, w.crossExists(coord))
		coord.x = 1 // dummy funcion to be removed - has no effect on anything
	}
}

func (w word) equalsByValue(wrd word) bool {
	// checks if two word cover the same positions - since
	lengthIsEqual := w.length == wrd.length
	directionIsEqual := w.direction.yDisplacement == wrd.direction.yDisplacement && w.direction.xDisplacement == wrd.direction.xDisplacement
	anchorIsEqual := w.anchor.y == wrd.anchor.y && w.anchor.x == wrd.anchor.x
	return lengthIsEqual && directionIsEqual && anchorIsEqual
}

//---------functions declaration---------

func readInput() []byte {
	data, err := os.ReadFile("./day_04/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func findAllXmadeOfMAS(w wordMap) []coords {
	// start at the center of each 'MAS' where 'A' would be the center of a cross
	// check at each of these 'A's, if there is a cross
	// confirmed crosses are collected in a Set - redundant find are eliminated this way
	crossCordSet := make(map[coords]struct{}) // recommended practise for creating a set

	for _, wrd := range w.words {
		centerDistanceFromStartOfWord := wrd.length / 2 // len 5 should give index 2
		centerOfWord := coords{
			x: wrd.anchor.x + wrd.direction.xDisplacement*centerDistanceFromStartOfWord,
			y: wrd.anchor.y + wrd.direction.yDisplacement*centerDistanceFromStartOfWord,
		}
		if w.crossExists(centerOfWord) {
			crossCordSet[centerOfWord] = struct{}{} // collect all unique centers
		}
	}
	// just parsing the set to an array
	allCrossCords := make([]coords, 0, len(crossCordSet))
	for coord, _ := range crossCordSet {
		allCrossCords = append(allCrossCords, coord)
	}

	return allCrossCords
}
