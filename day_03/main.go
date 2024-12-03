package day_03

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func Main() {
	input := readInput()

	fmt.Printf("\n====== DAY 03 ======\n")
	fmt.Printf("%d = Sum of all Multiplications\n", sumAllMultiplications(input))
	fmt.Printf("%d = Sum of all endabled Multiplications\n", calculateSumWithEnablern(input))

}

func calculateSumWithEnablern(input []byte) int {
	var p program
	p = newProgram(extractInstructionRaw(input))
	mem := memory{value: 0, flagExecOp: true}
	p.execute(&mem)

	return mem.value
}

func sumAllMultiplications(input []byte) int {
	sum := 0
	for _, s := range extractMulInstructions(input) {
		result := newMultiplyInstruction(s).calc()
		sum += result
	}
	return sum
}

func readInput() []byte {
	data, err := os.ReadFile("./day_03/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func extractMulInstructions(b []byte) [][]byte {
	// Define the regular expression of multiplication function
	pattern := `mul\(\d{1,3},\d{1,3}\)`
	re := regexp.MustCompile(pattern)
	mulStrings := re.FindAll(b, -1)
	return mulStrings
}

// maths-operation struct
type multiplyInstruction struct {
	x int
	y int
}

type disableInstruction struct {
}

type enableInstruction struct {
}

type noop struct{}

type program struct {
	instructions []instruction
}

type memory struct {
	value      int
	flagExecOp bool
}

type instruction interface {
	execute(m *memory)
}

func (_ disableInstruction) execute(m *memory) {
	m.flagExecOp = false
}

func (_ enableInstruction) execute(m *memory) {
	m.flagExecOp = true
}

func (mul multiplyInstruction) execute(mem *memory) {
	if mem.flagExecOp == true {
		mem.value += mul.x * mul.y
	}
}

func (_ noop) execute(mem *memory) {
	_ = mem.value
	// well, it does not change a thing
}

// constructor for multiplyInstruction out of e.g. mul(12,866) as []byte
func newMultiplyInstruction(b []byte) multiplyInstruction {
	mul := multiplyInstruction{}

	pattern := `\d{1,3}`
	re := regexp.MustCompile(pattern)
	arguments := re.FindAll(b, -1)

	mul.x, _ = strconv.Atoi(string(arguments[0]))
	mul.y, _ = strconv.Atoi(string(arguments[1]))
	return mul
}

// constructor for program
func newProgram(b [][]byte) program {
	p := program{}
	for _, instructionString := range b {
		ins := newInstructionFromRaw(instructionString)
		p.instructions = append(p.instructions, ins)
	}
	return p
}

func (p *program) execute(mem *memory) {
	if p.instructions == nil {
		fmt.Println("recieved a null pointer")
	}
	for _, instruction := range p.instructions {
		instruction.execute(mem)
		//fmt.Printf("Mem %08d flagExecOp %t\n", mem.value, mem.flagExecOp)
	}
}

// used for part one of day_03
func (mul multiplyInstruction) calc() int {
	return mul.x * mul.y
}

func extractInstructionRaw(input []byte) [][]byte {
	// hardcoded to collect (valid) multiply, do and don't instructionStrings
	pattern := `mul\(\d{1,3},\d{1,3}\)|don't\(\)|do\(\)`
	re := regexp.MustCompile(pattern)
	instructions := re.FindAll(input, -1)
	return instructions
}

func newInstructionFromRaw(raw []byte) instruction {
	reMultiply := regexp.MustCompile(`mul\(\d{1,3},\d{1,3}\)`)
	reDo := regexp.MustCompile(`do\(\)`)
	reDont := regexp.MustCompile(`don't\(\)`)

	if reMultiply.Match(raw) {
		return newMultiplyInstruction(raw)
	}
	if reDo.Match(raw) {
		return enableInstruction{}
	}
	if reDont.Match(raw) {
		return disableInstruction{}
	}
	return noop{}
}
