package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Hello, World!")

	buf := readInput()
	//print(string(buf[:]))

	splitIntoSlices(buf)

}

func readInput() []byte {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func splitIntoSlices(input []byte) ([]uint, []uint) {
	sizeOfLine := 14
	numberOfSpaces := 3
	lengthOfNumber := 5
	for i := 0; i < len(input)/sizeOfLine; i++ {
		startOfLine := i * sizeOfLine
		list1input, _ := strconv.Atoi(string(input[startOfLine : startOfLine+lengthOfNumber]))
		list2input, _ := strconv.Atoi(string(input[startOfLine+lengthOfNumber+numberOfSpaces : startOfLine+lengthOfNumber+numberOfSpaces+lengthOfNumber]))
		println(strconv.Itoa(list1input) + " " + strconv.Itoa(list2input))
	}
	return make([]uint, 1), make([]uint, 1)
}
