package day_09

import (
	"fmt"
	"io"
	"log"
	"os"
)

func Main() {
	input := readInput("./day_09/sample.txt")
	diskMap := interpretBytesAsNumerals(input)
	diskTsk1 := newStorageFromDiskMap(diskMap)
	//diskTsk1.compactTsk1()
	diskTsk2 := newStorageFromDiskMap(diskMap)
	diskTsk2.updateLookup()
	diskTsk2.compactTsk2()

	fmt.Printf("\n====== DAY 8 ======\n")
	fmt.Printf("%d = Checksum Of File After Compacting\n", diskTsk1.checksum())
	fmt.Printf("%d = Checksum Of File After Compacting Better\n", diskTsk2.checksum())
	//fmt.Println(string(diskTsk2.exportDataToText()))

}

//---------structs declaration---------

type storage struct {
	data   []datapoint
	lookup []block
}

type datapoint struct {
	// typ 0 = free space | typ 1 = file
	// datapoint is more of a datatype struct to distinct stored information
	id  int
	typ byte
}

type block struct {
	start  int
	length int
	data   datapoint
}

func newStorageFromDiskMap(in []byte) storage {
	// storage Constructor
	// worstCase: each space is 9 Bytes long --> 9*len/(
	out := make([]datapoint, 9*len(in))
	//var out = make([]datapoint, 0)
	var idxOut int = 0    // used as iterator
	var dataType byte = 1 // start with file bytes
	var fileIdx int = 0

	for _, segSize := range in {
		var dp datapoint
		switch dataType {
		case 0:
			dp = datapoint{id: 0, typ: dataType}
		case 1:
			dp = datapoint{id: fileIdx, typ: dataType}
			fileIdx++
		default:
			dp = datapoint{id: 0, typ: dataType}
		}
		for i := 0; byte(i) < segSize; i++ {
			out[idxOut] = dp
			idxOut++
		}
		dataType = dataType ^ 1 // alternating pattern of file and space bytes
	}
	return storage{data: out[:idxOut]}
}

func (s *storage) exportDataToText() string {
	var out string
	for _, dataPoint := range s.data {
		switch dataPoint.typ {
		case 0:
			out += ". "
		case 1:
			out += fmt.Sprintf("%02d", dataPoint.id)
		}
	}
	return out
}

//---------methods declaration---------

func (s *storage) compactTsk1() {
	var iterFront int = -1 // assign outOfBound Indicies cause nextIter function increments directly
	var iterBack int = len(s.data)
	for iterFront < iterBack {
		// next emptyStoragePoint typ 0 - iterFront is advanced
		if _, err := s.nextDatapoint(&iterFront, 0); err == io.EOF {
			panic(err)
		}
		// next fileDataPoint typ 1 - iterBack is advanced
		if _, err := s.prevDatapoint(&iterBack, 1); err == io.EOF {
			panic(err)
		}
		if iterFront < iterBack /* Condition because then all spaces have been filled*/ {
			s.swapDatapoints(iterFront, iterBack)
		}
		//fmt.Printf("%s\n", s.exportDataToText())
	}
}

func (s *storage) compactTsk2() {
	// to Tsk2 move file blocks one at a time from the end of the disk to the leftmost free space block
	// move each file exactly once and if there is not enough free space the file is left there
	currBlock, err := s.getBlockAt(len(s.data) - 1)
	if err != nil {
		panic(err)
	}
	for err != io.EOF {
		// the curr block is retrieved at the end of each loop iteration
		//s.updateLookup()
		//fmt.Println(s.exportLookupToText())
		fmt.Println(s.exportDataToText())
		if currBlock.data.typ == 0 {
			// id == 0 -> empty block -> dont process
			currBlock, err = s.nextBlock(currBlock.start, -1)
			continue
		}
		if destAddr, _ := s.allocateSpace(currBlock.length); destAddr != -1 {
			// value '-1' signals there is no free space
			if destAddr+currBlock.length > currBlock.start {
				// destAddr must be left of oldBlock according to tsk2 so skip instead
				currBlock, err = s.nextBlock(currBlock.start, -1)
				continue
			}
			err = s.moveBlock(currBlock, destAddr)
			if err != nil {
				panic(err)
			}
		} else {
			// no free space found -> skip that block
			// TODO process the error of 's.allocateSpace' somehow
			currBlock, err = s.nextBlock(currBlock.start, -1)
			continue
		}
		currBlock, err = s.nextBlock(currBlock.start, -1)
	}
}

func (s *storage) updateLookup() {
	// to Tsk2
	newLookup := make([]block, 0)
	currBlock, err := s.getBlockAt(0)
	if err != nil {
		panic(err)
	}
	for err != io.EOF {
		newLookup = append(newLookup, currBlock)
		startOfNextBlock := currBlock.start + currBlock.length
		currBlock, err = s.getBlockAt(startOfNextBlock)
		if err != nil && err != io.EOF {
			panic(err)
		}
	}
	// TODO: implement lookup validity and consistency check
	s.lookup = newLookup

}

func (s *storage) exportLookupToText() string {
	var out string = ""
	for _, currBlock := range s.lookup {
		out += fmt.Sprintf("%d", currBlock.length)
	}
	return out
}

func (s *storage) nextDatapoint(iter *int, typ byte) (datapoint, error) {
	//scan storage for the next dataPoint of type typ
	// returns empty datapoint when eof
	for {
		*iter = *iter + 1
		if *iter >= len(s.data) || *iter < 0 {
			*iter = len(s.data) // stop iter progress here
			return *new(datapoint), io.EOF
		}
		if s.data[*iter].typ == typ {
			return s.data[*iter], nil
		}
	}
}

func (s *storage) prevDatapoint(iter *int, typ byte) (datapoint, error) {
	// scan storage for the prev dataPoint of type typ
	// returns empty datapoint when eof
	for {
		*iter = *iter - 1
		if *iter >= len(s.data) || *iter < 0 {
			*iter = 0 // stop iter progress here
			return *new(datapoint), io.EOF
		}
		if s.data[*iter].typ == typ {
			return s.data[*iter], nil
		}
	}
}

func (s *storage) allocateSpace(size int) (int, error) {
	// or rather find continuous allocated Space - I don't bother to reserve it
	// size is the number of dataPoints and the search starts at the leftmost(first) dataPoint
	var freeBlockLen int = 0
	for i, dataPoint := range s.data {
		if dataPoint.typ == 0 {
			freeBlockLen++
		} else {
			// check if block long enough was found
			if freeBlockLen >= size {
				var startOfFreeBlock int = i - freeBlockLen
				return startOfFreeBlock, nil
			}
			freeBlockLen = 0 // reset free-counter
		}
	}
	return -1, fmt.Errorf("no free space of size %d", size)
}

func (s *storage) writeBlock(b block, destAddr int) error {
	// TODO return value: number of datapoints written
	if destAddr < 0 || (destAddr+b.length) > len(s.data) {
		return fmt.Errorf("destination address-range out of bound")
	}
	if b.start < 0 || (b.start+b.length) > len(s.data) {
		return fmt.Errorf("source block address range out of bound")
	}
	for i := 0; i < b.length; i++ {
		s.data[destAddr+i] = b.data
	}
	return nil
}

func (s *storage) moveBlock(b block, destAddr int) error {
	err := s.writeBlock(b, destAddr)
	emptyBlock := block{
		start:  b.start,
		length: b.length,
		data:   *new(datapoint),
	}
	if err != nil {
		return err
	}
	err = s.writeBlock(emptyBlock, b.start)
	return err
}

func (s *storage) getBlockAt(start int) (block, error) {
	if start > len(s.data)-1 || start < 0 {
		// on error return empty block
		return *new(block), io.EOF
	}
	var refferenceDP datapoint = s.data[start] // the datapoint the block contains
	var rightBorder int = start
	var leftBorder int = start
	for rightBorder = start; true; rightBorder++ {
		if rightBorder > len(s.data)-1 {
			rightBorder--
			break
		}
		if s.data[rightBorder] != refferenceDP {
			rightBorder-- // since in last iteration the start of the next block was found
			break
		}
	}
	for leftBorder = start; true; leftBorder-- {
		if leftBorder < 0 {
			leftBorder++
			break
		}
		if s.data[leftBorder] != refferenceDP {
			leftBorder++ // since in last iteration the end of the prev block was found
			break
		}
	}
	return block{
		start:  leftBorder,
		length: rightBorder - leftBorder + 1,
		data:   refferenceDP,
	}, nil
}

func (s *storage) nextBlock(pos int, direction int) (block, error) {
	// TODO: restric direction to be either -1 or 1
	var result block = *new(block)
	currBlock, err := s.getBlockAt(pos)
	if err != nil {
		return result, err
	}
	switch direction {
	case -1:
		result, err = s.getBlockAt(currBlock.start - 1)
		if err != nil {
			return result, err
		}
	default:
		result, err = s.getBlockAt(currBlock.start + currBlock.length)
		if err != nil {
			return result, err
		}
	}
	return result, err
}

func (s *storage) swapDatapoints(front int, back int) {
	tmp := s.data[front]
	s.data[front] = s.data[back]
	s.data[back] = tmp
}

func (s *storage) checksum() uint64 {
	var checksum uint64 = 0
	var pointPos int = 0
	for _, dataPoint := range s.data {
		if dataPoint.typ != 0 {
			checksum += uint64(pointPos * dataPoint.id)
		}
		pointPos++
	}
	return checksum
}

//---------functions declaration---------

func interpretBytesAsNumerals(input []byte) []byte {
	var output []byte = make([]byte, len(input))
	for i, byt := range input {
		// assumes only [0-9]
		output[i] = byte(int(rune(byt) - '0'))
	}
	return output
}

func readInput(name string) []byte {
	data, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
