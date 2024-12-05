package day_04

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

func Main() {
	input := readInput()
	var ridl wordMap
	ridl.readIntoDataFromText(input)
	ridl.findAllWords([]byte("XMAS"))

	fmt.Printf("\n====== DAY 04 ======\n")
	fmt.Printf("%d = Number of 'XMAS' in input\n", len(ridl.words))

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

	for _, word := range w.words {
		s, _ := w.wordToString(word)
		fmt.Printf("%s @[%3d|%3d] \n", s, word.anchor.x, word.anchor.y)
	}

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

//---------functions declaration---------

func readInput() []byte {
	data, err := os.ReadFile("./day_04/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	return data
}
