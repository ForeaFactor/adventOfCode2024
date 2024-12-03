package day_02

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func Main() {
	input := readInput()
	reportList := parseInputToReportData(input)

	fmt.Printf("\n====== DAY 02 ======\n")
	fmt.Printf("%d = Number of safe Reports\n", countSecureReports(reportList))
	fmt.Printf("%d = Number of safe Reports with Problem Dampener\n", countSecureReportsWithProblemDampener(reportList))

}

func readInput() []byte {
	data, err := os.ReadFile("./day_02/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func parseInputToReportData(input []byte) [][]int {
	/*
		only works for input where chars are the size of one byte. ExampleInput
		1 4 5 8 11 12 9
		7 8 9 10 12 15 17 17
	*/
	var reportList [][]int

	// find each of line and extract the line
	startOfLine := 0
	for i, potentialLFchar := range input {
		if potentialLFchar == '\n' {
			endOfLine := i //include the linefeed '\n' in slice selector to help the next for-loop
			reportRawData := input[startOfLine : endOfLine+1]

			// find each number and extract the number
			startOfNumber := 0
			report := make([]int, 0) // report for this iteration is empty
			for i2, potentialSpaceChar := range reportRawData {
				if potentialSpaceChar == ' ' || i2 == len(reportRawData)-1 {
					endOfNumber := i2
					numberRawData := reportRawData[startOfNumber:endOfNumber]
					number, _ := strconv.ParseInt(string(numberRawData), 10, 32)
					//fmt.Printf("%d ", number)
					report = append(report, int(number))
					startOfNumber = endOfNumber + 1 //for next iteration and next number excluding Whitespace
				}
			}

			//fmt.Println()
			reportList = append(reportList, report)
			startOfLine = endOfLine + 1 //for next iteration and next line excluding LF
		}
	}
	return reportList
}

func countSecureReports(reportList [][]int) int {
	var numOfSaveReports = 0

	for _, report := range reportList {
		isSave, _ := validateReport(report)
		if isSave {
			numOfSaveReports++
		}
	}
	//fmt.Printf("%2d ", reportList)
	return numOfSaveReports
}

func validateReport(report []int) (bool, int) {
	// follow the wierd rules in the day02 task, that specify weather as report is 'save'
	// just need to compare all numbers next to each other
	diffAtFront := report[1] - report[0]
	reportGrowthDirection := 1 // 1 OR -1
	if diffAtFront < 0 {
		reportGrowthDirection = -1
	}

	//fmt.Printf("%d %d \n", report, reportGrowthDirection)
	for i := 0; i < len(report)-1; i++ {
		diff := report[i+1] - report[i]
		dirCorDiff := diff * reportGrowthDirection // growth direction corrected diff
		if dirCorDiff <= 0 || dirCorDiff >= 4 {
			return false, i + 1 // return the place (number index in report) of the failure as error-code
		}
	}
	return true, -1
}

func countSecureReportsWithProblemDampener(reportList [][]int) int {
	var numOfSaveReports = 0

	for _, report := range reportList {
		isSave := isSafeWithDampener(report)
		if isSave {
			numOfSaveReports++
		}
		//fmt.Printf("%t %d %d \n", isSave, failedAtIndex, report)
	}
	//fmt.Printf("%2d ", reportList)
	return numOfSaveReports
}

func validateReportWithProblemDampener(report []int) (bool, int) {
	// follow the wierd rules in the day02 task, that specify weather as report is 'save'
	// just need to compare all numbers next to each other
	// 8 6 x4x 4 1 When Error: check if removing before, inside, after changes anything
	// 8     4 4 1 false
	// 8 6     4 1 true
	// 8 6  4    1 true

	valRes, errCode := validateReport(report)
	//	fmt.Printf("%t %d %d \n", valRes, errCode, report)
	if valRes == true {
		return true, -1
	}
	report1 := removeIndexInSlice(report, errCode-1)
	valRes, _ = validateReport(report1)
	//	fmt.Printf("%t %d %d \n", valRes, errCode, report1)
	if valRes == true {
		return true, errCode
	}
	report2 := removeIndexInSlice(report, errCode)
	valRes, _ = validateReport(report2)
	//	fmt.Printf("%t %d %d \n", valRes, errCode, report2)
	if valRes == true {
		return true, errCode
	}
	report3 := removeIndexInSlice(report, errCode+1)
	valRes, _ = validateReport(report3)
	//	fmt.Printf("%t %d %d \n", valRes, errCode, report3)
	if valRes == true {
		return true, errCode
	}
	return false, errCode
}

func removeIndexInSlice(slice []int, index int) []int {
	if index > len(slice)-1 || index < 0 {
		return slice // prevent access outside of slice
	}
	shorterSlice := make([]int, 0, len(slice)-1)
	shorterSlice = append(shorterSlice, slice[:index]...)
	shorterSlice = append(shorterSlice, slice[index+1:]...)

	return shorterSlice
}

func isSafe(report []int) bool {
	isIncreasing := true
	isDecreasing := true

	for i := 0; i < len(report)-1; i++ {
		diff := report[i+1] - report[i]
		if diff < 1 || diff > 3 {
			return false
		}
		if diff > 0 {
			isDecreasing = false
		}
		if diff < 0 {
			isIncreasing = false
		}
	}

	return isIncreasing || isDecreasing
}

// isSafeWithDampener checks if a report becomes safe after removing a single level
func isSafeWithDampener(report []int) bool {
	for i := 0; i < len(report); i++ {
		// Create a copy of the report with one level removed
		modifiedReport := make([]int, 0, len(report)-1)
		modifiedReport = append(modifiedReport, report[:i]...)
		modifiedReport = append(modifiedReport, report[i+1:]...)

		// Check if the modified report is safe
		if isSafe(modifiedReport) {
			return true
		}
	}
	return false
}
