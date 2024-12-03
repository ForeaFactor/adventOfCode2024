package day_01

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func Main() {
	buf := readInput()

	sliceOne, sliceTwo := splitIntoSlices(buf)

	fmt.Printf("\n====== DAY 01 ======\n")
	fmt.Printf("%d = Distance of Lists\n", distanceOfLists(sliceOne, sliceTwo))
	fmt.Printf("%d = Similarity Score of Lists\n", similarityScoreOfLists(sliceTwo, sliceOne))
}

func distanceOfLists(sliceOne []int, sliceTwo []int) int {
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
	data, err := os.ReadFile("./day_01/input.txt")
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

// the similarity score according to the tasks description
func similarityScoreOfLists(sliceOne []int, sliceTwo []int) int {
	score := 0
	for _, val := range sliceOne {
		countOfMatchesToVal := 0

		for _, potentialMatchToVal := range sliceTwo {
			if val == potentialMatchToVal {
				countOfMatchesToVal++
			}
		}

		score += countOfMatchesToVal * val
	}
	return score
}
