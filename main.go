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
	print(string(buf[:]))

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
	for i := 0; i < len(input)/sizeOfLine; i++ {
		startOfLine := i * sizeOfLine
		list1input, _ := strconv.Atoi(string(input[startOfLine : startOfLine+5]))
		println(list1input)
	}
	return make([]uint, 1), make([]uint, 1)
}
