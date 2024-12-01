package main

import (
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	buf := readInput()

	sliceOne, sliceTwo := splitIntoSlices(buf)

	println(calcDistance(sliceOne, sliceTwo))
}

func calcDistance(sliceOne []int, sliceTwo []int) int {
	sort.Ints(sliceOne)
	sort.Ints(sliceTwo)

	var sum int
	sum = 0
	for i, value := range sliceOne {
		sum += absDiffInt(value, sliceTwo[i])
	}
	return sum
}

func readInput() []byte {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func splitIntoSlices(input []byte) ([]int, []int) {
	sizeOfLine := 14
	numberOfSpaces := 3
	lengthOfNumber := 5
	numberOfLines := len(input) / sizeOfLine

	sliceOne := make([]int, 0)
	sliceTwo := make([]int, 0)
	for i := 0; i < numberOfLines; i++ {
		startOfLine := i * sizeOfLine
		list1input, _ := strconv.Atoi(string(input[startOfLine : startOfLine+lengthOfNumber]))
		list2input, _ := strconv.Atoi(string(input[startOfLine+lengthOfNumber+numberOfSpaces : startOfLine+lengthOfNumber+numberOfSpaces+lengthOfNumber]))
		//println(strconv.Itoa(list1input) + " " + strconv.Itoa(list2input))

		sliceOne = append(sliceOne, list1input)
		sliceTwo = append(sliceTwo, list2input)
	}

	return sliceOne, sliceTwo
}

// absDiffInt Source: https://stackoverflow.com/questions/57648933/why-doesnt-go-have-a-function-to-calculate-the-absolute-value-of-integers
func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}
