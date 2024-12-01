package main

import (
	"fmt"
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
	file, err := os.Open(`input.txt`)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	buf := make([]byte, file)
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			panic(err)
		}
		//fmt.Print(string(buf[:n]))
	}
	return buf
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
